package usecase

import (
	"context"
	"fmt"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

type CommentUseCase struct {
	commentRepo repository.Comment
}

func NewCommentUseCase(comment repository.Comment) *CommentUseCase {
	return &CommentUseCase{commentRepo: comment}
}

func (u *CommentUseCase) GetByPostID(postid string) (comments []*entity.Comment, err error) {
	comments, err = u.commentRepo.GetByPostID(postid)
	return
}
