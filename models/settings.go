package models

type GetSettings struct {
	OmUrl                string `json:"omUrl"`
	LcUrl                string `json:"lcUrl"`
	NewCustomerID        int64  `json:"Customer ID to replace (orders)"`
	CustomerIDForHistory int64  `json:"Customer ID to replace (history)"`
}

type PostSettings struct {
	OmUrl                string `json:"omUrl"`
	LcUrl                string `json:"lcUrl"`
	NewCustomerID        int64  `json:"customer_id"`
	CustomerIDForHistory int64  `json:"customer_id_for_history"`
}

type ResponseSettings struct {
	Message string      `json:"message"`
	Params  GetSettings `json:"params"`
}

type ResponseSettingsError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
