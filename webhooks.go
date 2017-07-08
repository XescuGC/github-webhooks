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

type webhooks struct {
	Port   int
	events map[string]bool
}

func NewWebhooks(p int, events []string) *webhooks {
	wh := &webhooks{
		Port:   p,
		events: make(map[string]bool),
	}

	for _, e := range events {
		wh.AddEvent(e)
	}

	return wh
}

func (wh *webhooks) AddEvent(e string) {
	if _, ok := wh.events[e]; !ok {
		wh.events[e] = true
	}
}

func (wh *webhooks) HasEvent(e string) bool {
	_, ok := wh.events[e]
	return ok
}

func (wh *webhooks) Events() (events []string) {
	for e, _ := range wh.events {
		events = append(events, e)
	}
	return
}

func (wh *webhooks) Start() {
	http.HandleFunc("/", wh.eventHandle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", wh.Port), nil))
}

func (wh *webhooks) eventHandle(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		return
	}

	e := req.Header.Get("X-GitHub-Event")

	if !wh.HasEvent(e) {
		return
	}

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
