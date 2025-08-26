module db-transaction

go 1.23

require (
	db-testing v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
	go.uber.org/zap v1.27.0
	gorm.io/gorm v1.25.5
)

replace db-testing => ../db-testing

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.5.4 // indirect
)
