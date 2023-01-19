module watchmen

go 1.16

require (
	github.com/go-playground/validator/v10 v10.11.1
	github.com/go-sql-driver/mysql v1.7.0
	github.com/golang-jwt/jwt/v4 v4.4.3
	github.com/jmoiron/sqlx v1.3.5
	github.com/labstack/echo-jwt/v4 v4.0.0
	github.com/labstack/echo/v4 v4.10.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/qustavo/sqlhooks/v2 v2.1.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.14.0
	golang.org/x/crypto v0.4.0
)

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab
