package main

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/yanyiwu/gojieba"
)

//go:generate cqcfg -c .
// cqp: 名称: Ikaros
// cqp: 版本: 1.0.0:1
// cqp: 作者: adorableparker
// cqp: 简介: 模板测试

type stagedSession struct {
	Group          int64                                       // 调用者所在群
	QQ             int64                                       // 调用者QQ
	Function       *func([]string, int32, int64, int64, uint8) // 执行的命令函数
	Parameter      []string                                    // 参数
	TryOpportunity uint8                                       // 已尝试次数
}

var stagedSessionPool = make(map[int32]*stagedSession, 30)

// Datedir 数据库位置
var Datedir string

// Appdir 配置文件位置
var Appdir string

// atMe @我 的CQ码
var atMe string

// Jb 结巴库对象
var Jb *gojieba.Jieba

// 线程调度
var wg sync.WaitGroup

// SauceNAO 初始化开关
var SauceNAO bool

func newStagedSession(group, qq int64, function func([]string, int32, int64, int64, uint8), parameter []string, try uint8) *stagedSession {
	return &stagedSession{
		Group:          group,     // 调用者所在群
		QQ:             qq,        // 调用者QQ
		Function:       &function, // 执行的命令函数
		Parameter:      parameter, // 参数
		TryOpportunity: try,       // 已尝试次数
	}
}

func main() { /*此处应当留空*/ }

func init() {
	cqp.AppID = "io.github.adorableparker.Ikaros" // TODO: 修改为这个插件的ID
	cqp.PrivateMsg = onPrivateMsg
	cqp.GroupMsg = onGroupMsg
	cqp.Enable = onEnable
	cqp.Disable = onDisable
	cqp.GroupMemberIncrease = onGroupMemberIncrease

}

func onEnable() int32 {
	SauceNAO = false
	Jb = gojieba.NewJieba()
	Appdir = cqp.GetAppDir()
	Datedir = filepath.Join(Appdir, "User.db")
	atMe = fmt.Sprintf("[CQ:at,qq=%d]", cqp.GetLoginQQ())

	// 每小时的 报时任务
	wg.Add(1)
	go callBellTask()

	// 每六分钟的 检查动态更新任务
	wg.Add(1)
	go updateCheckTask()

	// 每天的晚九点半 提醒任务
	wg.Add(1)
	go remindTask()

	wg.Wait()
	return 0
}

func onDisable() int32 {
	Jb.Free()
	// cron.Clear()
	wg.Done()
	wg.Done()
	wg.Done()
	return 0
}

func onPrivateMsg(subType, msgID int32, fromQQ int64, msg string, font int32) int32 {
	// cqp.SendPrivateMsg(fromQQ, msg) //复读机
	ok := parser(msgID, -1, fromQQ, msg)
	if !ok {
		tuling(msg, -1, fromQQ, true)
	}
	return 0
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	if atForMe(msg) {
		tuling(strings.NewReplacer(atMe, "").Replace(msg), fromGroup, fromQQ, true)
		return 0
	}
	ok := parser(msgID, fromGroup, fromQQ, msg)
	if !ok {
		if strings.Contains(msg, "射爆") {
			if fire(fromGroup, fromQQ) {
				return 0
			}
		}
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(10) <= 2 {
			tuling(msg, fromGroup, fromQQ, false)
			return 0
		}
		repeater(msg, fromGroup)
	}
	return 0
}

func parser(msgID int32, fromGroup, fromQQ int64, msg string) (ok bool) {
	instructionPacket := strings.Fields(msg) // 解析消息
	for i, s := range stagedSessionPool {    // 在会话池里面查找符合标志的
		if s.Group == fromGroup && s.QQ == fromQQ { // 存在会话任务
			if strings.Contains(msg, "算了") || strings.Contains(msg, "不用了") || strings.Contains(msg, "不搜了") || strings.Contains(msg, "不查了") {
				delete(stagedSessionPool, i) // 删除掉这个会话任务
				sendMsg(fromGroup, fromQQ, "命令已取消")
				return true
			}
			fun := *(s.Function)                                           // 获取功能函数对象
			fun(instructionPacket, msgID, s.Group, s.QQ, s.TryOpportunity) // 执行功能
			delete(stagedSessionPool, i)                                   // 删除掉这个会话任务
			return true
		}
	}

	if len(instructionPacket) != 0 { // 如果有前缀
		ok = functionList(instructionPacket, msgID, fromGroup, fromQQ) // 判定功能触发
		return
	}
	return true
}

