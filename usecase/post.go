package usecase

import (
	"context"
	"fmt"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

// PostUsecase は投稿に関するユースケースの構造体です
type PostUsecase struct {
	postRepo repository.Post
}

// NewPostUsecase は投稿に関するユースケースのポインタを生成します
func NewPostUsecase(postRepo repository.Post) *PostUsecase {
	return &PostUsecase{postRepo: postRepo}
}

// Create は投稿の情報を保存するというユースケースです
func (p *PostUsecase) Create(ctx context.Context, post *entity.Post) error {
	if err := post.IsValid(); err != nil {
		return fmt.Errorf("invalid post field: %w", err)
	}

	if err := p.postRepo.Insert(ctx, post); err != nil {
		return fmt.Errorf("failed Store Post entity: %w", err)
	}
	return nil
}
