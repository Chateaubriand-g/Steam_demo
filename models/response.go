package models

type ResponseDto struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data",omitempty`
}

type PageDto struct {
	Total     int64       `json:"total"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
	Data      interface{} `json:"data"`
}

const (
	SuccessCode      = 200
	BadRequestCode   = 400
	UnauthorizedCode = 401
	NotFoundCode     = 404
	ConflictCode     = 409
	ServerErrorCode  = 500
)

func SuccessResponse(data interface{}) ResponseDto {
	return ResponseDto{
		Code:    SuccessCode,
		Message: "success",
		Data:    data,
	}
}

func SuccessResponseWithMsg(data interface{}, msg string) ResponseDto {
	return ResponseDto{
		Code:    SuccessCode,
		Message: msg,
		Data:    data,
	}
}

func BadRequestResponse(data interface{}, msg string) ResponseDto {
	return ResponseDto{
		Code:    BadRequestCode,
		Message: msg,
		Data:    data,
	}
}

func UnauthorizedResponse(data interface{}, msg string) ResponseDto {
	return ResponseDto{
		Code:    UnauthorizedCode,
		Message: msg,
		Data:    data,
	}
}

func ConflictResponse(data interface{}, msg string) ResponseDto {
	return ResponseDto{
		Code:    ConflictCode,
		Message: msg,
		Data:    data,
	}
}

func NotFoundResponse(data interface{}, msg string) ResponseDto {
	return ResponseDto{
		Code:    NotFoundCode,
		Message: msg,
		Data:    data,
	}
}

func ServerErrorResponse(data interface{}, msg string) ResponseDto {
	return ResponseDto{
		Code:    ServerErrorCode,
		Message: msg,
		Data:    data,
	}
}
