package usecase

import (
	"context"
	"fmt"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

// PostUsecase は投稿に関するusecaseです
type Post struct {
	repo repository.Post
}

// NewPostUsecase はPostUsecaseのポインタを生成する関数です
func NewPostUsecase(repo repository.Post) *Post {
	return &Post{repo: repo}
}

// GetAll は保存されている全ての投稿を取得します
func (p *Post) GetAll(ctx context.Context) ([]*entity.Post, error) {
	posts, err := p.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to GetAll: %w", err)
	}

	return posts, nil
}
