package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

// PostUsecase は投稿に関するユースケースの構造体です
type PostUsecase struct {
	postRepo repository.Post
	userRepo repository.User
}

// NewPostUsecase は投稿に関するユースケースのポインタを生成します
func NewPostUsecase(postRepo repository.Post, userRepo repository.User) *PostUsecase {
	return &PostUsecase{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// GetAll は保存されている投稿を全て取得します
func (p *PostUsecase) GetAll(ctx context.Context) ([]*entity.Post, error) {
	posts, err := p.postRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to GetAll: %w", err)
	}

	return posts, nil
}

// Get はpostIDを満たす投稿を1つ取得します
func (p *PostUsecase) Get(ctx context.Context, postID int) (*entity.Post, error) {
	post, err := p.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed PostUsecase.Get: %w", err)
	}
	return post, nil
}

// Create は引数のpostエンティティをもとに投稿を1つ生成します
func (p *PostUsecase) Create(ctx context.Context, post *entity.Post) error {
	if err := p.postRepo.Insert(ctx, post); err != nil {
		return fmt.Errorf("failed Create Post entity: %w", err)
	}
	return nil
}

// Update は引数のpostエンティティをもとに投稿を1つ更新します
func (p *PostUsecase) Update(ctx context.Context, post *entity.Post) error {
	if _, err := p.userRepo.FindByID(post.UserID); err != nil {
		return fmt.Errorf("failed find user(userID: %s): %w", post.UserID, err)
	}

	fmt.Fprintf(os.Stderr, "OUT: %+v", post)

	if err := p.postRepo.Update(ctx, post); err != nil {
		return fmt.Errorf("failed Update Post: %w", err)
	}
	return nil
}
