package main

import (
	"context"
	"fmt"

	pb "rickroll/proto"
	x402 "github.com/becomeliminal/grpc-gateway-x402"
)

// RickRollServer implements the RickRollService
type RickRollServer struct {
	pb.UnimplementedRickRollServiceServer
	rickrollCount int64
}

// NewRickRollServer creates a new RickRollServer
func NewRickRollServer() *RickRollServer {
	return &RickRollServer{
		rickrollCount: 0,
	}
}

// GetRickRoll returns the content (requires payment)
func (s *RickRollServer) GetRickRoll(ctx context.Context, req *pb.GetRickRollRequest) (*pb.GetRickRollResponse, error) {
	// Extract payment info from context (injected by x402 middleware)
	payment, ok := x402.GetPaymentFromGRPCContext(ctx)
	if !ok {
		// This shouldn't happen if x402 middleware is working
		// The middleware will return 402 before this handler is called
		return nil, fmt.Errorf("payment required but not found in context")
	}

	// Increment counter
	s.rickrollCount++

	// Build payment receipt from the x402 payment context
	receipt := &pb.PaymentReceipt{
		TransactionHash: payment.TransactionHash,
		AmountPaid:      payment.Amount,
		TokenSymbol:     payment.TokenSymbol,
		PayerAddress:    payment.PayerAddress,
		SettledAt:       payment.SettledAt.Unix(),
	}

	// Return the response with payment receipt
	// TODO: Add actual content here (lyrics, video URL, GIF, etc.)
	return &pb.GetRickRollResponse{
		Lyrics:   "Premium content goes here - add your own text",
		VideoUrl: "https://www.youtube.com/embed/dQw4w9WgXcQ",
		Gif:      []byte{}, // Add GIF bytes if needed
		Receipt:  receipt,
	}, nil
}

// GetInfo returns service information (free endpoint, no payment required)
func (s *RickRollServer) GetInfo(ctx context.Context, req *pb.GetInfoRequest) (*pb.GetInfoResponse, error) {
	return &pb.GetInfoResponse{
		Description:    "Pay-per-use service powered by x402 protocol",
		PriceUsd:       "0.01",
		TotalRickrolls: s.rickrollCount,
		AcceptedTokens: []*pb.TokenInfo{
			{
				Symbol:          "USDC",
				Network:         "arbitrum-one",
				ContractAddress: "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
			},
		},
	}, nil
}
