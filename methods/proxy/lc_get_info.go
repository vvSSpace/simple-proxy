package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func LcGetInfo(w http.ResponseWriter, r *http.Request, condition <-chan bool, lcUrl string, oldCustomerID *int64) {
	log.Println("Received request:", r.URL)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		return
	}

	cond, ok := <-condition
	if !ok || !cond {
		log.Printf("Skip replacing \"customer_id\" in a request to the \"Lc\" service: could not determine the owner of the order")
		return
	} else {
		log.Printf("Replacing \"customer_id\" in a request to the \"Lc\" service")
	}

	var lcRequest map[string]interface{}

	err = json.Unmarshal(body, &lcRequest)
	if err != nil {
		return
	}

	lcRequest["params"].(map[string]interface{})["customer_id"] = oldCustomerID

	modifiedLcBodyRequest, err := json.Marshal(lcRequest)
	if err != nil {
		return
	}

	lcURL := fmt.Sprintf("%s%s", lcUrl, r.URL.Path)

	request, err := http.NewRequest(r.Method, lcURL, strings.NewReader(string(modifiedLcBodyRequest)))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Request to \"%s\" was not sent, status code: %d", lcUrl, response.StatusCode)
		w.WriteHeader(response.StatusCode)
		return
	}

	defer response.Body.Close()

	bodyResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading Lc response body:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(bodyResp)
}
