package entity

import "errors"

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Profile   string `json:"profile"`
	TwitterID string `json:"twitter_id"`
	IconURL   string `json:"icon_url"`
}

func NewUser(id, name, profile, twitterID, iconURL string) *User {
	return &User{
		ID:        id,
		Name:      name,
		Profile:   profile,
		TwitterID: twitterID,
		IconURL:   iconURL,
	}
}

// IsValid は各エンティティに問題がある場合はerrorを返すメソッドです
func (u *User) IsValid() error {
	if len(u.ID) == 0 {
		return errors.New("user ID must not be empty")
	}
	if len(u.Name) == 0 {
		return errors.New("user name must not be empty")
	}
	return nil
}
