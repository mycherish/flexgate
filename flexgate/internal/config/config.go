package config

type Config struct {
	Server ServerConfig `yaml:"server"`
	Routes []Route      `yaml:"routes"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Route struct {
	PathPrefix  string          `yaml:"path_prefix"`
	StripPrefix string          `yaml:"strip_prefix"`
	Upstream    string          `yaml:"upstream"`
	RateLimit   RateLimitConfig `yaml:"ratelimit"`
}

type RateLimitConfig struct {
	Rate     int64 `yaml:"rate"`
	Capacity int64 `yaml:"capacity"`
}
