module github.com/mattermost/mattermost-plugin-autotranslate/server

go 1.12

require (
	github.com/aws/aws-sdk-go v1.15.67
	github.com/fsnotify/fsnotify v1.4.7
	github.com/golang/protobuf v1.2.0
	github.com/google/uuid v1.0.0
	github.com/gorilla/websocket v0.0.0-20181030144553-483fb8d7c32f
	github.com/hashicorp/go-hclog v0.0.0-20180402200405-69ff559dc25f
	github.com/hashicorp/go-plugin v0.0.0-20180331002553-e8d22c780116
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d
	github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8
	github.com/magiconair/properties v1.8.0
	github.com/mattermost/mattermost-server v5.4.0+incompatible
	github.com/mattermost/viper v1.0.4 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/nicksnyder/go-i18n v1.10.0
	github.com/oklog/run v1.0.0
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/pelletier/go-toml v1.2.0
	github.com/pkg/errors v0.8.0
	github.com/spf13/afero v1.1.2
	github.com/spf13/cast v1.3.0
	github.com/spf13/jwalterweatherman v1.0.0
	github.com/spf13/pflag v1.0.3
	go.uber.org/atomic v1.3.2
	go.uber.org/multierr v1.1.0
	go.uber.org/zap v1.9.1
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9
	golang.org/x/net v0.0.0-20181102091132-c10e9556a7bc
	golang.org/x/sys v0.0.0-20181205085412-a5c9d58dba9a
	golang.org/x/text v0.3.0
	google.golang.org/genproto v0.0.0-20181101192439-c830210a61df
	google.golang.org/grpc v1.16.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
	gopkg.in/yaml.v2 v2.2.2
)
