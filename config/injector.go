// go:build wireinject
//go:build wireinject
// +build wireinject

package config

import (
	"github.com/budgetin-app/user-management-service/config/database"
	"github.com/budgetin-app/user-service/app/controller"
	"github.com/budgetin-app/user-service/app/repository"
	"github.com/google/wire"
)

// Databases
var db = wire.NewSet(database.ConnectDB)

// Repositories
var accountRepository = wire.NewSet(
	repository.NewAccountRepository,
	wire.Bind(new(repository.AccountRepository), new(*repository.AccountRepositoryImpl)),
)

var loginInfoRepository = wire.NewSet(
	repository.NewLoginInfoRepository,
	wire.Bind(new(repository.LoginInfoRepository), new(*repository.LoginInfoRepositoryImpl)),
)

var roleRepository = wire.NewSet(
	repository.NewRoleRepository,
	wire.Bind(new(repository.RoleRepository), new(*repository.RoleRepositoryImpl)),
)

var sessionRepository = wire.NewSet(
	repository.NewSessionRepository,
	wire.Bind(new(repository.SessionRepository), new(*repository.SessionRepositoryImpl)),
)

var emailVerificationRepository = wire.NewSet(
	repository.NewEmailVerificationRepository,
	wire.Bind(new(repository.EmailVerificationRepository), new(*repository.EmailVerificationRepositoryImpl)),
)

// Controllers
var authController = wire.NewSet(
	controller.NewAuthController,
	wire.Bind(new(controller.AuthController), new(*controller.AuthControllerImpl)),
)

// Configure initialized the dependency injection components
func Configure() *Configuration {
	wire.Build(
		NewConfiguration,
		db,
		accountRepository,
		loginInfoRepository,
		roleRepository,
		sessionRepository,
		emailVerificationRepository,
		authController,
	)
	return nil
}
