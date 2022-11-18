package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap/zapcore"
)

type Backend struct {
	S3       *BackendS3 `env:",prefix=S3_"`
	Type     string     `env:"TYPE,required"`
	RootPath string     `env:"ROOT_PATH,default=/tmp"`
}

func (b *Backend) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("type", b.Type)
	return nil
}

type BackendS3 struct {
	AccessKey    string `env:"ACCESS_KEY"`
	Bucket       string `env:"BUCKET"`
	Endpoint     string `env:"ENDPOINT"`
	SecretKey    string `env:"SECRET_KEY"`
	UsePathStyle bool   `env:"USE_PATH_STYLE, default=false"`
}

type Config struct {
	Backend  *Backend `env:",prefix=BACKEND_"`
	Log      *Log     `env:",prefix=LOG_"`
	Lock     *Lock    `env:",prefix=LOCK_"`
	Name     string   `env:"NAME,default=kerranamodb"`
	HTTPPort int      `env:"HTTP_PORT,default=5000"`
	GRPCPort int      `env:"GRPC_PORT,default=10020"`
	Trace    *Trace   `env:",prefix=TRACE_"`
}

type Lock struct {
	ServiceDiscoveryEndpoint  string `env:"SERVICE_DISCOVERY_ENDPOINT"`
	ServiceDiscoveryNodeCount int    `env:"SERVICE_DISCOVERY_NODE_COUNT"`
	ServiceDiscoveryTimeout   int    `env:"SERVICE_DISCOVERY_TIMEOUT"`
	ServiceDiscoveryPort   int    `env:"SERVICE_DISCOVERY_PORT"`
	HostIP                    string `env:"HOST_IP"`
	Nodes                     string `env:"NODES"`
}

func (l *Lock) GetNodes() []string {
	return strings.Split(l.Nodes, ",")
}

type Trace struct {
	Enable bool   `env:"ENABLE,default=false"`
	Name   string `env:"NAME,default=kerranamodb"`
	Type   string `env:"TYPE,default=console"`

	Jaeger *TraceJaeger `env:",prefix=JAEGER_"`
}

type TraceJaeger struct {
	Endpoint string `env:"ENDPOINT"`
}

func (cfg *Config) HTTPAddress() string {
	return fmt.Sprintf(":%d", cfg.HTTPPort)
}

type Log struct {
	Format string `env:"FORMAT,default=json"`
	Level  string `env:"LEVEL,default=info"`
}

func Load(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
