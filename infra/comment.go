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

var _ repository.Comment = (*CommentRepository)(nil)

type CommentRepository struct {
	dbMap *gorp.DbMap
}

func NewCommentRepository(dbMap *gorp.DbMap) *CommentRepository {
	dbMap.AddTableWithName(CommentDTO{}, "comments").SetKeys(true, "ID")
	dbMap.AddTableWithName(CommentInsertDTO{}, "comments").SetKeys(true, "ID")
	return &CommentRepository{dbMap: dbMap}
}

// FindByPostID は該当PostIDに属するコメントのスライスを返す
func (r *CommentRepository) FindByPostID(postID int) (comments []*entity.Comment, err error) {
	var commentDTOs []CommentDTO
	if _, err = r.dbMap.Select(&commentDTOs, "SELECT * FROM comments WHERE post_id = ?", postID); err != nil {
		return nil, err
	}
	for _, commentDTO := range commentDTOs {
		comment := &entity.Comment{
			ID:        commentDTO.ID,
			UserID:    commentDTO.UserID,
			PostID:    commentDTO.PostID,
			Type:      commentDTO.Type,
			Content:   commentDTO.Content,
			FirstLine: commentDTO.FirstLine,
			LastLine:  commentDTO.LastLine,
			Code:      commentDTO.Code,
			CreatedAt: service.ConvertTimeToStr(commentDTO.CreatedAt),
			UpdatedAt: service.ConvertTimeToStr(commentDTO.UpdatedAt),
		}
		comments = append(comments, comment)
	}
	if comments == nil {
		return nil, entity.NewErrorNotFound("comment")
	}
	return
}

// FindByUserID は該当IDのユーザのコメントをDBから取得して返す
func (r *CommentRepository) FindByUserID(ctx context.Context, uid string) ([]*entity.Comment, error) {
	var commentDTOs []CommentDTO
	if _, err := r.dbMap.Select(&commentDTOs, "SELECT * FROM comments WHERE user_id = ?", uid); err != nil {
		return nil, err
	}

	var comments []*entity.Comment
	for _, commentDTO := range commentDTOs {
		comment := &entity.Comment{
			ID:        commentDTO.ID,
			UserID:    commentDTO.UserID,
			PostID:    commentDTO.PostID,
			Type:      commentDTO.Type,
			Content:   commentDTO.Content,
			FirstLine: commentDTO.FirstLine,
			LastLine:  commentDTO.LastLine,
			Code:      commentDTO.Code,
			CreatedAt: service.ConvertTimeToStr(commentDTO.CreatedAt),
			UpdatedAt: service.ConvertTimeToStr(commentDTO.UpdatedAt),
		}
		comments = append(comments, comment)
	}
	if comments == nil {
		return nil, entity.NewErrorNotFound("comment")
	}
	return comments, nil
}

// Insert は該当ユーザーをDBに保存する
func (r *CommentRepository) Insert(comment *entity.Comment) error {
	commentDTO := &CommentInsertDTO{
		ID:        comment.ID,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		Type:      comment.Type,
		Content:   comment.Content,
		FirstLine: comment.FirstLine,
		LastLine:  comment.LastLine,
		Code:      comment.Code,
	}

	if err := r.dbMap.Insert(commentDTO); err != nil {
		if sqlerr, ok := err.(*mysql.MySQLError); ok {
			// 存在しないPostIDで登録した時のエラー
			if sqlerr.Number == mysqlerr.ER_NO_REFERENCED_ROW_2 && strings.Contains(sqlerr.Message, "post_id") {
				return entity.NewErrorNotFound("post")
			}
			// 存在しないUserIDで登録した時のエラー
			if sqlerr.Number == mysqlerr.ER_NO_REFERENCED_ROW_2 && strings.Contains(sqlerr.Message, "user_id") {
				return entity.NewErrorNotFound("user")
			}
		}
		return err
	}
	return nil
}

// FindByID はpostID, commentIDからコメントを取得します
func (r *CommentRepository) FindByID(postID, commentID int) (*entity.Comment, error) {
	var commentDTO CommentDTO
	if err := r.dbMap.SelectOne(&commentDTO, "SELECT * FROM comments WHERE post_id = ? AND id = ?", postID, commentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewErrorNotFound("comment")
		}
		return nil, err
	}

	return &entity.Comment{
		ID:        commentDTO.ID,
		UserID:    commentDTO.UserID,
		PostID:    commentDTO.PostID,
		Type:      commentDTO.Type,
		Content:   commentDTO.Content,
		FirstLine: commentDTO.FirstLine,
		LastLine:  commentDTO.LastLine,
		Code:      commentDTO.Code,
		CreatedAt: service.ConvertTimeToStr(commentDTO.CreatedAt),
		UpdatedAt: service.ConvertTimeToStr(commentDTO.UpdatedAt),
	}, nil
}

// CommentDTO はDBとやり取りするためのDataTransferObject
type CommentDTO struct {
	ID        int       `db:"id"`
	UserID    string    `db:"user_id"`
	PostID    int       `db:"post_id"`
	Type      string    `db:"type"`
	Content   string    `db:"content"`
	FirstLine int       `db:"first_line"`
	LastLine  int       `db:"last_line"`
	Code      string    `db:"code"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// CommentInsertDTO はInsert用のDataTransferObject
// timestamp系は参照しないようにしています
// ref: https://github.com/go-gorp/gorp/issues/125
type CommentInsertDTO struct {
	ID        int       `db:"id"`
	UserID    string    `db:"user_id"`
	PostID    int       `db:"post_id"`
	Type      string    `db:"type"`
	Content   string    `db:"content"`
	FirstLine int       `db:"first_line"`
	LastLine  int       `db:"last_line"`
	Code      string    `db:"code"`
	CreatedAt time.Time `db:"-"`
	UpdatedAt time.Time `db:"-"`
}
