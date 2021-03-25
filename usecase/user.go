package usecase

import (
	"context"
	"fmt"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

// UserUseCase はユーザに関するユースケースです
type UserUseCase struct {
	userRepo    repository.User
	authRepo    repository.Auth
	postRepo    repository.Post
	commentRepo repository.Comment
}

// NewUserUseCase はユーザに関するユースケースのポインタを生成します
func NewUserUseCase(user repository.User, auth repository.Auth, post repository.Post, comment repository.Comment) *UserUseCase {
	return &UserUseCase{
		userRepo:    user,
		authRepo:    auth,
		postRepo:    post,
		commentRepo: comment,
	}
}

// Get は引数のuidを満たすユーザを1つ取得します
func (u *UserUseCase) Get(ctx context.Context, uid string) (user *entity.User, err error) {
	user, err = u.userRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to Get User from DB: %w", err)
	}

	user.IconURL, err = u.authRepo.GetIconURL(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to Get User from Firebase: %w", err)
	}
	return
}

// GetComments は引数のuidを満たすユーザが行ったコメントを全て取得します
func (u *UserUseCase) GetComments(ctx context.Context, uid string) ([]*entity.Comment, error) {
	comments, err := u.commentRepo.FindByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetPosts は引数のuidを満たすユーザが行った投稿を全て取得します
func (u *UserUseCase) GetPosts(ctx context.Context, uid string) ([]*entity.Post, error) {
	posts, err := u.postRepo.FindByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed UserUseCase.GetPosts: %w", err)
	}
	return posts, nil
}

// Create は引数のユーザエンティティをもとにユーザを1つ生成します
func (u *UserUseCase) Create(ctx context.Context, user *entity.User) error {
	if err := user.IsValid(); err != nil {
		return fmt.Errorf("invalid user fields: %w", err)
	}
	user.Format()
	if err := u.userRepo.Insert(ctx, user); err != nil {
		return fmt.Errorf("failed to Insert User into DB: %w", err)
	}
	return nil
}

// Update は引数のユーザエンティティをもとに、同じuidを持つユーザがいればそのユーザを更新します
// 存在しないユーザの場合更新は行われません
func (u *UserUseCase) Update(ctx context.Context, user *entity.User) error {
	if err := user.IsValid(); err != nil {
		return fmt.Errorf("invalid user fields: %w", err)
	}
	user.Format()

	// Updateは存在しないユーザーの更新をしてもエラーにならないので，ここでユーザーの存在確認をする
	if _, err := u.userRepo.FindByID(ctx, user.ID); err != nil {
		return fmt.Errorf("not found user %s in DB: %w", user.ID, err)
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to Update User into DB: %w", err)
	}
	return nil
}
