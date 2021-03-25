//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

import "github.com/openhacku-saboten/OmnisCode-backend/domain/entity"

// Comment はコメントに関する永続化と再構築のためのリポジトリです
type Comment interface {
	FindByPostID(postid int) (comments []*entity.Comment, err error)
	Insert(comment *entity.Comment) error
}
