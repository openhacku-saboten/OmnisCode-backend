package infra

import (
	"github.com/go-gorp/gorp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
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
	if _, err = r.dbMap.Select(&commentDTOs, "SELECT * FROM comments WHERE postid = ?", postid); err != nil {
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
			CreatedAt: commentDTO.CreatedAt,
			UpdatedAt: commentDTO.UpdatedAt,
		}
		comments = append(comments, comment)
	}
	return
}

// CommentDTO はDBとやり取りするためのDataTransferObject
type CommentDTO struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	PostID    int    `json:"post_id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	FirstLine int    `json:"first_line"`
	LastLine  int    `json:"last_line"`
	Code      string `json:"code"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
