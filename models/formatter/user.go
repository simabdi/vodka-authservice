package formatter

type UserFormatter struct {
	Uuid           string `json:"uuid"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
	Status         string `json:"status"`
	Role           string `json:"role"`
}

type UserAvailableFormatter struct {
	Uuid       string `json:"uuid"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Status     string `json:"status"`
	LinkExpire string `json:"link_expire"`
}

type LoginFormatter struct {
	Uuid           string `json:"uuid"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
	Token          string `json:"token"`
	RefType        string `json:"ref_type"`
	Role           string `json:"role"`
}

type AccountFormatter struct {
	Uuid        string `json:"uuid"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}
