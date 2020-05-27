package main

import (
	"fmt"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jinzhu/gorm"
)

type roster struct {
	Code string `gorm:"column:code"`
	Name string `gorm:"column:name"`
}

func realName(msg []string, msgID int32, group, qq int64, try uint8) {

	if len(msg) == 0 { // 如果没有获取到参数
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入索引信息\nq(≧▽≦q)")
			} else {
				sendMsg(group, qq, "索引不能为空哦,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, construction, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}

	index := msg[0]

	// 读取数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		return
	}
	var roster1, roster2 []roster

	db.Table("Roster").Where("code GLOB ?", fmt.Sprintf("*%s*", index)).Find(&roster1)
	db.Table("Roster").Where("name GLOB ?", fmt.Sprintf("*%s*", index)).Find(&roster2)

	var r = make(map[string]string)

	for _, i := range roster1 {
		r[i.Code] = i.Name
	}
	for _, i := range roster2 {
		r[i.Code] = i.Name
	}

	// 格式化输出
	if len(r) == 0 {
		sendMsg(group, qq, "名字中包含有 %s 的舰船未收录")
	}
	var str string = fmt.Sprintf("名字中包含有 %s 的舰船有:", index)
	for code, name := range r {
		str += fmt.Sprintf("\n和谐名:%s    原名:%s", code, name)
	}
	sendMsg(group, qq, str)
}
