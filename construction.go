package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jinzhu/gorm"
)

type shipBuilding struct {
	Currentname string `gorm:"column:Currentname"`
	Usedname    string `gorm:"column:Usedname"`
	Time        string `gorm:"column:time"`
}

type ttn []shipBuilding

func (l ttn) Len() int           { return len(l) }
func (l ttn) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ttn) Less(i, j int) bool { return len(l[i].Currentname) < len(l[j].Currentname) }

type ntt []shipBuilding

func (l ntt) Len() int           { return len(l) }
func (l ntt) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ntt) Less(i, j int) bool { return s2i(l[i].Time) < s2i(l[j].Time) }
func s2i(s string) int {
	list := strings.Split(s, ":")
	h, _ := strconv.Atoi(list[0])
	m, _ := strconv.Atoi(list[1])
	return h*60 + m
}

func construction(msg []string, msgID int32, group, qq int64, try uint8) {

	if len(msg) == 0 { // 如果没有获取到参数
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入索引信息\nq(≧▽≦q)")
			} else {
				sendMsg(group, qq, "索引不能为空哦,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, construction, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}

	re := regexp.MustCompile(`\d:\d\d`)
	constructionTime := re.FindAllString(strings.Replace(msg[0], "：", ":", -1), 1) // 正则匹配查找索引
	var results []string
	if constructionTime == nil { // 没有找到
		results = nameToTime(msg[0]) // 由名字查找
	} else {
		results = timeToName(constructionTime[0]) // 由时间查找
	}
	for _, result := range results {
		sendMsg(group, qq, result)
	}
}

func nameToTime(index string) []string {
	index = strings.ToUpper(index)
	var shipInfos ntt
	// 读取数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		return nil
	}

	db.Table("AzurLane_construct_time").Where("CurrentName GLOB ?", fmt.Sprintf("*%s*", index)).Or("UsedName GLOB ?", fmt.Sprintf("*%s*", index)).Find(&shipInfos)

	// 格式化输出
	if len(shipInfos) == 0 {
		return []string{fmt.Sprintf("名字包含有 %s 的舰船不可建造或尚未收录", index)}
	}

	sort.Sort(shipInfos) // 排序

	var out = make([]string, 1)
	page := 0
	out[page] = fmt.Sprintf("名字包含有 %s 的舰船有:", index)
	for i, data := range shipInfos {
		if i%20 == 0 && i != 0 {
			out[page] += fmt.Sprintf("\n每页最多20条，当前是第%d页", page+1)
			page++
			out = append(out, "")
		}
		out[page] += fmt.Sprintf("\n原名:%s\t和谐名:%s\t建造时长:%s", data.Currentname, data.Usedname, data.Time)
	}

	out = append(out, fmt.Sprintf("结果共计%d条,已全部列出", len(shipInfos)))

	return out
}

func timeToName(index string) []string {
	// 读取数据库
	db, err := gorm.Open("sqlite3", Datedir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		return nil
	}
	var shipInfos ttn
	db.Table("AzurLane_construct_time").Where("time = ?", index).Find(&shipInfos)

	// 格式化输出
	if len(shipInfos) == 0 {
		return []string{fmt.Sprintf("没有建造时间为 %s 的舰船或尚未收录", index)}
	}
	sort.Sort(shipInfos) // 排序

	var out = make([]string, 1)
	page := 0
	out[page] = fmt.Sprintf("建造时间为 %s 的舰船有:", index)
	for i, data := range shipInfos {
		if i%50 == 0 && i != 0 {
			out[page] += fmt.Sprintf("\n每页最多50条，当前是第%d页", page+1)
			page++
			out = append(out, "")
		}
		out[page] += fmt.Sprintf("\n原名:%s\t和谐名:%s", data.Currentname, data.Usedname)
	}

	out = append(out, fmt.Sprintf("结果共计%d条,已全部列出", len(shipInfos)))

	return out
}
