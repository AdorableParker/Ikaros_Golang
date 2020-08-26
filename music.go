package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/buger/jsonparser"
)

func getSongID(songName string, source int) (songID int64) {
	client := &http.Client{}
	var req *http.Request
	var err error
	switch source {
	case 0:
		req, err = http.NewRequest("GET", fmt.Sprintf("http://music.163.com/api/search/get/web?type=1&s=%s", url.QueryEscape(songName)), nil)
	case 1:
		req, err = http.NewRequest("GET", fmt.Sprintf("https://c.y.qq.com/soso/fcgi-bin/client_search_cp?format=json&w=%s", url.QueryEscape(songName)), nil)
	}

	if err != nil {
		cqp.AddLog(30, "HTTP错误", fmt.Sprintf("错误信息:%v", err))
		// fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		cqp.AddLog(30, "HTTP错误", fmt.Sprintf("错误信息:%v", err))
		// fmt.Println(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cqp.AddLog(30, "HTTP错误", fmt.Sprintf("错误信息:%v", err))
		// fmt.Println(err)
	}
	switch source {
	case 0:
		songID, _ = jsonparser.GetInt([]byte(bodyText), "result", "songs", "[0]", "id")
	case 1:
		songID, _ = jsonparser.GetInt([]byte(bodyText), "data", "song", "list", "[0]", "songid")
	}
	return songID
}

// DocMusic 点歌功能文档
var DocMusic = &HelpDoc{
	Name:        "点歌姬",
	KeyWord:     []string{"点歌", "点歌姬"},
	Example:     "点歌姬 ハートの确率",
	Description: "点歌姬<空格><歌曲名>"}

func music(msg []string, msgID int32, group, qq int64, try uint8) {
	if len(msg) == 0 { // 如果没有获取到参数
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入歌曲名\nq(≧▽≦q)")
			} else {
				sendMsg(group, qq, "歌曲名不能为空哦,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, music, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}
	songID := getSongID(strings.Join(msg, " "), 1) // 网易云的分享失效
	// sendMsg(group, qq, "[CQ:music,id=4931944,type=163]")
	// sendMsg(group, qq, fmt.Sprintf("%d", songID))
	sendMsg(group, qq, fmt.Sprintf("[CQ:music,id=%d,type=qq]", songID))
}
