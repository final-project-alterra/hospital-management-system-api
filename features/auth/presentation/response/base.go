package response

import (
	"fmt"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	jsonformat "github.com/final-project-alterra/hospital-management-system-api/utils/json-format"
	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func Success(c echo.Context, code int, message string, data interface{}) error {
	resp := SuccessResponse{}
	resp.Meta.Code = code
	resp.Meta.Message = message
	resp.Data = data
	return c.JSON(code, resp)
}

func Error(c echo.Context, err error) error {
	resp := ErrorResponse{}
	resp.Error.Code = int(errors.Kind(err))
	resp.Error.Message = string(errors.ClientMessage(err))

	// log stack trace error
	if e, ok := err.(*errors.Error); ok {
		fmt.Printf("error trace: %+v\n", jsonformat.JSON(errors.Ops(e)))
	}
	fmt.Printf("error: %+v\n", err.Error())

	return c.JSON(resp.Error.Code, resp)
}
