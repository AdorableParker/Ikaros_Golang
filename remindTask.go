package main

import (
	"fmt"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type remindTaskInfo struct {
	GroupID int64 `gorm:"column:group_id"`
}

//remindTask 每天的晚九点半 提醒任务
func remindTask() {
	var calibrationCountdown uint8 = 10 // 首次运行 需要执行计时器校准
	for {
		if calibrationCountdown >= 10 { // 如果 校准倒数大于 10 执行校准操作
			calibrationCountdown = 0                                                        // 进入校准部分 计数归零
			t := time.Now().Round(24 * time.Hour).Add(13 * time.Hour).Add(30 * time.Minute) // 返回距离今天 或 明天 的的晚九点半 时间
			if t.Before(time.Now()) {                                                       // 如果时间晚于当前时间
				t = t.Add(24 * time.Hour) // 往后推24小时
			}
			cqp.AddLog(0, "提醒任务执行时间", fmt.Sprintln(t))
			time.Sleep(t.Sub(time.Now())) // 休眠直到该时间
		}
		azle, rgo := remindtext(int(time.Now().Weekday()))

		// 链接数据库
		db, err := gorm.Open("sqlite3", Datadir)
		if err != nil {
			cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
			return
		}

		// 每日提醒
		var remindTasks []remindTaskInfo

		// AzLe 版本
		// 查询数据库
		db.Table("group_info").Select("group_id").Where("Daily_remind_AzurLane = ?", "1.0").Find(&remindTasks)
		for _, grouInfo := range remindTasks {
			cqp.SendGroupMsg(grouInfo.GroupID, azle)
		}
		// FGO 版本
		// 查询数据库
		db.Table("group_info").Select("group_id").Where("Daily_remind_FGO = ?", "1.0").Find(&remindTasks)
		for _, grouInfo := range remindTasks {
			cqp.SendGroupMsg(grouInfo.GroupID, rgo)
		}

		db.Close()                 // 关闭数据库
		calibrationCountdown++     // 计数增加
		time.Sleep(24 * time.Hour) // 二十四小时后继续
	}
}

func remindtext(w int) (string, string) {
	var weekday = [7]string{"周一", "周二", "周三", "周四", "周五", "周六"}
	var AzurLane = [4]string{
		"今天是%s哦,今天开放的是「战术研修」「斩首行动」，困难也记得打呢。各位指挥官晚安咯\nο(=•ω＜=)ρ⌒☆",
		"今天是%s哦,今天开放的是「战术研修」「商船护送」，困难也记得打呢。各位指挥官晚安咯\nο(=•ω＜=)ρ⌒☆",
		"今天是%s哦,今天开放的是「战术研修」「海域突进」，困难也记得打呢。各位指挥官晚安咯\nο(=•ω＜=)ρ⌒☆",
		"今天是周日哦,每日全部模式开放，每周两次的破交作战记得打哦，困难模式也别忘了。各位指挥官晚安咯\nο(=•ω＜=)ρ⌒☆"}
	var FGO = [7]string{
		"晚上好,Master,今天是周日, 今天周回本开放「剑阶修炼场」,「收集火种(All)」。\nο(=•ω＜=)ρ⌒☆",
		"晚上好,Master,今天是周一, 今天周回本开放「弓阶修炼场」,「收集火种(枪杀)」。\nο(=•ω＜=)ρ⌒☆",
		"晚上好,Master,今天是周二, 今天周回本开放「枪阶修炼场」,「收集火种(剑骑)」。\nο(=•ω＜=)ρ⌒☆",
		"晚上好,Master,今天是周三, 今天周回本开放「狂阶修炼场」,「收集火种(弓术)」。\nο(=•ω＜=)ρ⌒☆",
		"晚上好,Master,今天是周四, 今天周回本开放「骑阶修炼场」,「收集火种(枪杀)」。\nο(=•ω＜=)ρ⌒☆",
		"晚上好,Master,今天是周五, 今天周回本开放「术阶修炼场」,「收集火种(剑骑)」。\nο(=•ω＜=)ρ⌒☆",
		"晚上好,Master,今天是周六, 今天周回本开放「杀阶修炼场」,「收集火种(弓术)」。\nο(=•ω＜=)ρ⌒☆"}

	if w == 0 {
		return AzurLane[3], FGO[6]
	}
	return fmt.Sprintf(AzurLane[w%3], weekday[w-1]), FGO[w]

}
