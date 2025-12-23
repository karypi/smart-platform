# -*- coding: utf-8 -*-
import os
import json
import requests
import arrow
import time
from flask import Flask
from flask import request
app = Flask(__name__)
def bytes2json(data_bytes):
    try:
      data = data_bytes.decode('utf8').replace("'", '"')
      return json.loads(data)
    except:
        return "";
def makealertdata(data):
    send_data = {
        "msgtype": "text",
        "text": {
            "content": "消息非prometheus标准格式"
        }
    }
    try:
        alerts_date = data['alerts']
    except:
        alerts_date = ""

    for output in alerts_date[:]:
        try:
            severity = output['labels']['severity']
        except:
            severity = 'null'
        try:
            alertname = output['labels']['alertname']
        except:
            alertname = 'null'
        try:
            instance = output['labels']['instance']
        except:
            instance = 'null'
        try:
            message = output['annotations']['summary']
        except:
            try:
                message = output['annotations']['description']
            except:
                message = 'null'
        try:
            status = output['status']
        except:
            status = 'null'
        try:
            startsAt = output['startsAt']
        except:
            startsAt = time.localtime()
        try:
            endsAt = output['endsAt']
        except:
            endsAt = time.localtime()

        if status == 'firing':
            status_zh = '报警'
            title = '【%s】' % (status_zh)
            send_data = {
                "msgtype": "markdown",
                "markdown": {
                    "content": "## %s \n\n" %title +
                            ">**告警级别**: %s \n\n" % severity +
                            ">**告警类型**: %s \n\n" % alertname +
                            ">**告警主机**: %s \n\n" % instance +
                            ">**告警详情**: %s \n\n" % message +
                            ">**告警状态**: %s \n\n" % status +
                            ">**触发时间**: %s \n\n" % arrow.get(startsAt).to('Asia/Shanghai').format(
                        'YYYY-MM-DD HH:mm:ss')
                }
            }
        elif status == 'resolved':
            status_zh = '恢复'
            title = '【%s】' % (status_zh)
            send_data = {
                "msgtype": "markdown",
                "markdown": {
                    "content": "## %s \n\n" %title +
                            ">**告警级别**: %s \n\n" % severity +
                            ">**告警类型**: %s \n\n" % alertname +
                            ">**告警主机**: %s \n\n" % instance +
                            ">**告警详情**: %s \n\n" % message +
                            ">**告警状态**: %s \n\n" % status +
                            ">**触发时间**: %s \n\n" % arrow.get(startsAt).to('Asia/Shanghai').format(
                        'YYYY-MM-DD HH:mm:ss') +
                            ">**结束时间**: %s \n" % arrow.get(endsAt).to('Asia/Shanghai').format(
                        'YYYY-MM-DD HH:mm:ss')
                }
            }
    return send_data
def send_alert(data):
  #此处获取环境变量“ROBOT_TOKEN”，会在docker-compose的配置文件中配置，docker-compose启动docker时向docker容器注入环境变量
    #token = "xxxx"
    token = os.getenv('ROBOT_TOKEN')
    if not token:
        print('you must set ROBOT_TOKEN env')
        return
    url = 'https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s' % token
    send_data = makealertdata(data)
    req = requests.post(url, json=send_data)
    result = req.json()
    if result['errcode'] != 0:
        print('notify webchat error: %s' % result['errcode'])

@app.route('/', methods=['POST', 'GET'])
def send():
    if request.method == 'POST':
        post_data = request.get_data()
        send_alert(bytes2json(post_data))
        return 'success'
    else:
        return 'weclome to use prometheus alertmanager webchat webhook server!'
if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
