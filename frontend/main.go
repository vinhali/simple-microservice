package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vinhali/simple-microservice/otelsetup"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	authenticated bool
	authMu        sync.Mutex

	authTotalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "web_banking_auth_total",
			Help: "Total number of auth attempts",
		},
		[]string{"result"},
	)

	authSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "web_banking_auth_success_total",
			Help: "Total number of successful auths",
		},
	)

	authFailureCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "web_banking_auth_failure_total",
			Help: "Total number of failed auths",
		},
	)
	tracer = otel.Tracer("auth")
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	serviceName := "digital-bank"
	serviceVersion := "1.0.0"
	otelShutdown, err := otelsetup.SetupOTelSDK(ctx, serviceName, serviceVersion)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	time.Sleep(2 * time.Second)

	prometheus.MustRegister(authTotalCounter)
	prometheus.MustRegister(authSuccessCounter)
	prometheus.MustRegister(authFailureCounter)

	go func() {
		defer wg.Done()
		startFrontend()
	}()

	wg.Wait()
}

func clearCookies(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "authpass",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

func startFrontend() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/auth", frontHandler)

	http.Handle("/style.css", http.FileServer(http.Dir("/bin/files/css")))
	http.Handle("/script.js", http.FileServer(http.Dir("/bin/files/js")))
	log.Println("Starting web-banking frontend server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func frontHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	renderauthForm(w)

	authTotalCounter.WithLabelValues("attempts").Inc()

	authpassCookie, err := r.Cookie("authpass")

	if err == nil {
		parentCtx, parentSpan := tracer.Start(r.Context(), "/auth")
		parentSpan.SetAttributes(semconv.HTTPStatusCodeKey.Int(200))
		authInfo := attribute.String("auth.info", "call sent successfully for authentication")
		parentSpan.SetAttributes(authInfo)
		defer parentSpan.End()
		span := parentSpan
		if authpassCookie.Value == "passed" {
			authMu.Lock()
			authenticated = true
			authMu.Unlock()

			authTotalCounter.WithLabelValues("success").Inc()
			authSuccessCounter.Inc()

			parentCtx, span = tracer.Start(parentCtx, "auth.success")
			span.SetAttributes(semconv.HTTPStatusCodeKey.Int(200))
			authInfo := attribute.String("auth.info", "authentication completed successfully")
			span.SetAttributes(authInfo)
			defer span.End()

			clearCookies(w)
			callTransferAPI(w, span, parentCtx)
			return

		} else if authpassCookie.Value == "rejected" {

			authTotalCounter.WithLabelValues("failure").Inc()
			authFailureCounter.Inc()

			parentCtx, span = tracer.Start(parentCtx, "auth.failure")
			span.SetAttributes(semconv.HTTPStatusCodeKey.Int(401))
			authInfo := attribute.String("auth.info", "authentication failed")
			span.SetAttributes(authInfo)
			defer span.End()

			return
		}
	}
	return
}

func renderauthForm(w http.ResponseWriter) {
	htmlContent, err := os.ReadFile("/bin/files/index.html")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to read auth form", http.StatusInternalServerError)
		return
	}

	jsContent, err := os.ReadFile("/bin/files/js/script.js")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to read JavaScript file", http.StatusInternalServerError)
		return
	}

	cssContent, err := os.ReadFile("/bin/files/css/style.css")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to read CSS file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n<script>%s</script>\n<style>%s</style>\n", htmlContent, jsContent, cssContent)
}

func callTransferAPI(w http.ResponseWriter, parentSpan trace.Span, parentCtx context.Context) {
	enableCors(&w)
	authMu.Lock()
	defer authMu.Unlock()

	if !authenticated {
		log.Print(http.StatusUnauthorized)
		return
	}

	req, err := http.NewRequest("GET", "http://web-banking-backend-service:8081/transfer", nil)

	// Use o contexto e span da autenticação como contexto pai
	ctx, span := tracer.Start(parentCtx, "/transfer")
	defer span.End()

	parentSpanContext := trace.SpanContextFromContext(ctx)

	req.Header.Set("authenticated", "true")
	req.Header.Set("traceId", parentSpanContext.TraceID().String())
	req.Header.Set("spanId", parentSpanContext.SpanID().String())

	if err != nil {
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(500))
		transferForward := attribute.String("transfer.forward", "server not found")
		span.SetAttributes(transferForward)
	} else {
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(200))
		transferForward := attribute.String("transfer.forward", "request sent successfully")
		span.SetAttributes(transferForward)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
}
