package main

import (
	"fmt"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

// DocWater 群活跃数据分享功能文档
var DocWater = &HelpDoc{
	Name:        "群活跃数据分享",
	KeyWord:     []string{"群活跃数据"},
	Description: "返回当前群的群活跃数据统计信息"}

func water(fromGroup int64) {
	cqp.SendGroupMsg(fromGroup, fmt.Sprintf("https://qqweb.qq.com/m/qun/activedata/active.html?gc=%d", fromGroup))
}
