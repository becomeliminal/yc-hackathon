package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "rickroll/proto"
	x402 "github.com/becomeliminal/grpc-gateway-x402"
	"github.com/becomeliminal/grpc-gateway-x402/evm"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcPort         = flag.Int("grpc-port", 9090, "gRPC server port")
	httpPort         = flag.Int("http-port", 8080, "HTTP gateway port")
	facilitatorURL   = flag.String("facilitator-url", "https://facilitator.x402.org", "x402 facilitator URL")
	recipientAddress = flag.String("recipient-address", "", "Wallet address to receive payments (required)")
)

func main() {
	flag.Parse()

	// Validate required flags
	if *recipientAddress == "" {
		log.Fatal("--recipient-address is required")
	}

	// Start gRPC server in background
	grpcServer := startGRPCServer(*grpcPort)
	defer grpcServer.GracefulStop()

	// Start HTTP gateway with x402 middleware
	if err := startHTTPGateway(*httpPort, *grpcPort); err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}
}

func startGRPCServer(port int) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register RickRoll service
	rickrollService := NewRickRollServer()
	pb.RegisterRickRollServiceServer(grpcServer, rickrollService)

	// Start serving in background
	go func() {
		log.Printf("gRPC server listening on :%d", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	return grpcServer
}

func startHTTPGateway(httpPort, grpcPort int) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create EVM verifier for x402
	verifier, err := evm.NewEVMVerifier(*facilitatorURL)
	if err != nil {
		return fmt.Errorf("failed to create EVM verifier: %w", err)
	}

	// Configure x402 middleware
	x402Config := x402.Config{
		Verifier: verifier,
		EndpointPricing: map[string]x402.PricingRule{
			// Paywall the /v1/rickroll endpoint
			"/v1/rickroll": {
				Amount:      "0.01", // $0.01 USD
				Description: "Get premium content",
				AcceptedTokens: []x402.TokenRequirement{
					{
						Network:       "arbitrum-one",
						AssetContract: "0xaf88d065e77c8cC2239327C5EDb3A432268e5831", // USDC on Arbitrum One
						Symbol:        "USDC",
						Recipient:     *recipientAddress,
						TokenName:     "USD Coin",
						TokenDecimals: 6,
					},
				},
			},
		},
		// Skip payment for info endpoint
		SkipPaths: []string{
			"/v1/info",
			"/health",
		},
		ValidityDuration: 5 * time.Minute,
	}

	// Validate config
	if err := x402Config.Validate(); err != nil {
		return fmt.Errorf("invalid x402 config: %w", err)
	}

	// Create grpc-gateway mux with payment metadata propagation
	mux := runtime.NewServeMux(
		x402.WithPaymentMetadata(), // Propagate payment info to gRPC context
	)

	// Connect to local gRPC server
	endpoint := fmt.Sprintf("localhost:%d", grpcPort)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register RickRoll service handler
	err = pb.RegisterRickRollServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return fmt.Errorf("failed to register handler: %w", err)
	}

	// Wrap mux with x402 payment middleware
	handler := x402.PaymentMiddleware(x402Config)(mux)

	// Add logging wrapper
	loggingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)

		// Check for payment header
		paymentHeader := r.Header.Get("X-Payment")
		if paymentHeader != "" {
			log.Printf("[PAYMENT] Received X-Payment header (length: %d bytes)", len(paymentHeader))
			log.Printf("[PAYMENT] Payment data preview: %s...", paymentHeader[:min(100, len(paymentHeader))])
		} else {
			log.Printf("[PAYMENT] No X-Payment header found")
		}

		// Wrap the response writer to capture status code
		wrappedWriter := &loggingResponseWriter{ResponseWriter: w, statusCode: 200}

		// Call the x402 middleware
		handler.ServeHTTP(wrappedWriter, r)

		log.Printf("[RESPONSE] Status: %d", wrappedWriter.statusCode)
	})

	// Add CORS middleware
	corsHandler := addCORS(loggingHandler)

	// Add health check endpoint (no payment required)
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	healthMux.Handle("/", corsHandler)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: healthMux,
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down HTTP server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		cancel()
	}()

	log.Printf("HTTP gateway listening on :%d", httpPort)
	log.Printf("x402 payment middleware enabled for /v1/rickroll")
	log.Printf("Free endpoints: /v1/info, /health")
	log.Printf("Recipient address: %s", *recipientAddress)
	log.Printf("Facilitator: %s", *facilitatorURL)
	log.Printf("\nTry:")
	log.Printf("  curl http://localhost:%d/health", httpPort)
	log.Printf("  curl http://localhost:%d/v1/info", httpPort)
	log.Printf("  curl http://localhost:%d/v1/rickroll  # Will return 402 Payment Required", httpPort)

	// Start serving (blocking)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// addCORS adds CORS headers to allow requests from the frontend
func addCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from localhost (development)
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Payment")
		w.Header().Set("Access-Control-Expose-Headers", "X-Payment-Response")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
