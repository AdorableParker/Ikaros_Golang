package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	mapset "github.com/deckarep/golang-set"
)

type aiQA struct {
	Answer   string `gorm:"column:answer"`
	Question string `gorm:"column:question"`
	Keys     string `gorm:"column:keys"`
	From     int64  `gorm:"column:from"`
}

type favorScore struct {
	ID int64 `gorm:"column:ID"`
	// Score string `gorm:"column:score"`  // 好感度分数
	Banned bool `gorm:"column:banned"`
}

var response = [...]string{"我记住了你的话，因为你的认真教导，好感度上升了",
	"我记住了你的话，你教学的时候太严厉了，好感度下降了",
	"这样的吗，我大概记住了\nฅ( ̳• ◡ • ̳)ฅ",
	"我喜欢学习\nヾ(◍°∇°◍)ﾉﾞ",
	"虽然不太懂，但是我还是把你教的知识记在了心里"}

// DocTuling 图灵对话功能文档
var DocTuling = &HelpDoc{
	Name:    "图灵AI",
	// KeyWord: []string{"我"},
	// Example: "我 你好啊"}
}

func tuling(msg string, group, qq int64, flag bool) {
	var ai []aiQA
	msg = strings.TrimSpace(msg)              // 移除前后端空格
	wordinfos := Jb.ExtractWithWeight(msg, 3) // 关键词提取
	compareSources := Jb.Cut(msg, true)       // 分词
	source := mapset.NewSet()                 // 建立集合
	for _, word := range compareSources {
		source.Add(word)
	}
	// cqp.AddLog(0, "测试信息", fmt.Sprintf("测试信息:%v", source))
	if !DBConn {
		if flag {
			sendMsg(group, qq, "数据库离线状态")
		}
		return
	}
	// 链接数据库
	db, err := gorm.Open("sqlite3", Appdir+"Ai.db")
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		sendMsg(group, qq, "数据库连接失败,智商已离线")
		return
	}

	for _, i := range wordinfos { // 第一次 关键词索引寻找
		// 查询数据库
		db.Table("universal_corpus").Select("answer, question").Where("keys = ?", i.Word).Find(&ai)
		answerList := filter(ai, source, 0.65) // 评分

		numAanswers := len(answerList)
		if numAanswers != 0 {
			rand.Seed(time.Now().UnixNano())
			sendMsg(group, qq, answerList[rand.Intn(numAanswers)])
			return
		}
	}

	db.Table("universal_corpus").Select("answer, question").Where("question = ?", msg).Find(&ai)
	answerList := filter(ai, source, 0.75) // 评分
	numAanswers := len(answerList)
	if numAanswers != 0 {
		rand.Seed(time.Now().UnixNano())
		sendMsg(group, qq, answerList[rand.Intn(numAanswers)])
		return
	}

	if flag {
		sendMsg(group, qq, "你在说什么，我怎么听不懂\n(○´･д･)ﾉ")
	}
}

func filter(ai []aiQA, source mapset.Set, maxScore float32) []string {
	var QAList = make(map[string]string, 0)
	for _, pair := range ai {
		QAList[pair.Answer] = pair.Question
	}
	var answerList = make([]string, 0)
	for a, q := range QAList { // 对每个问答组
		contrast := mapset.NewSet()
		for _, word := range Jb.Cut(q, true) { // 分词
			contrast.Add(word)
		}
		score := float32(source.Intersect(contrast).Cardinality()) / float32(source.Union(contrast).Cardinality())
		// cqp.AddLog(0, "测试信息", fmt.Sprintf("测试信息:%v\n%v\n%v\n%v", score, contrast, q, a))
		switch {
		case score > maxScore:
			maxScore = score
			answerList = []string{a}
		case score == maxScore:
			answerList = append(answerList, a)
		}

	}
	return answerList
}

// DocTraining 图灵教学功能文档
var DocTraining = &HelpDoc{
	Name:        "图灵教学",
	KeyWord:     []string{"教学", "训练", "调教"},
	Example:     "训练 生命、宇宙以及万物的答案是什么#42",
	Description: "我会完全信任你教给她的所有知识，她把你教给她的所有知识视作珍宝并会很认真的将其牢牢记住..所以请不要让她学坏哦！"}

func training(msgs []string, msgID int32, group, qq int64, try uint8) {
	if !DBConn {
		sendMsg(group, qq, "数据库离线状态")
		return
	}

	if len(msgs) == 0 {
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入问答哦,格式为:问题#回答\n例如:还记得我们的约定吗#我会永远记得的")
				sendMsg(group, qq, "注意：问答将录入至通用数据库,所有群组都能触发该问答,所以请不要教学带有针对特定对象、特定语境的问答")
			} else {
				sendMsg(group, qq, "我没有看懂,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, training, msgs, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}
	msgString := strings.Join(msgs, " ")
	if !strings.Contains(msgString, "#") {
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入问答哦,格式为:问题#回答\n例如:还记得我们的约定吗#我会永远记得的")
			} else {
				sendMsg(group, qq, "我没有看懂,检查一下格式,再发一次吧\n格式为:问题#回答\n例如：还记得我们的约定吗#我会永远记得的") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, training, msgs, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}

	// 链接数据库
	db, err := gorm.Open("sqlite3", Appdir+"Ai.db")
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		sendMsg(group, qq, "数据库连接失败,智商已离线")
		return
	}

	var blacklist favorScore
	db.Table("favor_score").Where("ID = ?", qq).Find(&blacklist)
	// cqp.AddLog(0,"测试",fmt.Sprintln(blacklist))
	if blacklist.Banned {
		sendMsg(group, qq, "很多哥哥姐姐说你是坏人，我不要听你的的ヾ(≧へ≦)〃\n（由于多人举报，你已经被加入黑名单了哦)")
		return
	}

	// 解析问答组
	var QAPair aiQA
	QA := strings.SplitN(msgString, "#", 2)
	keyWord := Jb.ExtractWithWeight(QA[0], 1) // 关键词提取
	if len(keyWord) <= 0 {
		QAPair.Keys = QA[0]
	} else {
		QAPair.Keys = keyWord[0].Word
	}

	QAPair.Question = QA[0]
	QAPair.Answer = QA[1]
	QAPair.From = qq

	// 写入数据库
	db.Table("universal_corpus").Create(&QAPair)

	rand.Seed(time.Now().UnixNano()) // 置随机数种子
	sendMsg(group, qq, response[rand.Intn(4)])
}
