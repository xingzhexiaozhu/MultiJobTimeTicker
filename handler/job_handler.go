package handler

import (
	"MultiJobTimeTicker/config"
	"MultiJobTimeTicker/dao"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type UserData struct {
	dao *dao.SQLDao
}

var MapUserData map[int]*UserData

type UserInfo struct {
	ID int64
}

func Init(conf *config.Config) error {
	if conf == nil {
		return fmt.Errorf("err: %s", "conf is nil")
	}

	// 里程活动数据库、数据表
	dbConf := conf.DB["push"]
	table1 := conf.Table["1"].Table
	data1, err := NewUserData(&dbConf, table1)
	if err != nil {
		return fmt.Errorf("init user data err: %s, table: %s, conf: %+v", err.Error(), table1, dbConf)
	}

	// 公交活动数据库、数据表
	table2 := conf.Table["2"].Table
	data2, err := NewUserData(&dbConf, table2)
	if err != nil {
		return fmt.Errorf("init user data err: %s, table: %s, conf: %+v", err.Error(), table2, dbConf)
	}

	MapUserData = map[int]*UserData{
		1: data1,
		2: data2,
	}
	return nil
}

// 获取用户数据库dao实例
func NewUserData(conf *config.DataBaseConf, table string) (*UserData, error) {
	dao, err := dao.NewSQLDao(conf, table)
	if err != nil {
		err = fmt.Errorf("init UserData err, %s", err.Error())
		return nil, err
	}

	return &UserData{
		dao: dao,
	}, nil
}

// 获取用户数量
func (userData *UserData) GetUserCount() (int64, error) {
	// 事务
	trans, err := userData.dao.GetDB().Begin()
	defer trans.Commit()
	if err != nil {
		err = fmt.Errorf("begin transaction err, %s", err.Error())
		return 0, err
	}

	count := int64(0)
	row := trans.QueryRow(`SELECT count(*) FROM` + userData.dao.GetTable())
	row.Scan(&count)

	return count, nil
}

func timeHepler() int64 {
	t, _ := strconv.ParseInt(time.Now().Format("20060102"), 10, 64)
	return t
}

// 给符合条件的用户发push消息
func (userData *UserData) SelectUserAndDoJob(job *config.Job) {
	// 事务
	trans, err := userData.dao.GetDB().Begin()
	if err != nil {
		log.Printf("begin transaction err, %s", err.Error())
		return
	}
	defer trans.Commit()

	count := int64(0)
	sql := `SELECT count(*) FROM ` + userData.dao.GetTable()
	traceId := fmt.Sprintf("%v%v", time.Now().UnixNano()/1e6, rand.Intn(1000))
	if job.Condition != "" {
		condition := fmt.Sprintf("%s and dn=%d", job.Condition, timeHepler())
		sql = sql + ` where ` + condition
	}
	row := trans.QueryRow(sql)
	row.Scan(&count)
	log.Printf("transQueryRow: traceId=%v sql=%v count=%d", traceId, sql, count)
	if count == 0 {
		log.Printf("sql=%v, count=0", sql)
		return
	}

	var id int64
	// 分批查找
	for i := int64(0); i <= int64(count/200 + 1); i++ {
		//for {
		var userList = make([]UserInfo, 0, 0)
		searchSql := `SELECT * success FROM ` + userData.dao.GetTable() + ` WHERE id>?`
		if job.Condition != "" {
			condition := fmt.Sprintf("%s and dn=%d", job.Condition, timeHepler())
			searchSql = searchSql + ` and ` + condition
		}
		searchSql = searchSql + ` ORDER BY id ASC LIMIT 200`
		rows, err := trans.Query(searchSql, id)
		log.Printf("pageHelper: traceId=%v searchSql=%s", traceId, searchSql)
		if err != nil {
			log.Printf("execute sql err, %s", err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			var info UserInfo
			err = rows.Scan(&info.ID)
			if err != nil {
				log.Printf("rows err, row=%v, err=%v", rows, err.Error())
				return
			}
			userList = append(userList, info)
			id = info.ID
		}

		err = rows.Err()
		if err != nil {
			log.Printf("rows err, row=%v, err=%v", rows, err.Error())
			return
		}

		if len(userList) > 0 {
			userData.doJob(userList, job, traceId)
			atomic.AddInt64(&job.Success, int64(len(userList)))
			log.Printf("pushSuccess: traceId=%v, job=%v, count=%v, successCnt=%v", traceId, job, count, job.Success)
		}
	}
	return
}

func (userData *UserData) doJob(userList []UserInfo, job *config.Job, traceId string) {
	var idStrs string
	var ids []string
	for _, v := range userList {
		ids = append(ids, strconv.FormatInt(v.ID, 10))
	}

	if len(ids) > 0 {
		idStrs = strings.Join(ids, ",")
	}

	// 基于id做一些具体的任务
	// ...

	// 更新数据库内容
	userData.updateJobSuccess(idStrs, job.Type, traceId)
}

func (userData *UserData) updateJobSuccess(idstr string, jobType int, traceId string) (err error) {
	//事务
	trans, err := userData.dao.GetDB().Begin()
	if err != nil {
		err = fmt.Errorf("get commute user by id begin transaction failed, %s traceId=%v", err.Error(), traceId)
		return err
	}
	defer trans.Commit()

	condition := fmt.Sprintf("id in (%s) and type=%d", idstr, jobType)
	updateSql := `update ` + userData.dao.GetTable() + ` set success = 1 where ` + condition
	updateRows, err := trans.Exec(updateSql)
	if err != nil {
		log.Printf("update user by id err, %v updateRows=%v traceId=%v", err, updateRows, traceId)
		return err
	}
	cnt, err := updateRows.RowsAffected()
	if err != nil {
		log.Printf("update user by id err, %v updateRows=%v traceId=%v", err, updateRows, traceId)
		return err
	}
	log.Printf("updateDB: traceId=%s update sql=%s success rows=%d", traceId, updateSql, cnt)
	return nil
}