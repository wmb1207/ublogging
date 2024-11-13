package models

type (
	User struct {
		UUID      string  `json:"uuid"`
		Username  string  `json:"username"`
		Email     string  `json:"email"`
		Followers []*User `json:"followers"`
		Following []*User `json:"following"`

		CreatedAt string `json:"created_at"`

		Feed []*Post `json:"feed,omitempty"`
	}

	option func(*User)
)

func withUUID(uuid string) option {
	return func(c *User) {
		c.UUID = uuid
	}
}

func withFollowers(followers []*User) option {
	return func(c *User) {
		c.Followers = followers
	}
}

func NewUser(username, email string, options ...option) *User {
	newUser := &User{
		UUID:     "",
		Username: username,
		Email:    email,
	}

	for _, option := range options {
		option(newUser)
	}

	return newUser
}
