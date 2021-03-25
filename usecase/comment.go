package usecase

import (
	"context"
	"fmt"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

// CommentUseCase はコメントに関するユースケースです
type CommentUseCase struct {
	commentRepo repository.Comment
	postRepo    repository.Post
}

// NewCommentUseCase はCommentUseCaseのポインタを生成する関数です
func NewCommentUseCase(comment repository.Comment, post repository.Post) *CommentUseCase {
	return &CommentUseCase{commentRepo: comment, postRepo: post}
}

// Get は引数のpostIDとcommentIDの両方を満たすコメントを1つ取得します
func (u *CommentUseCase) Get(postID, commentID int) (comment *entity.Comment, err error) {
	comment, err = u.commentRepo.FindByID(postID, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to Get comment from DB: %w", err)
	}
	return
}

// GetByPostID は引数のpostIDを満たす投稿にぶら下がるコメントを全て取得します
func (u *CommentUseCase) GetByPostID(postID int) (comments []*entity.Comment, err error) {
	comments, err = u.commentRepo.FindByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to GetByPostID from DB: %w", err)
	}
	return
}

// Create は引数のcommentエンティティをもとにコメントを1つ生成します
func (u *CommentUseCase) Create(ctx context.Context, comment *entity.Comment) error {
	// Postのオーナー以外によるcommitを弾く
	post, err := u.postRepo.FindByID(ctx, comment.PostID)
	if err != nil {
		return fmt.Errorf("not found post %d in DB: %w", comment.PostID, err)
	}
	if comment.Type == "commit" && comment.UserID != post.UserID {
		return entity.ErrCannotCommit
	}

	if err := u.commentRepo.Insert(comment); err != nil {
		return fmt.Errorf("failed to Insert Comment into DB: %w", err)
	}
	return nil
}

// Update は引数のCommentエンティティをもとにコメントを1つ更新します
func (u *CommentUseCase) Create(ctx context.Context, comment *entity.Comment) error {
	// リクエストにAPI仕様にないフィールドidが含まれていたら任意のcommentIDを
	// フロントでセットできてしまうので，ここらへんでcommentIDを初期化しておく
	comment.ID = 0
	if err := comment.IsValid(); err != nil {
		return fmt.Errorf("invalid Comment fields: %w", err)
	}

	// Postのオーナー以外によるcommitを弾く
	post, err := u.postRepo.FindByID(ctx, comment.PostID)
	if err != nil {
		return fmt.Errorf("not found post %d in DB: %w", comment.PostID, err)
	}
	if comment.Type == "commit" && comment.UserID != post.UserID {
		return entity.ErrCannotCommit
	}

	if err := u.commentRepo.Insert(comment); err != nil {
		return fmt.Errorf("failed to Insert Comment into DB: %w", err)
	}
	return nil
}
