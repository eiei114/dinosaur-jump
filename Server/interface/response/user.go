package response

type UserCreateResponse struct {
	Token string `json:"token"`
}

type UserGetResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	HighScore int    `json:"highScore"`
}

type UserRankingResponse struct {
	Name      string `json:"name"`
	HighScore int    `json:"highScore"`
}
