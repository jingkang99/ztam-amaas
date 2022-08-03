package global

var CFile	string
var Debug	bool
var CFObj	Config

type SERVER struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Seed string `mapstructure:"seed"`
}

type DBCONF struct {
	Type string `mapstructure:"type "`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Schm string `mapstructure:"schm"`
}

type Config struct {
	Server SERVER `mapstructure:"server"`
	DBConf DBCONF `mapstructure:"dbconf"`
}