package adapters

import (
	"net/http"
	"time"

	_helper "github.com/croin-app-project/package/pkg/utils/helpers"
	http_response "github.com/croin-app-project/package/pkg/utils/http-response"
	"github.com/croin-app-project/user-service/internal/adapters/dto"
	"github.com/croin-app-project/user-service/internal/domain"
	"github.com/croin-app-project/user-service/internal/usecases/iservices"
	"github.com/croin-app-project/user-service/middleware"
	"github.com/mitchellh/mapstructure"

	"github.com/gofiber/fiber/v2"
)

type AuthenticateControllerImpl struct {
	_userService iservices.IUserService
}

func NewAuthenticateController(f fiber.Router, userService iservices.IUserService) {
	handler := &AuthenticateControllerImpl{_userService: userService}
	v1 := f.Group("/v1")
	api := v1.Group("/authenticate")
	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)

}

func (impl *AuthenticateControllerImpl) Register(c *fiber.Ctx) error {
	var body dto.RegisterDto
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
	passwordHash, err := _helper.HashPassword(body.Password)
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
			Status: "Register success!",
			Result: user,
		})
	}

}

func (impl *AuthenticateControllerImpl) Login(c *fiber.Ctx) error {
	var body dto.LoginDto
	if err := c.BodyParser(&body); err != nil {
		errCode, errObj := http_response.HandleException(http_response.INVALID_INPUT_PARAMETER, err)
		return c.Status(errCode).JSON(errObj)
	}

	user, err := impl._userService.VerifyUsernamePassword(body.Username, body.Password)
	if err != nil {
		errCode, errObj := http_response.HandleException(http_response.USER_NOT_FOUND, err)
		return c.Status(errCode).JSON(errObj)
	}
	if user == nil {
		errCode, errObj := http_response.HandleException(http_response.INVALID_PASSWORD, nil)
		return c.Status(errCode).JSON(errObj)
	}

	token, errBuildeToken := middleware.BuildToken(*user)
	if errBuildeToken != nil {
		errCode, errObj := http_response.HandleException(http_response.INTERNAL_SYSTEM_ERROR, errBuildeToken)
		return c.Status(errCode).JSON(errObj)
	}

	result := dto.LoginPresenter{
		AccessToken: token,
	}
	errCode, errObj := http.StatusOK, http_response.SuccessReponse{Code: http.StatusOK,
		Status: "Login success!",
		Result: result,
	}
	return c.Status(errCode).JSON(errObj)
}
