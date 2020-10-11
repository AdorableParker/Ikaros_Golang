package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jinzhu/gorm"
)

type setoMod struct {
	Mod   uint8 `gorm:"column:SetoMod"` // 模式
	Quota int   `gorm:"column:quota"`   // 剩余份额
	Score int   `gorm:"column:score"`   // 审核积分
	Date  int64 `gorm:"column:date"`    // 最后更新时间
	Fine  int   `gorm:"column:fine"`    // 惩罚
}

var loger *log.Logger

// DocRandSeto 随机色图功能文档
var DocRandSeto = &HelpDoc{
	Name:        "随机色图功能",
	KeyWord:     []string{"随机色图"},
	Description: "本功能分为“安全模式”和“审核模式”\n\n处于安全模式时,发出的图一定不是露点图\n处于审核模式时,有概率出现露点图\n\n本功能可由群管理或以上权限者设定,相关控制命令详见“help 控制台”\n\n以群为单位，每天拥有20张图片份额,份额刷新时间为每日凌晨3点\n\n另外:本着`人人为我,我为人人`的集体主义精神,审核模式要求用户对于发出的图片进行正确的审核,审核内容主要为对于是否露点做出正确判断并反馈\n作为感谢对于净化社区环境所做出贡献的奖励\n做出正确判断可以获取更多每日份额(每10次反馈予以1张每日份额奖励)\n但若发现恶意错误反馈获取份额的行为\n将会被处永久降低每日份额的惩罚(降低量5张起步,下至永久关闭该群该功能)\n惩罚不予申诉"}

func init() {
	file := filepath.Join(Appdir, "logs", fmt.Sprintf("Icarus-%s.log", time.Now().Format("2006-1-2")))
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	loger = log.New(logFile, "[seTo]", log.LstdFlags) // 将文件设置为loger作为输出
	return
}

func randSeto(msg []string, msgID int32, group, qq int64, try uint8) {
	// 读取数据库
	db, err := gorm.Open("sqlite3", Datadir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		return
	}
	var seto setoMod
	if len(msg) >= 2 {
		if try > 0 && try <= 3 { // 如果已尝试次数不超过3次
			try++ // 已尝试次数+1
			judgment, err := strconv.ParseBool(msg[1])
			if err != nil {
				sendMsg(group, qq, "审核结果\n只接受 1/0、t/f、T/F 这几个输入哦\n(。・ω・。)")
				stagedSessionPool[msgID] = newStagedSession(group, qq, randSeto, []string{msg[0]}, try)
				return
			}
			markImg(msg[0], judgment)
			db.Table("rendSeto").Select("score").Where("groupID = ?", group).Find(&seto)   // 读取分数
			db.Table("rendSeto").Where("groupID = ?", group).Update("score", seto.Score+1) // 更新分数
			sendMsg(group, qq, fmt.Sprintf("输入审核结果为 %t,当前群奖励分数为 %d", judgment, seto.Score+1))
			writeLog(msg[0], group, judgment)
			return
		}
		sendMsg(group, qq, "错误次数过多,本次审核失败(失败结果将不计入奖励分数)(｀･ω･´)ゞ")
		return
	}

	db.Table("group_info").Select("SetoMod").Where("group_id = ?", group).Find(&seto)

	if seto.Mod == 0 {
		sendMsg(group, qq, "关闭状态")
		return
	}
	db.Table("rendSeto").Where("groupID = ?", group).Find(&seto)

	now := time.Now().Unix() // 现在时间戳

	today, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	today3clock := today.Add(3 * time.Hour).Unix() // 今三点时间戳

	if now >= today3clock {
		if seto.Date < today3clock {
			seto.Quota = 20 + seto.Score/10 - seto.Fine*5
			db.Table("rendSeto").Where("groupID = ?", group).Update("quota", seto.Quota)
			db.Table("rendSeto").Where("groupID = ?", group).Update("date", now)
			// db.Table("rendSeto").Where("groupID = ?", group).Find(&seto)
		}
	} else {
		yesterday3clock := today.Add(-21 * time.Hour).Unix()
		if seto.Date < yesterday3clock {
			seto.Quota = 20 + seto.Score/10 - seto.Fine*5
			db.Table("rendSeto").Where("groupID = ?", group).Update("quota", seto.Quota)
			db.Table("rendSeto").Where("groupID = ?", group).Update("date", now)
			// db.Table("rendSeto").Where("groupID = ?", group).Find(&seto)
		}
	}

	if seto.Quota <= 0 {
		sendMsg(group, qq, "本群今日份额已用尽,等凌晨3点的刷新吧(｡ŏ﹏ŏ)\n进行正确审核可以获得更多的每日份额哦ヾ(≧O≦)〃")
		return
	}
	var imgName string
	switch seto.Mod {
	case 1:
		sendMsg(group, qq, "当前为安全模式")
		imgName = filepath.Join(Appdir, "seto", "Safe", sendSeto(seto.Mod))
	case 2:
		sendMsg(group, qq, "当前为审核模式,请发送审核结果,是否露点?(仅接受 1/0、t/f)")
		imgName = sendSeto(seto.Mod)
		stagedSessionPool[msgID] = newStagedSession(group, qq, randSeto, []string{imgName}, 1) // 添加新的会话到会话池
		imgName = filepath.Join(Appdir, "seto", "Unsafe", imgName)
	}
	// cqp.AddLog(0, "CQCode", fmt.Sprintf("[CQ:image,file=%s]", imgName))
	sendMsg(group, qq, fmt.Sprintf("[CQ:image,file=%s]", imgName))
	// sendMsg(group, qq, `[CQ:image,file=00B42DD8A147B5CE5D88B88723B61797(1181×1181).jpg]`)
	db.Table("rendSeto").Where("groupID = ?", group).Update("quota", seto.Quota-1)
	sendMsg(group, qq, fmt.Sprintf(
		"本群今日剩余份额为 %d\n下次刷新份额为 %d\n扣除错误审核处罚份额 %d ,实得份额 %d",
		seto.Quota-1,                  // 剩余份额
		20+seto.Score/10,              // 刷新份额
		seto.Fine*5,                   // 处罚份额
		20+seto.Score/10-seto.Fine*5)) // 实际份额
}

func sendSeto(mod uint8) string {
	setoMod := [...]string{"Safe", "Unsafe"}
	// 取文件列表
	files, err := ioutil.ReadDir(filepath.Join(Appdir, "seto", setoMod[mod-1]))
	if err != nil {
		cqp.AddLog(30, "文件列表读取错误", fmt.Sprintln(err))
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return files[rand.Intn(len(files))].Name()
}

func markImg(name string, value bool) {
	var d string
	if value {
		d = "MarkUnsafe"
	} else {
		d = "MarkSafe"
	}
	os.Rename(filepath.Join(Appdir, "seto", "Unsafe", name), filepath.Join(Appdir, "seto", d, name))
}

func writeLog(name string, id int64, flag bool) {
	loger.Println(name, id, flag)
	// loger.Print("Hello, log file!")
	// fmt.Print(&buf)
}
