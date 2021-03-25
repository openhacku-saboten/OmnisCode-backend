package usecase

import (
	"context"
	"fmt"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

type UserUseCase struct {
	userRepo    repository.User
	authRepo    repository.Auth
	commentRepo repository.Comment
}

func NewUserUseCase(user repository.User, auth repository.Auth, comment repository.Comment) *UserUseCase {
	return &UserUseCase{
		userRepo:    user,
		authRepo:    auth,
		commentRepo: comment,
	}
}

func (u *UserUseCase) Get(ctx context.Context, uid string) (user *entity.User, err error) {
	user, err = u.userRepo.FindByID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to Get User from DB: %w", err)
	}

	user.IconURL, err = u.authRepo.GetIconURL(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to Get User from Firebase: %w", err)
	}
	return
}

func (u *UserUseCase) GetComments(ctx context.Context, uid string) ([]*entity.Comment, error) {
	comments, err := u.commentRepo.FindByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (u *UserUseCase) Create(user *entity.User) error {
	if err := user.IsValid(); err != nil {
		return fmt.Errorf("invalid user fields: %w", err)
	}
	user.Format()
	if err := u.userRepo.Insert(user); err != nil {
		return fmt.Errorf("failed to Insert User into DB: %w", err)
	}
	return nil
}

func (u *UserUseCase) Update(user *entity.User) error {
	if err := user.IsValid(); err != nil {
		return fmt.Errorf("invalid user fields: %w", err)
	}
	user.Format()

	// Updateは存在しないユーザーの更新をしてもエラーにならないので，ここでユーザーの存在確認をする
	if _, err := u.userRepo.FindByID(user.ID); err != nil {
		return fmt.Errorf("not found user %s in DB: %w", user.ID, err)
	}

	if err := u.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to Update User into DB: %w", err)
	}
	return nil
}
