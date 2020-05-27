import sqlite3
import time
import pickle
import os


def sql_read():  # 读
    conn = sqlite3.connect("Ai.db")
    pointer = conn.cursor()
    cursor = pointer.execute("SELECT * FROM universal_corpus")
    out = set()
    i = 0
    for para in cursor:  # 逐个打包
        out.add(para)
        i += 1
    print("数据库中原有:", i, "条数据")
    conn.commit()
    conn.close()
    return out  # 返回结果


def sql_write(info):  # 添加行
    conn = sqlite3.connect("Ai.db")
    pointer = conn.cursor()
    print(info)
    pointer.execute('''INSERT INTO universal_corpus ("answer", "question", "keys") VALUES {};'''.format(info[:3]))
    if info[3] != None:
        pointer.execute('''INSERT INTO universal_corpus ("from") VALUES {};'''.format(info[3]))

    conn.commit()
    conn.close()


def sql_delete():
    conn = sqlite3.connect("Ai.db")
    pointer = conn.cursor()

    pointer.execute('''DELETE FROM universal_corpus;''')

    conn.commit()
    conn.close()


def main():
    print("开始进行重复数据检查")
    a = sql_read()
    tasks = len(a)
    print("去重后共:", tasks, "条")
    inp = input(
        ">>>> ！！ 危险操作 ！！ <<<<\n是否进行数据清洗\n\n>>确定后将会删除数据库中所有数据,然后重写入去重后数据\n！！此操作不可逆！！\n开始后中断将会造成数据丢失\n【Yes】确认  【Exit】退出\n\n")
    while True:
        if inp == "Yes":          
            print("确认进行数据清洗,开始删除数据库记录")

            f = open('cache.pickle','wb')
            pickle.dump(a,f)
            f.close()

            sql_delete()
            print("已删除库中所有记录,开始写入数据")
            s = 0
            for i in a:
                print(i)
                sql_write(i)
                s += 1
                print("已写入", s, "条,共有", tasks, "条需要写入", end="\r")
            input("\n已完成,回车退出")
            os.remove("cache.pickle") 
            return
        elif inp == "Exit":
            return
        elif inp == "R":
            try:
                f=open('cache.pickle','rb')
                a=pickle.load(f)
                f.close()
                tasks = len(a)
                    
                s = 0
                for i in a:
                    sql_write(i)
                    s += 1
                    print("已写入", s, "条,共有", tasks, "条需要写入", end="\r")
            except Exception as err:
                print("执行出错",err)
                exit()
            input("\n已完成,回车退出")
            os.remove("cache.pickle") 
            return  
        else:
            inp = input()


if __name__ == "__main__":
    main()
