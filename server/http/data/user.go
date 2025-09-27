package data

type UserVO struct {
	ID       int    `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type UserActivateReq struct {
	Code string `json:"code" binding:"required,min=6,max=6"`
}
