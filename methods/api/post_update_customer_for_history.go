package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"simple-proxy/models"
	"strings"
)

func UpdateCustomerForHistory(w http.ResponseWriter, r *http.Request, CustomerIDForHistory *int64) {
	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var requestParams models.RequestUpdateCustomerForHistory

	err = json.Unmarshal(requestData, &requestParams)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}

	if requestParams.CustomerIDForHistory == 0 {
		errorResponse(w, "The `customer_id_for_history` parameter is missing or equal to 0")
		return
	}

	*CustomerIDForHistory = requestParams.CustomerIDForHistory

	responseParams := models.ResponseUpdateCustomerForHistory{
		Message: "Customer ID for history was successfully updated",
		Param: models.Param{
			CustomerIDForHistory: *CustomerIDForHistory,
		},
	}

	responseBody, err := json.Marshal(responseParams)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Customer ID for history was successfully updated")
	log.Println(strings.Repeat("-", 80))
	log.Println("Customer ID to replace (history): ", *CustomerIDForHistory)
	log.Println(strings.Repeat("-", 80))

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
