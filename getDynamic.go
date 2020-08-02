package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type crawlerUpdateInfo struct {
	UpdateTime int64 `gorm:"column:update_time"`
}

// DocDynamicByID B站动态查询功能文档
var DocDynamicByID = &HelpDoc{
	Name:        "B站动态查询",
	KeyWord:     []string{"B站动态"},
	Example:     "B站动态 114514",
	Description: "B站动态<空格><UID>\n用于查询指定UID的用户的B站动态"}

func dynamicByID(msg []string, msgID int32, group, qq int64, try uint8) {
	if len(msg) == 0 { // 如果没有获取到参数
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入up的UID\nq(≧▽≦q)")
			} else {
				sendMsg(group, qq, "不能为空哦,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, dynamicByID, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}

	uid, err := strconv.ParseUint(msg[0], 10, 0)
	if err != nil {
		sendMsg(group, qq, "不正常的UID,你怕不是打错了")
		return
	}
	sendDynamic([]string{}, group, qq, int(uid))
}

// DocSendDynamic B站动态快捷查询功能文档
var DocSendDynamic = &HelpDoc{
	Name: "B站动态快捷查询",
	KeyWord: []string{
		"小加加", "火星加", "B博更新", "b博更新",
		"转推姬", "碧蓝日推",
		"罗德岛线报", "方舟公告", "方舟B博", "阿米娅",
		"月球人公告", "FGO公告", "呆毛王"},
	Example:     "小加加 1\n碧蓝日推 5\n方舟功告\nFGO公告",
	Description: "命令空格后加数字可回溯指定条数的历史动态"}

func sendDynamic(tomsg []string, group, qq int64, id int) {
	var pages int = 0
	var err error
	if len(tomsg) != 0 {
		pages, err = strconv.Atoi(tomsg[0])
		if err != nil {
			sendMsg(group, qq, "请不要输入一些奇奇怪怪的东西\n＞︿＜")
			return
		}
		switch {
		case pages >= 10:
			pages = 9
			sendMsg(group, qq, "最多只能往前10条哦\n的(￣﹃￣)")
		case pages < 0:
			pages = 0
			sendMsg(group, qq, "未来的事情我怎么会知道\n=￣ω￣=")
		}
	}
	timeStamp, msg, img := getDynamic(id, pages, false)
	var endResult string
	var newmsg string
	if img != nil {
		// cqp.AddLog(0, "测试文本", fmt.Sprintf("图片列表:%v", img))
		if cqp.CanSendImage() {
			newmsg = msg + "\n附图：\n[CQ:image,file=" + strings.Join(img, "]\n[CQ:image,file=") + "]"
		} else {
			newmsg = msg + "\n附图：\n" + strings.Join(img, "\n")
			// cqp.AddLog(0, "测试文本", fmt.Sprintf("输出:%v", newmsg))
		}
	} else {
		newmsg = msg
	}

	if timeStamp > 0 {
		timeString := time.Unix(timeStamp, 0).Format("2006-01-02 15:04:05")
		endResult = fmt.Sprintf("%s\n发布时间: %s", newmsg, timeString)
	} else {
		endResult = newmsg
	}
	sendMsg(group, qq, endResult)
}

func getDynamic(id, pages int, flag bool) (int64, string, []string) {

	apiurl := "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history?host_uid=%d"
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(apiurl, id), nil)
	if err != nil {
		cqp.AddLog(30, "HTTP错误", fmt.Sprintf("错误信息:%v", err))
		return -1, "", nil
	}
	resp, err := client.Do(req)
	if err != nil {
		cqp.AddLog(30, "HTTP错误", fmt.Sprintf("错误信息:%v", err))
		return -1, "", nil
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cqp.AddLog(30, "HTTP错误", fmt.Sprintf("错误信息:%v", err))
		return -1, "", nil
	}

	index := fmt.Sprintf("[%d]", pages)

	timestamp, _ := jsonparser.GetInt([]byte(bodyText), "data", "cards", index, "desc", "timestamp")
	if flag {
		// 链接数据库
		db, err := gorm.Open("sqlite3", Datedir)
		defer db.Close()
		if err != nil {
			cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v", err))
			return -1, "", nil
		}
		// 查询数据库
		var crawlerUpdate crawlerUpdateInfo
		db.Table("Crawler_update_time").Select("update_time").Where("update_url = ?", id).First(&crawlerUpdate)
		// cqp.AddLog(30, "数据库错误", fmt.Sprintf("错误信息:%v\t%v\t%v", err, timestamp, crawlerUpdate.UpdateTime))
		if timestamp <= crawlerUpdate.UpdateTime {
			return -1, "", nil
		}
	}

	dynamicType, err := jsonparser.GetInt([]byte(bodyText), "data", "cards", index, "desc", "type")

	// fmt.Printf("Tpye:%T,value:%v,err:%v", dynamicType, dynamicType, err)
	cardold, _ := jsonparser.GetString([]byte(bodyText), "data", "cards", index, "card")
	card := []byte(cardold)
	switch dynamicType {
	case 0: // 无效数据
		return 0, "没有相关动态信息", nil
	case 1: // 转发
		content, _ := jsonparser.GetString(card, "item", "content") // 评论转发
		return timestamp, fmt.Sprintf("转发并评论：%s", content), nil
	case 2: // 含图动态
		description, _ := jsonparser.GetString(card, "item", "description")   // 描述
		picturesCount, _ := jsonparser.GetInt(card, "item", "pictures_count") // 图片数
		imgSrc := make([]string, 0)                                           // 容器
		for i := 0; i < int(picturesCount); i++ {
			url, _ := jsonparser.GetString(card, "item", "pictures", fmt.Sprintf("[%d]", i), "img_src") // 图片地址
			imgSrc = append(imgSrc, url)
		}
		return timestamp, fmt.Sprintf("更新动态：%s", description), imgSrc
	case 4: //无图动态
		content, _ := jsonparser.GetString(card, "item", "content") // 内容
		return timestamp, fmt.Sprintf("更新动态：%s", content), nil
	case 8: // 视频
		dynamic, _ := jsonparser.GetString(card, "dynamic") // 描述
		imgSrc, _ := jsonparser.GetString(card, "pic")      //封面图片
		return timestamp, dynamic, []string{imgSrc}
	case 64: // 专栏
		title, _ := jsonparser.GetString(card, "title")       // 标题
		summary, _ := jsonparser.GetString(card, "summary")   // 摘要
		imgSrc, _ := jsonparser.GetString(card, "banner_url") // 封面图片
		return timestamp, fmt.Sprintf("专栏标题:%s\n专栏摘要：\n%s…", title, summary), []string{imgSrc}
	case 2048:
		title, _ := jsonparser.GetString(card, "sketch", "title")          // 标题
		context, _ := jsonparser.GetString(card, "vest", "content")        // 内容
		targetURL, _ := jsonparser.GetString(card, "sketch", "target_url") // 相关链接
		return timestamp, fmt.Sprintf("动态标题:%s\n动态内容：\n%s\n相关链接:\n%s", title, context, targetURL), nil
	default:
		cqp.AddLog(30, "JSON错误", fmt.Sprintf("错误信息:未知的类型码 %v ", dynamicType))
		return 0, "是未知的动态类型,无法解析", nil
	}
}
