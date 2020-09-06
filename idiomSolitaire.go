package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jinzhu/gorm"
)

type mark struct {
	Word     string
	Phonetic string
}

type idioms struct {
	ID   int    `gorm:"column:ID"`
	Word string `gorm:"column:word"`
	Head string `gorm:"column:headPhonetic"`
	Tail string `gorm:"column:tailPhonetic"`
}

var solitaireGroupList = make(map[int64]*mark)

func solutaire(oldPhonetic, input string) (word, phonetic string, ok uint8) {
	db, err := gorm.Open("sqlite3", Datadir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		return
	}
	rand.Seed(time.Now().UnixNano())
	var idiom idioms
	if oldPhonetic == "" {
		db.Table("idiomsDictionary").Where("ID = ?", rand.Intn(29500)+1).Find(&idiom)
		return idiom.Word, idiom.Tail, 0
	}
	db.Table("idiomsDictionary").Where("Word = ?", input).Find(&idiom)
	if idiom.Word == "" {
		return word, phonetic, 1
	}
	if idiom.Head != oldPhonetic {
		return word, phonetic, 2
	}
	var idiomList []idioms
	db.Table("idiomsDictionary").Where("headPhonetic = ?", idiom.Tail).Find(&idiomList)
	if len(idiomList) != 0 {
		idiom = idiomList[rand.Intn(len(idiomList))]
		return idiom.Word, idiom.Tail, 0
	}
	return word, phonetic, 3
}

// DocSolitaire 成语接龙功能文档
var DocSolitaire = &HelpDoc{
	Name:        "成语接龙",
	KeyWord:     []string{"成语接龙", "接龙"},
	Example:     "接龙 不越雷池",
	Description: "成语接龙,同音即可,回答前需加前缀"}

func solitaire(msg []string, group, qq int64) {
	oldPhonetic, ok := solitaireGroupList[group]
	if !ok {
		word, phonetic, _ := solutaire("", "")
		solitaireGroupList[group] = &mark{Word: word, Phonetic: phonetic}
		sendMsg(group, qq, fmt.Sprintf("开始成语接龙咯,题目是 %s", word))
		return
	}
	if len(msg) == 0 {
		sendMsg(group, qq, "你什么也没答啊")
		return
	}
	if msg[0] == "不玩了" {
		delete(solitaireGroupList, group)
		sendMsg(group, qq, "主人说不玩了那咱就不玩了")
		return
	}
	word, phonetic, code := solutaire(oldPhonetic.Phonetic, msg[0])
	switch code {
	case 0:
		oldPhonetic.Word = word
		oldPhonetic.Phonetic = phonetic
		sendMsg(group, qq, fmt.Sprintf("我接 %s", word))
	case 1:
		sendMsg(group, qq, "这个成语我没听过呢,我觉得你在骗我")
	case 2:
		sendMsg(group, qq, fmt.Sprintf("错了,题目是 %s,看清楚再答啊", oldPhonetic.Word))
	case 3:
		delete(solitaireGroupList, group)
		sendMsg(group, qq, "不玩啦,你肯定是作弊,这题我不会辣ε(┬┬﹏┬┬)3")
	}
}
