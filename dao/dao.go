package dao

import (
	"MultiJobTimeTicker/config"
	"database/sql"
	"fmt"
)

// 数据库连接管理
type DBConnManager struct {
	mapDBConns map[string](*sql.DB) // key:DBName value:DBConn
}

// 数据库连接管理实例
var dbConnManager *DBConnManager

// 初始化连接管理实例
func init() {
	dbConnManager = newDBConnManager()
}

// 创建数据库连接管理实例
func newDBConnManager() *DBConnManager {
	return &DBConnManager{mapDBConns: make(map[string](*sql.DB))}
}

// 从map维护的连接中按名字获取，如果取不到则新建立连接
func (manager *DBConnManager) getDBConn(conf *config.DataBaseConf) (*sql.DB, error) {
	dbConn, exist := manager.mapDBConns[conf.DB]
	if exist {
		return dbConn, nil
	}

	dbConn, err := initDBConn(conf)
	if err != nil {
		return nil, err
	}

	manager.mapDBConns[conf.DB] = dbConn
	return manager.mapDBConns[conf.DB], nil
}

// 创建连接
func initDBConn(conf *config.DataBaseConf) (*sql.DB, error) {
	if conf == nil {
		return nil, fmt.Errorf("db conf is empty")
	}

	// 数据库编码设置
	if conf.CharSet == "" {
		conf.CharSet = "utf8mb4"
	}
	// 连接超时设置
	if conf.ConnTimeOut <= 0 {
		conf.ConnTimeOut = 5000
	}
	// 超时
	if conf.TimeOut <= 0 {
		conf.TimeOut = 20000
	}

	// 创建连接
	connUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=%s&timeout=%dms&readTimeout=%dms&writeTimeout=%dms",
		conf.User, conf.Pwd, conf.Host, conf.Port, conf.DB, conf.CharSet, conf.ConnTimeOut, conf.TimeOut, conf.TimeOut)
	dbConn, err := sql.Open("mysql", connUrl)
	dbConn.SetMaxOpenConns(conf.MaxConnNums)
	dbConn.SetMaxIdleConns(conf.MaxIdleNums)
	if err != nil {
		return nil, fmt.Errorf("init sql db failed, %s", err.Error())
	}

	return dbConn, err
}

// SQLDao 为sql提供dao服务
type SQLDao struct {
	dbConn *sql.DB // 连接
	table  string  // 表名
}

// 生成Dao实例
func NewSQLDao(conf *config.DataBaseConf, table string)(*SQLDao, error) {
	dbConn, err := dbConnManager.getDBConn(conf)
	if err != nil {
		err = fmt.Errorf("get sql dao err, %s", err.Error())
		return nil, err
	}

	return &SQLDao{
		dbConn: dbConn,
		table:  table,
	}, nil
}

// 获取Dao的连接实例
func (dao *SQLDao) GetDB() *sql.DB {
	return dao.dbConn
}

// 获取Dao的表实例
func (dao *SQLDao) GetTable() string {
	return dao.table
}