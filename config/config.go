package config

var gconfig = new(Config)

type Config struct {
	// chain server
	ChainServerAddr string

	//
	GrpcAddr  string
}



func GetConfig() *Config {
	return gconfig
}