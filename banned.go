package main

import (
	"strconv"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

// DocBanned 快捷禁言功能文档
var DocBanned = &HelpDoc{
	Name:        "快捷禁言",
	KeyWord:     []string{"求口", "自助禁言"},
	Example:     "自助禁言 30",
	Description: "求口<空格><时间分钟>\n用于满足群员自己禁言自己的需求,需要提供管理员权限"}

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

// DocAllBanned 全群禁言功能文档
var DocAllBanned = &HelpDoc{
	Name:        "全群禁言",
	KeyWord:     []string{"全员禁言", "全员自闭"},
	Description: "开启当前群的全员禁言,需要提供管理员权限"}

func allBanned(group, qq int64) {
	groupInfo := cqp.GetGroupMemberInfo(group, qq, false)
	if groupInfo.Auth >= 2 {
		cqp.SetGroupWholeBan(group, true)
	} else {
		cqp.SendGroupMsg(group, "我才不听你的呢\n(￣﹃￣)")
	}
}

// DocAllNotBanned 解除全群禁言功能文档
var DocAllNotBanned = &HelpDoc{
	Name:        "解除全群禁言",
	KeyWord:     []string{"解禁"},
	Description: "关闭当前群的全员禁言,需要提供管理员权限"}

func allNotBanned(group int64) {
	cqp.SetGroupWholeBan(group, false)
}
