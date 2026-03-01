module github.com/ApeironFoundation/axle/bff

go 1.26

replace github.com/ApeironFoundation/axle/contracts => ../../contracts/generated

replace github.com/ApeironFoundation/axle/db => ../../db

require (
	connectrpc.com/connect v1.19.1
	github.com/ApeironFoundation/axle/contracts v0.0.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/nats-io/nats.go v1.39.1
	github.com/redis/go-redis/v9 v9.7.3
	github.com/rs/cors v1.11.1
	github.com/rs/zerolog v1.33.0
	golang.org/x/net v0.42.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/stretchr/testify v1.11.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
