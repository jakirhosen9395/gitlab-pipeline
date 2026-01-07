package main

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// tiny template so calculatorHandler can ExecuteTemplate("index.html", ...)
func setTestTemplate(t *testing.T) {
	t.Helper()
	var err error
	tpl, err = template.New("index.html").Parse(`
		{{- if .Result -}}{{.Result}}{{- end -}}
		{{- if .Error -}}{{.Error}}{{- end -}}
	`)
	if err != nil {
		t.Fatalf("template parse failed: %v", err)
	}
}

func TestCalculatorHandler_GET_OK(t *testing.T) {
	setTestTemplate(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	calculatorHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
}

func TestCalculatorHandler_Add(t *testing.T) {
	setTestTemplate(t)

	form := url.Values{
		"num1":     {"3"},
		"num2":     {"4"},
		"operator": {"add"},
	}
	req := httptest.NewRequest(http.MethodPost, "/calculator", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	calculatorHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "7") {
		t.Fatalf("want result 7 in body, got: %q", body)
	}
}

func TestCalculatorHandler_DivideByZero(t *testing.T) {
	setTestTemplate(t)

	form := url.Values{
		"num1":     {"10"},
		"num2":     {"0"},
		"operator": {"divide"},
	}
	req := httptest.NewRequest(http.MethodPost, "/calculator", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	calculatorHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Cannot divide by zero") {
		t.Fatalf("expected divide-by-zero error, got: %q", rr.Body.String())
	}
}

func TestCalculatorHandler_MethodNotAllowed(t *testing.T) {
	setTestTemplate(t)

	req := httptest.NewRequest(http.MethodPut, "/calculator", nil)
	rr := httptest.NewRecorder()

	calculatorHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Method not allowed") {
		t.Fatalf("expected method-not-allowed message, got: %q", rr.Body.String())
	}
}

func TestWithReqID_SetsHeader(t *testing.T) {
	h := withReqID(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rr := httptest.NewRecorder()

	h(rr, req)

	if got := rr.Header().Get("X-Request-ID"); got == "" {
		t.Fatalf("expected X-Request-ID header to be set")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
}

func TestHealthz_ReturnsServiceNameHealthy(t *testing.T) {
	handler := healthzHandler("calculator")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "calculator healthy") {
		t.Fatalf("want 'calculator healthy', got: %q", body)
	}
}
