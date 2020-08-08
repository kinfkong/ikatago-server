module github.com/kinfkong/ikatago-server

go 1.14

require (
	github.com/aliyun/aliyun-oss-go-sdk v2.1.4+incompatible
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatedier/beego v0.0.0-20171024143340-6c6a4f5bd5eb
	github.com/fatedier/frp v0.0.0-00010101000000-000000000000
	github.com/fatedier/golib v0.0.0-20181107124048-ff8cd814b049
	github.com/gliderlabs/ssh v0.3.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/mergermarket/go-pkcs7 v0.0.0-20170926155232-153b18ea13c9
	github.com/rakyll/statik v0.1.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.7.1
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	moul.io/http2curl v1.0.0
)

replace github.com/fatedier/frp => github.com/kinfkong/frp v1.33.8
