import { http, createConfig } from 'wagmi'
import { arbitrum } from 'wagmi/chains'
import { coinbaseWallet, metaMask, walletConnect } from 'wagmi/connectors'

export const config = createConfig({
  chains: [arbitrum],
  connectors: [
    metaMask(),
    coinbaseWallet({ appName: 'Rickroll x402' }),
    walletConnect({ projectId: import.meta.env.VITE_WC_PROJECT_ID || 'demo' }),
  ],
  transports: {
    [arbitrum.id]: http(),
  },
})
