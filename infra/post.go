package infra

import (
	"context"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

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
func (p *PostRepository) Store(ctx context.Context, post *entity.Post) error {
	postDTO := &PostDTO{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Code:      post.Code,
		Language:  post.Language,
		Content:   post.Content,
		Source:    post.Source,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	// if err := p.dbMap.In

	_ = postDTO
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
