package model

type Quote struct {
	Text     string `json:"quote"`
	ServedBy string `json:"ipAddress"`
	Language string `json:"language"`
}
