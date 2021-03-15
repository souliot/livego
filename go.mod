module livego

go 1.15

replace (
	public => ../public
)

require (
	github.com/astaxie/beego v1.12.2
	github.com/auth0/go-jwt-middleware v0.0.0-20190805220309-36081240882b
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis/v7 v7.2.0
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/kr/pretty v0.1.0
	github.com/orcaman/concurrent-map v0.0.0-20190826125027-8c72a8bb44f6
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	github.com/urfave/negroni v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.4.1
	public v1.0.0
)
