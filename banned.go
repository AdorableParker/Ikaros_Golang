package main

import (
	"strconv"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

func banned(msg []string, msgID int32, group, qq int64, try uint8) {
	if len(msg) != 0 {
		bantime, err := strconv.Atoi(msg[0])
		if err != nil {
			cqp.SendGroupMsg(group, "请不要输入一些奇奇怪怪的东西\n＞︿＜")
			return
		}

		cqp.SetGroupBan(group, qq, int64(bantime*60))

	} else {
		cqp.SendGroupMsg(group, "请选择套餐时长(单位:分钟)\n(。・ω・。)")
		try++
		stagedSessionPool[msgID] = newStagedSession(group, qq, banned, msg, try) // 将一个 待跟进会话 加入 会话池
		return
	}
}

func allBanned(group, qq int64) {
	groupInfo := cqp.GetGroupMemberInfo(group, qq, false)
	if groupInfo.Auth >= 2 {
		cqp.SetGroupWholeBan(group, true)
	} else {
		cqp.SendGroupMsg(group, "我才不听你的呢\n(￣﹃￣)")
	}
}

func allNotBanned(group int64) {
	cqp.SetGroupWholeBan(group, false)
}
