package config

// DataBaseConf 定义了sql数据库相关的配置
type DataBaseConf struct {
	DB          string `toml:"db"`
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	User        string `toml:"user"`
	Pwd         string `toml:"pwd"`
	CharSet     string `toml:"charset"`
	ConnTimeOut int64  `toml:"conn_time_out"`
	TimeOut     int64  `toml:"time_out"`
	MaxConnNums int    `toml:"max_conn_nums"`
	MaxIdleNums int    `toml:"max_idle_nums"`
}

// TableConf 定义了sql表名的配置
type TableConf struct {
	Table string `toml:"table"`
}

// GlobalConfig 定义全局push信息
type GlobalConf struct {
	JobNum         int `toml:"job_num"`
	ConnTimeoutMs  int `toml:"conn_timeout_ms"`
	ServeTimeoutMs int `toml:"serve_timeout_ms"`
}

// JobContent 配置化的消息内容
type JobContent struct {
	Type      int    `toml:"type"`      // 标识不同的活动，根据活动id区分所数据库
	Title     string `toml:"title"`     // 主题
	Content   string `toml:"content"`   // 具体任务
	Url       string `toml:"url"`       // url信息
	Freq      string `toml:"freq"`      // 活动频率
	Condition string `toml:"condition"` // 活动条件
	Force     int    `toml:"force"`
}

// Job 落入sql中的job信息
type Job struct {
	Type      int    `json:"type"`      // 标识不同的活动，根据活动id区分所数据库
	Title     string `json:"title"`     // 主题
	Content   string `json:"content"`   // 具体任务
	Url       string `json:"url"`       // url信息
	Freq      string `json:"freq"`      // 活动频率
	Condition string `json:"condition"` // 活动条件
	Success   int64  `json:"success"`   // 执行状态
}
