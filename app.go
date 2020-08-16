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
	"gopkg.in/ini.v1"
)

//go:generate cqcfg -c .
// cqp: 名称: Ikaros
// cqp: 版本: 1.2.0:1
// cqp: 作者: adorableparker
// cqp: 简介: 模板测试
// cqp: 菜单

type stagedSession struct {
	Group          int64                                       // 调用者所在群
	QQ             int64                                       // 调用者QQ
	Function       *func([]string, int32, int64, int64, uint8) // 执行的命令函数
	Parameter      []string                                    // 参数
	TryOpportunity uint8                                       // 已尝试次数
}

// HelpDoc 帮助文档类
type HelpDoc struct {
	Name        string   // 功能名
	KeyWord     []string // 关键字
	Example     string   // 例子
	Description string   // 说明
}

type keyConf struct {
	AdminAccount int64
	SaucenaoKey  string
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

// LoadingFinished 初始化完成标识
var LoadingFinished bool

// DBConn 数据库连接控制标记
var DBConn bool

// AuthorizedGroupList 群申请授权名录
var AuthorizedGroupList [5]int64

// AdminConfig 关键配置项
var AdminConfig = new(keyConf)

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
	cqp.AppID = "com.adorableparker.github.ikaros_golang" // TODO: 修改为这个插件的ID
	cqp.PrivateMsg = onPrivateMsg
	cqp.GroupMsg = onGroupMsg
	cqp.Enable = onEnable
	cqp.Disable = onDisable
	cqp.GroupMemberIncrease = onGroupMemberIncrease
	cqp.GroupRequest = onGroupRequest
}

func onGroupRequest(subType, sendTime int32, fromGroup, fromQQ int64, msg, responseFlag string) int32 {
	// subType 1: 加群请求 2: 被邀请入群
	// sendTime 消息时间
	// msg 验证问答,被邀请时为空
	// responseFlag 请求回馈密钥
	cqp.AddLog(0, "GroupRequest", fmt.Sprintln(subType, sendTime, fromGroup, fromQQ, msg, responseFlag))
	if subType == 2 {
		for i, authorizedgroup := range AuthorizedGroupList {
			if authorizedgroup == fromGroup {
				cqp.SetGroupAddRequest(responseFlag, subType, 1, "")
				AuthorizedGroupList[i] = 0
				// cqp.AddLog(0, "GroupRequest", "执行同意语句")
				return 0
			}
		}
		cqp.SetGroupAddRequest(responseFlag, subType, 2, "未授权的请求")
		// cqp.AddLog(0, "GroupRequest", "执行拒绝语句")
	}
	return 0
}

func onEnable() int32 {
	SauceNAO = false
	DBConn = true
	Jb = gojieba.NewJieba()
	Appdir = cqp.GetAppDir()
	Datedir = filepath.Join(Appdir, "User.db")
	// cqp.AddLog(0, "Dir", Appdir)
	atMe = fmt.Sprintf("[CQ:at,qq=%d]", cqp.GetLoginQQ())
	err := ini.MapTo(AdminConfig, Appdir+"MainConf.ini")
	if err != nil {
		cqp.AddLog(10, "关键配置文件读取异常", fmt.Sprintln(err))
		panic("关键配置文件读取异常")
	}
	// cqp.AddLog(0, "", fmt.Sprintln(AdminConfig))
	// 每小时的 报时任务
	// wg.Add(1)
	go callBellTask()

	// 每六分钟的 检查动态更新任务
	// wg.Add(1)
	go updateCheckTask()

	// 每天的晚九点半 提醒任务
	// wg.Add(1)
	go remindTask()
	LoadingFinished = true
	// wg.Wait()
	return 0
}

func onDisable() int32 {
	Jb.Free()
	// cron.Clear()
	// wg.Done()
	// wg.Done()
	// wg.Done()
	return 0
}

func onPrivateMsg(subType, msgID int32, fromQQ int64, msg string, font int32) int32 {
	// cqp.AddLog(0, "调试编码", fmt.Sprintln("私聊部分", msg))
	if !LoadingFinished {
		cqp.AddLog(0, "初始化中", "初始化完成前不处理消息")
		return 0
	}
	// cqp.SendPrivateMsg(fromQQ, msg) //复读机
	ok := parser(msgID, -1, fromQQ, msg)
	if !ok {
		tuling(msg, -1, fromQQ, true)
	}
	return 0
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	if !LoadingFinished {
		cqp.AddLog(0, "初始化中", "初始化完成前不处理消息")
		return 0
	}
	// cqp.AddLog(0, "调试编码", fmt.Sprintln("群消息", msg))

	// cqp.AddLog(0, "t", fmt.Sprintln("msgid:", msgID))
	if atForMe(msg) {
		tuling(strings.Join(strings.Split(msg, atMe), ""), fromGroup, fromQQ, true)
		return 0
	}
	ok := parser(msgID, fromGroup, fromQQ, msg)
	if !ok {
		if autoTrigger(fromGroup) {
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(10) <= 2 {
				tuling(msg, fromGroup, fromQQ, false)
				return 0
			}
			return 0
		}
		repeater(msg, fromGroup)
	}

	return 0
}

