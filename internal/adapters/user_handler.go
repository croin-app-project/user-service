package adapters

import (
	"net/http"
	"time"

	_helper "github.com/croin-app-project/package/pkg/utils/helpers"
	http_response "github.com/croin-app-project/package/pkg/utils/http-response"
	"github.com/croin-app-project/user-service/internal/adapters/dto"
	"github.com/croin-app-project/user-service/internal/domain"
	"github.com/croin-app-project/user-service/internal/usecases/iservices"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
)

type UserControllerImpl struct {
	_userService iservices.IUserService
}

func NewUserController(f fiber.Router, userService iservices.IUserService) {
	handler := &UserControllerImpl{_userService: userService}
	v1 := f.Group("/v1")
	api := v1.Group("/user")
	api.Get("/", handler.GetAllUsers)
	api.Post("/", handler.AddUser)
	api.Put("/", handler.UpdateUser)
	api.Post("/ban", handler.BanUser)
}

func (impl *UserControllerImpl) GetAllUsers(c *fiber.Ctx) error {
	users, err := impl._userService.GetAllUsers()
	if err != nil {
		errCode, errObj := http_response.HandleException(http_response.INTERNAL_SYSTEM_ERROR, err)
		return c.Status(errCode).JSON(errObj)
	}
	code, result := http.StatusOK, http_response.SuccessReponse{Code: http.StatusOK,
		Status: "Success!",
		Result: users,
	}
	return c.Status(code).JSON(result)
}

func (impl *UserControllerImpl) AddUser(c *fiber.Ctx) error {
	var body dto.UserDto
	if err := c.BodyParser(&body); err != nil {
		errCode, errObj := http_response.HandleException(http_response.INVALID_INPUT_PARAMETER, err)
		return c.Status(errCode).JSON(errObj)
	}

	isExists, err := impl._userService.IsAlreadyExistsByUsername(body.Username)
	if err != nil && err.Error() != "record not found" {
		errCode, errObj := http_response.HandleException(http_response.DATABASE_ERROR, err)
		return c.Status(errCode).JSON(errObj)
	}
	if isExists {
		errCode, errObj := http_response.HandleException(http_response.DATA_ALREADY_EXISTS, nil)
		return c.Status(errCode).JSON(errObj)
	}
	passwordHash, err := _helper.HashPassword(body.Username + "1234")
	if err != nil {
		errCode, errObj := http_response.HandleException(http_response.INTERNAL_SYSTEM_ERROR, err)
		return c.Status(errCode).JSON(errObj)
	}
	user := domain.User{
		PasswordHash:     passwordHash,
		RegistrationDate: time.Now(),
	}
	mapstructure.Decode(body, &user)

	if err := impl._userService.SaveNewUser(user); err != nil {
		errCode, errObj := http_response.HandleException(http_response.INTERNAL_SYSTEM_ERROR, err)
		return c.Status(errCode).JSON(errObj)
	} else {
		return c.Status(fiber.StatusCreated).JSON(http_response.SuccessReponse{
			Code:   http.StatusCreated,
			Status: "Add success!",
			Result: user,
		})
	}
}

func (impl *UserControllerImpl) UpdateUser(c *fiber.Ctx) error {
	var body dto.UserDto
	if err := c.BodyParser(&body); err != nil {
		errCode, errObj := http_response.HandleException(http_response.INVALID_INPUT_PARAMETER, err)
		return c.Status(errCode).JSON(errObj)
	}

	isExists, err := impl._userService.IsAlreadyExistsByUsername(body.Username)
	if err != nil && err.Error() != "record not found" {
		errCode, errObj := http_response.HandleException(http_response.DATABASE_ERROR, err)
		return c.Status(errCode).JSON(errObj)
	}
	if isExists {
		user := domain.User{}

		mapstructure.Decode(body, &user)

		if err := impl._userService.UpdateUser(user); err != nil {
			errCode, errObj := http_response.HandleException(http_response.INTERNAL_SYSTEM_ERROR, err)
			return c.Status(errCode).JSON(errObj)
		}

		code, result := http.StatusOK, http_response.SuccessReponse{Code: http.StatusOK,
			Status: "Update user success!",
			Result: user,
		}
		return c.Status(code).JSON(result)
	} else {
		errCode, errObj := http_response.HandleException(http_response.USER_NOT_FOUND, nil)
		return c.Status(errCode).JSON(errObj)
	}
}

func (impl *UserControllerImpl) BanUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessReponse{
		Status: "Ban user success!",
		Result: nil,
	})
}
