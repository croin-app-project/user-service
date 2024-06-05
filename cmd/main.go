package main

import (
	_helpers "github.com/croin-app-project/package/pkg/utils/helpers"
	_config "github.com/croin-app-project/user-service/config"
	_middleware "github.com/croin-app-project/user-service/middleware"

	"github.com/croin-app-project/user-service/internal/adapters"
	"github.com/croin-app-project/user-service/internal/domain"
	"github.com/croin-app-project/user-service/internal/domain/repositories"
	"github.com/croin-app-project/user-service/internal/usecases"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	config := _config.ReadConfiguration()
	dbConfig := _helpers.Filter(config.Databases, func(s _config.DatabaseSetting) bool {
		return s.DbName == "user"
	})[0]
	configUserService := _helpers.Filter(config.Server.Services, func(s _config.ServiceSetting) bool {
		return s.Name == "user-service"
	})[0]
	app := fiber.New()
	api := app.Group("/api/" + configUserService.Name)
	dsn := dbConfig.Url
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.UserRole{})
	db.AutoMigrate(&domain.Role{})
	db.AutoMigrate(&domain.RolePermission{})
	db.AutoMigrate(&domain.Permission{})

	userRepository := repositories.NewUserGormRepository(db)
	userService := usecases.NewUserService(userRepository)
	adapters.NewAuthenticateController(api, userService)

	middL := _middleware.InitMiddleware()
	app.Use(middL.CORS())

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("ok!")
	})

	app.Listen(":" + configUserService.Port)
}
