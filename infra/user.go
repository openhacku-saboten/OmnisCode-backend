package infra

import (
	"database/sql"
	"errors"
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
func (r *UserRepository) FindByID(uid string) (user *entity.User, err error) {
	var userDTO UserDTO
	err = r.dbMap.SelectOne(&userDTO, "SELECT * FROM users WHERE id = ?", uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrUserNotFound
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

// Insert は該当ユーザーをDBに保存する
func (r *UserRepository) Insert(user *entity.User) error {
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

// Update は該当ユーザーをDBに保存する
func (r *UserRepository) Update(user *entity.User) error {
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

// UserDTO はDBとやり取りするためのDataTransferObject
type UserDTO struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Profile   string `db:"profile"`
	TwitterID string `db:"twitter_id"`
}
