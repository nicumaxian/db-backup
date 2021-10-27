package configuration

type DbConfiguration struct {
	Host string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Configuration struct {
	Databases map[string]DbConfiguration `mapstructure:"databases"`
}
