module github.com/ApeironFoundation/axle/gateway

go 1.26

require (
	github.com/ApeironFoundation/axle/contracts/generated v0.0.0
	connectrpc.com/connect v1.19.1
	github.com/go-chi/chi/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/nats-io/nats.go v1.39.1
	github.com/redis/go-redis/v9 v9.7.3
	github.com/rs/cors v1.11.1
	github.com/rs/zerolog v1.33.0
	golang.org/x/net v0.38.0
	google.golang.org/protobuf v1.36.9
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace github.com/ApeironFoundation/axle/contracts/generated => ../../contracts/generated
