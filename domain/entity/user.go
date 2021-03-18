package entity

type User struct {
	ID        string
	Name      string
	Profile   string
	TwitterID string
}

func NewUser(id, name, profile, twitterID string) *User {
	return &User{
		ID:        id,
		Name:      name,
		Profile:   profile,
		TwitterID: twitterID,
	}
}
