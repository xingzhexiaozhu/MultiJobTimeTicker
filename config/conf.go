package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	DB         map[string]DataBaseConf `toml："database"`
	Table      map[string]TableConf    `toml:"sql_table"`
	GlobalConf GlobalConf              `toml:"global_conf"`
	JobContent map[string]JobContent   `toml:"job_content"`
}

var conf *Config

func init() {
	conf = new(Config)
	_, err := toml.Decode("conf.toml", &conf)
	if err != nil {
		panic(err)
	}

	//body, _ := json.Marshal(Conf)
	//fmt.Println(body)
}

// 获取配置
func GetConfig() *Config {
	return conf
}