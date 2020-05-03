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
	for {
		if calibrationCountdown >= 3 { // 如果 校准倒数大于 120 执行校准操作
			calibrationCountdown = 0         // 进入校准部分 计数归零
			t := time.Now().Round(time.Hour) // 返回最近的整小时数
			if t.Before(time.Now()) {        // 如果时间晚于当前时间
				t = t.Add(time.Hour) // 往后推一小时
			}
			cqp.AddLog(0, "整点报时时间", fmt.Sprintln(t))
			time.Sleep(t.Sub(time.Now())) // 休眠直到该时间
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
			cqp.SendGroupMsg(checkGroup.GroupID, fmt.Sprintf("现在%d点咯", nowTime))
		}

		msg := getScript(nowTime)
		if msg == "" {
			msg = fmt.Sprintf("现在%d点咯，好啦是我忘词啦，你好烦欸╰（‵□′）╯", nowTime)
		}
		// 查询数据库
		db.Table("group_info").Select("group_id").Where("Call_bell_AZ = ?", "1.0").Find(&checkList)
		for _, checkGroup := range checkList {
			cqp.SendGroupMsg(checkGroup.GroupID, msg)
		}

		db.Close()
		calibrationCountdown++ // 计数增加
		time.Sleep(time.Hour)  // 一小时后继续
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
	return script[h] + "        ————" + fileName
}
