package infra

import (
	"github.com/go-gorp/gorp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

type UserRepository struct {
	dbMap *gorp.DbMap
}

func NewUserRepository(dbMap *gorp.DbMap) *UserRepository {
	dbMap.AddTableWithName(entity.User{}, "users").SetKeys(false, "ID")
	return &UserRepository{dbMap: dbMap}
}

func (r *UserRepository) FindByID(uid string) (user *entity.User, err error) {
	user = &entity.User{}
	err = r.dbMap.SelectOne(user, "SELECT * FROM users WHERE id = ?", uid)
	return
}
