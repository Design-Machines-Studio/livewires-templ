package main

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Design-Machines-Studio/livewires-templ/component"
	"github.com/a-h/templ"
)

type session struct {
	Actor, CSRF string
	Cases       map[string]string
	Expires     time.Time
}
type app struct {
	store    *store
	files    map[string][]byte
	lock     assetLock
	receipt  sourceReceipt
	mu       sync.Mutex
	sessions map[string]*session
}

func newApp(s *store, files map[string][]byte, lock assetLock, receipt sourceReceipt) *app {
	return &app{store: s, files: files, lock: lock, receipt: receipt, sessions: make(map[string]*session)}
}

func grant(actor, action, list string) bool {
	if list != "reading" && list != "requirements" {
		return false
	}
	switch actor {
	case "editor":
		return action == "read" || action == "reorder"
	case "reading-editor":
		return action == "read" || (action == "reorder" && list == "reading")
	case "viewer":
		return action == "read"
	default:
		return false
	}
}

func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func (a *app) getSession(r *http.Request) (string, session, bool) {
	c, err := r.Cookie("sortable-session")
	if err != nil {
		return "", session{}, false
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	s, ok := a.sessions[c.Value]
	if !ok || time.Now().After(s.Expires) {
		delete(a.sessions, c.Value)
		return "", session{}, false
	}
	// Do not expose the mutable case map outside its lock.
	return c.Value, session{Actor: s.Actor, CSRF: s.CSRF, Expires: s.Expires}, true
}

func (a *app) issueSession(w http.ResponseWriter, old, actor string) (session, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for id, s := range a.sessions {
		if time.Now().After(s.Expires) {
			delete(a.sessions, id)
		}
	}
	if len(a.sessions) >= 1024 {
		return session{}, errors.New("demo session limit reached")
	}
	delete(a.sessions, old)
	id := randomToken()
	s := session{Actor: actor, CSRF: randomToken(), Cases: make(map[string]string), Expires: time.Now().Add(8 * time.Hour)}
	a.sessions[id] = &s
	http.SetCookie(w, &http.Cookie{Name: "sortable-session", Value: id, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, MaxAge: 8 * 60 * 60})
	return s, nil
}

func sameOrigin(r *http.Request) bool {
	host, port, err := net.SplitHostPort(r.Host)
	if err != nil || host != "127.0.0.1" {
		return false
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return false
	}
	if site := r.Header.Get("Sec-Fetch-Site"); site != "" && site != "same-origin" && site != "none" {
		return false
	}
	if origin := r.Header.Get("Origin"); origin != "" {
		return origin == "http://"+r.Host
	}
	// Native navigation may omit Origin. An exact Referer, or browser
	// same-origin fetch metadata, plus the server-held CSRF token is required.
	ref, err := url.Parse(r.Referer())
	return (err == nil && ref.Scheme == "http" && ref.Host == r.Host) || r.Header.Get("Sec-Fetch-Site") == "same-origin"
}

// Both transports use a bounded, closed set of string fields. JSON token
// decoding detects duplicate members rather than silently taking the last.
func fields(w http.ResponseWriter, r *http.Request, allowed ...string) (map[string]string, error) {
	if r.URL.RawQuery != "" {
		return nil, errors.New("mutation query fields are forbidden")
	}
	valid := map[string]bool{}
	for _, key := range allowed {
		valid[key] = true
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	typeName, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	switch typeName {
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		for key, values := range r.PostForm {
			if !valid[key] || len(values) != 1 {
				return nil, errors.New("unknown or repeated field")
			}
			out[key] = values[0]
		}
	case "application/json":
		d := json.NewDecoder(r.Body)
		token, err := d.Token()
		if err != nil || token != json.Delim('{') {
			return nil, errors.New("object required")
		}
		for d.More() {
			t, err := d.Token()
			if err != nil {
				return nil, err
			}
			key, ok := t.(string)
			if !ok || !valid[key] {
				return nil, errors.New("unknown field")
			}
			if _, exists := out[key]; exists {
				return nil, errors.New("repeated field")
			}
			value, err := d.Token()
			if err != nil {
				return nil, err
			}
			text, ok := value.(string)
			if !ok {
				return nil, errors.New("string field required")
			}
			out[key] = text
		}
		if _, err := d.Token(); err != nil {
			return nil, err
		}
		var extra any
		if err := d.Decode(&extra); err != io.EOF {
			return nil, errors.New("trailing data")
		}
	default:
		return nil, errors.New("unsupported content type")
	}
	for _, key := range allowed {
		if _, ok := out[key]; !ok {
			return nil, fmt.Errorf("missing field %s", key)
		}
	}
	return out, nil
}

func csrfOK(r *http.Request, s session, f map[string]string) bool {
	return sameOrigin(r) && len(f["csrf"]) == len(s.CSRF) && subtle.ConstantTimeCompare([]byte(f["csrf"]), []byte(s.CSRF)) == 1
}

func (a *app) routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /session", a.sessionPage)
	m.HandleFunc("POST /session", a.selectPersona)
	m.HandleFunc("GET /sortable", a.page)
	m.HandleFunc("POST /sortable/lists/{list}/move", a.move)
	m.HandleFunc("POST /sortable/lists/{list}/case", a.setCase)
	m.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		a.store.mu.Lock()
		healthy := !a.store.unavailable
		a.store.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if !healthy {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(struct {
			Healthy  bool            `json:"healthy"`
			Source   sourceReceipt   `json:"source"`
			Producer string          `json:"producerCommit"`
			Assets   []assetIdentity `json:"assets"`
		}{healthy, a.receipt, a.lock.ProducerCommit, a.lock.Assets})
	})
	m.HandleFunc("GET /bridge.js", func(w http.ResponseWriter, r *http.Request) {
		b, _ := authored.ReadFile("bridge.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write(b)
	})
	m.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Match the original escaped path too: no encoded separators, dot paths,
		// directories or run-file discovery through this static mapping.
		if r.URL.EscapedPath() != r.URL.Path {
			http.NotFound(w, r)
			return
		}
		b, ok := a.files[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", mime.TypeByExtension(assetExtension(r.URL.Path)))
		w.Header().Set("X-Asset-SHA256", digest(b))
		w.Write(b)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		m.ServeHTTP(w, r)
	})
}
func assetExtension(path string) string {
	i := strings.LastIndex(path, ".")
	if i < 0 {
		return ""
	}
	return path[i:]
}

