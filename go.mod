module github.com/uroborosq-go-dfs/server

go 1.19

require (
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/uroborosq-go-dfs/models v0.0.0-20230120113109-1d73a1f01e9b
)

require (
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/uroborosq-go-dfs/models => ../models
