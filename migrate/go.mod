module migrate

go 1.24.5

require (
	github.com/golang-migrate/migrate/v4 v4.18.1
	github.com/spf13/cobra v1.10.1
	stock-automation/database v0.0.0
	stock-automation/helper v0.0.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gorm.io/driver/mysql v1.6.0 // indirect
	gorm.io/gorm v1.31.0 // indirect
	stock-automation/schema v0.0.0 // indirect
)

replace (
	stock-automation/database => ../database
	stock-automation/helper => ../helper
	stock-automation/schema => ../schema
)
