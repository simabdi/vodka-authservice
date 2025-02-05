package request

type (
	RegisterRequest struct {
		FullName    string `json:"full_name" validate:"required,max=50"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required,min=11,max=13"`
		Password    string `json:"password" validate:"required,min=8"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	CreateAccountRequest struct {
		Uuid     string `json:"uuid" validate:"required"`
		Password string `json:"password" validate:"required,min=8"`
		TypeLink string `json:"type_link" validate:"required"`
	}

	UpdatePasswordRequest struct {
		PasswordExist      string `json:"password_exist" validate:"required"`
		NewPassword        string `json:"new_password" validate:"required,min=8"`
		NewPasswordConfirm string `json:"new_password_confirm" validate:"required,eqfield=NewPassword"`
	}

	ResetPasswordRequest struct {
		Email string `json:"email" validate:"required"`
	}

	UuidRequest struct {
		Uuid string `json:"uuid" validate:"required"`
	}
)
