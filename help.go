package main

import (
	"fmt"
	"strings"
)

var helpdocset = make(map[string]*HelpDoc, 25)

func (funcName HelpDoc) readme(group, qq int64) {
	sendMsg(group, qq, fmt.Sprintf("功能名：%s\n命令关键字:%v\n范例:\n%s\n说明:\n%s", funcName.Name, funcName.KeyWord, funcName.Example, funcName.Description))
}

func help(nameList []string, fromGroup, fromQQ int64) {
	if len(helpdocset) == 0 {
		helpdocset[DocActivity.Name] = DocActivity
		helpdocset[DocAllBanned.Name] = DocAllBanned
		helpdocset[DocAllNotBanned.Name] = DocAllNotBanned
		helpdocset[DocApproveAuthorization.Name] = DocApproveAuthorization
		helpdocset[DocBanned.Name] = DocBanned
		helpdocset[DocCalculationExp.Name] = DocCalculationExp
		helpdocset[DocCalculato.Name] = DocCalculato
		helpdocset[DocConsole.Name] = DocConsole
		helpdocset[DocConstruction.Name] = DocConstruction
		helpdocset[DocDynamicByID.Name] = DocDynamicByID
		helpdocset[DocEquipmentRanking.Name] = DocEquipmentRanking
		helpdocset[DocPixivRanking.Name] = DocPixivRanking
		helpdocset[DocRealName.Name] = DocRealName
		helpdocset[DocSaucenao.Name] = DocSaucenao
		helpdocset[DocSendDynamic.Name] = DocSendDynamic
		helpdocset[DocShipMap.Name] = DocShipMap
		helpdocset[DocSolitaire.Name] = DocSolitaire
		helpdocset[DocSrengthRanking.Name] = DocSrengthRanking
		helpdocset[DocTraining.Name] = DocTraining
		helpdocset[DocTuling.Name] = DocTuling
		helpdocset[DocWater.Name] = DocWater
		helpdocset[DocMusic.Name] = DocMusic
		helpdocset[DocRandSeto.Name] = DocRandSeto
	}
	if len(nameList) != 0 {
		funName := strings.Join(nameList, " ")
		doc, ok := helpdocset[funName]
		if !ok {
			sendMsg(fromGroup, fromQQ, fmt.Sprintf("没有找到名为<%s>的命令，你是不是打错了\n(●'◡'●)", nameList[0]))
			return
		}
		doc.readme(fromGroup, fromQQ)

	} else {
		text := make([]string, 0, len(helpdocset))
		for functionName := range helpdocset {
			text = append(text, functionName)
		}
		sendMsg(fromGroup, fromQQ, "使用样例格式输入命令查看详细帮助内容\n使用帮助<空格><命令名>\n例:\n帮助 以图搜图\n帮助 控制台\n\n以下为命令名单(命令关键字 与 命令名并不相同)")
		sendMsg(fromGroup, fromQQ, strings.Join(text, "\n"))
	}
}
