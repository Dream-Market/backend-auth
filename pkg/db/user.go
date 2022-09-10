package db

import (
	"github.com/Dream-Market/backend-auth/pkg/models"
)

type UserHandler struct {
	*Handler
}

func InitUserHandler(h *Handler) *UserHandler {
	return &UserHandler{
		Handler: h,
	}
}

func (h *UserHandler) FindByEmailOrPhone(email string, phone string) (user models.User, err error) {
	err = h.DB.Where(&models.User{Email: email}).Or(&models.User{Phone: phone}).First(&user).Error
	return
}

func (h *UserHandler) Create(user models.User) (created models.User, err error) {
	err = h.DB.Create(&user).Error
	return user, err
}
