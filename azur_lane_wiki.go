package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

// DocEquipmentRanking 碧蓝航线Wiki装备榜单查询功能文档
var DocEquipmentRanking = &HelpDoc{
	Name:        "碧蓝航线Wiki装备榜单查询",
	KeyWord:     []string{"装备榜单", "装备榜", "装备排行榜"},
	Description: "爬取碧蓝航线Wiki装备强度测评榜",
}

func equipmentRanking(group, qq int64) {
	imgURL := getWikiImg("装备一图榜", 0)
	if imgURL != "" {
		sendPhoto(group, qq, imgURL)
	} else {
		sendMsg(group, qq, "访问Wiki失败惹\nε(┬┬﹏┬┬)3")
	}
}

// DocSrengthRanking 碧蓝航线Wiki强度榜查询功能文档
var DocSrengthRanking = &HelpDoc{
	Name: "碧蓝航线Wiki PVE用舰船综合性能强度榜查询",
	KeyWord: []string{
		"强度榜单", "强度榜", "舰娘强度榜", "舰娘排行榜",
		"强度副榜", "舰娘强度副榜", "舰娘排行副榜"},
	Description: "爬取碧蓝航线Wiki PVE用舰船综合性能强度榜",
}

func srengthRanking(group, qq int64) {
	imgURL := getWikiImg("PVE用舰船综合性能强度榜", 1)
	if imgURL != "" {
		sendPhoto(group, qq, imgURL)
	} else {
		sendMsg(group, qq, "访问Wiki失败惹\nε(┬┬﹏┬┬)3")
	}
}

func srengthRankingEXC(group, qq int64) {
	imgURL := getWikiImg("PVE用舰船综合性能强度榜", 2)
	if imgURL != "" {
		sendPhoto(group, qq, imgURL)
	} else {
		sendMsg(group, qq, "访问Wiki失败惹\nε(┬┬﹏┬┬)3")
	}
}

// DocPixivRanking 碧蓝航线Wiki P站搜索结果一览榜查询功能文档
var DocPixivRanking = &HelpDoc{
	Name:        "碧蓝航线Wiki P站搜索结果一览榜查询",
	KeyWord:     []string{"社保榜", "射爆榜", "P站榜", "p站榜"},
	Description: "爬取碧蓝航线Wiki P站搜索结果一览榜",
}

func pixivRanking(group, qq int64) {
	imgURL := getWikiImg("P站搜索结果一览榜（社保榜）", 0)
	if imgURL != "" {
		sendPhoto(group, qq, imgURL)
	} else {
		sendMsg(group, qq, "访问Wiki失败惹\nε(┬┬﹏┬┬)3")
	}
}

func getWikiImg(index string, exception int) string {

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
	if index == "PVE用舰船综合性能强度榜" {
		s := soup.Find("img").Map(findImg)
		return s[exception]
	}
	s, ok := soup.Find("img").Attr("src")
	if !ok {
		cqp.AddLog(30, "搜索错误", "未能找到目标")
		return ""
	}
	return s
}
func findImg(i int, j *goquery.Selection) string {
	s, ok := j.Attr("src")
	if !ok {
		cqp.AddLog(30, "搜索错误", "未能找到目标")
		return ""
	}
	return s
}

func sendPhoto(group, qq int64, imgURL string) {
	// if cqp.CanSendImage() {
	sendMsg(group, qq, fmt.Sprintf("[CQ:image,url=%s]", imgURL))

	// } else {
	// s := fmt.Sprintf("[CQ:at,qq=%d]\n图片链接:%s", qq, imgURL)
	// sendMsg(group, qq, s)
	// }
}
