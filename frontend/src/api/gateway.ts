import { createClient } from '@connectrpc/connect'
import { StreamingService } from '@axle/contracts/gateway/v1/streaming_pb'
import { createConnectTransport } from '@connectrpc/connect-web'

const gatewayTransport = createConnectTransport({
  baseUrl: import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:8081',
})

export const streamingClient = createClient(StreamingService, gatewayTransport)
