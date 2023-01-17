package api

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
