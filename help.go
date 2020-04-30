package main

import (
	"fmt"
	"strings"
)

var helpText = map[string]string{
	"群活跃数据": `------water------
命令关键字: "群活跃数据"

效果: 返回当前群的群活跃数据统计信息

######################`,
	">快捷禁言": `------banned------
命令关键字: "求口", "自助禁言"
命令输入格式: 

求口<空格><时间分钟>
例:
自助禁言 30

效果: 满足你奇怪的要求
######################

------all_banned------
命令关键字: "全员禁言", "全员自闭"

&群管理权限者使用有效&

命令关键字: 

开启当前群的全员禁言

效果: 全都给我安静
######################

------all_not_banned------
命令关键字: "解禁"

解除当前群的全员禁言

效果: 阿拉霍洞开，封印解除
######################`,
	">wiki榜单": `------equipment_ranking------
命令关键字: "装备榜单", "装备榜", "装备排行榜"

效果: 装备强度测评榜
######################

------srength_ranking------
命令关键字: "强度榜单", "强度榜", "舰娘强度榜", "舰娘排行榜"

效果: 舰娘强度测评榜
######################

------pixiv_ranking------
命令关键字: "社保榜", "射爆榜", "P站榜", "p站榜"

效果: P站搜索结果排行榜

######################`,
	">B站动态获取": `------update_bilibili------
关键字: "小加加", "火星加", "B博更新", "b博更新"

效果: 锉刀怪又gū了？让我康康

######################

---bilibili_jp_Twitter---
关键字: "转推姬", "碧蓝日推"

效果: 这特么是什么东西

######################

--update_bili_Arknights--
关键字: "罗德岛线报", "方舟公告", "方舟B博", "阿米娅"

效果: 博士，您还有许多事情需要处理，现在还不能休息哦。

# 欲启用自动更新模式请使用命令 "控制台" 获取更多信息`,
	"#以图搜图": `------img_saucenao------
命令关键字: "saucenao", "图片搜索", "搜图"
命令输入格式: 

图片搜索<空格><图片>
例:
搜图 [图片]

效果: 根据输入的图片，使用saucenao搜图引擎进行图源搜索，返回相似度最高的。
######################`,
	"建造时间查询": `------construction_------
命令关键字: "建造时间查询", "建造时间", "建造查询"
命令输入格式: 

建造时间查询<空格><时间|船名>
例: 
建造时间 0:27
建造时间 萨拉托加

效果: 返回数据库中，符合的船名和时间
######################`,
	"打捞定位": `------ship_map------
命令关键字: "打捞定位"
命令输入格式: 

打捞定位<空格><船名|地图坐标>
例: 
打捞定位 萨拉托加
打捞定位 3-4

效果: 返回对应的船名列表或者地图坐标列表
######################`,
	"#碧蓝航线活动进度": `------activity------
命令关键字: "活动进度", "进度计算", "奖池计算"
命令输入格式: 

活动进度<空格><已刷点数>#[目标点数]
注: #号分隔可选参数
例: 
活动进度 12345
活动进度 12345#67890

效果: 根据已刷点数，返回活动进度报告
######################`,
	">控制台": `&需群管理员以上权限才能触发&
控制命令:
1.改变复读姬状态
2.改变开火许可状态
3.设定新入群禁言时间
4.改变火星时报订阅状态
5.改变标枪快讯订阅状态
6.改变罗德岛线报订阅状态
7.改变报时鸟状态
8.改变报时鸟_舰C版状态
9.改变迎新功能状态
10.改变每日提醒_舰B版功能状态
11.改变每日提醒_FGO版功能状态

[序号只是用来看的]`,
	">图灵AI": `------tuling------
关键字："伊卡洛斯"

"我是娱乐用人造天使，α型号「伊卡洛斯」，My Master。"
\t\t\t————伊卡洛斯

------training------
命令关键字："教学", "训练", "调教"
命令输入格式：

训练<空格><问题>#<回答>
例: 
训练 生命、宇宙以及万物的答案是什么#42

$ 伊卡洛斯会完全信任你教给她的所有知识，她把你教给她的所有知识视作珍宝并会很认真的将其牢牢记住..所以请不要让她学坏哦！`}

func help(nameList []string, fromGroup, fromQQ int64) {
	if len(nameList) != 0 {
		if helpText[nameList[0]] != "" {
			sendMsg(fromGroup, fromQQ, helpText[nameList[0]])

		} else {
			sendMsg(fromGroup, fromQQ, fmt.Sprintf("没有找到名为<%s>的命令集，你是不是打错了\n(●'◡'●)", nameList[0]))

		}
	} else {
		text := make([]string, 0, len(helpText))
		for i := range helpText {
			text = append(text, i)
		}
		sendMsg(fromGroup, fromQQ, strings.Join(text, "\n"))
		sendMsg(fromGroup, fromQQ, "带有>标志的为命令集\n带有#标志的为非同名命令,即命令关键词与命令名不一致\n具体命令关键字都请查看详细内容获知")
		sendMsg(fromGroup, fromQQ, "查看详细帮助内容\n使用帮助<空格><命令名>\n例:\n帮助 >wiki榜单\n效果: 根据输入的命令名，返回帮助信息\n######################")
	}
}
