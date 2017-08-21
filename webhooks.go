package webhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

// Webhooks defines the basic struct for configuring the Webhooks server
type Webhooks struct {
	Port int

	events map[string]bool

	projectCards chan *github.ProjectCardEvent
	issues       chan *github.IssuesEvent
}

// New initializes the webhooks with the port and the events that need to take care of
func New(p int, events []string) *Webhooks {
	wh := &Webhooks{
		Port:   p,
		events: make(map[string]bool),

		projectCards: make(chan *github.ProjectCardEvent, 10),
		issues:       make(chan *github.IssuesEvent, 10),
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

// ProjectCards defines the public channel to listen for *github.ProjectCardEvent
func (wh *Webhooks) ProjectCards() <-chan *github.ProjectCardEvent { return wh.projectCards }

// Issues defines the public channel to listen for *github.IssuesEvent
func (wh *Webhooks) Issues() <-chan *github.IssuesEvent { return wh.issues }

// Start starts the webhook server
func (wh *Webhooks) Start() {
	http.HandleFunc("/", wh.eventHandle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", wh.Port), nil))
}

func (wh *Webhooks) eventHandle(w http.ResponseWriter, req *http.Request) {
	var err error
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
		return
	}

	switch e {
	case "project_card":
		err = wh.newProjectCardEvent(b)
	case "issues":
		err = wh.newIssuesEvent(b)
	}

	if err != nil {
		log.Println(err)
	}
}

func (wh *Webhooks) newProjectCardEvent(b []byte) error {
	var e github.ProjectCardEvent
	err := json.Unmarshal(b, &e)
	if err != nil {
		return err
	}
	wh.projectCards <- &e
	return nil
}

func (wh *Webhooks) newIssuesEvent(b []byte) error {
	var e github.IssuesEvent
	err := json.Unmarshal(b, &e)
	if err != nil {
		return err
	}
	wh.issues <- &e
	return nil
}
