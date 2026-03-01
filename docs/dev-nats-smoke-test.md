# Dev-only NATS smoke test (BFF ↔ Gateway ↔ LLM)

This scenario validates inter-service communication through NATS using protobuf payloads.

> Scope: development only. Do not enable in production.

## Preconditions

- Infra is running: `just dev-up`
- Contracts are generated: `just contracts-generate`
- `ENABLE_DEV_ENDPOINTS=true` for BFF and LLM (default is `true`)
- Services are running in separate terminals:
  - `just bff-run`
  - `just gateway-run`
  - `just llm-run`

## Dev-only subjects and endpoint

- Event publish subject (BFF -> Gateway/LLM): `axle.events.test.ping`
- Request-reply subject (BFF -> LLM): `axle.test.ping.rpc`
- Dev-only ConnectRPC endpoint on BFF: `test.v1.TestService/Ping`

## Protobuf payloads

- `test.v1.PingRequest` is serialized with protobuf and used in:
  - BFF -> NATS request-reply payload
  - `gateway.v1.Event.payload` bytes for pub/sub
- `test.v1.PingReply` is serialized with protobuf and sent back by LLM over NATS request-reply
- Gateway streams `gateway.v1.Event` over ConnectRPC streaming (`Subscribe`)

## Test scenario

1. Open a streaming ConnectRPC connection to Gateway:
   - Service: `gateway.v1.StreamingService`
   - Method: `Subscribe`
2. Call BFF unary ConnectRPC method:
   - Service: `test.v1.TestService`
   - Method: `Ping`
   - Example JSON body: `{ "message": "ping" }`
3. BFF publishes `gateway.v1.Event` to `axle.events.test.ping`.
4. Gateway receives the event from NATS and pushes it to active `Subscribe` streams.
5. LLM receives:
   - the pub/sub event on `axle.events.test.ping` (logs it)
   - the request-reply call on `axle.test.ping.rpc` (returns `PingReply`, logs it)

## Expected logs

### BFF

- `DEV-ONLY endpoint enabled: test.v1.TestService/Ping`
- `dev-only ping event published`
- `dev-only ping reply received`

### Gateway

- `dev-only ping event received by gateway`

### LLM

- `DEV-ONLY NATS subscriptions enabled: axle.events.test.ping, axle.test.ping.rpc`
- `dev-only ping event received by llm`
- `dev-only ping rpc request handled by llm`

## Production safety

Disable the dev-only behavior with:

- `ENABLE_DEV_ENDPOINTS=false`

When disabled:

- BFF does not register `test.v1.TestService/Ping`
- LLM does not register dev-only NATS subscriptions
