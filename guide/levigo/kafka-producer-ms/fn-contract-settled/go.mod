module fn-contract-settled

go 1.24.0

toolchain go1.24.11

require (
	// Core
	github.com/cloudevents/sdk-go/v2 v2.16.2
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/google/uuid v1.6.0

	// 🔒 Deps explícitas do Viper (ESSENCIAL)
	github.com/spf13/afero v1.15.0 // indirect

	// Config
	github.com/spf13/viper v1.21.0
	github.com/subosito/gotenv v1.6.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
)

require (
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.37.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace (
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20220503193339-ba3ae3f07e29
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto v0.0.0-20220503193339-ba3ae3f07e29
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto v0.0.0-20220503193339-ba3ae3f07e29
)
