package request

type UserCreateRequest struct {
	Name string `json:"name"`
}

type UserGetRequest struct {
	Token string `json:"auth_token"`
}
