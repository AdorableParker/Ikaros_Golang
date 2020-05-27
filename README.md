# Ikaros_Golang #

[![GitHub issues](https://img.shields.io/github/issues/AdorableParker/Ikaros_Golang?style=plastic)](https://github.com/AdorableParker/Ikaros_Golang/issues)
[![GitHub stars](https://img.shields.io/github/stars/AdorableParker/Ikaros_Golang)](https://github.com/AdorableParker/Ikaros_Golang/stargazers)
[![GitHub license](https://img.shields.io/github/license/AdorableParker/Ikaros_Golang?style=plastic)](https://github.com/AdorableParker/Ikaros_Golang/blob/master/LICENSE)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/AdorableParker/Ikaros_Golang?style=plastic)
![GitHub last commit](https://img.shields.io/github/last-commit/AdorableParker/Ikaros_Golang?style=plastic)

自用QQ机器人伊卡洛斯的Golang版本

因为Python太慢，所以从Python版本迁移过来

这个语言版本不知道我能坚持维护多久，说起来Python前前后后也更新了2年了...

~~另外，因为gojieba库的编译报错问题，也不知道是我做错了什么，或许哪天实在忍不了了会转C++~~

## 如何使用 ##
如果只是想直接用的话:

1.  下载`app.dll`和`app.json`文件
2.  打开 Cool Q 的开发者模式
3.  把`app.dll`和`app.json`放到`..\dev\io.github.adorableparker.Ikaros`目录下
4.  下载`nfiguration_data`文件夹里面的文件
5.  放置到`..\data\app\io.github.adorableparker.Ikaros`目录下

---

如果是改写了代码想自己编译:

1.  运行`build.bat`来进行编译生成`app.dll`和`app.json`

2.  然后从上面的第2步开始操作，就可以了

    >   `app.json`应该会要根据你的修改稍微改改内容，具体的修改请参考SDK文档
    >
    >   `https://pkg.go.dev/github.com/Tnze/CoolQ-Golang-SDK/cqp?tab=doc`



## 目前功能有 ##

* 基于saucenao的搜图
* 自动检测bilibili的动态更新
    > 只写了`碧蓝航线`、`火星bot小黄瓜`和`明日方舟`的
* 主动查询bilibili的动态更新
    > 一样只写了`碧蓝航线`、`火星bot小黄瓜`和`明日方舟`，但是模板套入B站的UID就可以增加对象了
* 碧蓝航线的建造时间查询
    > 由于游戏更新，需要时不时的往数据库添加数据
* 碧蓝航线的打捞地点查询
    > 由于游戏更新，一样需要时不时的往数据库添加数据
* 碧蓝航线的活动计算器
* 智障聊天
    > ~~我仍不知道gojieba报错的原因~~问题大概解决了，是该死的初始化时长太久的问题
* 自助禁言
    > 解放管理员的双手，让闲的蛋疼的群员自己解决自闭问题
* 碧蓝航线的Wiki榜单图片爬取
    > 只爬取了`装备排行榜`、`P站榜`和`PVE节奏版`
* 整点报时
    > 花了些功夫去搞舰C的台词呢
* 游戏每日提醒
    > 这个开了就知道是什么了
* 船名查询
    > 就是和谐名对照表了

### 还有很多功能没有迁移，不想写... ###
~~都2020年了，你是个成熟的IDE了，该学会自己把项目写好了.jpg~~
