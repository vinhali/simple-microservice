package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vinhali/simple-microservice/otelsetup"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var accountsMu sync.Mutex
var authenticatedAccounts = make(map[string]bool)

type Account struct {
	ID      string
	Balance float64
	mu      sync.Mutex
}

var accounts = map[string]*Account{
	"123456": {ID: "123456", Balance: 100000000.0},
	"789012": {ID: "789012", Balance: 10.0},
}

var transactions []Transaction

type Transaction struct {
	Message                   string  `json:"message"`
	SourceAccountBalance      float64 `json:"source_account_balance"`
	DestinationAccountBalance float64 `json:"destination_account_balance"`
	Amount                    float64 `json:"amount"`
}

var (
	accountBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "web_banking_account_balance",
			Help: "Current balance of an account",
		},
		[]string{"account_id"},
	)

	transactionStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "web_banking_transaction_status",
			Help: "Transaction status (success or failure)",
		},
		[]string{"account_id", "status"},
	)

	transactionAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "web_banking_transaction_amount",
			Help: "Amount transferred in a transaction",
		},
		[]string{"source_account_id", "destination_account_id"},
	)
	tracer = otel.Tracer("transfer")
)

func main() {
	http.HandleFunc("/transfer", transactionHandler)
	http.Handle("/output", http.HandlerFunc(outputHandler))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	serviceName := "digital-bank"
	serviceVersion := "1.0.0"
	otelShutdown, err := otelsetup.SetupOTelSDK(ctx, serviceName, serviceVersion)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	prometheus.MustRegister(accountBalance)
	prometheus.MustRegister(transactionStatus)
	prometheus.MustRegister(transactionAmount)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Starting web-banking backend server at port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, authenticated")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func transactionHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authHeader := r.Header.Get("authenticated")
	traceId := r.Header.Get("traceId")
	spanId := r.Header.Get("spanId")

	if authHeader != "true" {
		log.Println(http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		log.Println(http.StatusNotFound)
		return
	}

	spanContext, err := createSpanContext(traceId, spanId)
	if err != nil {
		log.Printf("Erro ao criar SpanContext: %v\n", err)
		http.Error(w, "Erro ao criar SpanContext", http.StatusInternalServerError)
		return
	}

	sourceAccountID := "123456"
	destinationAccountID := "789012"
	amount := rand.Float64() * 20

	ctx, span := tracer.Start(
		trace.ContextWithRemoteSpanContext(r.Context(), spanContext),
		"transfer.status",
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("source_account_id", sourceAccountID),
		attribute.String("destination_account_id", destinationAccountID),
	)

	if err := transferAmount(ctx, sourceAccountID, destinationAccountID, amount); err != nil {
		transactionStatus.WithLabelValues(sourceAccountID, "failure").Inc()

		http.Error(w, fmt.Sprintf("Failed to transfer amount: %v", err), http.StatusInternalServerError)
		span.SetAttributes(attribute.String("transfer.status", "failure"))

		return
	}

	accountBalance.WithLabelValues(sourceAccountID).Set(getAccountBalance(sourceAccountID))
	accountBalance.WithLabelValues(destinationAccountID).Set(getAccountBalance(destinationAccountID))
	transactionStatus.WithLabelValues(sourceAccountID, "success").Inc()
	transactionAmount.WithLabelValues(sourceAccountID, destinationAccountID).Set(amount)

	span.SetAttributes(
		attribute.String("transfer.status", "success"),
		attribute.Float64("transfer.amount", amount),
	)

	transaction := Transaction{
		Message:                   "Transaction completed successfully",
		SourceAccountBalance:      getAccountBalance(sourceAccountID),
		DestinationAccountBalance: getAccountBalance(destinationAccountID),
		Amount:                    amount,
	}

	transactions = append(transactions, transaction)
}

func outputHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(transactions)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func createSpanContext(traceId, spanId string) (trace.SpanContext, error) {
	traceID, err := trace.TraceIDFromHex(traceId)
	if err != nil {
		return trace.SpanContext{}, err
	}

	spanID, err := trace.SpanIDFromHex(spanId)
	if err != nil {
		return trace.SpanContext{}, err
	}

	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	}), nil
}

func transferAmount(ctx context.Context, sourceID, destinationID string, amount float64) error {
	sourceAccount, ok := accounts[sourceID]
	if !ok {
		transactionStatus.WithLabelValues(sourceID, "failure").Inc()
		return fmt.Errorf("source account not found")
	}

	destinationAccount, ok := accounts[destinationID]
	if !ok {
		transactionStatus.WithLabelValues(sourceID, "failure").Inc()
		return fmt.Errorf("destination account not found")
	}

	sourceAccount.mu.Lock()
	defer sourceAccount.mu.Unlock()

	destinationAccount.mu.Lock()
	defer destinationAccount.mu.Unlock()

	if sourceAccount.Balance < amount {
		transactionStatus.WithLabelValues(sourceID, "failure").Inc()
		return fmt.Errorf("insufficient funds")
	}

	sourceAccount.Balance -= amount
	destinationAccount.Balance += amount

	return nil
}

func getAccountBalance(accountID string) float64 {
	accountsMu.Lock()
	defer accountsMu.Unlock()

	account, ok := accounts[accountID]
	if !ok {
		return 0.0
	}

	return account.Balance
}
