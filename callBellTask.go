package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type updateCallInfo struct {
	GroupID int64 `gorm:"column:group_id"`
}

func callBellTask() {
	var calibrationCountdown uint8 = 3 // 首次运行 进行计时器校准
	var flag bool = false              // 逾期旗帜
	for {
		if calibrationCountdown >= 3 { // 如果 校准倒数大于 3 执行校准操作

			calibrationCountdown = 0                                  // 进入校准部分 计数归零
			nextTime := time.Now().Truncate(time.Hour).Add(time.Hour) // 返回下一个整点时间
			jetLag := nextTime.Sub(time.Now())                        // 时差

			switch {
			// 时差 大于或等于 59分钟 说明逾期1分钟以内 || 时差 小于或等于 0 说明时间刚过去
			case jetLag >= 59*time.Minute || jetLag <= 0:
				cqp.AddLog(0, "整点报时", fmt.Sprintln("计划下次执行时间:", nextTime))

			// 时差 大于或等于 57分钟 说明逾期3分钟以内
			case jetLag >= 57*time.Minute:
				flag = true
				cqp.AddLog(0, "整点报时逾期", fmt.Sprintln("需要校准,计划下次执行时间:", nextTime))

			// 其他情况 ∈ 时差 在 57~0 之间
			default:
				cqp.AddLog(0, "整点报时", fmt.Sprintln("等待时长:", jetLag, "计划下次执行时间:", nextTime))
				time.Sleep(jetLag) // 休眠直到该时间
			}
		}

		// 链接数据库
		db, err := gorm.Open("sqlite3", Datedir)
		if err != nil {
			cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
			return
		}
		nowTime := time.Now().Hour()
		var checkList []updateCallInfo
		// 查询数据库
		db.Table("group_info").Select("group_id").Where("Call_bell = ?", "1.0").Find(&checkList)
		for _, checkGroup := range checkList {
			if flag {
				cqp.SendGroupMsg(checkGroup.GroupID, fmt.Sprintf("啊，迟到了，现在是%d点", nowTime))
			} else {
				cqp.SendGroupMsg(checkGroup.GroupID, fmt.Sprintf("现在%d点咯", nowTime))
			}
		}

		msg := getScript(nowTime)
		if msg == "" {
			msg = fmt.Sprintf("现在%d点咯，好啦是我忘词啦，你好烦欸╰（‵□′）╯", nowTime)
		}
		// 查询数据库
		db.Table("group_info").Select("group_id").Where("Call_bell_AZ = ?", "1.0").Find(&checkList)
		for _, checkGroup := range checkList {
			if flag {
				cqp.SendGroupMsg(checkGroup.GroupID, "啊，迟到了")
			}
			cqp.SendGroupMsg(checkGroup.GroupID, msg)
		}

		db.Close()
		calibrationCountdown++ // 计数增加
		if flag {
			flag = false
			time.Sleep(time.Now().Truncate(time.Hour).Add(time.Hour).Sub(time.Now())) // 虽然看起来很长,但是功能就是等待到下一个整点而已
		} else {
			time.Sleep(time.Hour) // 一小时后继续
		}
		// cqp.AddLog(0, "test", msg)
		// time.Sleep(2 * time.Minute) // 六分钟后继续
	}
}

func getScript(h int) string {
	files, err := ioutil.ReadDir(filepath.Join(Appdir, "time_txt"))
	if err != nil {
		cqp.AddLog(30, "文件列表读取错误", fmt.Sprintln(err))
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	fileName := files[rand.Intn(len(files))].Name()

	f, err := os.Open(filepath.Join(Appdir, "time_txt", fileName))
	defer f.Close()
	if err != nil {
		cqp.AddLog(30, "台词文件读取错误", fmt.Sprintln(err))
		return ""
	}
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		cqp.AddLog(30, "文件流读取错误", fmt.Sprintln(err))
		return ""
	}
	script := strings.Split(string(fd), "\n")
	return script[h] + "                  ————" + fileName
}
