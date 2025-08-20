package entitie

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	ExpiresIn   int32        `json:"expires_in"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
