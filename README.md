# Ikaros_Golang #

[![GitHub issues](https://badgen.net/github/issues/AdorableParker/Ikaros_Golang)](https://github.com/AdorableParker/Ikaros_Golang/issues)
[![GitHub stars](https://badgen.net/github/stars/AdorableParker/Ikaros_Golang)](https://github.com/AdorableParker/Ikaros_Golang/stargazers)
[![GitHub latest tag](https://badgen.net/github/tag/AdorableParker/Ikaros_Golang)](https://github.com/AdorableParker/Ikaros_Golang/tags)
![GitHub last commit](https://badgen.net/github/last-commit/AdorableParker/Ikaros_Golang)


自用QQ机器人伊卡洛斯的Golang版本

因为Python太慢，所以从Python版本迁移过来

这个语言版本不知道我能坚持维护多久，说起来Python前前后后也更新了2年了...

~~另外，因为gojieba库的编译报错问题，也不知道是我做错了什么，或许哪天实在忍不了了会转C++~~

## 如何使用 ##
### 使用者请看这里: ###
1. 在[![release](https://badgen.net/github/release/AdorableParker/Ikaros_Golang/stable)](https://github.com/AdorableParker/Ikaros_Golang/releases)下载最新稳定版本
2. 将下载的放入`.cpk`文件放入酷Q的应用目录
3. 载入并启用插件(报错属于正常情况)
4. 下载`https://github.com/AdorableParker/Ikaros_Golang/tree/master/Configuration_data`里面的文件并将其放置到`..\data\app\com.adorableparker.github.ikaros_golang`目录下
5. 根据自己情况修改`MainConf.ini`的值并保存

此时已经可以正常使用

---
### 二次开发请看这里: ###
0.  二次开发代码
1.  运行`build.bat`来进行编译生成`app.dll`和`app.json`
2.  下载`app.dll`和`app.json`文件
3.  打开 Cool Q 的开发者模式
4.  把`app.dll`和`app.json`放到`..\dev\com.adorableparker.github.ikaros_golang`目录下
    >   `app.json`应该会要根据你的修改稍微改改内容，具体的修改请参考SDK文档
    >   `https://pkg.go.dev/github.com/Tnze/CoolQ-Golang-SDK/cqp?tab=doc`
5.  下载`Configuration_data`文件夹里面的文件
6.  放置到`..\data\app\com.adorableparker.github.ikaros_golang`目录下

## 使用须知 ##
* 因为代码太渣,退出后会有遗留的野线程需要手动结束
* 因为一些数据库的数据并不能自动更新,所以请手动更新,下面是相关项目说明:
    * `activityConfig.ini`
        > 活动进度计算器的相关数据
    * `time_txt`文件夹
        > 里面的无后缀文件皆为台词数据,记事本直接打开即可，可根据已有文件的格式自行添加台词内容
    * `User.db`
        > sqlite数据库,主要数据库,其中三个表有更新需求：
        > * `AzurLane_construct_time` 碧蓝航线舰船建造时间表
        > * `Roster` 和谐名对照表
        > * `ship_map` 舰船打捞地点表  
        >
        > 特别说明:  
        > `ship_map`中,使用的是四位二进制数来表示该章节相关海域打捞情况  
        > *假设：*  
        > *1、4海域不可打捞,2、3海域可以打捞*  
        > *则,四位二进制表示为 0110 转化为十进制数则为 6*  
        > *即,在数据表中储存为 6*


## 目前功能有 ##

* 基于saucenao的搜图
* 主动查询bilibili的动态更新
    > 通过封装好的快捷查询命令或者输入指定UID来查询
* 自动检测bilibili的动态更新并推送
    > 只写了`碧蓝航线`、`火星bot小黄瓜`和`明日方舟`的
    >> *`V.d647584`* 追加了`命运-冠位指定`
* 碧蓝航线的建造时间查询
    > 由于游戏更新，需要时不时的往数据库添加数据
* 碧蓝航线的打捞地点查询
    > 由于游戏更新，一样需要时不时的往数据库添加数据
* 碧蓝航线的活动计算器
* 碧蓝航线的舰船经验计算器
* 智障聊天
    > ~~我仍不知道gojieba报错的原因~~问题大概解决了，是该死的初始化时长太久的问题
* 自助禁言
    > 解放管理员的双手，让闲的蛋疼的群员自己解决自闭问题
* 碧蓝航线的Wiki榜单图片爬取
    > 只爬取了`装备排行榜`、`P站榜`和`PVE节奏版`
* 整点报时
    > 花了些功夫去搞舰C的台词呢
* 游戏每日提醒
    > `每日`.tg 的 `每日`.nz 提醒
* 船名查询
    > 就是和谐名对照表了
* 计算器
    > 数据结构的练习
* 半自动化同意入群邀请
    > 命令关键字：`approveAuthorization`,`授权批准`  
    > 使用时需系统管理员权限  
    > 主要用于系统管理员批准邀请加群的申请
* 成语接龙 
    > 计划实现 同字和同音两个版本

### 还有很多功能没有迁移，不想写... ###
~~都2020年了，你是个成熟的IDE了，该学会自己把项目写好了.jpg~~
