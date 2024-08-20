package model

type User struct {
	id   string
	name string
}

func NewUser(id, name string) *User {
	return &User{
		id:   id,
		name: name,
	}
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}
