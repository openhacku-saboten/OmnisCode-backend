package entity

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
