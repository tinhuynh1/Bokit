package response

import "net/http"

type Response struct {
	Status     string      `json:"status"`
	Code       int         `json:"code"`
	MessageKey string      `json:"message_key"`
	ErrorCode  string      `json:"error_code"`
	Data       interface{} `json:"data,omitempty"`
	Paging     *Paging     `json:"paging,omitempty"`
}

func Success(data interface{}, messageKey string) Response {
	return Response{
		Status:     "success",
		Code:       http.StatusOK,
		MessageKey: messageKey,
		Data:       data,
	}
}

func SuccessWithPaging(data interface{}, paging Paging, messageKey string) Response {
	return Response{
		Status:     "success",
		Code:       http.StatusOK,
		MessageKey: messageKey,
		Data:       data,
		Paging:     &paging,
	}
}

func Error(code int, errorCode, messageKey string) Response {
	return Response{
		Status:     "error",
		Code:       code,
		MessageKey: messageKey,
		ErrorCode:  errorCode,
	}
}