func (a *app) sessionPage(w http.ResponseWriter, r *http.Request) {
	_, s, ok := a.getSession(r)
	if !ok {
		var err error
		s, err = a.issueSession(w, "", "")
		if err != nil {
			http.Error(w, err.Error(), 503)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	personaPage(s).Render(r.Context(), w)
}
func (a *app) selectPersona(w http.ResponseWriter, r *http.Request) {
	id, s, ok := a.getSession(r)
	if !ok {
		http.Error(w, "Session required.", 403)
		return
	}
	f, err := fields(w, r, "csrf", "persona")
	if err != nil {
		http.Error(w, "Invalid fields.", 400)
		return
	}
	if !csrfOK(r, s, f) || !grant(f["persona"], "read", "reading") {
		http.Error(w, "Denied.", 403)
		return
	}
	if _, err := a.issueSession(w, id, f["persona"]); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	http.Redirect(w, r, "/sortable", http.StatusSeeOther)
}

type presentation struct{ Variant, Density, Controls string }

func choose(variant, density, controls string) (presentation, bool) {
	if variant == "" {
		variant = "cards"
	}
	if density == "" {
		density = "default"
	}
	if controls == "" {
		controls = "both"
	}
	return presentation{variant, density, controls}, (variant == "cards" || variant == "dividers") && (density == "compact" || density == "default") && (controls == "handle" || controls == "arrows" || controls == "both")
}
func (a *app) page(w http.ResponseWriter, r *http.Request) {
	id, s, ok := a.getSession(r)
	if !ok || !grant(s.Actor, "read", "reading") || !grant(s.Actor, "read", "requirements") {
		http.Redirect(w, r, "/session", 303)
		return
	}
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Invalid presentation.", 400)
		return
	}
	for key, v := range q {
		if (key != "variant" && key != "density" && key != "controls") || len(v) != 1 {
			http.Error(w, "Invalid presentation.", 400)
			return
		}
	}
	p, ok := choose(q.Get("variant"), q.Get("density"), q.Get("controls"))
	if !ok {
		http.Error(w, "Invalid presentation.", 400)
		return
	}
	reading, requirements := a.store.snapshot("reading"), a.store.snapshot("requirements")
	if reading.Revision == 0 || requirements.Revision == 0 {
		http.Error(w, "Saving could not be confirmed. Restart and reload required.", http.StatusServiceUnavailable)
		return
	}
	a.mu.Lock()
	if current := a.sessions[id]; current != nil {
		clear(current.Cases)
	}
	a.mu.Unlock()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	sortablePage(a, s, p, reading, requirements).Render(r.Context(), w)
}

var exampleCases = []string{"normal", "reject", "delay", "stale", "missing", "invalid", "unrelated"}

func (a *app) setCase(w http.ResponseWriter, r *http.Request) {
	id, s, ok := a.getSession(r)
	list := r.PathValue("list")
	if !ok || !grant(s.Actor, "reorder", list) {
		http.Error(w, "Denied.", 403)
		return
	}
	f, err := fields(w, r, "csrf", "case")
	if err != nil {
		http.Error(w, "Invalid fields.", 400)
		return
	}
	if !csrfOK(r, s, f) {
		http.Error(w, "Denied.", 403)
		return
	}
	valid := false
	for _, value := range exampleCases {
		valid = valid || value == f["case"]
	}
	if !valid {
		http.Error(w, "Invalid case.", 400)
		return
	}
	a.mu.Lock()
	current := a.sessions[id]
	if current != nil {
		current.Cases[list] = f["case"]
	}
	a.mu.Unlock()
	if current == nil {
		http.Error(w, "Session expired.", 403)
		return
	}
	if r.Header.Get("Datastar-Request") == "true" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{list + "Case": f["case"]})
		return
	}
	http.Redirect(w, r, "/sortable", 303)
}

