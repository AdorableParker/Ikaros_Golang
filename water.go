package main

import (
	"fmt"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

func water(fromGroup int64) {
	cqp.SendGroupMsg(fromGroup, fmt.Sprintf("https://qqweb.qq.com/m/qun/activedata/active.html?gc=%d", fromGroup))
}
