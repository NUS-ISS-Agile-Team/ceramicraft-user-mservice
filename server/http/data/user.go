package data

type UserLoginVO struct {
	ID       int    `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type UserActivateReq struct {
	Code string `json:"code" binding:"required,min=6,max=6"`
}

type UserProfileVO struct {
	ID             int            `json:"id"`
	Email          string         `json:"email"`
	Name           string         `json:"name"`
	Avatar         string         `json:"avatar"`
	DefaultAddress *UserAddressVO `json:"default_address,omitempty"`
}

type UserAddressVO struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	ZipCode      string `json:"zip_code" binding:"required"`
	Country      string `json:"country" binding:"required"`
	Province     string `json:"province"`
	City         string `json:"city"`
	Detail       string `json:"detail" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	ContactPhone string `json:"contact_phone" binding:"required,e164"`
	IsDefault    bool   `json:"is_default"`
}
