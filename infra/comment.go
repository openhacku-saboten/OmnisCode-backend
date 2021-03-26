package infra

import (
	"context"
	"database/sql"
	"errors"
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

// CommentRepository は認証情報の永続化と再構成のためのリポジトリです
type CommentRepository struct {
	dbMap *gorp.DbMap
}

// NewCommentRepository は投稿情報のリポジトリのポインタを生成する関数です
func NewCommentRepository(dbMap *gorp.DbMap) *CommentRepository {
	dbMap.AddTableWithName(CommentDTO{}, "comments").SetKeys(true, "ID")
	dbMap.AddTableWithName(CommentInsertDTO{}, "comments").SetKeys(true, "ID")
	return &CommentRepository{dbMap: dbMap}
}

// FindByID はpostID, commentIDからコメントを取得します
func (r *CommentRepository) FindByID(ctx context.Context, postID, commentID int) (*entity.Comment, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
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
}

// FindByPostID は該当PostIDに属するコメントのスライスを返す
func (r *CommentRepository) FindByPostID(ctx context.Context, postID int) (comments []*entity.Comment, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var commentDTOs []CommentDTO
		if _, err = r.dbMap.Select(&commentDTOs, "SELECT * FROM comments WHERE post_id = ?", postID); err != nil {
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
}

// FindByUserID は該当IDのユーザのコメントをDBから取得して返す
func (r *CommentRepository) FindByUserID(ctx context.Context, uid string) ([]*entity.Comment, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
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
}

// Insert は該当ユーザーをDBに保存する
func (r *CommentRepository) Insert(ctx context.Context, comment *entity.Comment) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := comment.IsValid(); err != nil {
			return fmt.Errorf("invalid Comment fields: %w", err)
		}

		commentDTO := &CommentInsertDTO{
			ID:        0, // auto incrementされるのでこれで良い
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
		comment.ID = commentDTO.ID
		return nil
	}
}

// Update は引数で渡したエンティティのコメントでDBに保存されている情報を更新します
// コメントした人以外が更新する場合、更新は行われません
func (r *CommentRepository) Update(ctx context.Context, comment *entity.Comment) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// 該当するコメントが存在するか確認
		gotComment, err := r.FindByID(ctx, comment.PostID, comment.ID)
		if err != nil {
			// そもそも取得できない場合は権限がないのと同義なのでErrNotFoundを返す
			return entity.NewErrorNotFound("comment")
		}

		// 所有者でないなら、更新処理は行わない
		if gotComment.UserID != comment.UserID {
			return entity.ErrIsNotAuthor
		}

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

		if _, err := r.dbMap.Update(commentDTO); err != nil {
			return err
		}
	}
	return nil
}

// Delete は該当コメントをDBから削除する
func (r *CommentRepository) Delete(ctx context.Context, userID string, postID, commentID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// 該当するコメントが存在するか確認
		gotComment, err := r.FindByID(ctx, postID, commentID)
		if err != nil {
			return entity.NewErrorNotFound("comment")
		}

		// 所有者でないなら、削除処理は行わない
		if gotComment.UserID != userID {
			return entity.ErrIsNotAuthor
		}

		commentDTO := &CommentInsertDTO{
			ID:     commentID,
			PostID: postID,
		}

		if _, err := r.dbMap.Delete(commentDTO); err != nil {
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
