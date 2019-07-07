package dao

import (
	"database/sql"
	"github.com/LucasGao67/blog001/util"
	"time"
)

var Award award

type award struct{}

type AwardInfo struct {
	Id     int64  `json:"id" db:"id"`
	UserId int64  `json:"user_id" db:"user_id"`
	Remark string `json:"remark"`
	Ct     int64  `json:"ct"`
	Ut     int64  `json:"ut"`
}

const creatTable = `create table blog1 (id bigint primary key auto_increment,user_id bigint default 0,ct bigint default 0,ut bigint default 0,remark varchar(255) default "")`

func (award) FindOne(userId int64) (item *AwardInfo, err error) {
	item = &AwardInfo{}
	sqlStr := " select id,user_id,remark,ct,ut from blog1 where user_id = ?"
	err = util.SqlClient.Get(item, sqlStr, userId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return item, err
}

func (award) InsertOne(userId int64, remark string) (int64, error) {
	now := time.Now()
	id := int64(0)
	sqlStr := "insert into blog1 (user_id,remark,ct,ut) values(?,?,?,?)"
	result, err := util.SqlClient.Exec(sqlStr, userId, remark, now.Unix(), now.Unix())
	if err == nil {
		id, err = result.LastInsertId()
	}
	return id, err
}
