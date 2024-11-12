package models

type (
	User struct {
		UUID      string
		Username  string
		Email     string
		Followers []*User
		Following []*User

		Feed []*Post
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
