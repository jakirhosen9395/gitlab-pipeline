package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const AppVersion = "v1"

type CalculatorResult struct {
	Num1     string
	Num2     string
	Operator string
	Result   string
	Error    string
	Version  string
}

var tpl *template.Template

func main() {
	// parse template once
	var err error
	tpl, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("parse template: %v", err)
	}

	// configurable bind address
	host := strings.TrimSpace(os.Getenv("HOST"))
	if host == "" {
		host = "0.0.0.0" // all interfaces; use 127.0.0.1 for local-only
	}
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "9000"
	}
	addr := host + ":" + port

	// service name for /healthz output
	serviceName := strings.TrimSpace(os.Getenv("SERVICE_NAME"))
	if serviceName == "" {
		serviceName = "calculator"
	}

	// routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", withReqID(calculatorHandler))
	mux.HandleFunc("/calculator", withReqID(calculatorHandler))
	mux.HandleFunc("/healthz", withReqID(healthzHandler(serviceName)))

	log.Printf("%s %s listening on %s", serviceName, AppVersion, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// middleware: add request id + latency log
func withReqID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		start := time.Now()
		w.Header().Set("X-Request-ID", id)
		next(w, r)
		log.Printf("[%s] %s %s %v", id, r.Method, r.URL.Path, time.Since(start))
	}
}

func healthzHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "%s healthy\n", serviceName)
	}
}

// GET: form, POST: compute and render
func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	data := CalculatorResult{Version: AppVersion}

	switch r.Method {
	case http.MethodGet:
		// just render
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			data.Error = "Failed to parse form"
			break
		}
		data.Num1 = strings.TrimSpace(r.FormValue("num1"))
		data.Num2 = strings.TrimSpace(r.FormValue("num2"))
		data.Operator = strings.TrimSpace(r.FormValue("operator"))

		if data.Num1 == "" || data.Num2 == "" || data.Operator == "" {
			data.Error = "Missing parameters"
			break
		}
		n1, err1 := strconv.ParseFloat(data.Num1, 64)
		n2, err2 := strconv.ParseFloat(data.Num2, 64)
		if err1 != nil || err2 != nil {
			data.Error = "Invalid numbers"
			break
		}

		switch data.Operator {
		case "add":
			data.Result = strconv.FormatFloat(n1+n2, 'f', -1, 64)
		case "subtract":
			data.Result = strconv.FormatFloat(n1-n2, 'f', -1, 64)
		case "multiply":
			data.Result = strconv.FormatFloat(n1*n2, 'f', -1, 64)
		case "divide":
			if n2 == 0 {
				data.Error = "Cannot divide by zero"
			} else {
				data.Result = strconv.FormatFloat(n1/n2, 'f', -1, 64)
			}
		default:
			data.Error = "Invalid operator (use: add, subtract, multiply, divide)"
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		data.Error = "Method not allowed"
	}

	if err := tpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("template execute: %v", err)
	}
}
