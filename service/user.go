package service

import (
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	gomail "github.com/simabdi/go-mail"
	"github.com/simabdi/vodka-authservice/helper"
	"github.com/simabdi/vodka-authservice/models"
	"github.com/simabdi/vodka-authservice/repository"
	"github.com/simabdi/vodka-authservice/request"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	STATUS_ACTIVE             = "ACTIVE"
	STATUS_WAITING_ACTIVATION = "WAITING ACTIVATION"
)

type UserService interface {
	GetAll(ctx *fiber.Ctx) ([]models.User, error)
	Login(input request.LoginRequest) (models.User, error)
	GetById(userId int) (models.User, error)
	GetByUuid(uuid string) (models.User, error)
	ActivationAccount(ctx *fiber.Ctx, input request.CreateAccountRequest) error
	UpdatePassword(ctx *fiber.Ctx, input request.UpdatePasswordRequest) error
	ResetPassword(ctx *fiber.Ctx, input request.ResetPasswordRequest) error
	ResendActivation(ctx *fiber.Ctx, input request.UuidRequest) error
}

type service struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *service {
	return &service{repository}
}

func (s *service) GetAll(ctx *fiber.Ctx) ([]models.User, error) {
	payload, err := s.repository.GetAll(ctx)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func (s *service) Login(input request.LoginRequest) (models.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.GetByEmail(email)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) GetById(userId int) (models.User, error) {
	user, err := s.repository.GetById(userId)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) GetByUuid(uuid string) (models.User, error) {
	user, err := s.repository.GetByUuid(uuid)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) ActivationAccount(ctx *fiber.Ctx, input request.CreateAccountRequest) error {
	userAccount, err := s.repository.GetByUuid(input.Uuid)
	if err != nil {
		return err
	}

	if input.TypeLink == "activation" {
		if userAccount.Status == STATUS_ACTIVE {
			return errors.New("Account is already activated.")
		}
	}

	passwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	err = s.repository.UpdateColumn(userAccount, "password", string(passwd))
	if err != nil {
		return err
	}
	if input.TypeLink == "activation" {
		err = s.repository.UpdateColumn(userAccount, "status", STATUS_ACTIVE)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) UpdatePassword(ctx *fiber.Ctx, input request.UpdatePasswordRequest) error {
	uuidLogin := ctx.Locals("uuid").(string)
	account, err := s.repository.GetByUuid(uuidLogin)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.PasswordExist))
	if err != nil {
		return errors.New("Your password was not updated, since the provided current password does not match.")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
	if err != nil {
		return err
	}

	err = s.repository.UpdateColumn(account, "password", string(password))
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ResetPassword(ctx *fiber.Ctx, input request.ResetPasswordRequest) error {
	account, err := s.repository.GetByEmail(input.Email)
	if err != nil {
		return err
	}

	var now = time.Now().AddDate(0, 0, 1)
	dateTime := now.Format(time.DateTime)

	err = s.repository.UpdateColumn(account, "link_expire", dateTime)
	if err != nil {
		return err
	}

	html := getTemplateResetPassword(input.Email, account.Uuid)
	gomail.SendMail(input.Email, "Reset Your Password", html)

	return nil
}

func (s *service) ResendActivation(ctx *fiber.Ctx, input request.UuidRequest) error {
	varUser, err := s.repository.GetByUuid(input.Uuid)
	if err != nil {
		return err
	}

	html := getTemplateResendActivation(ctx, varUser.Email, varUser.Uuid)
	gomail.SendMail(varUser.Email, "Activation Account", html)

	return nil
}

func getTemplateResetPassword(email string, uuidUser string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("[Error Mail] : Error loading .env file")
	}
	var templateBuffer bytes.Buffer

	type EmailData struct {
		Email    string
		Link     string
		UrlImage string
	}

	data := EmailData{
		Email:    email,
		Link:     os.Getenv("ACTIVATION_LINK") + helper.Std64Encode(uuidUser+":"+"forgot"),
		UrlImage: os.Getenv("URL_IMAGE"),
	}

	fileHtml := filepath.Join("..", "internal/view/reset_password.html")
	htmlData, err := os.ReadFile(fileHtml)
	if err != nil {
		log.Fatal(err)
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))
	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}

func getTemplateResendActivation(ctx *fiber.Ctx, email string, uuidUser string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("[Error Mail] : Error loading .env file")
	}
	var templateBuffer bytes.Buffer

	type EmailData struct {
		Email    string
		Link     string
		UrlImage string
	}

	data := EmailData{
		Email:    email,
		Link:     os.Getenv("ACTIVATION_LINK") + helper.Std64Encode(uuidUser+":"+"create"),
		UrlImage: os.Getenv("URL_IMAGE"),
	}

	fileHtml := filepath.Join("..", "internal/view/activation_account.html")
	htmlData, err := os.ReadFile(fileHtml)
	if err != nil {
		log.Fatal(err)
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))
	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}
