package webhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

// All the public channels where the evens that you have registered will be published
var (
	ProjectCards = make(chan *github.ProjectCardEvent, 10)
)

type Webhooks struct {
	Port   int
	events map[string]bool
}

// New initializes the webhooks with the port and the events that need to take care of
func New(p int, events []string) *Webhooks {
	wh := &webhooks{
		Port:   p,
		events: make(map[string]bool),
	}

	for _, e := range events {
		wh.AddEvent(e)
	}

	return wh
}

// AddEvent adds a new event to the regeistered events
func (wh *Webhooks) AddEvent(e string) {
	if _, ok := wh.events[e]; !ok {
		wh.events[e] = true
	}
}

// HasEvent checks if an event is registered or not
func (wh *Webhooks) HasEvent(e string) bool {
	_, ok := wh.events[e]
	return ok
}

// Events returns all the events that are registered
func (wh *Webhooks) Events() (events []string) {
	for e, _ := range wh.events {
		events = append(events, e)
	}
	return
}

// Start starts the webhook server
func (wh *Webhooks) Start() {
	http.HandleFunc("/", wh.eventHandle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", wh.Port), nil))
}

func (wh *Webhooks) eventHandle(w http.ResponseWriter, req *http.Request) {
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
	var pc github.ProjectCardEvent
	err := json.Unmarshal(b, &pc)
	if err != nil {
		return err
	}
	ProjectCards <- &pc
	return nil
}
