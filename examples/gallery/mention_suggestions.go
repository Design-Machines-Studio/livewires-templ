package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Design-Machines-Studio/livewires-templ/component"
	"github.com/a-h/templ"
)

// Fictional fixture data only. This is not an application directory or API.
type mentionPerson struct {
	ID, Name, Detail string
	Avatar           *component.AvatarProps
	Disabled         bool
}

var mentionPeople = []mentionPerson{
	{"ada", "Ada Okafor", "Finance circle", &component.AvatarProps{Name: "Ada Okafor", Src: "/assets/comment-image.jpg"}, false},
	{"bo", "Bo Lindqvist", "Kitchen", &component.AvatarProps{Name: "Bo Lindqvist"}, false},
	{"sam-k", "Sam Rivera", "Kitchen", &component.AvatarProps{Name: "Sam Rivera"}, false},
	{"sam-f", "Sam Rivera", "Finance circle", &component.AvatarProps{Name: "Sam Rivera"}, true},
	{"di", "Di Ng", "", nil, false},
	{"long", "Maria-Fernanda Oyelaran-Castellanos de la Vega", "Festival submissions working group", &component.AvatarProps{Name: "Maria-Fernanda Oyelaran-Castellanos de la Vega"}, false},
	{"word", "Wolfeschlegelsteinhausenbergerdorff", "Archive", nil, false},
	{"ed", "Ed Tanaka", "Board", nil, false},
	{"fa", "Fatima Haddad", "Membership", nil, false},
	{"gu", "Gus Pereira", "Kitchen", nil, false},
	{"ha", "Hana Kowalski", "Finance circle", nil, false},
	{"ik", "Ike Mensah", "Delivery", nil, false},
	{"jo", "Jo Whitfield", "Delivery", nil, false},
	{"ka", "Kai Nakamura", "Board", nil, false},
	{"le", "Lea Fontaine", "Membership", nil, false},
	{"mo", "Mo Rahman", "Kitchen", nil, false},
	{"ni", "Nia Adeyemi", "Delivery", nil, false},
	{"ol", "Ola Brandt", "Finance circle", nil, false},
	{"pi", "Pia Moreau", "Kitchen", nil, false},
	{"qu", "Quinn Doyle", "Delivery", nil, false},
	{"ro", "Rosa Jiménez", "Board", nil, false},
}

func mentionMatches(query string) []mentionPerson {
	q := strings.ToLower(query)
	var result []mentionPerson
	for _, p := range mentionPeople {
		for _, word := range strings.FieldsFunc(strings.ToLower(p.Name), func(r rune) bool { return r == ' ' || r == '-' }) {
			if strings.HasPrefix(word, q) {
				result = append(result, p)
				break
			}
		}
	}
	return result
}

func mentionStatus(state string, count int) string {
	switch state {
	case "loading":
		return "Searching people…"
	case "empty":
		return "No matching people"
	case "error":
		return "People could not be loaded. Keep typing to try again."
	case "ready":
		if count == 1 {
			return "1 person found"
		}
		return fmt.Sprintf("%d people found", count)
	default:
		return ""
	}
}

func mentionReply(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("list")
	if id != "mention-list" && id != "mention-rtl-list" {
		http.Error(w, "Unknown demonstration", http.StatusBadRequest)
		return
	}
	dir := "ltr"
	if id == "mention-rtl-list" {
		dir = "rtl"
	}
	state := q.Get("state")
	busy := false
	query := q.Get("q")
	switch state {
	case "closed":
	case "loading":
		busy = q.Get("stale") == "true"
	case "ready":
		delay := 120 * time.Millisecond
		if q.Get("mode") == "slow" {
			if len(query) <= 1 {
				delay = 1600 * time.Millisecond
			} else {
				delay = 500 * time.Millisecond
			}
		}
		select {
		case <-time.After(delay):
		case <-r.Context().Done():
			return
		}
		if q.Get("mode") == "error" {
			state = "error"
		} else if len(mentionMatches(query)) == 0 {
			state = "empty"
		}
	default:
		http.Error(w, "Unknown state", http.StatusBadRequest)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	templ.Handler(mentionFragment(id, dir, state, query, busy)).ServeHTTP(w, r)
}
