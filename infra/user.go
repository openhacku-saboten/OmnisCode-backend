package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-gorp/gorp"
	"github.com/go-sql-driver/mysql"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

var _ repository.User = (*UserRepository)(nil)

// UserRepository ユーザー情報の永続化と再構成のためのリポジトリです
type UserRepository struct {
	dbMap *gorp.DbMap
}

// NewUserRepository はユーザー情報のリポジトリのポインタを生成する関数です
func NewUserRepository(dbMap *gorp.DbMap) *UserRepository {
	dbMap.AddTableWithName(UserDTO{}, "users").SetKeys(false, "ID")
	return &UserRepository{dbMap: dbMap}
}

// FindByID は該当IDのユーザーの情報をDBから取得して返す
func (r *UserRepository) FindByID(ctx context.Context, uid string) (user *entity.User, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var userDTO UserDTO
		err = r.dbMap.SelectOne(&userDTO, "SELECT * FROM users WHERE id = ?", uid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, entity.NewErrorNotFound("user")
			}
			return nil, err
		}
		user = entity.NewUser(
			userDTO.ID,
			userDTO.Name,
			userDTO.Profile,
			userDTO.TwitterID,
			"",
		)
		return
	}
}

// Insert は該当ユーザーをDBに保存する
func (r *UserRepository) Insert(ctx context.Context, user *entity.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		userDTO := &UserDTO{
			ID:        user.ID,
			Name:      user.Name,
			Profile:   user.Profile,
			TwitterID: user.TwitterID,
		}

		if err := r.dbMap.Insert(userDTO); err != nil {
			if sqlerr, ok := err.(*mysql.MySQLError); ok {
				// userIDが重複したときのエラー
				if sqlerr.Number == mysqlerr.ER_DUP_ENTRY && strings.Contains(sqlerr.Message, "users.PRIMARY") {
					return entity.ErrDuplicatedUser
				}
				// twitterIDが重複したときのエラー
				if sqlerr.Number == mysqlerr.ER_DUP_ENTRY && strings.Contains(sqlerr.Message, "twitter_id") {
					return entity.ErrDuplicatedTwitterID
				}
			}
			return err
		}
		return nil
	}
}

// Update は該当ユーザーのデータを更新するDBに保存する
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		userDTO := &UserDTO{
			ID:        user.ID,
			Name:      user.Name,
			Profile:   user.Profile,
			TwitterID: user.TwitterID,
		}

		if _, err := r.dbMap.Update(userDTO); err != nil {
			if sqlerr, ok := err.(*mysql.MySQLError); ok {
				// twitterIDが重複したときのエラー
				if sqlerr.Number == mysqlerr.ER_DUP_ENTRY && strings.Contains(sqlerr.Message, "twitter_id") {
					return entity.NewErrorDuplicated("user TwitterID")
				}
			}
			return err
		}
		return nil
	}
}

// Delete は該当ユーザIDを満たすユーザをDBから削除します
func (r *UserRepository) Delete(ctx context.Context, user *entity.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// ユーザの存在確認は参照のみなのでトランザクションには入れない
		// 該当ユーザの存在確認
		_, err := r.FindByID(ctx, user.ID)
		if err != nil {
			return entity.NewErrorNotFound("user")
		}

		userDTO := &UserDTO{
			ID: user.ID,
		}

		// トランザクションオブジェクトをctxから取得する
		dao, ok := getTx(ctx)
		if !ok {
			// 見つからなかったら、dbMapをそのまま設定する
			dao = r.dbMap
		}
		if _, err := dao.Delete(userDTO); err != nil {
			return err
		}

		return nil
	}
}

// DoInTx はトランザクションの中でDBにアクセスするためのラッパー関数です
func (r *UserRepository) DoInTx(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	tx, err := r.dbMap.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed dbMap.Begin(): %w", err)
	}

	// トランザクションをctxに埋め込む
	ctx = context.WithValue(ctx, &txKey, tx)
	// 中身の処理を実行する
	v, err := f(ctx)
	if err != nil {
		_ = tx.Rollback()
		return v, fmt.Errorf("rollbacked: %w", err)
	}

	// コミット時に失敗してもロールバック
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return v, fmt.Errorf("failed to commit: rollbacked: %w", err)
	}
	return v, nil
}

// UserDTO はDBとやり取りするためのDataTransferObject
type UserDTO struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Profile   string `db:"profile"`
	TwitterID string `db:"twitter_id"`
}
