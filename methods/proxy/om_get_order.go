package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func OrdersGetOrder(w http.ResponseWriter, r *http.Request, condition chan<- bool, omUrl string, oldCustomerID *int64, newCustomerID int64) {
	if len(r.URL.Query()) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Missing query input parameters")
		return
	}
	log.Println("Received request:", r.URL)

	ordersURL := fmt.Sprintf("%s%s?%s", omUrl, r.URL.Path, r.URL.RawQuery)

	request, err := http.NewRequest(r.Method, ordersURL, r.Body)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("keep-alive", " */*")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		condition <- false
		log.Println("Error sending request:", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		condition <- false
		log.Printf("Request to \"%s\" was not sent, status code: %d", omUrl, response.StatusCode)
		w.WriteHeader(response.StatusCode)
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	var ordersResp map[string]interface{}

	err = json.Unmarshal(body, &ordersResp)
	if err != nil {
		return
	}

	if ordersResp["error"] != nil {
		condition <- false
		log.Printf("Order \"%s\" not found", r.URL.Query().Get("order_nr"))
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(body)
		return
	}

	*oldCustomerID = int64(ordersResp["result"].(map[string]interface{})["customer"].(map[string]interface{})["customer_id"].(float64))
	condition <- true

	ordersResp["result"].(map[string]interface{})["customer"].(map[string]interface{})["customer_id"] = newCustomerID

	modifiedOrdersBody, err := json.Marshal(ordersResp)
	if err != nil {
		return
	}
	log.Printf("The owner of the \"%s\" order has successfully changed to customer %d", r.URL.Query().Get("order_nr"), newCustomerID)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(modifiedOrdersBody)
}
