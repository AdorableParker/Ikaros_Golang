package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"

	"gopkg.in/ini.v1"
)

func activity(msg []string, msgID int32, group, qq int64, try uint8) {
	if len(msg) == 0 { // 如果没有获取到参数
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入已刷点数\nq(≧▽≦q)")
			} else {
				sendMsg(group, qq, "不能为空哦,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, activity, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}
	cfg, err := ini.Load(Appdir + "config.ini")
	if err != nil {
		cqp.AddLog(10, "INI文件读取异常", fmt.Sprintln(err))
		sendMsg(group, qq, "读取活动信息失败了\n/(ㄒoㄒ)/~~")
		return
	}

	on, _ := cfg.Section("ongoing_activities").Key("form").Bool()
	if !on {
		sendMsg(group, qq, "暂时没有开启的活动哦\n(＃°Д°)")
		return
	}

	name := cfg.Section("ongoing_activities").Key("name").String()
	report := fmt.Sprintf("活动名：%s", name)
	var point, aims int

	parameter := strings.Split(msg[0], "#")

	point, err = strconv.Atoi(parameter[0])
	if err != nil || point < 0 {
		sendMsg(group, qq, "计算被迫结束,原因:\n (°д°) 奇怪的参数增加了")
		return
	}

	if len(parameter) > 1 {
		aims, _ = strconv.Atoi(parameter[1])
		if err != nil || aims < 0 {
			sendMsg(group, qq, "计算被迫结束,原因:\n (°д°) 奇怪的参数增加了.jpg")
			return
		}

	} else {
		aims, _ = cfg.Section("shop").Key("all").Int()
	}

	if point >= aims { // 已完成
		sendMsg(group, qq, "目标已完成\nお疲れ様です\n(p≧w≦q)")
		return
	}

	mapIDList := cfg.Section("mapid").KeyStrings()
	schedule := strconv.FormatFloat(float64(point)/float64(aims)*100, 'f', 2, 32)

	Remaining := aims - point
	for _, mapID := range mapIDList {
		divisor, _ := cfg.Section("mapid").Key(mapID).Int()
		quotient := Remaining / divisor
		if Remaining%divisor != 0 {
			quotient++
		}
		report += fmt.Sprintf("\n若只出击%s还需%d次", mapID, quotient)
	}

	stopTime := cfg.Section("time").Key("stoptime").String()
	t, _ := time.Parse("2006-1-2 15:4 MST", stopTime)
	if time.Now().After(t) { // 如果已经结束
		report += fmt.Sprintf("\n当前已获得 %d 积分\n已完成进度\n%s%%\n活动已经结束了哦", point, schedule)
		sendMsg(group, qq, report)

		// 更改保存
		cfg.Section("ongoing_activities").Key("form").SetValue("false")
		cfg.SaveTo(Appdir + "config.ini")
		return
	}
	remainingTime := t.Sub(time.Now().Truncate(time.Second)).String()
	report += fmt.Sprintf("\n当前已获得 %d 积分\n已完成进度\n%s%%\n距离活动结束还有 %s", point, schedule, remainingTime)
	sendMsg(group, qq, report)
}