func atForMe(msg string) bool {

	if strings.Contains(msg, atMe) {
		return true
	}
	return false
}

// sendMsg 自动判断消息来源并发送消息
// sendMsg(group, qq int64, msg)
func sendMsg(group, qq int64, msg string) {
	if group != -1 {
		cqp.SendGroupMsg(group, msg)
	} else {
		cqp.SendPrivateMsg(qq, msg)
	}
}

func functionList(msg []string, msgID int32, fromGroup, fromQQ int64) bool {
	switch msg[0] {
	case "响应池":
		cqp.AddLog(0, "调试输出", fmt.Sprintln(stagedSessionPool))
	case "training", "训练", "调教", "教学":
		training(msg[1:], msgID, fromGroup, fromQQ, 0)
	case "activity", "活动进度", "进度计算", "奖池计算":
		activity(msg[1:], msgID, fromGroup, fromQQ, 0)
	case "shipMap", "打捞定位":
		shipMap(msg[1:], msgID, fromGroup, fromQQ, 0)
	case "construction", "建造时间查询", "建造时间", "建造查询":
		construction(msg[1:], msgID, fromGroup, fromQQ, 0)
	case "saucenao", "图片搜索", "搜图":
		saucenao(msg[1:], msgID, fromGroup, fromQQ, 0)
	case "help", "使用说明", "使用帮助", "帮助", "使用方法":
		help(msg[1:], fromGroup, fromQQ)
	case "equipmentRanking", "装备榜单", "装备榜", "装备排行榜":
		equipmentRanking(fromGroup, fromQQ)
	case "srengthRanking", "强度榜单", "强度榜", "舰娘强度榜", "舰娘排行榜":
		srengthRanking(fromGroup, fromQQ)
	case "pixivRanking", "社保榜", "射爆榜", "P站榜", "p站榜":
		pixivRanking(fromGroup, fromQQ)
	case "小加加", "火星加", "B博更新", "b博更新":
		sendDynamic(msg[1:], fromGroup, fromQQ, 233114659)
	case "转推姬", "碧蓝日推":
		sendDynamic(msg[1:], fromGroup, fromQQ, 300123440)
	case "罗德岛线报", "方舟公告", "方舟B博", "阿米娅":
		sendDynamic(msg[1:], fromGroup, fromQQ, 161775300)
	case "伊卡洛斯":
		words := strings.Join(msg[1:], " ") // 拼接字符串
		tuling(words, fromGroup, fromQQ, true)
	default: // 只有群消息有效
		if fromGroup == -1 {
			return false
		}
		switch msg[0] {
		case "water", "群活跃数据":
			water(fromGroup)
		case "allNotBanned", "解禁":
			allNotBanned(fromGroup)
		case "banned", "求口", "自助禁言":
			banned(msg[1:], msgID, fromGroup, fromQQ, 0)
		default:
			groupInfo := cqp.GetGroupMemberInfo(fromGroup, fromQQ, true)
			if groupInfo.Auth >= 2 {
				switch msg[0] {
				case "allBanned", "全员禁言", "全员自闭":
					allBanned(fromGroup, fromQQ)
				case "改变报时鸟状态":
					callBellAlter(fromGroup, true)
				case "改变报时鸟_舰C版状态":
					callBellAlter(fromGroup, false)
				case "改变每日提醒_舰B版功能状态":
					dailyRemindAlter(fromGroup)
				case "改变每日提醒_FGO版功能状态":
					dailyRemindAlterFGO(fromGroup)
				case "改变火星时报订阅状态":
					saraNewsAlter(fromGroup)
				case "改变标枪快讯订阅状态":
					javelinNewsAlter(fromGroup)
				case "改变罗德岛线报订阅状态":
					arknightsAlter(fromGroup)
				case "改变开火许可状态":
					fireAlter(fromGroup)
				case "改变复读姬状态":
					repeatAlter(fromGroup)
				case "改变迎新功能状态":
					newAddAlter(fromGroup)
				case "设定新入群禁言时间":
					groupPolicy(msg[1:], msgID, fromGroup, fromQQ, 0)
				default:
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}
