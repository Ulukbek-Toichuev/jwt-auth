package model

type GeneralMessageResponse struct {
	Message string `json:"message"`
}

func NewGeneralMessageResponse(mssg string) *GeneralMessageResponse {
	return &GeneralMessageResponse{mssg}
}
