package webhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ProjectCards = make(chan projectCardEvent, 10)
)

func Start(port int) {
	startServer(port)
}

func startServer(port int) {
	http.HandleFunc("/", eventHandle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func eventHandle(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		return
	}

	e := req.Header.Get("X-GitHub-Event")

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}

	switch e {
	case "project_card":
		newProjectCardEvent(b)
	}
}

func newProjectCardEvent(b []byte) error {
	var pc projectCardEvent
	err := json.Unmarshal(b, &pc)
	if err != nil {
		return err
	}
	ProjectCards <- pc
	return nil
}
