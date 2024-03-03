package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"simple-proxy/models"
	"strings"
)

func PostSettings(w http.ResponseWriter, r *http.Request, omUrl *string, lcUrl *string, newCustomerID *int64, CustomerIDForHistory *int64) {
	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var requestParams models.PostSettings

	err = json.Unmarshal(requestData, &requestParams)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}

	if requestParams.OmUrl == "" ||
		requestParams.LcUrl == "" {
		errorResponse(w, "Missing required parameters omUrl or lcUrl")
		return
	}

	if strings.HasPrefix(requestParams.OmUrl, "http://") == false ||
		strings.HasPrefix(requestParams.LcUrl, "http://") == false {
		errorResponse(w, "The omUrl and lcUrl parameters must begin with http://")
		return
	}

	if requestParams.NewCustomerID == 0 {
		errorResponse(w, "The `customer_id` parameter is missing or equal to 0")
		return
	}

	if requestParams.CustomerIDForHistory == 0 {
		errorResponse(w, "The `customer_id_for_history` parameter is missing or equal to 0")
		return
	}

	*omUrl = requestParams.OmUrl
	*lcUrl = requestParams.LcUrl
	*newCustomerID = requestParams.NewCustomerID
	*CustomerIDForHistory = requestParams.CustomerIDForHistory

	responseParams := models.ResponseSettings{
		Message: "Settings were successfully updated",
		Params: models.GetSettings{
			OmUrl:                *omUrl,
			LcUrl:                *lcUrl,
			NewCustomerID:        *newCustomerID,
			CustomerIDForHistory: *CustomerIDForHistory,
		},
	}

	responseBody, err := json.Marshal(responseParams)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Settings were successfully updated")
	log.Println(strings.Repeat("-", 80))
	log.Println("Orders service URL:               ", *omUrl)
	log.Println("History service URL:              ", *lcUrl)
	log.Println("Customer ID to replace (orders):  ", *newCustomerID)
	log.Println("Customer ID to replace (history): ", *CustomerIDForHistory)
	log.Println(strings.Repeat("-", 80))

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func errorResponse(w http.ResponseWriter, errorText string) {
	errorParams := models.ResponseSettingsError{
		Message: "Settings not updated",
		Error:   errorText,
	}

	responseBody, err := json.Marshal(errorParams)
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(responseBody)
}
