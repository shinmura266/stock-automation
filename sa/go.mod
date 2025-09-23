module sa

go 1.24.5

require (
	github.com/spf13/cobra v1.8.0
	gorm.io/gorm v1.31.0
	stock-automation/database v0.0.0
	stock-automation/schema v0.0.0
)

replace (
	stock-automation/database => ../database
	stock-automation/schema => ../schema
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/text v0.20.0 // indirect
	gorm.io/driver/mysql v1.6.0 // indirect
)
