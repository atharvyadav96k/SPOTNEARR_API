package services

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
)

type base_service struct{}

func res(message string, statusCode int, data interface{}) response.Res {
	return response.Res{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}
}

func (b *base_service) ResponseOK(message string, data interface{}) response.Res {
	return res(message, http.StatusOK, data)
}

func (b *base_service) ResponseNotFound(message string) response.Res {
	return res(message, http.StatusNotFound, nil)
}

func (b *base_service) ResponseInternalServer(message string) response.Res {
	return res(message, http.StatusInternalServerError, nil)
}

func (b *base_service) ResponseBadRequest(message string) response.Res {
	return res(message, http.StatusBadRequest, nil)
}

func (b *base_service) ResponseConflict(message string) response.Res {
	return res(message, http.StatusConflict, nil)
}

func (b *base_service) ResponseCreated(message string, data interface{}) response.Res {
	return res(message, http.StatusCreated, data)
}
