module github.com/krafton-hq/version-helper

go 1.18

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/spf13/cobra v1.4.0
	github.com/thediveo/enumflag v0.10.1
	github.krafton.com/xtrm/fox/client v1.1.10
	github.krafton.com/xtrm/fox/core v0.0.0-coreunpublished
	go.uber.org/zap v1.21.0
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506 // indirect
	google.golang.org/grpc v1.43.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.krafton.com/xtrm/fox/core v0.0.0-coreunpublished => github.krafton.com/xtrm/fox/core v1.1.10

replace github.krafton.com/xtrm/fox/server v0.0.0-serverunpublished => github.krafton.com/xtrm/fox/server v1.1.10
