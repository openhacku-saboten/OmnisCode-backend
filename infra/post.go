package infra

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-gorp/gorp"
	"github.com/go-sql-driver/mysql"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

var _ repository.Post = (*PostRepository)(nil)

// PostRepository は投稿情報の永続化と再構築のためのリポジトリです
type PostRepository struct {
	dbMap *gorp.DbMap
}

// NewPostRepository は投稿情報のリポジトリのポインタを生成する関数です
func NewPostRepository(dbMap *gorp.DbMap) *PostRepository {
	dbMap.AddTableWithName(PostDTO{}, "posts").SetKeys(false, "ID")
	return &PostRepository{dbMap: dbMap}
}

// Insert は引数で渡したエンティティの投稿をDBに保存します
func (p *PostRepository) Insert(ctx context.Context, post *entity.Post) error {
	postDTO := &PostDTO{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Code:      post.Code,
		Language:  post.Language,
		Content:   post.Content,
		Source:    post.Source,
		CreatedAt: time.Now(), // 空だとエラーになるので、ひとまず現在時刻を入れる
		UpdatedAt: time.Now(),
	}

	if err := p.dbMap.Insert(postDTO); err != nil {
		if sqlerr, ok := err.(*mysql.MySQLError); ok {
			// 存在しないユーザIDで登録した時のエラー
			if sqlerr.Number == mysqlerr.ER_NO_REFERENCED_ROW_2 && strings.Contains(sqlerr.Message, "user_id") {
				return errors.New("unexisted user")
			}
			// postIDが重複したときのエラー
			if sqlerr.Number == mysqlerr.ER_DUP_ENTRY && strings.Contains(sqlerr.Message, "posts.PRIMARY") {
				return errors.New("post ID is duplicated")
			}
		}

		return err
	}

	return nil
}

// PostDTO はDBとやりとりするためのDataTransferObjectです
// ref: migrations/20210319141439-CreatePosts.sql
type PostDTO struct {
	ID        int       `db:"id"`
	UserID    string    `db:"user_id"`
	Title     string    `db:"title"`
	Code      string    `db:"code"`
	Language  string    `db:"language"`
	Content   string    `db:"content"`
	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
