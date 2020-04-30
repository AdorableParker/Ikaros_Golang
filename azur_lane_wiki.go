package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

func equipmentRanking(group, qq int64) {
	imgURL := getWikiImg("装备一图榜")
	if imgURL != "" {
		sendPhoto(group, qq, imgURL)
	} else {
		sendMsg(group, qq, "访问Wiki失败惹\nε(┬┬﹏┬┬)3")
	}
}

func srengthRanking(group, qq int64) {
	imgURL := getWikiImg("PVE用舰船综合性能强度榜")
	if imgURL != "" {
		sendPhoto(group, qq, imgURL)
	} else {
		sendMsg(group, qq, "访问Wiki失败惹\nε(┬┬﹏┬┬)3")
	}
}

func pixivRanking(group, qq int64) {
	imgURL := getWikiImg("P站搜索结果一览榜（社保榜）")
	if imgURL != "" {
		sendPhoto(group, qq, imgURL)
	} else {
		sendMsg(group, qq, "访问Wiki失败惹\nε(┬┬﹏┬┬)3")
	}
}

func getWikiImg(index string) string {

	// 请求html页面
	res, err := http.Get(fmt.Sprintf("https://wiki.biligame.com/blhx/%s", index))
	if err != nil {
		// 错误处理
		cqp.AddLog(30, "HTTP错误", fmt.Sprintf("错误信息:%v", err))
		return ""
	}
	if res.StatusCode != 200 {
		cqp.AddLog(30, "HTTP错误", fmt.Sprintf("StatusCode:%d", res.StatusCode))
		return ""
	}
	// 载入HTML文件
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		cqp.AddLog(30, "HTTP错误", fmt.Sprintf("错误信息:%v", err))
		return ""
	}
	// 查找图片
	soup := doc.Find("#mw-content-text")
	s, ok := soup.Find("img").Attr("src")
	if !ok {
		cqp.AddLog(30, "搜索错误", "未能找到目标")
		return ""
	}
	return s
}

func sendPhoto(group, qq int64, imgURL string) {
	if cqp.CanSendImage() {
		sendMsg(group, qq, fmt.Sprintf("[CQ:image,file=%s]", imgURL))
	} else {
		s := fmt.Sprintf("[CQ:at,qq=%d]\n图片链接:%s", qq, imgURL)
		sendMsg(group, qq, s)
	}
}
