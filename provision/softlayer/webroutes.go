package softlayer

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
}

func parseBody(r *http.Request) (ACPOpts, error) {
	//get the post data with credentials and setting
	bag := ACPOpts{}
	//	var conf Config
	//	var nodeData serverData
	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body", err)
	}

	decoded, decerr := b64.StdEncoding.DecodeString(string(bodyData))
	if decerr != nil {
		log.Println("Something wrong with the post data 64bit encoding", decerr)
		return bag, decerr
	}

	unmarshErr := json.Unmarshal(decoded, &bag)
	if unmarshErr != nil {
		log.Println("Cannot deserialize ket bag", unmarshErr)
		return bag, unmarshErr
	}
	fmt.Println("Bag", bag)
	defer r.Body.Close()
	return bag, nil
}

func BootUp(w http.ResponseWriter, r *http.Request) {
	log.Println("Checking if boot server is up")
	w.Header().Set("Content-Type", "application/json")
	resp := Response{Status: "up"}
	json.NewEncoder(w).Encode(resp)
}

func NodeUp(w http.ResponseWriter, r *http.Request) {
	log.Println("Node Up called")
	q := r.URL.Query()
	fmt.Println("received", q)
	nodeType := q["type"][0]
	nodeIP := q["ip"][0]
	nodeName := q["name"][0]
	log.Println("Parsed vals:", nodeType, nodeIP, nodeName)

	bag, bodyErr := parseBody(r)
	if bodyErr != nil {
		log.Println("Body error", bodyErr)
	}

	log.Println("Bag", bag)

	w.Header().Set("Content-Type", "application/json")
	resp := Response{Status: "Received node"}
	json.NewEncoder(w).Encode(resp)
}
