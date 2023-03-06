package config

var gconfig = &defaultConfig

type Config struct {
	// chain server
	ChainServerAddr string
	//
	GrpcAddr string
}

var defaultConfig = Config{
	ChainServerAddr: ":3801",
	GrpcAddr:        "0.0.0.0:3802",
}

func GetConfig() *Config {
	return gconfig
}

type NodeConfig struct {
	Generate     bool
	GivenPrivate string
	GrpcPort     int
	NodeDir      string
	ChainServer  string
}
