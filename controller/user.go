package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"watchmen/config"

	"watchmen/api"
	"watchmen/repository"
	"watchmen/utils"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type SignUpBody struct {
	FullName  string `json:"fullname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Cellphone string `json:"cellphone"`
	Password  string `json:"password" validate:"required"`
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

		err := validate.Struct(body)
		if err != nil {
			errs := make([]string, 0)
			for _, err := range err.(validator.ValidationErrors) {
				errs = append(errs, err.Error())
			}

			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   errs,
			})
		}

		body.Cellphone, err = utils.NormalizeCellphone(body.Cellphone)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		err = utils.ValidatePassword(body.Password)
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
			return err
		}

		exists, err = userRepo.EmailExists(ctx.Request().Context(), body.Email)
		if err != nil {
			return fmt.Errorf("SignUp.EmailExists: %w", err)
		} else if exists {
			return err
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
				Message: fmt.Sprintf("Signed up successfully. user id: %d", user.ID),
			},
		})
	}
}

type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Login(userRepo repository.UserRepository) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		body := new(LoginBody)
		if err := ctx.Bind(body); err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		err := validate.Struct(body)
		if err != nil {
			errs := make([]string, 0)
			for _, err := range err.(validator.ValidationErrors) {
				errs = append(errs, err.Error())
			}

			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   errs,
			})
		}

		err = utils.ValidatePassword(body.Password)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		user, err := userRepo.FindByEmail(ctx.Request().Context(), body.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				return ctx.JSON(http.StatusBadRequest, api.Response{
					Status: http.StatusBadRequest,
					Data:   api.MessageResponse{Message: "email and/or password is wrong"},
				})
			}

			return fmt.Errorf("Login.FindByEmail: %w", err)
		}

		err = user.ValidatePassword(body.Password)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data: api.MessageResponse{
					Message: "email and/or password is wrong",
				},
			})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(int(user.ID)),
		})

		tokenString, err := token.SignedString(config.C.User.JWTSecret)

		return ctx.JSON(http.StatusOK, echo.Map{
			"token": tokenString,
		})
	}
}

type GetAlertBody struct {
	UserID uint `param:"id"`
}

// GetAlert return all links of a user, that error count of them greater than or equal to link threshold
func GetAlert(userRepo repository.UserRepository) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		body := new(GetAlertBody)
		if err := ctx.Bind(body); err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		alerts, err := userRepo.GetAlerts(ctx.Request().Context(), body.UserID)
		if err != nil {
			return fmt.Errorf("GetAlert.GetAlerts: %w", err)
		}

		return ctx.JSON(http.StatusOK, api.Response{
			Status: http.StatusOK,
			Data:   alerts,
		})
	}
}
