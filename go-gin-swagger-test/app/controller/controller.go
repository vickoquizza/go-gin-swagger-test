package controller

import (
	"go-gin-swagger-test/core/domain/repository"
)

type Controller struct {
	Persistence repository.Repository
}

func NewController(persistence repository.Repository) *Controller {
	return &Controller{Persistence: persistence}
}
