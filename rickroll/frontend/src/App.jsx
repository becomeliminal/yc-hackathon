import { useState } from 'react'
import { useAccount } from 'wagmi'
import { ConnectButton } from '@rainbow-me/rainbowkit'
import RickrollButton from './RickrollButton'
import './App.css'

function App() {
  const { isConnected } = useAccount()
  const [videoUrl, setVideoUrl] = useState(null)

  return (
    <div className="App">
      <header>
        <h1>🎁 Premium Content Access</h1>
        <p>Unlock exclusive content for just $0.01 USDC</p>
      </header>

      <main>
        <div className="connect-section">
          <ConnectButton />
        </div>

        {isConnected && !videoUrl && (
          <div className="payment-section">
            <RickrollButton onSuccess={setVideoUrl} />
          </div>
        )}

        {videoUrl && (
          <div className="video-section">
            <h2>Enjoy your content! 🎉</h2>
            <iframe
              width="560"
              height="315"
              src={videoUrl}
              title="Premium Content"
              frameBorder="0"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowFullScreen
            ></iframe>
            <button
              onClick={() => setVideoUrl(null)}
              style={{ marginTop: '20px' }}
            >
              Get More Content
            </button>
          </div>
        )}

        {!isConnected && (
          <div className="info-section">
            <p>Connect your wallet to get started</p>
            <p className="small">Make sure you're on Base Mainnet with USDC</p>
          </div>
        )}
      </main>

      <footer>
        <p>Powered by x402 protocol + your custom facilitator</p>
      </footer>
    </div>
  )
}

export default App
