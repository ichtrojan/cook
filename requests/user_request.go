package requests

type PreForgot struct {
	Email string `json:"email"`
}

type PostForgot struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}
