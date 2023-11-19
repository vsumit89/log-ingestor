package config

type AppConfig struct {
	Server *ServerConfig `yaml:"server"`
	DB     []DBConfig    `yaml:"db"`
	Search *SearchConfig `yaml:"search_index"`
	MQCfg  *MQConfig     `yaml:"messagequeue"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type MQConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type SearchConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Name    string `yaml:"index_name"`
	API_KEY string `yaml:"api_key"`
}
