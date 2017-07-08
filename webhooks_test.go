package webhooks

import (
	"reflect"
	"testing"
)

func TestAddEvent(t *testing.T) {
	w := New(3000, nil)
	e := "test"
	w.AddEvent(e)

	if _, ok := w.events[e]; !ok {
		t.Errorf("Expected to have event %q but was not found", e)
	}
}

func TestHasEvent(t *testing.T) {
	w := New(3000, nil)
	e := "test"
	w.AddEvent(e)

	if !w.HasEvent(e) {
		t.Errorf("Expected to have event %q but was not found", e)
	}
}

func TestEvents(t *testing.T) {
	w := New(3000, nil)
	e := "test"
	w.AddEvent(e)

	if !reflect.DeepEqual(w.Events(), []string{"test"}) {
		t.Errorf("Expected the list of events to be %v and found %v", []string{"test"}, w.Events())
	}
}

func TestNew(t *testing.T) {
	w := New(3000, []string{"test"})

	if !reflect.DeepEqual(w.Events(), []string{"test"}) {
		t.Errorf("Expected the list of events to be %v and found %v", []string{"test"}, w.Events())
	}

	if w.Port != 3000 {
		t.Errorf("Expected the port be 3000 and found %v", w.Port)
	}
}
