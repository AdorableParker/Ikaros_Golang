package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type groupInfo struct {
	ID      int   `gorm:"column:id"`
	GroupID int64 `gorm:"column:group_id"`

	Repeat      uint8 `gorm:"column:repeat"`      // 复读
	AutoTrigger uint8 `gorm:"column:autoTrigger"` // 开火限制

	// 动态更新
	Arknights      uint8 `gorm:"column:Arknights"`
	SaraNews       uint8 `gorm:"column:Sara_news"`
	JavelinNews    uint8 `gorm:"column:Javelin_news"`
	FateGrandOrder uint8 `gorm:"column:FateGrandOrder"`

	// 报时
	CallBell uint8 `gorm:"column:Call_bell"`

	// 每日提醒
	DailyRemindAzurLane uint8 `gorm:"column:Daily_remind_AzurLane"`
	DailyRemindFGO      uint8 `gorm:"column:Daily_remind_FGO"`

	// 群策略
	NewAdd uint8 `gorm:"column:New_add"`
	Policy int64 `gorm:"column:policy"`
}

type repeatInfo struct {
	Info    string `gorm:"column:info"`
	Flag    uint8  `gorm:"column:flag"`
	GroupID int64  `gorm:"column:groupid"`
}

var real = [2]bool{false, true}

// DocConsole 控制台功能文档
var DocConsole = &HelpDoc{
	Name: "控制台",
	KeyWord: []string{
		"改变复读姬状态\n", "改变主动对话许可状态\n", "设定新入群禁言时间\n",
		"改变火星时报订阅状态\n", "改变标枪快讯订阅状态\n", "改变罗德岛线报订阅状态\n",
		"改变FGO订阅状态\n", "改变报时鸟模式\n",
		"改变迎新功能状态\n", "改变每日提醒_舰B版功能状态\n", "改变每日提醒_FGO版功能状态\n"},
	Example:     "改变复读姬状态\n改变报时鸟模式 1\n设定新入群禁言时间 5",
	Description: "需群管理员以上权限才能触发\n报时模式目前共四种:\n0\t关闭\n1\t标准\n2\t舰C\n3\t明日方舟"}

func fireAlter(group int64) {
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("autoTrigger", 1^g.AutoTrigger)
	cqp.SendGroupMsg(group, fmt.Sprintf("主动对话原状态为 %t\n现状态已改为 %t", real[g.AutoTrigger], real[1^g.AutoTrigger]))
}

func repeatAlter(group int64) {
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("repeat", 1^g.Repeat)
	cqp.SendGroupMsg(group, fmt.Sprintf("复读姬原状态为 %t\n现状态已改为 %t", real[g.Repeat], real[1^g.Repeat]))
}

func arknightsAlter(group int64) {
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("Arknights", 1^g.Arknights)
	cqp.SendGroupMsg(group, fmt.Sprintf("明日方舟　B站动态订阅原状态为 %t\n现状态已改为 %t", real[g.Arknights], real[1^g.Arknights]))
}

func fgoAlter(group int64) {
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("FateGrandOrder", 1^g.FateGrandOrder)
	cqp.SendGroupMsg(group, fmt.Sprintf("命运－冠位指定 B站动态订阅原状态为 %t\n现状态已改为 %t", real[g.FateGrandOrder], real[1^g.FateGrandOrder]))
}

func saraNewsAlter(group int64) {
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("Sara_news", 1^g.SaraNews)
	cqp.SendGroupMsg(group, fmt.Sprintf("碧蓝航线　B站动态订阅原状态为 %t\n现状态已改为 %t", real[g.SaraNews], real[1^g.SaraNews]))
}

func javelinNewsAlter(group int64) {
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("Javelin_news", 1^g.JavelinNews)
	cqp.SendGroupMsg(group, fmt.Sprintf("火星bot小黄瓜　B站动态订阅原状态为 %t\n现状态已改为 %t", real[g.JavelinNews], real[1^g.JavelinNews]))
}

// callBellAlter(group int64, flag uint8)
// group 群号码
// flag = 0 关闭
// flag = 1 标准
// flag = 2 舰C
// flag = 3 明日方舟
func callBellAlter(group int64, msg []string) {
	var g groupInfo
	var mod = [...]string{"关闭", "标准", "舰C", "明日方舟"}
	if len(msg) == 0 {
		cqp.SendGroupMsg(group, "无参数,执行失败")
		return
	}
	flag, err := strconv.ParseUint(msg[0], 0, 8)
	if err != nil || flag >= 4 {
		cqp.SendGroupMsg(group, "无效参数,执行失败")
		return
	}
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("Call_bell", flag)
	cqp.SendGroupMsg(group, fmt.Sprintf("报时鸟原模式为 %s\n现已改为 %s", mod[g.CallBell], mod[flag]))
}

