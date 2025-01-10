package config

type JWT struct {
	Issuer              string `mapstructure:"issuer"`
	UserAccessTokenKey  string `mapstructure:"user_access_token_key"`
	UserRefreshTokenKey string `mapstructure:"user_refresh_token_key"`
}

type Server struct {
	Port      int    `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`
	SentryDNS string `mapstructure:"sentry_dns"`
}

type GRPC struct {
	Port int `mapstructure:"port"`
}
