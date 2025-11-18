# x402 Payment Gateway Demo

**YC Agentic Payments Hackathon - November 16, 2025**
**Stripe Track Submission**

This demo showcases how enterprise merchants can accept HTTP 402 payments using the x402 protocol. Built with [Liminal](https://becomeliminal.com), a global stablecoin neobank, and powered by our open-source [grpc-gateway-x402 middleware](https://github.com/becomeliminal/grpc-gateway-x402) and the [Liminal Facilitator](https://facilitator.liminal.cash).

## Overview

This project demonstrates what a merchant wanting to accept x402 payments can do with Stripe-like payment infrastructure. It shows the complete end-to-end flow:

1. **Content Gating** - Merchant protects premium content behind a paywall
2. **USDC Payment** - User pays $0.01 USDC via gasless EIP-3009 signature
3. **Settlement** - Payment settles directly into merchant's Liminal wallet
4. **Offramp** - Merchant can convert to fiat via Bridge using Privy-signed transactions

[Liminal](https://becomeliminal.com) is a global stablecoin neobank that uses Privy as a signer on Safe wallets and Bridge for on/offramping. At Liminal, we run a chain-agnostic, ERC-20 agnostic facilitator at `facilitator.liminal.cash` that can settle any EIP-3009 compatible token on any EVM chain.

## The Stack

- **Frontend**: React + Vite + wagmi + RainbowKit (Arbitrum One)
- **Backend**: gRPC service with grpc-gateway + x402 middleware
- **Payment Protocol**: x402 with EIP-3009 (USDC transferWithAuthorization)
- **Wallet**: [Liminal](https://becomeliminal.com) - global stablecoin neobank (Privy signer on Safe)
- **Facilitator**: Chain-agnostic, ERC-20 agnostic payment verification at facilitator.liminal.cash
- **Onramp/Offramp**: Bridge for fiat conversion
- **Teleport**: Bank-to-bank instant transfers between Liminal wallets

## How It Works

1. User connects wallet (MetaMask, Coinbase Wallet, or WalletConnect)
2. Clicks "Unlock Content" and receives `402 Payment Required` response
3. Signs EIP-3009 authorization for $0.01 USDC on Arbitrum (gasless, no transaction)
4. Liminal facilitator verifies the payment signature and settles directly into merchant's Liminal wallet
5. Content unlocks with autoplay
6. Merchant can:
   - Offramp to fiat via Bridge using Privy-signed transactions
   - Teleport funds bank-to-bank to other Liminal wallets
   - Earn high-yield on deposited stablecoins

## Why This Matters for Enterprise Merchants

Enterprise merchants using high-performance microservice architectures can integrate x402 payments with:

- **Production-ready middleware**: Open source [grpc-gateway-x402](https://github.com/becomeliminal/grpc-gateway-x402) for gRPC services
- **Liminal wallet as main business bank**: Direct settlement with high-yield savings
- **No intermediaries**: Facilitator settles directly into merchant's wallet
- **Chain-agnostic**: Facilitator supports any EIP-3009 compatible token on any EVM chain
- **Seamless fiat offramp**: Convert to fiat via Bridge
- **Instant transfers**: Bank-to-bank transfers via Liminal Teleport

## Running the Demo

### Prerequisites

- Go 1.21+
- Node.js 18+
- Please build system installed

### Backend

```bash
# For production (using Liminal facilitator)
plz run //rickroll/service:service -- \
  --recipient-address 0x37BaE31Bd75020DF9C4A7EE86c9A53Ba2B27F10E \
  --facilitator-url https://facilitator.liminal.cash/v1/x402

# For local testing (hackathon setup)
plz run //rickroll/service:service -- \
  --recipient-address 0x37BaE31Bd75020DF9C4A7EE86c9A53Ba2B27F10E \
  --facilitator-url http://localhost:8080/v1/x402 \
  --http-port 8081
```

### Frontend

```bash
cd rickroll/frontend
npm install
npm run dev
```

### Environment Variables

- `VITE_SERVICE_URL` - Backend URL (default: `http://localhost:8081`)
- `VITE_WC_PROJECT_ID` - WalletConnect project ID (optional)

## Project Structure

```
rickroll/
├── proto/          # gRPC/Protobuf definitions
├── service/        # Go backend with x402 middleware
│   ├── main.go     # HTTP gateway + gRPC server
│   └── handler.go  # Service implementation
└── frontend/       # React app
    └── src/
        ├── App.jsx
        ├── RickrollButton.jsx  # Payment flow + EIP-3009 signing
        └── wagmi.js            # Wallet config
```

## Technical Details

- **Network**: Arbitrum One (chainId: 42161)
- **Token**: USDC (0xaf88d065e77c8cC2239327C5EDb3A432268e5831)
- **Price**: $0.01 USD per access
- **Payment Method**: EIP-3009 transferWithAuthorization (gasless, no transaction required)
- **Signature Standard**: EIP-712 typed data
- **Facilitator**: Chain-agnostic, supports any EIP-3009 compatible ERC-20 token on any EVM chain

## Links

- [Liminal](https://becomeliminal.com) - Global stablecoin neobank
- [grpc-gateway-x402 Middleware](https://github.com/becomeliminal/grpc-gateway-x402)
- [Liminal Facilitator](https://facilitator.liminal.cash)

## Built For

**YC Agentic Payments Hackathon - November 16, 2025**
**Stripe Track**

Demonstrating enterprise-grade x402 payment infrastructure for merchants.

## License

MIT
