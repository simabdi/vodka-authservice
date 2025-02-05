package resource

import (
	"github.com/simabdi/vodka-authservice/config"
	"github.com/simabdi/vodka-authservice/models"
	"github.com/simabdi/vodka-authservice/models/formatter"
)

func LoginResource(user models.User, token string) formatter.LoginFormatter {
	var profilePicture string
	if len(user.ProfilePicture) > 0 {
		profilePicture = config.UrlImage + "" + user.ProfilePicture[17:len(user.ProfilePicture)]
	}

	Resource := formatter.LoginFormatter{
		Uuid:           user.Uuid,
		FullName:       user.FullName,
		Email:          user.Email,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: profilePicture,
		Token:          token,
		RefType:        user.RefType,
		Role:           user.Role,
	}

	return Resource
}

func UserResource(user models.User) formatter.UserFormatter {
	var profilePicture string
	if len(user.ProfilePicture) > 0 {
		profilePicture = config.UrlImage + "" + user.ProfilePicture[17:len(user.ProfilePicture)]
	}

	Resource := formatter.UserFormatter{
		Uuid:           user.Uuid,
		FullName:       user.FullName,
		Email:          user.Email,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: profilePicture,
		Status:         user.Status,
		Role:           user.Role,
	}

	return Resource
}

func AccountResource(user models.User) formatter.AccountFormatter {
	Resource := formatter.AccountFormatter{
		Uuid:        user.Uuid,
		FullName:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}

	return Resource
}

func UserActivationResource(user models.User) formatter.UserAvailableFormatter {
	Resource := formatter.UserAvailableFormatter{
		Uuid:       user.Uuid,
		FullName:   user.FullName,
		Email:      user.Email,
		LinkExpire: user.LinkExpire,
	}

	return Resource
}

func UserCollectionResource(users []models.User) []formatter.UserFormatter {
	resourceCollection := []formatter.UserFormatter{}

	for _, user := range users {
		format := UserResource(user)
		resourceCollection = append(resourceCollection, format)
	}

	return resourceCollection
}
