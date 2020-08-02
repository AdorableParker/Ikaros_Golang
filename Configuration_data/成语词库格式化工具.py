import ujson
import csv

with open('idiom.json', 'r', encoding='utf8')as fp:
    json_data = ujson.load(fp)

with open('idiom.csv', 'w', newline='') as csvfile:
    fieldnames = ['word', 'headPhonetic', 'tailPhonetic']
    writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
    writer.writeheader()

    for i in json_data:
        s = i['pinyin'].split()
        print(s, s[0], s[-1])
        writer.writerow(
            {'word': i['word'], 'headPhonetic': s[0], 'tailPhonetic': s[-1]})
input()
