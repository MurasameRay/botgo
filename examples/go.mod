module github.com/tencent-connect/botgo/examples

go 1.21

require (
	github.com/google/uuid v1.3.0
	github.com/tencent-connect/botgo v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
	golang.org/x/oauth2 v0.23.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/go-resty/resty/v2 v2.6.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/tidwall/gjson v1.9.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
)

replace github.com/tencent-connect/botgo => ../
