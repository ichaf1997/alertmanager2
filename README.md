# Alertmanager2

Alertmanager Webhook with integration `AliSms` `AliVms` `WxworkRobot` `ElasticSearchAPI` etc.

## Build

```bash
make
```

## How to use

```bash
bin/alertmanager2
```

## API examples

```bash
# 发送短信消息 (Prometheus notification)
curl -X POST "http://127.0.0.1:8080/channel/ali/sendsms?phone=手机号码&templatecode=短信模板代码&signname=短信签名&ak=AK密钥&as=AS密钥&tpname=渲染模板名称" -d @alertmanager_webhook_payload_example.json -v

# 发送语音消息 (Prometheus notification)
curl -X POST "http://127.0.0.1:8080/channel/ali/sendvms?phone=手机号码&templatecode=语音模板代码&ak=AK密钥&as=AS密钥&tpname=渲染模板名称" -d @alertmanager_webhook_payload_example.json -v

# 发送短信消息 (normal notification)
curl 'http://127.0.0.1:8080/channel/ali/sendsimplesms?phone=手机号码&templatecode=短信模板代码' -v -d '{"content": "这是一条测试的短信告警消息"}'

# 发送语音消息 (normal notification)
curl 'http://127.0.0.1:8080/channel/ali/sendsimplevms?phone=手机号码&templatecode=语音模板代码' -v -d '{"content": "这是一条测试的语音告警消息"}'

```

阿里云语音接口默认使用`公共模式`呼出,如需要指定`固定号码`呼出,需要额外指定一个查询参数`callnumber=xxx`