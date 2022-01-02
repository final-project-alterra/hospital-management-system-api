package presentation

import (
	"net/http"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/auth"
	"github.com/final-project-alterra/hospital-management-system-api/features/auth/presentation/request"
	"github.com/final-project-alterra/hospital-management-system-api/features/auth/presentation/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthPresetation struct {
	business auth.IBusiness
	validate *validator.Validate
}

func NewAuthPresentation(authBusiness auth.IBusiness) *AuthPresetation {
	return &AuthPresetation{
		business: authBusiness,
		validate: validator.New(),
	}
}

func (p *AuthPresetation) PostLogin(c echo.Context) error {
	status := http.StatusOK
	message := "Login success"
	const op errors.Op = "auth.presentation.PostLogin"
	var errMessage errors.ErrClientMessage

	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		errMessage = "Unable to parse request payload"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindBadRequest))
	}

	if err := p.validate.Struct(req); err != nil {
		errMessage = "Invalid email or password"
		return response.Error(c, errors.E(err, op, errMessage, errors.KindUnauthorized))
	}

	token, err := p.business.Login(req.Email, req.Password)
	if err != nil {
		return response.Error(c, errors.E(err, op))
	}
	return response.Success(c, status, message, response.Token(token))
}
