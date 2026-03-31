package config

import "time"

// Config is a main structure of global config
type Config struct {
	App      App      `yaml:"app" validate:"required"`
	Server   Server   `yaml:"server" validate:"required"`
	Postgres Postgres `yaml:"postgres" validate:"required"`
	Redis    Redis    `yaml:"redis" validate:"required"`
	Session  Session  `yaml:"session" validate:"required"`
}

type App struct {
	Env     string `yaml:"env" validate:"required,oneof=dev prod local"`
	Name    string `yaml:"name" validate:"required"`
	Version string `yaml:"version" validate:"required,semver"`
}

type TLS struct {
	Enable         bool   `yaml:"enable"`
	ServerCertPath string `yaml:"server-cert-path"`
	ServerKeyPath  string `yaml:"server-key-path"`
}

type HTTP struct {
	Addr string `yaml:"addr" validate:"required,hostname_port"`
	TLS  TLS    `yaml:"tls" validate:"required"`
}

type Server struct {
	HTTP         HTTP          `yaml:"http" validate:"required"`
	ReadTimeout  time.Duration `yaml:"read-timeout" validate:"required,min=100ms"`
	WriteTimeout time.Duration `yaml:"write-timeout" validate:"required,min=100ms"`
	IdleTimeout  time.Duration `yaml:"idle-timeout" validate:"required,min=100ms"`
}

type PostgresAuth struct {
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	DBName   string `yaml:"dbname" validate:"required"`
}

type PostgresConn struct {
	MaxIdles    int           `yaml:"max-idles"`
	MaxOpens    int           `yaml:"max-opens"`
	MaxIdleTime time.Duration `yaml:"max-idle-time"`
	MaxLifetime time.Duration `yaml:"max-lifetime"`
}

type Postgres struct {
	Host    string       `yaml:"host" validate:"required,hostname"`
	Port    int          `yaml:"port" validate:"required,gte=1,lte=65535"`
	SSLMode string       `yaml:"sslmode" validate:"required"`
	Auth    PostgresAuth `yaml:"auth" validate:"required"`
	Conn    PostgresConn `yaml:"conn" validate:"required"`
}

type Session struct {
	TokenTTL  time.Duration `yaml:"token-ttl" validate:"required,min=1h"`
	SecretKey string        `yaml:"secret-key" validate:"required"`
}

type RedisAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Redis struct {
	Addr    string        `json:"addr" validate:"required,hostname_port"`
	Auth    RedisAuth     `json:"auth"`
	UserTTL time.Duration `json:"user-ttl" validate:"required,min=100ms"`
}
