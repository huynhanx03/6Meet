package settings

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
	Logger  LoggerConfig  `mapstructure:"logger"`
	Redis   RedisConfig   `mapstructure:"redis"`
}

// ServerConfig is the configuration for the server
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	Host string `mapstructure:"host"`
}

// MongoDBConfig is the configuration for MongoDB
type MongoDBConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Timeout         int    `mapstructure:"timeout"`
	MaxPoolSize     uint64 `mapstructure:"max_pool_size"`
	MinPoolSize     uint64 `mapstructure:"min_pool_size"`
	MaxConnIdleTime uint64 `mapstructure:"max_conn_idle_time"`
}

// LoggerConfig is the configuration for the logger
type LoggerConfig struct {
	LogLevel    string `mapstructure:"log_level"`
	FileLogName string `mapstructure:"file_log_name"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAge      int    `mapstructure:"max_age"`
	MaxSize     int    `mapstructure:"max_size"`
	Compress    bool   `mapstructure:"compress"`
}

// RedisConfig is the configuration for Redis
type RedisConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Password        string `mapstructure:"password"`
	Database        int    `mapstructure:"database"`
	DialTimeout     int    `mapstructure:"dial_timeout"`
	ReadTimeout     int    `mapstructure:"read_timeout"`
	WriteTimeout    int    `mapstructure:"write_timeout"`
	PoolSize        int    `mapstructure:"pool_size"`
	MinIdleConns    int    `mapstructure:"min_idle_conns"`
	PoolTimeout     int    `mapstructure:"pool_timeout"`
	MaxRetries      int    `mapstructure:"max_retries"`
	MinRetryBackoff int    `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff int    `mapstructure:"max_retry_backoff"`
}
