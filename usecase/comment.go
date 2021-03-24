package usecase

import (
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

func (u *CommentUseCase) GetByPostID(postid int) (comments []*entity.Comment, err error) {
	comments, err = u.commentRepo.GetByPostID(postid)
	if err != nil {
		return nil, fmt.Errorf("failed to GetByPostID from DB: %w", err)
	}
	return
}

func (u *CommentUseCase) Create(comment *entity.Comment) error {
	// リクエストにAPI仕様にないフィールドidが含まれていたら任意のcommentIDを
	// フロントでセットできてしまうので，ここらへんでcommentIDを初期化しておく
	comment.ID = 0
	if err := comment.IsValid(); err != nil {
		return fmt.Errorf("invalid Comment fields: %w", err)
	}

	// Postの投稿者以外によるtype:commitを弾く

	if err := u.commentRepo.Insert(comment); err != nil {
		return fmt.Errorf("failed to Insert Comment into DB: %w", err)
	}
	return nil
}
