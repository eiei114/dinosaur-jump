package domain

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	TransformX int    `json:"transformX"`
	TransformY int    `json:"transformY"`
}

func NewUser(id, name string) (*User, error) {

	return &User{
		Id:         id,
		Name:       name,
		TransformX: 0,
		TransformY: 0,
	}, nil
}
