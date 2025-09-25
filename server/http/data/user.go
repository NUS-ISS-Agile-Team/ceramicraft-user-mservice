package data

type UserVO struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserActivateReq struct {
	Code string `json:"code" binding:"required"`
}
