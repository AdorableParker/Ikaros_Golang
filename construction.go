package main

import (
	"regexp"
	"strings"
)

func construction(msg []string, msgID int32, group, qq int64, try uint8) {

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
	input := strings.ToUpper(msg[0])

	re := regexp.MustCompile(`\d:\d\d`)
	constructionTime := re.FindAllString(strings.Replace(input, "：", ":", -1), 1) // 正则匹配查找索引
	var results []string
	if constructionTime == nil { // 没有找到
		results = nameToTime(msg[0]) // 由名字查找
	} else {
		results = timeToName(constructionTime[0]) // 由时间查找
	}
	for _, result := range results {
		sendMsg(group, qq, result)
	}
}

func nameToTime(name string) []string {
	return []string{name}
}

func timeToName(time string) []string {
	return []string{time}
}
