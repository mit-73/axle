/**
 * ConnectRPC transport for the BFF service.
 * Import clients from here throughout the app.
 */
import { createConnectTransport } from '@connectrpc/connect-web'

export const transport = createConnectTransport({
  baseUrl: import.meta.env.VITE_BFF_URL ?? 'http://localhost:9001',
})
