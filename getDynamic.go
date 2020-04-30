package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type crawlerUpdateInfo struct {
	UpdateTime int64 `gorm:"column:update_time"`
}

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
	_, msg, img := getDynamic(id, pages, false)
	if img != nil {
		var newmsg string
		// cqp.AddLog(0, "测试文本", fmt.Sprintf("图片列表:%v", img))
		if cqp.CanSendImage() {
			newmsg = msg + "\n附图：\n[CQ:image,file=" + strings.Join(img, "]\n[CQ:image,file=") + "]"
		} else {
			newmsg = msg + "\n附图：\n" + strings.Join(img, "\n")
			// cqp.AddLog(0, "测试文本", fmt.Sprintf("输出:%v", newmsg))
		}
		sendMsg(group, qq, newmsg)
		return
	}
	sendMsg(group, qq, msg)
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
	default:
		cqp.AddLog(30, "JSON错误", fmt.Sprintf("错误信息:未知的类型码 %v ", dynamicType))
		return 0, "是未知的动态类型,无法解析", nil
	}
}