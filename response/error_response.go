package response

import (
	"net/http"

	"github.com/anyufly/gin_common/apierr"
	"github.com/anyufly/gin_common/renders"
	"github.com/anyufly/gin_common/trans"
	"github.com/anyufly/stack_err/stackerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type errorResponse struct {
	*ErrorResponse
	Cause string `json:"cause"`
}

type ErrorResponse struct {
	*Response
	err     error
	realErr error
}

func NewErrorResponse(statusCode int, code string, msg string) *ErrorResponse {
	resp := NewResponse(statusCode, code, nil, msg)

	return &ErrorResponse{
		Response: resp,
	}
}

func (er *ErrorResponse) clone() *ErrorResponse {
	c := *er
	return &c
}

func (er *ErrorResponse) WithErr(err error) renders.ErrorRender {
	c := er.clone()

	c.err = err

	switch e := err.(type) {
	case stackerr.ErrorWithStack:
		c.realErr = e.Unwrap()
	default:
		c.realErr = err
	}

	return c
}

func (er *ErrorResponse) Err() (err error) {
	return er.err
}

func (er *ErrorResponse) Render(ctx *gin.Context) {

	if ve, ok := er.realErr.(validator.ValidationErrors); ok {
		data := make(map[string]interface{})
		for _, fe := range ve {
			errMsg := fe.Translate(trans.Trans())
			data[fe.Field()] = errMsg
		}
		er.Data = data
	}

	if ae, ok := er.realErr.(*apierr.APIError); ok {
		er.Code = ae.Code()
	}

	if gin.IsDebugging() {
		ctx.Render(er.StatusCode(), renders.JSON{Data: errorResponse{
			ErrorResponse: er,
			Cause:         er.err.Error(),
		}})
		return
	}

	ctx.Render(er.StatusCode(), renders.JSON{Data: er})
}

func (er *ErrorResponse) WithMsg(msg string) renders.ErrorRender {
	cer := er.clone()
	resp := cer.Response.WithMsg(msg)
	cer.Response = resp
	return cer
}

func (er *ErrorResponse) WithData(data interface{}) renders.ErrorRender {
	cer := er.clone()
	resp := cer.Response.WithData(data)
	cer.Response = resp
	return cer
}

func (er *ErrorResponse) WithStatusCode(statusCode int) renders.ErrorRender {
	cer := er.clone()
	resp := cer.Response.WithStatusCode(statusCode)
	cer.Response = resp
	return cer
}

func (er *ErrorResponse) WithCode(code string) renders.ErrorRender {
	cer := er.clone()
	resp := cer.Response.WithCode(code)
	cer.Response = resp
	return cer
}

var UnknownError = NewErrorResponse(http.StatusInternalServerError, "UnknownError", "????????????")
