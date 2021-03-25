package infra

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-gorp/gorp"
	"github.com/go-sql-driver/mysql"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/service"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

var _ repository.Post = (*PostRepository)(nil)

// PostRepository は投稿情報の永続化と再構築のためのリポジトリです
type PostRepository struct {
	dbMap *gorp.DbMap
}

// NewPostRepository は投稿情報のリポジトリのポインタを生成する関数です
func NewPostRepository(dbMap *gorp.DbMap) *PostRepository {
	dbMap.AddTableWithName(PostDTO{}, "posts").SetKeys(true, "id")
	dbMap.AddTableWithName(PostInsertDTO{}, "posts").SetKeys(true, "id")
	return &PostRepository{dbMap: dbMap}
}

// GetAll はMySQLサーバに接続して、全てのPostを取得して返すメソッドです
func (p *PostRepository) GetAll(context.Context) ([]*entity.Post, error) {
	var postDTOs []PostDTO
	if _, err := p.dbMap.Select(&postDTOs, "SELECT * FROM posts"); err != nil {
		return nil, err
	}

	var posts []*entity.Post
	for _, dto := range postDTOs {
		posts = append(posts, &entity.Post{
			ID:        dto.ID,
			UserID:    dto.UserID,
			Title:     dto.Title,
			Code:      dto.Code,
			Language:  dto.Language,
			Content:   dto.Content,
			Source:    dto.Source,
			CreatedAt: service.ConvertTimeToStr(dto.CreatedAt),
			UpdatedAt: service.ConvertTimeToStr(dto.UpdatedAt),
		})
	}
	if posts == nil {
		return nil, entity.NewErrorNotFound("post")
	}

	return posts, nil
}

// FindByID はpostIDから投稿を取得します
func (p *PostRepository) FindByID(ctx context.Context, postID int) (*entity.Post, error) {
	var postDTO PostDTO
	if err := p.dbMap.SelectOne(&postDTO, "SELECT * FROM posts WHERE id = ?", postID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewErrorNotFound("post")
		}
		return nil, err
	}

	return &entity.Post{
		ID:        postDTO.ID,
		UserID:    postDTO.UserID,
		Title:     postDTO.Title,
		Code:      postDTO.Title,
		Language:  postDTO.Language,
		Content:   postDTO.Content,
		Source:    postDTO.Source,
		CreatedAt: service.ConvertTimeToStr(postDTO.CreatedAt),
		UpdatedAt: service.ConvertTimeToStr(postDTO.UpdatedAt),
	}, nil
}

// FindByUserID はユーザの投稿をDBから取得します
func (r *PostRepository) FindByUserID(ctx context.Context, uid string) ([]*entity.Post, error) {
	var postDTOs []PostDTO
	if _, err := r.dbMap.Select(&postDTOs, "SELECT * FROM posts WHERE user_id = ?", uid); err != nil {
		return nil, err
	}

	var posts []*entity.Post

	for _, dto := range postDTOs {
		posts = append(posts, &entity.Post{
			ID:        dto.ID,
			UserID:    dto.UserID,
			Title:     dto.Title,
			Code:      dto.Code,
			Language:  dto.Language,
			Content:   dto.Content,
			Source:    dto.Source,
			CreatedAt: service.ConvertTimeToStr(dto.CreatedAt),
			UpdatedAt: service.ConvertTimeToStr(dto.UpdatedAt),
		})
	}
	if posts == nil {
		return nil, entity.NewErrorNotFound("post")
	}

	return posts, nil
}

// Insert は引数で渡したエンティティの投稿をDBに保存します
func (p *PostRepository) Insert(ctx context.Context, post *entity.Post) error {
	postDTO := &PostInsertDTO{
		ID:       post.ID,
		UserID:   post.UserID,
		Title:    post.Title,
		Code:     post.Code,
		Language: post.Language,
		Content:  post.Content,
		Source:   post.Source,
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

// Update は引数で渡したエンティティの投稿でDBに保存されている情報を更新します
// 投稿の所有者以外が更新する場合、更新は行われません
func (p *PostRepository) Update(ctx context.Context, post *entity.Post) error {
	// 該当するポストがあるか確認
	post, err := p.FindByID(ctx, post.ID)
	if err != nil {

	}

	postDTO := &PostInsertDTO{
		ID:       post.ID,
		UserID:   post.UserID,
		Title:    post.Title,
		Code:     post.Code,
		Language: post.Language,
		Content:  post.Content,
		Source:   post.Source,
	}

	if _, err := p.dbMap.Update(postDTO); err != nil {
		if sqlerr, ok := err.(*mysql.MySQLError); ok {
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

// PostInsertDTO はInsert用のDataTransferObjectです
// timestamp系は参照しないようにしています
// ref: https://github.com/go-gorp/gorp/issues/125
type PostInsertDTO struct {
	ID        int       `db:"id"`
	UserID    string    `db:"user_id"`
	Title     string    `db:"title"`
	Code      string    `db:"code"`
	Language  string    `db:"language"`
	Content   string    `db:"content"`
	Source    string    `db:"source"`
	CreatedAt time.Time `db:"-"`
	UpdatedAt time.Time `db:"-"`
}