func waitRequest(r *http.Request, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-r.Context().Done():
		return false
	}
}
func validMarker(value string) bool {
	if len(value) < 1 || len(value) > 128 {
		return false
	}
	for _, c := range value {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_') {
			return false
		}
	}
	return true
}

func (a *app) move(w http.ResponseWriter, r *http.Request) {
	id, s, ok := a.getSession(r)
	list := r.PathValue("list")
	if !ok || !grant(s.Actor, "read", list) || !grant(s.Actor, "reorder", list) {
		http.Error(w, "Denied.", 403)
		return
	}
	enhanced := r.Header.Get("Datastar-Request") == "true"
	keys := []string{"itemId", "before", "revision", "csrf"}
	if enhanced {
		keys = append(keys, "requestId", "marker", "variant", "density", "controls")
	}
	f, err := fields(w, r, keys...)
	if err != nil {
		http.Error(w, "Invalid fields.", 400)
		return
	}
	if !csrfOK(r, s, f) {
		http.Error(w, "Denied.", 403)
		return
	}
	p, valid := choose(f["variant"], f["density"], f["controls"])
	if enhanced && (!valid || !validMarker(f["requestId"]) || !validMarker(f["marker"])) {
		http.Error(w, "Invalid correlation or presentation.", 400)
		return
	}
	rev, err := strconv.ParseUint(f["revision"], 10, 64)
	if err != nil || rev == 0 || strconv.FormatUint(rev, 10) != f["revision"] {
		http.Error(w, "Invalid revision.", 400)
		return
	}
	// Membership and self-before checks precede all writes, with revision
	// revalidation inside the store mutex after any deliberately bounded delay.
	if !contains(initialOrders[list], f["itemId"]) || (f["before"] != "" && !contains(initialOrders[list], f["before"])) || f["before"] == f["itemId"] {
		err = errInvalidMove
	}
	a.mu.Lock()
	mode := "normal"
	if current := a.sessions[id]; current != nil {
		if v := current.Cases[list]; v != "" {
			mode = v
		}
		delete(current.Cases, list)
	}
	a.mu.Unlock()
	state := a.store.snapshot(list)
	if state.Revision == 0 {
		http.Error(w, "Saving could not be confirmed. Restart and reload required.", http.StatusServiceUnavailable)
		return
	}
	if enhanced {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("X-Accel-Buffering", "no")
	}
	if enhanced && err == nil {
		switch mode {
		case "stale", "missing", "invalid":
			status := "accepted"
			request := f["requestId"]
			if mode == "stale" {
				request += "-old"
			}
			if mode == "missing" {
				status = ""
			}
			if mode == "invalid" {
				status = "unknown"
			}
			a.result(w, list, request, f["marker"], status, "Example result; real reply follows.", state.Revision)
			if !waitRequest(r, 3*time.Second) {
				return
			}
		case "unrelated":
			if a.patch(w, r, s, list, p, f["marker"], "unrelated", state) != nil {
				return
			}
			if !waitRequest(r, 3*time.Second) {
				return
			}
		}
	}
	if mode == "delay" && !waitRequest(r, 5*time.Second) {
		return
	}
	// A session rotated/expired during the delay cannot authorize a later write.
	_, live, stillValid := a.getSession(r)
	if !stillValid || live.CSRF != s.CSRF || !grant(live.Actor, "reorder", list) {
		http.Error(w, "Session changed. Reload before trying again.", http.StatusForbidden)
		return
	}
	if err == nil && (mode == "reject" || mode == "unrelated") {
		err = errors.New("Example write rejected before saving.")
	}
	state = a.store.snapshot(list)
	if err == nil {
		state, err = a.store.move(list, f["itemId"], f["before"], rev)
	}
	if state.Revision == 0 {
		// No authoritative file could be recovered after rename. Leave enhanced
		// pending state intact; never patch the old cache as authoritative state.
		http.Error(w, "Saving could not be confirmed. Restart and reload required.", http.StatusServiceUnavailable)
		return
	}
	if !enhanced {
		if err == nil {
			http.Redirect(w, r, "/sortable", 303)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(409)
		nativeRejection(err.Error(), listView{List: list, State: state, Session: s, Presentation: p, CanReorder: true}).Render(r.Context(), w)
		return
	}
	if a.patch(w, r, s, list, p, f["marker"], f["requestId"], state) != nil {
		return
	}
	status, message := "accepted", "Order saved."
	if err != nil {
		status, message = "rejected", err.Error()
	}
	a.result(w, list, f["requestId"], f["marker"], status, message, state.Revision)
}
func contains(values []string, id string) bool {
	for _, v := range values {
		if v == id {
			return true
		}
	}
	return false
}

// Selector values are closed ASCII markers and are still CSS string encoded
// byte-by-byte. The host/request marker gates even the child replacement.
func cssString(value string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, c := range []byte(value) {
		fmt.Fprintf(&b, "\\%x ", c)
	}
	b.WriteByte('"')
	return b.String()
}
func patchSelector(list, marker string) string {
	return "#sortable-" + list + "[data-active-request=" + cssString(marker) + "] #order-" + list
}
func (a *app) patch(w http.ResponseWriter, r *http.Request, s session, list string, p presentation, marker, delivery string, state listState) error {
	var b bytes.Buffer
	if err := orderFragment(listView{List: list, State: state, Session: s, Presentation: p, CanReorder: grant(s.Actor, "reorder", list), Pending: true, Delivery: delivery}).Render(r.Context(), &b); err != nil {
		return err
	}
	fmt.Fprintf(w, "event: datastar-patch-elements\ndata: selector %s\ndata: mode replace\n", patchSelector(list, marker))
	for _, line := range strings.Split(b.String(), "\n") {
		fmt.Fprintf(w, "data: elements %s\n", line)
	}
	fmt.Fprint(w, "\n")
	return http.NewResponseController(w).Flush()
}
func (a *app) result(w http.ResponseWriter, list, request, marker, status, message string, revision uint64) {
	result := map[string]any{"list": list, "requestId": request, "marker": marker, "message": message, "revision": strconv.FormatUint(revision, 10)}
	if status != "" {
		result["status"] = status
	}
	// Patch one immutable JSON string, not a recursively merged signal object.
	// Otherwise an omitted status could inherit "accepted" from an older reply.
	envelope, _ := json.Marshal(result)
	b, _ := json.Marshal(map[string]string{list + "Result": string(envelope)})
	fmt.Fprintf(w, "event: datastar-patch-signals\ndata: signals %s\n\n", b)
	http.NewResponseController(w).Flush()
}

type listView struct {
	List                string
	State               listState
	Session             session
	Presentation        presentation
	CanReorder, Pending bool
	Delivery            string
}

func (v listView) revision() string { return strconv.FormatUint(v.State.Revision, 10) }
func (v listView) label() string {
	if v.List == "reading" {
		return "Reading list"
	}
	return "Member requirements"
}
func (v listView) controls() templ.Attributes {
	if v.Pending {
		return templ.Attributes{"aria-disabled": "true"}
	}
	return nil
}
func (v listView) host() component.SortableListProps {
	variant := ""
	if v.List == "reading" && v.Presentation.Variant == "cards" {
		variant = "cards"
	}
	attrs := templ.Attributes{"data-list": v.List, "data-csrf": v.Session.CSRF, "data-variant": v.Presentation.Variant, "data-density": v.Presentation.Density, "data-controls": v.Presentation.Controls}
	if v.CanReorder {
		attrs["data-on:sortable-move-request"] = "const payload = SortableBridge.payload(el, evt.detail); if (payload) @post('/sortable/lists/" + v.List + "/move', {payload, retry: 'never', retryMaxCount: 0, requestCancellation: 'disabled'}).catch(error => SortableBridge.failure(el, payload, error))"
		attrs["data-effect"] = "SortableBridge.result(el, $" + v.List + "Result)"
	}
	return component.SortableListProps{ID: "sortable-" + v.List, Label: v.label(), Variant: variant, Density: v.Presentation.Density, Attrs: attrs}
}
func (v listView) moveProps(index int, direction string) component.SortableMoveProps {
	id := v.State.Order[index]
	before := ""
	disabled := false
	if direction == "up" {
		before = id
		disabled = index == 0
		if index > 0 {
			before = v.State.Order[index-1]
		}
	} else {
		disabled = index == len(v.State.Order)-1
		if index+2 < len(v.State.Order) {
			before = v.State.Order[index+2]
		}
	}
	return component.SortableMoveProps{Label: itemLabel(id), Direction: direction, Action: "/sortable/lists/" + v.List + "/move", ItemID: id, Before: before, Disabled: disabled, Metadata: moveMetadata(v), Attrs: v.controls()}
}
func itemLabel(id string) string {
	return map[string]string{"essay-01": "Field notes", "essay-02": "On paper", "essay-03": "An index", "orientation": "Member orientation", "agreements": "Shared agreements", "governance": "Governance participation"}[id]
}
func purpose(id string) string {
	return map[string]string{"orientation": "Welcome people into the co-op and how members work together.", "agreements": "Review the agreements members use to work together.", "governance": "Take part in a member meeting or decision round."}[id]
}
func caseAction(list, mode string) string {
	return "@post('/sortable/lists/" + list + "/case', {payload: {csrf: el.form.elements.csrf.value, case: '" + mode + "'}, retry: 'never', retryMaxCount: 0, requestCancellation: 'disabled'})"
}
