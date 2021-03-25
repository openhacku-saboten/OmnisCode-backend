package infra

import (
	"fmt"
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

// CommentRepository は認証情報の永続化と再構築のためのリポジトリです
type CommentRepository struct {
	dbMap *gorp.DbMap
}

// NewCommentRepository は投稿情報のリポジトリのポインタを生成する関数です
func NewCommentRepository(dbMap *gorp.DbMap) *CommentRepository {
	dbMap.AddTableWithName(CommentDTO{}, "comments").SetKeys(true, "ID")
	dbMap.AddTableWithName(CommentInsertDTO{}, "comments").SetKeys(true, "ID")
	return &CommentRepository{dbMap: dbMap}
}

// FindByPostID は該当PostIDに属するコメントのスライスを返す
func (r *CommentRepository) FindByPostID(postid int) (comments []*entity.Comment, err error) {
	var commentDTOs []CommentDTO
	if _, err = r.dbMap.Select(&commentDTOs, "SELECT * FROM comments WHERE post_id = ?", postid); err != nil {
		return nil, fmt.Errorf("failed CommentRepository.FindByPostID: %w", err)
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
