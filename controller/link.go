package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"watchmen/api"
	"watchmen/config"
	"watchmen/repository"
)

var allowedMethods map[string]struct{}

func init() {

	validate = validator.New()

	allowedMethods = map[string]struct{}{
		http.MethodGet: {},
	}
}

type AddLinkBody struct {
	URL            string  `json:"url" validate:"required"`
	Method         *string `json:"method"`
	ErrorThreshold *uint   `json:"error_threshold"`
	UserID         uint    `param:"id"`
}

func AddLink(linkRepo repository.LinkRepository) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		body := new(AddLinkBody)
		if err := ctx.Bind(body); err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		user := ctx.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		if claims["sub"] != fmt.Sprintf("%d", body.UserID) {
			return ctx.JSON(http.StatusUnauthorized, api.MessageResponse{Message: "permission denied"})
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

		method := http.MethodGet
		if body.Method != nil {
			if _, ok := allowedMethods[*body.Method]; !ok {
				return ctx.JSON(http.StatusBadRequest, api.Response{
					Status: http.StatusBadRequest,
					Data:   "method: GET and POST allowed",
				})
			}

			method = *body.Method
		}

		threshold := config.C.Link.ErrorThreshold
		if body.ErrorThreshold != nil {
			threshold = *body.ErrorThreshold
		}

		link := &repository.Link{
			UserID:         body.UserID,
			ErrorThreshold: threshold,
			URL:            body.URL,
			Method:         method,
		}

		err = linkRepo.CreateLink(ctx.Request().Context(), link)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, api.Response{
			Status: http.StatusOK,
			Data: api.MessageResponse{
				Message: "link added",
			},
		})
	}
}

type GetLinkBody struct {
	UserID uint `param:"id"`
}

func GetLink(linkRepo repository.LinkRepository) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		body := new(GetLinkBody)
		if err := ctx.Bind(body); err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		user := ctx.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		if claims["sub"] != fmt.Sprintf("%d", body.UserID) {
			return ctx.JSON(http.StatusUnauthorized, api.MessageResponse{Message: "permission denied"})
		}

		links, err := linkRepo.GetLink(ctx.Request().Context(), body.UserID)
		if err != nil {
			return fmt.Errorf("GetLink.GetLink: %w", err)
		}

		return ctx.JSON(http.StatusOK, api.Response{
			Status: http.StatusOK,
			Data:   links,
		})
	}
}

type RetrieveLinkBody struct {
	UserID uint `param:"id"`
	LinkID uint `param:"link_id"`
}

// RetrieveLink used for check link status. counting successful and unsuccessful requests
func RetrieveLink(linkRepo repository.LinkRepository) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		body := new(RetrieveLinkBody)
		if err := ctx.Bind(body); err != nil {
			return ctx.JSON(http.StatusBadRequest, api.Response{
				Status: http.StatusBadRequest,
				Data:   api.MessageResponse{Message: err.Error()},
			})
		}

		user := ctx.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		if claims["sub"] != fmt.Sprintf("%d", body.UserID) {
			return ctx.JSON(http.StatusUnauthorized, api.MessageResponse{Message: "permission denied"})
		}

		status, err := linkRepo.RetrieveLinkData(ctx.Request().Context(), body.UserID, body.LinkID)
		if err != nil {
			if err == sql.ErrNoRows {
				return ctx.JSON(http.StatusNotFound, api.MessageResponse{Message: "link not found"})
			}

			return fmt.Errorf("GetLink.GetLink: %w", err)
		}

		return ctx.JSON(http.StatusOK, api.Response{
			Status: http.StatusOK,
			Data: echo.Map{
				"ok count":    status.OK,
				"error count": status.ERR,
			},
		})
	}
}
