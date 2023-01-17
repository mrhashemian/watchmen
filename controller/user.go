package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"watchmen/api"
	"watchmen/repository"
	"watchmen/utils"
)

type SignUpBody struct {
	FullName  string `json:"fullname"`
	Email     string `json:"email"`
	Cellphone string `json:"cellphone"`
	Password  string `json:"password"`
}

func SignUp(userRepo repository.UserRepository) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		body := new(SignUpBody)
		if err := ctx.Bind(body); err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		var err error
		body.Cellphone, err = utils.NormalizeCellphone(body.Cellphone)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		body.Email, err = utils.NormalizeEmail(body.Email)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		exists, err := userRepo.CellphoneExists(ctx.Request().Context(), body.Cellphone)
		if err != nil {
			return fmt.Errorf("SignUp.CellphoneExists: %w", err)
		} else if exists {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data: api.MessageResponse{
					Message: "Duplicate Cellphone",
				},
			})
		}

		exists, err = userRepo.EmailExists(ctx.Request().Context(), body.Email)
		if err != nil {
			return fmt.Errorf("SignUp.EmailExists: %w", err)
		} else if exists {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data: api.MessageResponse{
					Message: "Duplicate Email",
				},
			})
		}

		user := &repository.User{
			Email:     body.Email,
			Cellphone: body.Cellphone,
			FullName:  body.FullName,
		}

		err = user.SetPassword(body.Password)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data: api.MessageResponse{
					Message: "Password is not appropriate",
				},
			})
		}

		err = userRepo.CreateUser(ctx.Request().Context(), user)
		if err != nil {
			return fmt.Errorf("SignUp.CreateUser: %w", err)
		}

		return ctx.JSON(http.StatusOK, api.Response{
			Status: http.StatusOK,
			Data: api.MessageResponse{
				Message: fmt.Sprintf("Signed up: %d", user.ID),
			},
		})
	}
}
