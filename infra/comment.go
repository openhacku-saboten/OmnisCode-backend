package infra

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/service"
)

type CommentRepository struct {
	dbMap *gorp.DbMap
}

func NewCommentRepository(dbMap *gorp.DbMap) *CommentRepository {
	dbMap.AddTableWithName(CommentDTO{}, "comments").SetKeys(false, "ID")
	return &CommentRepository{dbMap: dbMap}
}

// GetByPostID は該当PostIDに属するコメントのスライスを返す
func (r *CommentRepository) GetByPostID(postid int) (comments []*entity.Comment, err error) {
	var commentDTOs []CommentDTO
	if _, err = r.dbMap.Select(&commentDTOs, "SELECT * FROM comments WHERE post_id = ?", postid); err != nil {
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
	return
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
