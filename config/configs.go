package config

type Config struct {
	RedisHost string
	RedisPort string
	Port      string
}

var AppConfig = Config{
	RedisHost: "localhost",
	RedisPort: "6379",
	Port:      ":8080",
}
