package main

import (
	"log"
	"net/http"

	"github.com/sashajeltuhin/bluemix-provisioner/provision/softlayer"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bootup", softlayer.BootUp)
	mux.HandleFunc("/nodeup", softlayer.NodeUp)
	log.Println("Listening on port 8013")
	err := http.ListenAndServe(":8013", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
