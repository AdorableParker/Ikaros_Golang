package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/jozsefsallai/gophersauce"
)

// Client SauceNAO 搜图引擎客户端
var Client *gophersauce.Client

// ,
// 	"#以图搜图": `------img_saucenao------
// 命令关键字: "saucenao", "图片搜索", "搜图"
// 命令输入格式:

// 图片搜索<空格><图片>
// 例:
// 搜图 [图片]

// 效果: 根据输入的图片，使用saucenao搜图引擎进行图源搜索，返回相似度最高的。
// ######################`

// DocSaucenao 以图搜图功能文档
var DocSaucenao = &HelpDoc{
	Name:        "以图搜图",
	KeyWord:     []string{"saucenao", "图片搜索", "搜图"},
	Example:     "搜图 [图片]",
	Description: "图片搜索<空格><图片>\n根据输入的图片,使用saucenao搜图引擎进行图源搜索,返回相似度最高的,部分被隐蔽的结果将会被抛弃\n支持的平台有:\nPixiv、IMDb、DeviantArt、Bcy、AniDBA、Pawoo、Seiga、Sankaku、Danbooru"}

func saucenao(msg []string, msgID int32, group, qq int64, try uint8) {

	if len(msg) == 0 { // 如果没有获取到参数
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "要搜索的图片是哪张呢\nq(≧▽≦q)")
			} else {
				sendMsg(group, qq, "没有收到图片哦,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, saucenao, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}
	re := regexp.MustCompile(`\[CQ:image,file=.*?\]`)
	cqcode := re.FindAllString(msg[0], 1) // 正则匹配查找图片CQ码
	if cqcode == nil {                    // 没有找到
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			sendMsg(group, qq, "没有收到图片哦,再发一次吧\n(。・ω・。)")                               // 发送提示消息
			stagedSessionPool[msgID] = newStagedSession(group, qq, saucenao, msg, try) // 添加新的会话到会话池
			// cqp.AddLog(0,"t",fmt.Sprintln(stagedSessionPool))
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}
	cq := cqcode[0]
	fileDir := cqp.GetImage(cq[15 : len(cq)-1]) // 获取图片

	if !SauceNAO { // 如果初始化阶段出错了
		var err error
		Client, err = gophersauce.NewClient(&gophersauce.Settings{
			APIUrl:     `https://saucenao.com/search.php`,
			APIKey:     AdminConfig.SaucenaoKey,
			MaxResults: 1})
		if err != nil {
			cqp.AddLog(20, "初始化异常", fmt.Sprintln("搜图引擎初始化出现错误\n\v", err))
			sendMsg(group, qq, "搜图引擎初始化异常,请联系维护 ≧ ﹏ ≦")
			return
		}
		SauceNAO = true
	}
	sendMsg(group, qq, "引擎初始化完成, 开始搜图")
	response, err := Client.FromFile(fileDir)
	defer os.Remove(fileDir)
	if err != nil {
		if strings.Contains(err.Error(), "Daily Search Limit Exceeded.") {
			cqp.AddLog(20, "搜图异常", "查询次数超限")
			sendMsg(group, qq, "API超出200次/日,请改日重试 ≧ ﹏ ≦")
			return
		}
		if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
			cqp.AddLog(20, "搜图异常", "连接被远程主机强制关闭")
			sendMsg(group, qq, "An existing connection was forcibly closed by the remote host. ≧ ﹏ ≦")
			return
		}
		if strings.Contains(err.Error(), "A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.") {
			cqp.AddLog(20, "搜图异常", "服务器正常未响应")
			sendMsg(group, qq, "A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond. ≧ ﹏ ≦")
			return
		}
		cqp.AddLog(20, "搜图异常", fmt.Sprintln(err))
		sendMsg(group, qq, "搜图引擎运行异常,请稍后重试 ≧ ﹏ ≦")
		return
	}
	first := response.First()
	title := first.Data.Title // 标题
	if title == "" {
		sendMsg(group, qq, "没有找到符合标准的结果 ≧ ﹏ ≦")
		return
	}
	similarity := first.Header.Similarity            // 相似度
	thumbnail := first.Header.Thumbnail              // 缩略信息
	RemainingQuotas := response.Header.LongRemaining // 剩余配额
	var text = "%s\n结果来自于 %s\n%s ID:\t%d\n相似度:\t%s%%\n缩略信息:%s\n24小时内剩余可查询次数: %d"
	switch {
	case first.IsPixiv():
		pixivID := first.Data.PixivID
		text = fmt.Sprintf(text, title, "Pixiv", "Pixiv", pixivID, similarity, thumbnail, RemainingQuotas)

	case first.IsIMDb():
		imdbID := first.Data.IMDbID
		text = fmt.Sprintf(text, title, "IMDb", "IMDb", imdbID, similarity, thumbnail, RemainingQuotas)

	case first.IsDeviantArt():
		deviantartID := first.Data.DeviantArtID
		text = fmt.Sprintf(text, title, "DeviantArt", "DeviantArt", deviantartID, similarity, thumbnail, RemainingQuotas)

	case first.IsBcy():
		bcyID := first.Data.BcyID
		text = fmt.Sprintf(text, title, "Bcy", "Bcy", bcyID, similarity, thumbnail, RemainingQuotas)

	case first.IsAniDB():
		anidbaID := first.Data.AniDBAID
		text = fmt.Sprintf(text, title, "AniDBA", "AniDBA", anidbaID, similarity, thumbnail, RemainingQuotas)

	case first.IsPawoo():
		pawooID := first.Data.PawooID
		text = fmt.Sprintf(text, title, "Pawoo", "Pawoo", pawooID, similarity, thumbnail, RemainingQuotas)

	case first.IsSeiga():
		seigaID := first.Data.SeigaID
		text = fmt.Sprintf(text, title, "Seiga", "Seiga", seigaID, similarity, thumbnail, RemainingQuotas)

	case first.IsSankaku():
		sankakuID := first.Data.SankakuID
		text = fmt.Sprintf(text, title, "Sankaku", "Sankaku", sankakuID, similarity, thumbnail, RemainingQuotas)

	case first.IsDanbooru():
		danbooruID := first.Data.DanbooruID
		text = fmt.Sprintf(text, title, "Danbooru", "Danbooru", danbooruID, similarity, thumbnail, RemainingQuotas)
	default:
		// text = fmt.Sprintf("%+v\n", first)
		text = "服务器返回的是被隐藏的低相似度结果\n(ノω<。)ノ"
	}
	sendMsg(group, qq, text)
}
