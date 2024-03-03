package api

import (
	"encoding/json"
	"log"
	"net/http"
	"simple-proxy/models"
)

func GetSettings(w http.ResponseWriter, omUrl string, lcUrl string, newCustomerID int64, CustomerIDForHistory int64) {

	body := models.GetSettings{
		OmUrl:                omUrl,
		LcUrl:                lcUrl,
		NewCustomerID:        newCustomerID,
		CustomerIDForHistory: CustomerIDForHistory,
	}

	responseBody, err := json.Marshal(body)
	if err != nil {
		log.Println(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(responseBody)
}