func dailyRemindAlter(group int64) {
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("Daily_remind_AzurLane", 1^g.DailyRemindAzurLane)
	cqp.SendGroupMsg(group, fmt.Sprintf("每日提醒_舰B版功能原状态为 %t\n现状态已改为 %t", real[g.DailyRemindAzurLane], real[1^g.DailyRemindAzurLane]))
}

func dailyRemindAlterFGO(group int64) {
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("Daily_remind_FGO", 1^g.DailyRemindFGO)
	cqp.SendGroupMsg(group, fmt.Sprintf("每日提醒_FGO版功能原状态为 %t\n现状态已改为 %t", real[g.DailyRemindFGO], real[1^g.DailyRemindFGO]))
}

func newAddAlter(group int64) {
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("New_add", 1^g.NewAdd)
	cqp.SendGroupMsg(group, fmt.Sprintf("迎新功能原状态为 %t\n现状态已改为 %t", real[g.NewAdd], real[1^g.NewAdd]))
}

func groupPolicy(msg []string, msgID int32, group, qq int64, try uint8) {
	var bantime int
	var err error
	if len(msg) != 0 {
		bantime, err = strconv.Atoi(msg[0])
		if err != nil {
			cqp.SendGroupMsg(group, "请不要输入一些奇奇怪怪的东西\n＞︿＜")
			return
		}
	} else {
		cqp.SendGroupMsg(group, "请设定时长(单位:分钟)\n(。・ω・。)")
		try++
		stagedSessionPool[msgID] = newStagedSession(group, qq, groupPolicy, msg, try) // 将一个 待跟进会话 加入 会话池
		return
	}
	var g groupInfo
	// 链接数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	db.Table("group_info").Where("group_id = ?", group).Update("policy", bantime)
	cqp.SendGroupMsg(group, fmt.Sprintf("入群禁言时长原为 %d\n现状态已改为 %d", g.Policy, bantime))
}

func repeater(msg string, group int64) {
	var g groupInfo
	var r repeatInfo
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		cqp.SendGroupMsg(group, "数据库连接异常\n×_×")
		return
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	if g.ID == 0 {
		addg := groupInfo{GroupID: group}
		r = repeatInfo{GroupID: group}
		db.Table("group_info").Create(&addg)
		db.Table("repeat_info").Create(&r)
		return
	}
	if g.Repeat == 1 {
		db.Table("repeat_info").Where("groupid = ?", group).First(&r)
		if msg == r.Info {
			r.Flag++
			db.Table("repeat_info").Where("groupid = ?", group).Update("flag", r.Flag)
			if r.Flag == 2 {
				cqp.SendGroupMsg(group, msg)
			}
		} else {
			db.Table("repeat_info").Where("groupid = ?", group).Updates(map[string]interface{}{"info": msg, "flag": 0})
		}
	}
}

func autoTrigger(group int64) bool {
	var g groupInfo
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		return true
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", group).First(&g)
	if g.AutoTrigger == 1 {
		return true
	}
	return false
}

var nmcm = [...]string{"欢迎新人,能表演一下退群吗",
	"群地位-1",
	"新来的别客气,把自己当成群主就行",
	"是大佬!啊,大佬!啊!我死了",
	"你已经是群大佬了,快和萌新们打个招呼吧"}

func onGroupMemberIncrease(subType, sendTime int32, fromGroup, fromQQ, beingOperateQQ int64) int32 {
	var g groupInfo
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		return 0
	}
	// 查询数据库
	db.Table("group_info").Where("group_id = ?", fromGroup).First(&g)
	if g.NewAdd == 1 {
		rand.Seed(time.Now().UnixNano()) // 置随机数种子
		cqp.SendGroupMsg(fromGroup, nmcm[rand.Intn(4)])
	}
	if g.Policy > 0 {
		cqp.SetGroupBan(fromGroup, beingOperateQQ, g.Policy*60)
		cqp.SendGroupMsg(fromGroup, fmt.Sprintf("根据你群规定,新人禁言 %d 分钟", g.Policy))
	}
	return 0
}
