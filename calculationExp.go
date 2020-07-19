package main

import (
	"fmt"
	"strconv"
)

type lv2Exp struct {
	Lv          int `json:"lv,omitempty"`
	Coefficient int `json:"coefficient,omitempty"`
}

var expList []lv2Exp = []lv2Exp{
	lv2Exp{40, 1}, lv2Exp{60, 2}, lv2Exp{70, 3},
	lv2Exp{80, 4}, lv2Exp{90, 5}, lv2Exp{92, 10},
	lv2Exp{94, 20}, lv2Exp{95, 40}, lv2Exp{97, 50},
	lv2Exp{98, 200}, lv2Exp{99, 720}, lv2Exp{100, -620},
	lv2Exp{104, 20}, lv2Exp{105, 70}, lv2Exp{110, 120},
	lv2Exp{115, 180}, lv2Exp{119, 210}, lv2Exp{120, 0}}

func calculateParts(lowLv, highLv, existing int, flag bool) (totalExp int) {
	for ; lowLv < highLv; lowLv++ {
		needExp := 0
		lastLv := 0
		for _, item := range expList {
			rankDifference := lowLv - item.Lv
			if rankDifference <= 0 {
				needExp += item.Coefficient * (rankDifference + item.Lv - lastLv)
				break
			}
			needExp += item.Coefficient * (item.Lv - lastLv)
			lastLv = item.Lv
		}
		if flag {
			if 90 <= lowLv && lowLv < 100 {
				totalExp += needExp * 13
			} else {
				totalExp += needExp * 12
			}
		} else {
			totalExp += needExp * 10
		}
	}
	return totalExp*10 - existing
}

func calculationExp(msg []string, msgID int32, group, qq int64, try uint8) {
	var lowLv, highLv, existingExp int = 0, 0, 0
	var shipType bool = false
	var err error
	switch len(msg) {
	case 4:
		existingExp, err = strconv.Atoi(msg[3])
		if err != nil {
			sendMsg(group, qq, "已有经验参数\n只接受正整数输入哦\n(。・ω・。)")
			return
		}
		fallthrough

	case 3:
		shipType, err = strconv.ParseBool(msg[2])
		if err != nil {
			sendMsg(group, qq, "是否为决战方案参数\n只接受 1/0、t/f、T/F 这几个输入哦\n(。・ω・。)")
			return
		}
		fallthrough

	case 2:
		lowLv, err = strconv.Atoi(msg[0])
		highLv, err = strconv.Atoi(msg[1])
		if err != nil {
			sendMsg(group, qq, "等级参数\n只接受正整数输入哦\n(。・ω・。)")
			return
		}
		if lowLv >= highLv {
			sendMsg(group, qq, "你有问题,小老弟(¬_¬\")")
		}

	case 1:
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			sendMsg(group, qq, "再输入目标等级\nq(≧▽≦q)")
			stagedSessionPool[msgID] = newStagedSession(group, qq, calculationExp, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "什么都不告诉我可没办法计算呢,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return

	case 0: // 如果参数为0个
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入当前等级\nq(≧▽≦q)")
			} else {
				sendMsg(group, qq, "当前等级不能留空哦,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, calculationExp, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "什么都不告诉我可没办法计算呢,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}

	balance := calculateParts(lowLv, highLv, existingExp, shipType)
	if balance <= 0 {
		sendMsg(group, qq, fmt.Sprintf("当前等级:%d,目标等级:%d\n是否为决战方案:%t\n已有经验:%d\n最终计算结果: 达成目标等级后将溢出 %d EXP", lowLv, highLv, shipType, existingExp, -balance)) // 发送提示消息
	}
	sendMsg(group, qq, fmt.Sprintf("当前等级:%d,目标等级:%d\n是否为决战方案:%t\n已有经验:%d\n最终计算结果: 还需 %d EXP 可以达成目标等级", lowLv, highLv, shipType, existingExp, balance)) // 发送提示消息
}
