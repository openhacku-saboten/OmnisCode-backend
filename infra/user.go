package infra

import (
	"database/sql"
	"errors"

	"github.com/go-gorp/gorp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

type UserRepository struct {
	dbMap *gorp.DbMap
}

func NewUserRepository(dbMap *gorp.DbMap) *UserRepository {
	dbMap.AddTableWithName(UserDTO{}, "users").SetKeys(false, "ID")
	return &UserRepository{dbMap: dbMap}
}

// FindByID は該当IDのユーザーの情報をDBから取得して返す
func (r *UserRepository) FindByID(uid string) (user *entity.User, err error) {
	var userDTO UserDTO
	err = r.dbMap.SelectOne(&userDTO, "SELECT * FROM users WHERE id = ?", uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	user = entity.NewUser(
		userDTO.ID,
		userDTO.Name,
		userDTO.Profile,
		userDTO.TwitterID,
		"",
	)
	return
}

// UserDTO はDBとやり取りするためのDataTransferObject
type UserDTO struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Profile   string `db:"profile"`
	TwitterID string `db:"twitter_id"`
}
