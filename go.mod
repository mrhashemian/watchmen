module watchmen

go 1.16

require (
	github.com/alicebob/miniredis/v2 v2.30.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.7.0
	github.com/golang-jwt/jwt/v4 v4.4.3 // indirect
	github.com/jmoiron/sqlx v1.3.5
	github.com/labstack/echo/v4 v4.10.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/qustavo/sqlhooks/v2 v2.1.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.14.0
	golang.org/x/crypto v0.2.0
)

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab
