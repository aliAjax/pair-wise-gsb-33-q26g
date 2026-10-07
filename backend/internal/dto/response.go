package dto

// Response is the unified API envelope {code, message, data}.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK builds a success response.
func OK(data interface{}) Response {
	return Response{Code: 0, Message: "ok", Data: data}
}

// Fail builds an error response.
func Fail(code int, message string) Response {
	return Response{Code: code, Message: message}
}

// PageData is a pagination wrapper.
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"page_size"`
}
