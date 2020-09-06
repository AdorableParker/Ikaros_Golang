package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type shipSalvageMap struct {
	UsedName    string `gorm:"column:UsedName"`
	CurrentName string `gorm:"column:CurrentName"`
	Chapter1    uint8  `gorm:"column:Chapter1"`
	Chapter2    uint8  `gorm:"column:Chapter2"`
	Chapter3    uint8  `gorm:"column:Chapter3"`
	Chapter4    uint8  `gorm:"column:Chapter4"`
	Chapter5    uint8  `gorm:"column:Chapter5"`
	Chapter6    uint8  `gorm:"column:Chapter6"`
	Chapter7    uint8  `gorm:"column:Chapter7"`
	Chapter8    uint8  `gorm:"column:Chapter8"`
	Chapter9    uint8  `gorm:"column:Chapter9"`
	Chapter10   uint8  `gorm:"column:Chapter10"`
	Chapter11   uint8  `gorm:"column:Chapter11"`
	Chapter12   uint8  `gorm:"column:Chapter12"`
	Chapter13   uint8  `gorm:"column:Chapter13"`
}

// DocShipMap 碧蓝航线舰船打捞定位功能文档
var DocShipMap = &HelpDoc{
	Name:        "碧蓝航线舰船打捞定位",
	KeyWord:     []string{"打捞定位"},
	Example:     "打捞定位 萨拉托加\n打捞定位 3-4",
	Description: "打捞定位<空格><船名|地图坐标>\n用于查询指定地图的掉落情况或是舰船打捞地点"}

func shipMap(msg []string, msgID int32, group, qq int64, try uint8) {
	if len(msg) == 0 { // 如果没有获取到参数
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入索引信息\nq(≧▽≦q)")
			} else {
				sendMsg(group, qq, "索引不能为空哦,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, shipMap, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}
	re := regexp.MustCompile(`\d*?\-\d`)
	mapID := re.FindAllString(strings.Replace(msg[0], "—", "-", -1), 1) // 正则匹配查找索引
	var results []string
	if mapID == nil { // 没有找到
		results = nameToMap(msg[0]) // 由名字查找
	} else {
		results = mapToName(mapID[0]) // 由地图坐标查找
	}
	for _, result := range results {
		sendMsg(group, qq, result)
	}
}

func nameToMap(index string) []string {
	index = strings.ToUpper(index) // 格式化为大写
	var shipInfos []shipSalvageMap

	// 读取数据库
	db, err := gorm.Open("sqlite3", Datadir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		return nil
	}

	db.Table("ship_map").Where("CurrentName GLOB ?", fmt.Sprintf("*%s*", index)).Or("UsedName GLOB ?", fmt.Sprintf("*%s*", index)).Find(&shipInfos)

	// 格式化输出
	if len(shipInfos) == 0 {
		return []string{fmt.Sprintf("名字中包含有 %s 的舰船无法打捞或未收录", index)}
	}

	var shipMapIndex = [13]string{"Chapter 1", "Chapter 2", "Chapter 3",
		"Chapter 4", "Chapter 5", "Chapter 6", "Chapter 7", "Chapter 8",
		"Chapter 9", "Chapter 10", "Chapter 11", "Chapter 12", "Chapter 13"}

	report := make([]string, 1)
	inx := 0
	report[0] = fmt.Sprintf("名字中包含有 %s 的舰船有：", index)
	for i, shipInfo := range shipInfos { // 对每条数据
		v := reflect.ValueOf(shipInfo)
		if (i+1)%5 == 0 { // 每5条分页
			report[inx] += fmt.Sprintf("\n每页至多5条,当前第 %d 页", inx+1)
			report = append(report, "")
			inx++
		}
		report[inx] += fmt.Sprintf("\n======\n原名: %s\t和谐名：%s\n可在以下地点打捞", shipInfo.UsedName, shipInfo.CurrentName)
		for j := 2; j <= 14; j++ {
			if v.Field(j).Uint() != 0 {
				report[inx] += fmt.Sprintf("\n%s\n", shipMapIndex[j-2])
				if v.Field(j).Uint()&1 == 1 { // 第一节
					report[inx] += "Part.1\t"
				}
				if v.Field(j).Uint()&2 == 2 { // 第二节
					report[inx] += "Part.2\t"
				}
				if v.Field(j).Uint()&4 == 4 { // 第三节
					report[inx] += "Part.3\t"
				}
				if v.Field(j).Uint()&8 == 8 { // 第四节
					report[inx] += "Part.4\t"
				}
				report[inx] += "\n------------"
			}

		}

	}
	return report
}

func mapToName(index string) []string {
	var shipInfos []shipSalvageMap
	xy := strings.Split(index, "-")

	// 读取数据库
	db, err := gorm.Open("sqlite3", Datadir)
	defer db.Close()
	if err != nil {
		cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
		return nil
	}
	y, _ := strconv.Atoi(xy[1])
	db.Table("ship_map").Select("UsedName, CurrentName").Where(fmt.Sprintf("Chapter%s & %d = ?", xy[0], 1<<(y-1)), 1<<(y-1)).Find(&shipInfos)

	// 格式化
	if len(shipInfos) == 0 {
		return []string{fmt.Sprintf("数据库中没有找到可以在 %s 打捞的舰船", index)}
	}

	report := []string{fmt.Sprintf("在 %s 可以打捞的舰船有：", index)}
	for _, shipName := range shipInfos {
		report[0] += fmt.Sprintf("\n原名: %s\t和谐名: %s", shipName.UsedName, shipName.CurrentName)
	}
	return report
}
