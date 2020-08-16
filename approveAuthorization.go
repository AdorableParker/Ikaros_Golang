package main

import (
	"strconv"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

// DocApproveAuthorization 受邀入群授权功能文档
var DocApproveAuthorization = &HelpDoc{
	Name:        "受邀入群授权",
	KeyWord:     []string{"授权批准"},
	Example:     "授权批准 123456789",
	Description: "授权批准<空格><批准群号>\n用于批准同意受邀加群申请\n需要机器人系统管理员权限"}

func approveAuthorization(msg []string, msgID int32, group, qq int64, try uint8) {
	if qq != AdminConfig.AdminAccount {
		sendMsg(group, qq, "非系统管理员,权限不足,授权失败")
		return
	}

	if len(msg) != 0 {
		authorizedGroup, err := strconv.ParseInt(msg[0], 10, 0)
		if err != nil {
			cqp.SendGroupMsg(group, "非法输入,授权失败")
			return
		}
		for i, j := range AuthorizedGroupList {
			if j == 0 {
				AuthorizedGroupList[i] = authorizedGroup
				sendMsg(group, qq, "完成")
				return
			}
		}
		sendMsg(group, qq, "授权数已达上限,授权失败")
	} else {
		cqp.SendGroupMsg(group, "请提供需授权群号")
		try++
		stagedSessionPool[msgID] = newStagedSession(group, qq, approveAuthorization, msg, try) // 将一个 待跟进会话 加入 会话池
		sendMsg(group, qq, "完成")
		return
	}
}
