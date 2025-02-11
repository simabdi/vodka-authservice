package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simabdi/gofiber-exception/exception"
	"github.com/simabdi/vodka-authservice/helper"
	"github.com/simabdi/vodka-authservice/middleware"
	"github.com/simabdi/vodka-authservice/models/resource"
	"github.com/simabdi/vodka-authservice/request"
	"github.com/simabdi/vodka-authservice/service"
	"net/http"
)

type userHandler struct {
	service           service.UserService
	middlewareService middleware.Service
}

func NewUserHandler(service service.UserService, middlewareService middleware.Service) *userHandler {
	return &userHandler{service, middlewareService}
}

func (h *userHandler) Login(ctx *fiber.Ctx) error {
	var input request.LoginRequest

	if err := ctx.BodyParser(&input); err != nil {
		JsonResponse := helper.JsonResponse(http.StatusUnprocessableEntity, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(JsonResponse)
	}

	userLogin, err := h.service.Login(input)
	if err != nil {
		JsonResponse := helper.JsonResponse(http.StatusBadRequest, "email or password incorrect.", false, "", nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	token, err := h.middlewareService.GenerateToken(userLogin)
	if err != nil {
		JsonResponse := helper.JsonResponse(http.StatusBadRequest, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	JsonResponse := helper.JsonResponse(http.StatusOK, "", true, "", resource.LoginResource(userLogin, token))
	return ctx.Status(fiber.StatusOK).JSON(JsonResponse)
}

func (h *userHandler) Index(ctx *fiber.Ctx) error {
	payload, err := h.service.GetAll(ctx)
	if err != nil {
		JsonResponse := helper.JsonResponse(http.StatusBadRequest, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	jsonResponse := helper.JsonResponse(http.StatusOK, "", true, "", resource.UserPaginationResources(payload))
	return ctx.Status(fiber.StatusOK).JSON(jsonResponse)
}

func (h *userHandler) ActivationAccount(ctx *fiber.Ctx) error {
	var input request.CreateAccountRequest

	if err := ctx.BodyParser(&input); err != nil {
		JsonResponse := helper.JsonResponse(http.StatusBadRequest, "", false, exception.Validation(input), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	err := exception.Validate.Struct(input)
	if err != nil {
		JsonResponse := helper.JsonResponse(http.StatusUnprocessableEntity, "", false, exception.Validation(input), nil)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(JsonResponse)
	}

	err = h.service.ActivationAccount(ctx, input)
	if err != nil {
		JsonResponse := helper.JsonResponse(fiber.StatusBadRequest, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	JsonResponse := helper.JsonResponse(http.StatusOK, "", true, "", true)
	return ctx.Status(http.StatusOK).JSON(JsonResponse)
}

func (h *userHandler) CheckUserAvailable(ctx *fiber.Ctx) error {
	userAccount, err := h.service.GetByUuid(ctx.Params("uuid"))
	if err != nil {
		JsonResponse := helper.JsonResponse(fiber.StatusBadRequest, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	if userAccount.Status == "ACTIVE" && ctx.Query("password_type") == "create" {
		JsonResponse := helper.JsonResponse(http.StatusBadRequest, "", false, "Akun user sudah diaktivasi.", nil)
		return ctx.Status(http.StatusBadRequest).JSON(JsonResponse)
	}

	JsonResponse := helper.JsonResponse(http.StatusOK, "", true, "", resource.UserActivationResource(userAccount))
	return ctx.Status(http.StatusOK).JSON(JsonResponse)
}

func (h *userHandler) ResendEmailVerification(ctx *fiber.Ctx) error {
	var input request.UuidRequest

	if err := ctx.BodyParser(&input); err != nil {
		JsonResponse := helper.JsonResponse(http.StatusUnprocessableEntity, "", false, exception.Validation(input), nil)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(JsonResponse)
	}

	err := h.service.ResendActivation(ctx, input)
	if err != nil {
		JsonResponse := helper.JsonResponse(fiber.StatusBadRequest, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	JsonResponse := helper.JsonResponse(http.StatusOK, "", true, "", true)
	return ctx.Status(http.StatusOK).JSON(JsonResponse)
}

func (h *userHandler) ResetPassword(ctx *fiber.Ctx) error {
	var input request.ResetPasswordRequest

	if err := ctx.BodyParser(&input); err != nil {
		JsonResponse := helper.JsonResponse(http.StatusUnprocessableEntity, "", false, exception.Validation(input), nil)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(JsonResponse)
	}

	err := h.service.ResetPassword(ctx, input)
	if err != nil {
		JsonResponse := helper.JsonResponse(fiber.StatusBadRequest, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	JsonResponse := helper.JsonResponse(http.StatusOK, "", true, "", true)
	return ctx.Status(http.StatusOK).JSON(JsonResponse)
}

func (h *userHandler) ChangePassword(ctx *fiber.Ctx) error {
	var input request.UpdatePasswordRequest

	if err := ctx.BodyParser(&input); err != nil {
		JsonResponse := helper.JsonResponse(http.StatusBadRequest, "", false, exception.Validation(input), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	err := h.service.UpdatePassword(ctx, input)
	if err != nil {
		JsonResponse := helper.JsonResponse(fiber.StatusBadRequest, "", false, exception.Error(err), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse)
	}

	JsonResponse := helper.JsonResponse(http.StatusOK, "", true, "", nil)
	return ctx.Status(http.StatusOK).JSON(JsonResponse)
}
