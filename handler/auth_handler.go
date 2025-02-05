package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simabdi/gofiber-exception/exception"
	"github.com/simabdi/vodka-authservice/config"
	"github.com/simabdi/vodka-authservice/helper"
	"github.com/simabdi/vodka-authservice/middleware"
	"github.com/simabdi/vodka-authservice/models"
	"github.com/simabdi/vodka-authservice/models/resource"
	"github.com/simabdi/vodka-authservice/request"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Register(ctx *fiber.Ctx) error {
	var input request.RegisterRequest

	if err := ctx.BodyParser(&input); err != nil {
		JsonResponse := helper.JsonResponse(http.StatusUnprocessableEntity, "", false, exception.Validation(input), nil)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(JsonResponse)
	}

	// Hash Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	if err := config.DB.Create(&input).Error; err != nil {
		JsonResponse := helper.JsonResponse(http.StatusInternalServerError, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusInternalServerError).JSON(JsonResponse)
	}

	JsonResponse := helper.JsonResponse(http.StatusOK, "", true, "", true)
	return ctx.Status(http.StatusOK).JSON(JsonResponse)
}

func Login(ctx *fiber.Ctx) error {
	var input request.LoginRequest

	if err := ctx.BodyParser(&input); err != nil {
		JsonResponse := helper.JsonResponse(http.StatusUnprocessableEntity, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(JsonResponse)
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		JsonResponse := helper.JsonResponse(http.StatusUnauthorized, "Invalid credentials", false, "", nil)
		return ctx.Status(fiber.StatusUnauthorized).JSON(JsonResponse)
	}

	// Compare Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		JsonResponse := helper.JsonResponse(http.StatusUnauthorized, "Invalid credentials", false, "", nil)
		return ctx.Status(fiber.StatusUnauthorized).JSON(JsonResponse)
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		JsonResponse := helper.JsonResponse(http.StatusBadRequest, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	JsonResponse := helper.JsonResponse(http.StatusOK, "", true, "", resource.LoginResource(user, token))
	return ctx.Status(fiber.StatusOK).JSON(JsonResponse)
}
