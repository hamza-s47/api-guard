package config

type Config struct {
	Server ServerConfig `yaml:"server"`
	Routes []Route      `yaml:"routes"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Route struct {
	Path      string      `yaml:"path"`
	Backend   string      `yaml:"backend"`
	RateLimit RateLimiter `yaml:"rate_limit"`
}

type RateLimiter struct {
	Capacity     int     `yaml:"capacity"`
	RefillPerSec float64 `yaml:"refill_per_sec"`
}