func parser(msgID int32, fromGroup, fromQQ int64, msg string) (ok bool) {
	instructionPacket := strings.Fields(msg) // 解析消息
	// cqp.AddLog(0, "t", fmt.Sprintln("poll:", stagedSessionPool))
	for i, s := range stagedSessionPool { // 在会话池里面查找符合标志的
		// cqp.AddLog(0, "t", fmt.Sprintln(i, s, fromGroup, fromQQ))
		if s.Group == fromGroup && s.QQ == fromQQ { // 存在会话任务
			if strings.Contains(msg, "算了") || strings.Contains(msg, "不用了") || strings.Contains(msg, "不搜了") || strings.Contains(msg, "不查了") {
				delete(stagedSessionPool, i) // 删除掉这个会话任务
				sendMsg(fromGroup, fromQQ, "命令已取消")
				return true
			}
			fun := *(s.Function)                                           // 获取功能函数对象
			instructionPacket = append(s.Parameter, instructionPacket...)  // 拼接参数
			fun(instructionPacket, msgID, s.Group, s.QQ, s.TryOpportunity) // 执行功能
			delete(stagedSessionPool, i)                                   // 删除掉这个会话任务
			return true
		}
	}

	if len(instructionPacket) != 0 { // 如果有前缀
		ok = functionList(instructionPacket, msgID, fromGroup, fromQQ) // 判定功能触发
		return ok
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
	case "approveAuthorization", "授权批准":
		approveAuthorization(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "calculationExp", "舰船经验", "经验计算":
		calculationExp(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "dbCM", "离线数据库", "连接数据库":
		dbCM(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "calculato", "计算", "计算器":
		calculato(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "training", "训练", "调教", "教学":
		training(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "activity", "活动进度", "进度计算", "奖池计算":
		activity(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "shipMap", "打捞定位":
		shipMap(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "realName", "船名查询", "和谐名":
		realName(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "construction", "建造时间查询", "建造时间", "建造查询":
		construction(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "saucenao", "图片搜索", "搜图":
		saucenao(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "dynamicByID", "B站动态":
		dynamicByID(msg[1:], msgID, fromGroup, fromQQ, 0)

	case "solitaire", "成语接龙", "接龙":
		solitaire(msg[1:], fromGroup, fromQQ)

	case "小加加", "火星加", "B博更新", "b博更新":
		sendDynamic(msg[1:], fromGroup, fromQQ, 233114659)

	case "转推姬", "碧蓝日推":
		sendDynamic(msg[1:], fromGroup, fromQQ, 300123440)

	case "罗德岛线报", "方舟公告", "方舟B博", "阿米娅":
		sendDynamic(msg[1:], fromGroup, fromQQ, 161775300)

	case "月球人公告", "FGO公告", "呆毛王":
		sendDynamic(msg[1:], fromGroup, fromQQ, 233108841)

	case "伊卡洛斯项目地址":
		sendMsg(fromGroup, fromQQ, "项目地址: https://adorableparker.github.io/Ikaros_Golang/\n欢迎前来送star、造轮子、提issues")

	case "help", "使用说明", "使用帮助", "帮助", "使用方法":
		help(msg[1:], fromGroup, fromQQ)

	case "equipmentRanking", "装备榜单", "装备榜", "装备排行榜":
		equipmentRanking(fromGroup, fromQQ)

	case "srengthRanking", "强度榜单", "强度榜", "舰娘强度榜", "舰娘排行榜":
		srengthRanking(fromGroup, fromQQ)

	case "pixivRanking", "社保榜", "射爆榜", "P站榜", "p站榜":
		pixivRanking(fromGroup, fromQQ)

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
				case "改变FGO订阅状态":
					fgoAlter(fromGroup)
				case "改变主动对话许可状态":
					fireAlter(fromGroup)
				case "改变复读姬状态":
					repeatAlter(fromGroup)
				case "改变迎新功能状态":
					newAddAlter(fromGroup)
				case "改变报时鸟模式":
					callBellAlter(fromGroup, msg[1:])
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
