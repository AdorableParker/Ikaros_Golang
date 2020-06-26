package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type updateCheckInfo struct {
	GroupID int64 `gorm:"column:group_id"`
}

// updateCheckTask 每六分钟的 检查动态更新任务
func updateCheckTask() {
	var calibrationCountdown uint8 = 120 // 首次运行 进行计时器校准
	for {
		if calibrationCountdown >= 120 { // 如果 校准倒数大于 120 执行校准操作
			calibrationCountdown = 0                // 进入校准部分 计数归零
			t := time.Now().Round(10 * time.Minute) // 返回最近的整十分钟数
			if t.Before(time.Now()) {               // 如果时间晚于当前时间
				t = t.Add(10 * time.Minute) // 往后推10分钟
			}
			cqp.AddLog(0, "动态更新执行时间", fmt.Sprintln(t))
			time.Sleep(t.Sub(time.Now())) // 休眠直到该时间
		}

		// 链接数据库
		db, err := gorm.Open("sqlite3", Datedir)
		if err != nil {
			cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
			return
		}

		// 火星加加
		msg, ok := updateCheckt(233114659)
		var checkList []updateCheckInfo
		if ok {
			// 查询数据库
			db.Table("group_info").Select("group_id").Where("Sara_news = ?", "1.0").Find(&checkList)
			for _, checkGroup := range checkList {
				cqp.SendGroupMsg(checkGroup.GroupID, msg)
			}
		}

		// 标枪快讯
		msg, ok = updateCheckt(300123440)
		if ok {
			// 查询数据库
			db.Table("group_info").Select("group_id").Where("Javelin_news = ?", "1.0").Find(&checkList)
			for _, checkGroup := range checkList {
				cqp.SendGroupMsg(checkGroup.GroupID, msg)
			}
		}

		// B站明日方舟
		msg, ok = updateCheckt(161775300)
		if ok {
			// 查询数据库
			db.Table("group_info").Select("group_id").Where("Arknights = ?", "1.0").Find(&checkList)
			for _, checkGroup := range checkList {
				cqp.SendGroupMsg(checkGroup.GroupID, msg)
			}
		}

		// B站FGO
		msg, ok = updateCheckt(233108841)
		if ok {
			// 查询数据库
			db.Table("group_info").Select("group_id").Where("FateGrandOrder = ?", "1.0").Find(&checkList)
			for _, checkGroup := range checkList {
				cqp.SendGroupMsg(checkGroup.GroupID, msg)
			}
		}
		db.Close()
		calibrationCountdown++      // 计数增加
		time.Sleep(6 * time.Minute) // 六分钟后继续
		// cqp.AddLog(0, "test", msg)
		// time.Sleep(2 * time.Minute) // 六分钟后继续
	}
}

func updateCheckt(id int) (string, bool) {

	// return "推送测试", true
	timestamp, msg, img := getDynamic(id, 0, true)
	switch timestamp {
	case -1:
		return "", false
	case 0:
		return msg, true
	default:
		// 链接数据库
		db, err := gorm.Open("sqlite3", Datedir)
		if err != nil {
			cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
			return "", false
		}
		db.Table("Crawler_update_time").Where("update_url = ?", id).Update("update_time", timestamp)
		// cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v\t%v\t%v", err, timestamp, "改写执行后"))
		var newmsg string
		if img != nil {
			if cqp.CanSendImage() {
				newmsg = msg + "\n附图：\n[CQ:image,file=" + strings.Join(img, "]\n[CQ:image,file=") + "]"
			} else {
				newmsg = msg + "\n附图：\n" + strings.Join(img, "\n")
			}
			return newmsg, true
		}
		return msg, true
	}
}
