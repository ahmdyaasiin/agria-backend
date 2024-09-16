package request

type FinishRegisterOAuth struct {
	Email       string  `json:"email" validate:"required,email,unique"`
	Token       string  `json:"token" validate:"required"`
	Username    string  `json:"username" validate:"required,min=3,max=20,unique"`
	PhoneNumber string  `json:"phone_number" validate:"required,e164,unique"`
	Address     string  `json:"address" validate:"required"`
	District    string  `json:"district" validate:"required"`
	City        string  `json:"city" validate:"required"`
	State       string  `json:"state" validate:"required"`
	PostalCode  string  `json:"postal_code" validate:"required,numeric"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
}

type PreRegister struct {
	Name     string `json:"name" validate:"required,min=3"`
	Username string `json:"username" validate:"required,min=3,max=20,unique"`
	Email    string `json:"email" validate:"required,email,unique"`
	Password string `json:"password" validate:"required,min=8"`
}

type Register struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Username    string  `json:"username" validate:"required,min=3,max=20,unique"`
	Email       string  `json:"email" validate:"required,email,unique"`
	Password    string  `json:"password" validate:"required,min=8"`
	PhoneNumber string  `json:"phone_number" validate:"required,e164,unique"`
	Address     string  `json:"address" validate:"required"`
	District    string  `json:"district" validate:"required"`
	City        string  `json:"city" validate:"required"`
	State       string  `json:"state" validate:"required"`
	PostalCode  string  `json:"postal_code" validate:"required,numeric"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
}

type PostRegister struct {
	Email string `json:"email" validate:"required,email"`
}

type FinishRegister struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
