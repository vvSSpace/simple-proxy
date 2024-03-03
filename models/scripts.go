package models

type RequestUpdateCustomerForHistory struct {
	CustomerIDForHistory int64 `json:"customer_id_for_history"`
}

type ResponseUpdateCustomerForHistory struct {
	Message string `json:"message"`
	Param   Param  `json:"param"`
}

type Param struct {
	CustomerIDForHistory int64 `json:"customer_id_for_history"`
}
type ResponseUpdateError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
