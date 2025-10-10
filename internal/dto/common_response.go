package dto

type CommonResponseSuccess struct {
	Data    interface{} `json:"data,omitempty"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}

type CommonRespponseError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"error,omitempty"`
}
