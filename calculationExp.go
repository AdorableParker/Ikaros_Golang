package main

import (
	"fmt"
	"strconv"
)

// n*100 + (n-40)*100 + (n-60)*100 + (n-70)*100 + (n-80)*100 + (n-90)*500 + (n-92)*1000 + (n-94)*2000 + (n-95)*1000 + (n-97)*15000 + (n-98)*52000
func calculateParts(lowLv, highLv int, flag bool) int {
	highLv--
	if highLv < lowLv {
		return 0
	} else if highLv >= 100 {
		return calculatePartsPro(lowLv, highLv, flag) + calculateParts(lowLv, 100, flag)
	}
	totalExp := 0
	switch {
	case highLv >= 98:
		totalExp += (highLv - 98) * 5200
		fallthrough
	case highLv >= 97:
		totalExp += (highLv - 97) * 1500
		fallthrough
	case highLv >= 95:
		totalExp += (highLv - 95) * 100
		fallthrough
	case highLv >= 94:
		totalExp += (highLv - 94) * 200
		fallthrough
	case highLv >= 92:
		totalExp += (highLv - 92) * 100
		fallthrough
	case highLv >= 90:
		totalExp += (highLv - 90) * 50
		fallthrough
	case highLv >= 80:
		totalExp += (highLv - 80) * 10
		fallthrough
	case highLv >= 70:
		totalExp += (highLv - 70) * 10
		fallthrough
	case highLv >= 60:
		totalExp += (highLv - 60) * 10
		fallthrough
	case highLv >= 40:
		totalExp += (highLv - 40) * 10
		fallthrough
	default:
		totalExp += highLv * 10
	}
	if flag {
		if highLv < 90 {
			totalExp *= 12
			// totalExp = totalExp * 12 / 10
		} else {
			totalExp *= 13
		}
	}

	return totalExp + calculateParts(lowLv, highLv, flag)

}

// 70000 + (n-100)*2000 + (n-104)*5000 + (n-105)*5000 + (n-110)*6000 + (n-115)*3000
func calculatePartsPro(lowLv, highLv int, flag bool) int {
	if highLv < lowLv || highLv < 100 {
		return 0
	}
	var totalExp int
	switch {
	case highLv > 120:
		totalExp = 0
	case highLv > 119:
		totalExp = 3000000
	case highLv >= 115:
		totalExp += (highLv - 115) * 300
		fallthrough
	case highLv >= 110:
		totalExp += (highLv - 110) * 600
		fallthrough
	case highLv >= 105:
		totalExp += (highLv - 105) * 500
		fallthrough
	case highLv >= 104:
		totalExp += (highLv - 104) * 500
		fallthrough
	default:
		totalExp += (highLv-100)*200 + 7000
	}
	fmt.Println(totalExp, highLv)
	if flag && highLv <= 119 {
		totalExp *= 12
	}
	highLv--
	return totalExp + calculatePartsPro(lowLv, highLv, flag)

}

// DocCalculationExp 碧蓝航线舰船经验计算器功能文档
var DocCalculationExp = &HelpDoc{
	Name:        "碧蓝航线舰船经验计算器",
	KeyWord:     []string{"舰船经验", "经验计算"},
	Example:     "舰船经验 10 20 F 0\n经验计算 10 20 F\n经验计算 10 20",
	Description: "命令输入格式:\n舰船经验<空格><当前等级><空格><目标等级><空格>[是否为决战方案]<空格>[已有经验]\n根据输入的参数，返回达成目标等级需要的经验或是溢出的经验\n"}

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
			return
		}

	case 1:
		try++         // 已尝试次数+1
		if try <= 5 { // 如果已尝试次数不超过3次
			sendMsg(group, qq, "再输入目标等级\nq(≧▽≦q)")
			stagedSessionPool[msgID] = newStagedSession(group, qq, calculationExp, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "什么都不告诉我可没办法计算呢,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return

	case 0: // 如果参数为0个
		try++         // 已尝试次数+1
		if try <= 5 { // 如果已尝试次数不超过3次
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

	balance := calculateParts(lowLv, highLv, shipType)
	balance -= existingExp
	if balance < 0 {
		sendMsg(group, qq, fmt.Sprintf("当前等级:%d,目标等级:%d\n是否为决战方案:%t\n已有经验:%d\n最终计算结果: 达成目标等级后将溢出 %d EXP", lowLv, highLv, shipType, existingExp, -balance)) // 发送提示消息
		return
	}
	sendMsg(group, qq, fmt.Sprintf("当前等级:%d,目标等级:%d\n是否为决战方案:%t\n已有经验:%d\n最终计算结果: 还需 %d EXP 可以达成目标等级", lowLv, highLv, shipType, existingExp, balance)) // 发送提示消息
}
