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

// GetAll は保存されている全ての投稿を取得します
func (p *PostUsecase) GetAll(ctx context.Context) ([]*entity.Post, error) {
	posts, err := p.postRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to GetAll: %w", err)
	}

	return posts, nil
}

// Get はpost IDをもとに投稿情報を取得するというユースケースです
func (p *PostUsecase) Get(ctx context.Context, postID int) (*entity.Post, error) {
	post, err := p.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// Create は投稿の情報を保存するというユースケースです
func (p *PostUsecase) Create(ctx context.Context, post *entity.Post) error {
	// リクエストにAPI仕様にないフィールドidが含まれていたら任意のpostIDを
	// フロントでセットできてしまうので，ここらへんでpostIDを初期化しておく
	post.ID = 0
	if err := post.IsValid(); err != nil {
		return fmt.Errorf("invalid post field: %w", err)
	}

	if err := p.postRepo.Insert(ctx, post); err != nil {
		return fmt.Errorf("failed Store Post entity: %w", err)
	}
	return nil
}
