package entity

// User はユーザを表します
type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Profile   string `json:"profile"`
	TwitterID string `json:"twitter_id"`
	IconURL   string `json:"icon_url"`
}

// NewUser はUserのポインタを生成する関数です
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
		// Authenticate時にuser IDを確認しているので想定しないエラー
		return NewErrorEmpty("user ID")
	}
	if len(u.Name) == 0 {
		return ErrEmptyUserName
	}
	// MySQLのVARCHARはマルチバイト文字も１と数えるので，それに合わせてバイト数ではなく文字数を数える
	if len([]rune(u.ID)) > 128 {
		return NewErrorTooLong("user ID")
	}
	if len([]rune(u.Name)) > 128 {
		return NewErrorTooLong("user Name")
	}
	if len([]rune(u.TwitterID)) > 15 {
		return NewErrorTooLong("user TwitterID")
	}

	return nil
}

// Format は各エンティティの表記ゆれを整形するメソッドです
func (u *User) Format() {
	// TwitterIDに@が含まれていたら取りのぞく
	if len(u.TwitterID) > 0 && u.TwitterID[0] == '@' {
		u.TwitterID = u.TwitterID[1:]
	}
}
