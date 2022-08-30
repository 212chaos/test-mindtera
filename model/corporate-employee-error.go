package model

type CorporateEmployeeError struct {
	Number       int    `json:"no"`
	Email        string `json:"email"`
	ErrorMessage string `json:"error"`
}
