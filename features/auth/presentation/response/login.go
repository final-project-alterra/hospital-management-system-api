package response

type LoginResponse struct {
	Token string `json:"token"`
}

func Token(token string) LoginResponse {
	return LoginResponse{
		Token: token,
	}
}
