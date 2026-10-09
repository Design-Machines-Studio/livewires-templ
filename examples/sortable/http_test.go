package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func testApp(t *testing.T) *app {
	t.Helper()
	return newApp(testStore(t), map[string][]byte{}, assetLock{}, sourceReceipt{Commit: strings.Repeat("a", 40), Fingerprint: strings.Repeat("b", 64)})
}
func request(a *app, method, path, body string, cookie *http.Cookie, enhanced bool) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://127.0.0.1:8080"+path, strings.NewReader(body))
	r.Header.Set("Origin", "http://127.0.0.1:8080")
	if enhanced {
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Datastar-Request", "true")
	} else {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if cookie != nil {
		r.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	a.routes().ServeHTTP(w, r)
	return w
}
func login(t *testing.T, a *app, actor string) (*http.Cookie, string) {
	t.Helper()
	w := request(a, "GET", "/session", "", nil, false)
	cookie := w.Result().Cookies()[0]
	_, s, ok := a.getSession(withCookie(cookie))
	if !ok {
		t.Fatal("anonymous session missing")
	}
	w = request(a, "POST", "/session", url.Values{"csrf": {s.CSRF}, "persona": {actor}}.Encode(), cookie, false)
	if w.Code != 303 {
		t.Fatal("login failed", w.Code, responseBody(w))
	}
	next := w.Result().Cookies()[0]
	if next.Value == cookie.Value {
		t.Fatal("session not rotated")
	}
	if _, _, ok := a.getSession(withCookie(cookie)); ok {
		t.Fatal("old session survived")
	}
	_, s, _ = a.getSession(withCookie(next))
	return next, s.CSRF
}
func withCookie(cookie *http.Cookie) *http.Request {
	r := httptest.NewRequest("GET", "http://127.0.0.1:8080/sortable", nil)
	r.AddCookie(cookie)
	return r
}
func form(csrf, item, before, revision string) url.Values {
	return url.Values{"csrf": {csrf}, "itemId": {item}, "before": {before}, "revision": {revision}}
}
func payload(csrf, item, before, list string, revision uint64) map[string]string {
	return map[string]string{"csrf": csrf, "itemId": item, "before": before, "revision": strconv.FormatUint(revision, 10), "requestId": "request-1", "marker": "connection-1", "variant": "cards", "density": "default", "controls": "both"}
}
func encode(value any) string { b, _ := json.Marshal(value); return string(b) }

// Decode only explicitly emitted result envelopes for HTTP assertions. This
// does not establish browser effect/delivery coverage.
func responseBody(w *httptest.ResponseRecorder) string {
	body := w.Body.String()
	for _, line := range strings.Split(body, "\n") {
		if !strings.HasPrefix(line, "data: signals ") {
			continue
		}
		var signals map[string]json.RawMessage
		if json.Unmarshal([]byte(strings.TrimPrefix(line, "data: signals ")), &signals) != nil {
			continue
		}
		for key, value := range signals {
			if !strings.HasSuffix(key, "Result") {
				continue
			}
			var envelope string
			if json.Unmarshal(value, &envelope) == nil {
				body += "\n" + envelope
			}
		}
	}
	return body
}

func TestSessionValidationAndOrigin(t *testing.T) {
	a := testApp(t)
	cookie, csrf := login(t, a, "editor")
	before, _ := os.ReadFile(a.store.path)
	for _, path := range []string{"/session", "/sortable/lists/reading/move", "/sortable/lists/reading/case"} {
		values := form(csrf, "essay-02", "essay-01", "1")
		if path == "/session" {
			values = url.Values{"csrf": {csrf}, "persona": {"viewer"}}
		}
		if strings.HasSuffix(path, "/case") {
			values = url.Values{"csrf": {csrf}, "case": {"reject"}}
		}
		for _, origin := range []string{"http://evil.test", "null", "http://127.0.0.1:9999", ""} {
			r := httptest.NewRequest("POST", "http://127.0.0.1:8080"+path, strings.NewReader(values.Encode()))
			r.AddCookie(cookie)
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Set("Origin", origin)
			w := httptest.NewRecorder()
			a.routes().ServeHTTP(w, r)
			if w.Code != 403 {
				t.Fatalf("%s accepted origin %q: %d", path, origin, w.Code)
			}
		}
		values.Del("csrf")
		if w := request(a, "POST", path, values.Encode(), cookie, false); w.Code != 400 {
			t.Fatal("missing csrf accepted", path, w.Code)
		}
	}
	if w := request(a, "POST", "/session", url.Values{"csrf": {csrf}, "persona": {"unknown"}}.Encode(), cookie, false); w.Code != 403 {
		t.Fatal("unknown persona accepted")
	}
	if w := request(a, "POST", "/session", url.Values{"csrf": {csrf}, "persona": {"viewer"}}.Encode(), nil, false); w.Code != 403 {
		t.Fatal("persona without session accepted")
	}
	if w := request(a, "POST", "/sortable/lists/reading/case", url.Values{"csrf": {csrf}, "case": {"forever"}}.Encode(), cookie, false); w.Code != 400 {
		t.Fatal("unbounded case accepted")
	}
	after, _ := os.ReadFile(a.store.path)
	if string(before) != string(after) {
		t.Fatal("origin/session failures wrote state")
	}
}

func TestSessionAndGrantBoundary(t *testing.T) {
	for _, actor := range []string{"editor", "reading-editor", "viewer"} {
		t.Run(actor, func(t *testing.T) {
			a := testApp(t)
			cookie, csrf := login(t, a, actor)
			page := request(a, "GET", "/sortable", "", cookie, false).Body.String()
			for _, list := range []string{"reading", "requirements"} {
				permitted := grant(actor, "reorder", list)
				if strings.Contains(page, `action="/sortable/lists/`+list+`/move"`) != permitted {
					t.Fatal("UI grants diverge", list)
				}
				item, before := "essay-02", "essay-01"
				if list == "requirements" {
					item, before = "governance", "orientation"
				}
				w := request(a, "POST", "/sortable/lists/"+list+"/move", form(csrf, item, before, "1").Encode(), cookie, false)
				if permitted && w.Code != 303 || !permitted && w.Code != 403 {
					t.Fatal("POST grants diverge", list, w.Code)
				}
				w = request(a, "POST", "/sortable/lists/"+list+"/case", url.Values{"csrf": {csrf}, "case": {"reject"}}.Encode(), cookie, false)
				if permitted && w.Code != 303 || !permitted && w.Code != 403 {
					t.Fatal("case grants diverge", list, w.Code)
				}
			}
		})
	}
	a := testApp(t)
	cookie, csrf := login(t, a, "editor")
	for _, c := range []*http.Cookie{nil, {Name: "sortable-session", Value: "unknown"}} {
		if w := request(a, "POST", "/sortable/lists/reading/move", form(csrf, "essay-02", "essay-01", "1").Encode(), c, false); w.Code != 403 {
			t.Fatal("missing session admitted")
		}
	}
	a.mu.Lock()
	a.sessions[cookie.Value].Actor = "unknown"
	a.mu.Unlock()
	if w := request(a, "POST", "/sortable/lists/reading/move", form(csrf, "essay-02", "essay-01", "1").Encode(), cookie, false); w.Code != 403 {
		t.Fatal("unknown actor admitted")
	}
	a.mu.Lock()
	a.sessions[cookie.Value].Actor = "editor"
	a.sessions[cookie.Value].Expires = time.Now().Add(-time.Second)
	a.mu.Unlock()
	if w := request(a, "POST", "/sortable/lists/reading/move", form(csrf, "essay-02", "essay-01", "1").Encode(), cookie, false); w.Code != 403 {
		t.Fatal("expired session admitted")
	}
	if grant("editor", "unknown", "reading") || grant("editor", "read", "unknown") {
		t.Fatal("default deny failed")
	}
}
func TestStrictValidationAndCSRF(t *testing.T) {
	a := testApp(t)
	cookie, csrf := login(t, a, "editor")
	before, _ := os.ReadFile(a.store.path)
	for _, bad := range []string{"missing-before", "duplicate", "cross-list", "self", "revision", "csrf", "unknown-field", "missing-item"} {
		t.Run(bad, func(t *testing.T) {
			f := form(csrf, "essay-02", "essay-01", "1")
			switch bad {
			case "missing-before":
				f.Del("before")
			case "duplicate":
				f.Add("itemId", "essay-03")
			case "cross-list":
				f.Set("before", "governance")
			case "self":
				f.Set("before", "essay-02")
			case "revision":
				f.Set("revision", "01")
			case "csrf":
				f.Set("csrf", "bad")
			case "unknown-field":
				f.Set("actor", "editor")
			case "missing-item":
				f.Set("itemId", "missing")
			}
			if w := request(a, "POST", "/sortable/lists/reading/move", f.Encode(), cookie, false); w.Code < 400 {
				t.Fatal("invalid request accepted", bad, responseBody(w))
			}
		})
	}
	for _, origin := range []string{"https://evil.test", "null", "http://127.0.0.1:9999"} {
		r := withCookie(cookie)
		r.Method = "POST"
		r.URL.Path = "/sortable/lists/reading/move"
		r.Body = http.NoBody
		r.Header.Set("Origin", origin)
		if sameOrigin(r) {
			t.Fatal("disallowed origin accepted")
		}
	}
	jsonBody := encode(payload(csrf, "essay-02", "essay-01", "reading", 1))
	jsonBody = strings.TrimSuffix(jsonBody, "}") + `,"itemId":"essay-03"}`
	if w := request(a, "POST", "/sortable/lists/reading/move", jsonBody, cookie, true); w.Code != 400 {
		t.Fatal("duplicate JSON admitted", w.Code)
	}
	for _, path := range []string{"/sortable/lists/unknown/move", "/sortable/lists/reading/move?itemId=essay-03"} {
		if w := request(a, "POST", path, form(csrf, "essay-02", "essay-01", "1").Encode(), cookie, false); w.Code < 400 {
			t.Fatal("unknown path/query admitted")
		}
	}
	after, _ := os.ReadFile(a.store.path)
	if string(before) != string(after) {
		t.Fatal("invalid requests wrote state")
	}
}
func nodes(t *testing.T, body string, match func(*html.Node) bool) []*html.Node {
	t.Helper()
	root, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	out := []*html.Node{}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if match(n) {
			out = append(out, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)
	return out
}
func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
func TestNativeFormsAndPresentation(t *testing.T) {
	a := testApp(t)
	cookie, csrf := login(t, a, "editor")
	for _, list := range []string{"reading", "requirements"} {
		item, before := "essay-02", "essay-01"
		if list == "requirements" {
			item, before = "governance", "orientation"
		}
		w := request(a, "POST", "/sortable/lists/"+list+"/move", form(csrf, item, before, "1").Encode(), cookie, false)
		if w.Code != 303 || w.Header().Get("Location") != "/sortable" {
			t.Fatal("native redirect missing")
		}
		if w := request(a, "POST", "/sortable/lists/"+list+"/move", form(csrf, item, before, "1").Encode(), cookie, false); w.Code != 409 {
			t.Fatal("replay admitted")
		}
	}
	page := request(a, "GET", "/sortable?controls=arrows", "", cookie, false).Body.String()
	forms := nodes(t, page, func(n *html.Node) bool { return n.Data == "form" && strings.HasSuffix(attr(n, "action"), "/move") })
	if len(forms) != 12 || strings.Contains(page, "data-sortable-handle") {
		t.Fatal("arrows fallback forms missing")
	}
	for _, f := range forms {
		inputs := map[string]string{}
		var walk func(*html.Node)
		walk = func(n *html.Node) {
			if n.Data == "input" {
				inputs[attr(n, "name")] = attr(n, "value")
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				walk(c)
			}
		}
		walk(f)
		if len(inputs) != 4 || inputs["revision"] != "2" || inputs["csrf"] != csrf {
			t.Fatal("native metadata stale", inputs)
		}
		if inputs["itemId"] == "essay-02" && attr(f.LastChild, "data-sortable-move") == "up" && inputs["before"] != "essay-02" {
			t.Fatal("upper boundary not updated")
		}
	}
	page = request(a, "GET", "/sortable?variant=dividers&density=compact&controls=handle", "", cookie, false).Body.String()
	if strings.Contains(page, "data-sortable-variant") || strings.Contains(page, `action="/sortable/lists/reading/move"`) || !strings.Contains(page, `data-sortable-density="compact"`) || !strings.Contains(page, "requires JavaScript") {
		t.Fatal("presentation choices invalid")
	}
	if w := request(a, "GET", "/sortable?controls=bad", "", cookie, false); w.Code != 400 {
		t.Fatal("unallowlisted choice accepted")
	}
}
func TestOrderedSSEAndCorrelation(t *testing.T) {
	a := testApp(t)
	cookie, csrf := login(t, a, "editor")
	for _, list := range []string{"reading", "requirements"} {
		item, before := "essay-02", "essay-01"
		if list == "requirements" {
			item, before = "governance", "orientation"
		}
		body := encode(payload(csrf, item, before, list, 1))
		w := request(a, "POST", "/sortable/lists/"+list+"/move", body, cookie, true)
		s := responseBody(w)
		elements, signals := strings.Index(s, "event: datastar-patch-elements"), strings.Index(s, "event: datastar-patch-signals")
		if w.Code != 200 || elements < 0 || signals <= elements || !strings.Contains(s, "data: mode replace") || !strings.Contains(s, patchSelector(list, "connection-1")) || !strings.Contains(s, `data-delivery="request-1"`) || !strings.Contains(s, `data-revision="2"`) || !strings.Contains(s, `aria-disabled="true"`) || !strings.Contains(s, `"status":"accepted"`) {
			t.Fatal("invalid SSE", s)
		}
		if strings.Contains(s, "<lw-sortable-list") || strings.Contains(s, "data-sortable-status") {
			t.Fatal("patch replaces host/status")
		}
		w = request(a, "POST", "/sortable/lists/"+list+"/move", body, cookie, true)
		if !strings.Contains(responseBody(w), `"status":"rejected"`) || a.store.snapshot(list).Revision != 2 {
			t.Fatal("enhanced replay not stale")
		}
	}
	if cssString(`a"b]`) != `"\61 \22 \62 \5d "` {
		t.Fatal("unsafe selector quoting")
	}
}

func TestNativeAppendAndNoOp(t *testing.T) {
	a := testApp(t)
	cookie, csrf := login(t, a, "editor")
	if w := request(a, "POST", "/sortable/lists/reading/move", form(csrf, "essay-01", "", "1").Encode(), cookie, false); w.Code != 303 {
		t.Fatal("explicit empty append rejected", w.Code)
	}
	state := a.store.snapshot("reading")
	if state.Order[2] != "essay-01" || state.Revision != 2 || a.store.snapshot("requirements").Revision != 1 {
		t.Fatal("append changed wrong list", state)
	}
	before, _ := os.ReadFile(a.store.path)
	if w := request(a, "POST", "/sortable/lists/reading/move", encode(payload(csrf, "essay-01", "", "reading", 2)), cookie, true); !strings.Contains(responseBody(w), `"status":"accepted"`) {
		t.Fatal("enhanced no-op rejected")
	}
	after, _ := os.ReadFile(a.store.path)
	if string(before) != string(after) {
		t.Fatal("current-revision no-op wrote file")
	}
}
func TestConcurrentEndpointsAndExampleCases(t *testing.T) {
	a := testApp(t)
	cookie, csrf := login(t, a, "editor")
	w := request(a, "POST", "/sortable/lists/reading/case", encode(map[string]string{"csrf": csrf, "case": "delay"}), cookie, true)
	if w.Code != 200 {
		t.Fatal(responseBody(w))
	}
	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		defer wg.Done()
		done <- request(a, "POST", "/sortable/lists/reading/move", encode(payload(csrf, "essay-02", "essay-01", "reading", 1)), cookie, true)
	}()
	deadline := time.Now().Add(time.Second)
	for {
		a.mu.Lock()
		mode := a.sessions[cookie.Value].Cases["reading"]
		a.mu.Unlock()
		if mode == "" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("delayed request not started")
		}
		time.Sleep(time.Millisecond)
	}
	start := time.Now()
	w = request(a, "POST", "/sortable/lists/requirements/move", encode(payload(csrf, "governance", "orientation", "requirements", 1)), cookie, true)
	if time.Since(start) > time.Second || !strings.Contains(responseBody(w), `"status":"accepted"`) || a.store.snapshot("reading").Revision != 1 {
		t.Fatal("delay blocked independent list")
	}
	wg.Wait()
	if !strings.Contains(responseBody(<-done), `"status":"accepted"`) {
		t.Fatal("delay never completed")
	}
	for _, mode := range []string{"reject", "stale", "missing", "invalid", "unrelated"} {
		t.Run(mode, func(t *testing.T) {
			request(a, "POST", "/sortable/lists/reading/case", encode(map[string]string{"csrf": csrf, "case": mode}), cookie, true)
			rev := a.store.snapshot("reading").Revision
			w := request(a, "POST", "/sortable/lists/reading/move", encode(payload(csrf, "essay-03", "essay-02", "reading", rev)), cookie, true)
			body := responseBody(w)
			if mode == "reject" || mode == "unrelated" {
				if a.store.snapshot("reading").Revision != rev || !strings.Contains(body, `"status":"rejected"`) {
					t.Fatal("rejection wrote state")
				}
			}
			if mode == "stale" && !strings.Contains(body, `"requestId":"request-1-old"`) || mode == "invalid" && !strings.Contains(body, `"status":"unknown"`) || mode == "unrelated" && !strings.Contains(body, `data-delivery="unrelated"`) {
				t.Fatal("case delivery missing", body)
			}
			if mode == "missing" {
				first := strings.Split(body, "\n\n")[0]
				if strings.Contains(first, `"status"`) {
					t.Fatal("missing-status case fabricated status")
				}
			}
			if strings.LastIndex(body, "event: datastar-patch-signals") <= strings.LastIndex(body, "event: datastar-patch-elements") {
				t.Fatal("real result not after authoritative patch")
			}
		})
	}
}
func TestReferenceBytesAndAllowlist(t *testing.T) {
	a := testApp(t)
	b, err := os.ReadFile("../../plans/sortable-components-9/producer/sortable-list.html")
	if err != nil {
		t.Fatal(err)
	}
	locked, _ := authored.ReadFile("assets.lock.json")
	var lock assetLock
	json.Unmarshal(locked, &lock)
	for _, v := range lock.Assets {
		if v.Serve == "/manual/components/sortable-list.html" && digest(b) != v.SHA256 {
			t.Fatal("reference evidence/lock mismatch")
		}
	}
	a.files["/manual/components/sortable-list.html"] = b
	w := request(a, "GET", "/manual/components/sortable-list.html", "", nil, false)
	if responseBody(w) != string(b) || w.Header().Get("X-Asset-SHA256") != digest(b) {
		t.Fatal("reference bytes altered")
	}
	for _, path := range []string{"/orders.json", "/data/orders.json", "/assets.lock.json", "/../orders.json", "/manual/components%2fsortable-list.html", "/%2e%2e/orders.json"} {
		w := request(a, "GET", path, "", nil, false)
		if w.Code == 200 {
			t.Fatal("unsafe path served", path)
		}
	}
	if w := request(a, "POST", "/manual/components/sortable-list.html", "", nil, false); w.Code != 405 {
		t.Fatal("asset writable")
	}
	if _, _, err := loadAssets(t.TempDir()); err == nil {
		t.Fatal("missing assets accepted")
	}
}
func TestHealthReceiptAndRestartSessions(t *testing.T) {
	a := testApp(t)
	cookie, csrf := login(t, a, "editor")
	w := request(a, "GET", "/healthz", "", nil, false)
	if !strings.Contains(responseBody(w), a.receipt.Fingerprint) {
		t.Fatal("receipt absent")
	}
	s, err := openStore(a.store.path)
	if err != nil {
		t.Fatal(err)
	}
	restarted := newApp(s, nil, assetLock{}, a.receipt)
	if w := request(restarted, "POST", "/sortable/lists/reading/move", form(csrf, "essay-02", "essay-01", "1").Encode(), cookie, false); w.Code != 403 {
		t.Fatal("restart retained process-local sessions")
	}
	if !strings.Contains(request(a, "GET", "/sortable", "", cookie, false).Body.String(), fmt.Sprintf(`data-source-head="%s"`, a.receipt.Commit)) {
		t.Fatal("visible source marker absent")
	}
}
