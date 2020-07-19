package main

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

var checkCode string

func dbCM(msg []string, msgID int32, group, qq int64, try uint8) {
	if len(msg) == 0 {
		cqp.SendGroupMsg(group, "提供32位效验码")
		checkCode = getCheckCode()
		cqp.AddLog(10, "checkCode", checkCode)
		try++
		stagedSessionPool[msgID] = newStagedSession(group, qq, dbCM, msg, try) // 将一个 待跟进会话 加入 会话池
		return
	}

	if checkCode != msg[0] {
		cqp.SendGroupMsg(group, "校验失败")
		return
	}
	DBConn = !DBConn
	sendMsg(group, qq, "操作完成")
}

func getCheckCode() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))
}
