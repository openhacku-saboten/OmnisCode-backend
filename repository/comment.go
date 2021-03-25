//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

import "github.com/openhacku-saboten/OmnisCode-backend/domain/entity"

type Comment interface {
	FindByPostID(postID int) (comments []*entity.Comment, err error)
	Insert(comment *entity.Comment) error
	FindByID(postID, commentID int) (comment *entity.Comment, err error)
}
