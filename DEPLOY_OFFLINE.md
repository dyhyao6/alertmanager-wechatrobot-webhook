# Alertmanager Webhook 多平台离线部署指南

## 一、构建镜像

在有网络的机器上执行：

```bash
# 进入项目目录
cd alertmanager-wechatrobot-webhook

# 构建镜像（指定 amd64 平台）
docker build --platform linux/amd64 -t registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1 .

# 验证镜像
docker images registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1
```

## 二、导出镜像（离线传输）

```bash
# 导出镜像
docker save --platform linux/amd64 -o wechat-webhook-ems1.1.tar registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1

# 压缩（可选）
gzip wechat-webhook-ems1.1.tar
```

## 三、离线导入镜像

将 tar 文件拷贝到目标机器后执行：

```bash
# 停止并删除旧容器
docker stop alertmanager-webhook && docker rm alertmanager-webhook

# 加载镜像
docker load -i wechat-webhook-ems1.1.tar

# 验证
docker images registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1
```

## 四、启动服务

创建 `docker-compose.yaml`：

```yaml
services:
  webhook:
    image: registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1
    container_name: alertmanager-webhook
    network_mode: host
    volumes:
      - ./logs:/var/log/wechat-webhook
    command:
      - /usr/bin/wechat-webhook
      - -wecom.key=your-wechat-key
      - -feishu.key=your-feishu-key
      - -dingtalk.key=your-dingtalk-key
      - -addr=:8999
      - -log.dir=/var/log/wechat-webhook
      - -log.file=webhook.log
```

启动服务：

```bash
docker compose up -d

# 查看日志
docker compose logs -f
```

## 五、启动参数

| 参数 | 说明 |
|------|------|
| -wecom.key | 企业微信机器人 Webhook Key |
| -feishu.key | 飞书机器人 Webhook Key |
| -dingtalk.key | 钉钉机器人 Webhook Key |
| -addr | 服务监听地址（默认 :8999） |

## 六、API 接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/webhook/wecom` | POST | 企业微信 |
| `/webhook/feishu` | POST | 飞书 |
| `/webhook/dingtalk` | POST | 钉钉 |
| `/health` | GET | 健康检查 |

## 七、Alertmanager 配置

在 Alertmanager 配置文件中添加对应平台的 receiver：

```yaml
receivers:
  # 企业微信
  - name: 'wecom'
    webhook_configs:
      - url: 'http://目标机器IP:8999/webhook/wecom'
        send_resolved: true

  # 飞书
  - name: 'feishu'
    webhook_configs:
      - url: 'http://目标机器IP:8999/webhook/feishu'
        send_resolved: true

  # 钉钉
  - name: 'dingtalk'
    webhook_configs:
      - url: 'http://目标机器IP:8999/webhook/dingtalk'
        send_resolved: true
```

## 八、验证部署

```bash
# 查看容器状态
docker ps | grep alertmanager-webhook

# 查看日志
docker logs alertmanager-webhook

# 测试企业微信
curl -X POST http://localhost:8999/webhook/wecom \
  -H "Content-Type: application/json" \
  -d '{"version":"4","groupKey":"{}:test","status":"firing","receiver":"webhook","groupLabels":{},"commonLabels":{},"commonAnnotations":{"summary":"测试告警：服务实例宕机"},"externalURL":"http://alertmanager.local","alerts":[{"status":"firing","labels":{"alertname":"ServiceDown","severity":"critical","service":"api","instance":"192.168.1.100:9100"},"annotations":{"description":"服务实例 192.168.1.100:9100 指标采集失败超过 1 分钟"},"startsAt":"2026-05-29T12:00:00Z","endsAt":"0001-01-01T00:00:00Z"}]}'

# 测试飞书
curl -X POST http://localhost:8999/webhook/feishu \
  -H "Content-Type: application/json" \
  -d '{"version":"4","groupKey":"{}:test","status":"firing","receiver":"webhook","groupLabels":{},"commonLabels":{},"commonAnnotations":{"summary":"测试告警：服务实例宕机"},"externalURL":"http://alertmanager.local","alerts":[{"status":"firing","labels":{"alertname":"ServiceDown","severity":"critical","service":"api","instance":"192.168.1.100:9100"},"annotations":{"description":"服务实例 192.168.1.100:9100 指标采集失败超过 1 分钟"},"startsAt":"2026-05-29T12:00:00Z","endsAt":"0001-01-01T00:00:00Z"}]}'
```

## 九、服务管理

```bash
# 启动服务
docker compose up -d

# 停止服务
docker compose down

# 重启服务
docker compose restart

# 查看日志
docker compose logs -f
```

## 十、卸载

```bash
docker stop alertmanager-webhook
docker rm alertmanager-webhook
```