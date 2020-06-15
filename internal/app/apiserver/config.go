package apiserver

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	SocketPort  string `toml:"socket_port"`
	FcmKey      string `toml:"fcm_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":3052",
		LogLevel: "debug",
		//DatabaseURL: "host=localhost dbname=mitso_chat sslmode=disable password=62egegos user=postgres",
		DatabaseURL: "host=localhost port=3053 dbname=mitsodb sslmode=require password=mitsodbpass user=dbuser",
		SocketPort:  ":3051",
		FcmKey:      "AAAAbgOPVb4:APA91bGTZhVFprUUAGEAOVnzEV_H34OOL_GXHGruuz2supgUUR9pryzDGvj_70OS4iTkJgOU15SdDjL4P7jsXI-8znwlJPEoB2Na_DIhSZI-WqXfaY1BOEnT5xfkNrwxTWid2zC9t3AQ",
	}

}
