package config

import "github.com/Budgetin-Project/user-service/app/controller"

type Configuration struct {
	AuthController controller.AuthController
}

func NewConfiguration(
	authController controller.AuthController,
) *Configuration {
	return &Configuration{
		AuthController: authController,
	}
}
