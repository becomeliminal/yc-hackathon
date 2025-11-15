import { useState } from 'react'
import { useAccount, useSignTypedData } from 'wagmi'
import { parseUnits, encodePacked, keccak256 } from 'viem'

const SERVICE_URL = import.meta.env.VITE_SERVICE_URL || 'http://localhost:8081'

function RickrollButton({ onSuccess }) {
  const { address } = useAccount()
  const { signTypedDataAsync } = useSignTypedData()
  const [status, setStatus] = useState('ready') // ready, fetching, signing, paying, error
  const [error, setError] = useState(null)

  const handleGetRickrolled = async () => {
    try {
      setError(null)
      setStatus('fetching')

      // Step 1: Fetch payment requirements (will get 402)
      const response = await fetch(`${SERVICE_URL}/v1/rickroll`)

      if (response.status !== 402) {
        throw new Error(`Expected 402, got ${response.status}`)
      }

      const paymentRequired = await response.json()
      console.log('Payment requirements:', paymentRequired)

      if (!paymentRequired.paymentRequirements || paymentRequired.paymentRequirements.length === 0) {
        throw new Error('No payment requirements returned')
      }

      const requirement = paymentRequired.paymentRequirements[0]

      // Step 2: Sign EIP-3009 payment
      setStatus('signing')
      const payment = await signEIP3009Payment(requirement, address, signTypedDataAsync)

      // Step 3: Retry request with payment
      setStatus('paying')
      const paidResponse = await fetch(`${SERVICE_URL}/v1/rickroll`, {
        method: 'GET',
        headers: {
          'X-PAYMENT': btoa(JSON.stringify(payment)),
          'Content-Type': 'application/json',
        },
      })

      if (!paidResponse.ok) {
        const errorData = await paidResponse.json()
        throw new Error(errorData.message || `Payment failed: ${paidResponse.status}`)
      }

      const result = await paidResponse.json()
      console.log('Success!', result)

      setStatus('ready')
      onSuccess(result.videoUrl)

    } catch (err) {
      console.error('Error:', err)
      setError(err.message)
      setStatus('error')
    }
  }

  const getButtonText = () => {
    switch (status) {
      case 'fetching': return 'Fetching payment details...'
      case 'signing': return 'Sign payment in wallet...'
      case 'paying': return 'Processing payment...'
      case 'error': return 'Try Again'
      default: return 'Unlock Content - $0.01 USDC'
    }
  }

  return (
    <div className="rickroll-button-container">
      <button
        onClick={handleGetRickrolled}
        disabled={status !== 'ready' && status !== 'error'}
        className={`unlock-button ${status}`}
      >
        {getButtonText()}
      </button>
      {error && (
        <div className="error-message">
          ❌ {error}
        </div>
      )}
      {status === 'signing' && (
        <div className="info-message">
          Please check your wallet and approve the signature request
        </div>
      )}
    </div>
  )
}

// Sign EIP-3009 transferWithAuthorization payment
async function signEIP3009Payment(requirement, userAddress, signTypedDataAsync) {
  // Generate random nonce
  const nonce = keccak256(encodePacked(['uint256'], [BigInt(Date.now())]))

  // Convert amount to token units (USDC has 6 decimals)
  const amountInUnits = parseUnits(requirement.maxAmountRequired, 6)

  // EIP-712 domain for USDC on Base Mainnet
  const domain = {
    name: 'USD Coin',
    version: '2',
    chainId: 8453, // Base Mainnet
    verifyingContract: requirement.assetContract,
  }

  // EIP-3009 transferWithAuthorization types
  const types = {
    TransferWithAuthorization: [
      { name: 'from', type: 'address' },
      { name: 'to', type: 'address' },
      { name: 'value', type: 'uint256' },
      { name: 'validAfter', type: 'uint256' },
      { name: 'validBefore', type: 'uint256' },
      { name: 'nonce', type: 'bytes32' },
    ],
  }

  // Message to sign
  const message = {
    from: userAddress,
    to: requirement.recipient,
    value: amountInUnits.toString(),
    validAfter: '0',
    validBefore: requirement.validBefore.toString(),
    nonce,
  }

  console.log('Signing message:', message)

  // Sign the typed data
  const signature = await signTypedDataAsync({
    domain,
    types,
    primaryType: 'TransferWithAuthorization',
    message,
  })

  // Return x402 payment format
  return {
    x402Version: 1,
    scheme: 'exact',
    network: requirement.network,
    payload: {
      signature,
      authorization: {
        from: userAddress,
        to: requirement.recipient,
        value: amountInUnits.toString(),
        validAfter: 0,
        validBefore: parseInt(requirement.validBefore),
        nonce,
        tokenContract: requirement.assetContract,
      },
    },
  }
}

export default RickrollButton
