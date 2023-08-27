package auth

import (
	"github.com/google/uuid"
	a "hilmy.dev/store/src/modules/account/account_entity"
)

type signupReq struct {
	Name     *string `json:"name" validate:"required,gt=0"`
	Username *string `json:"username" validate:"required,gt=0"`
	Password *string `json:"password" validate:"required,gt=0"`
}

type signinReq struct {
	Username *string `json:"username" validate:"required,gt=0"`
	Password *string `json:"password" validate:"required,gt=0"`
}

type signinRes struct {
	Token *string    `json:"token"`
	ID    *uuid.UUID `json:"id"`
	Name  *string    `json:"name"`
	Role  *a.Role    `json:"role"`
}

type accountRes struct {
	Token *string    `json:"token"`
	ID    *uuid.UUID `json:"id"`
	Name  *string    `json:"name"`
	Role  *a.Role    `json:"role"`
}
