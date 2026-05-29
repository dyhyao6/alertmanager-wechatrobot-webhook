# Alertmanager 企业微信机器人 Webhook 离线部署指南

## 一、构建镜像

在有网络的机器上执行：

```bash
# 进入项目目录
cd alertmanager-wechatrobot-webhook

# 构建镜像（指定 amd64 平台，避免目标机器兼容性问题）
docker build --platform linux/amd64 -t registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1 .

# 验证镜像
docker images registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1
```

## 二、导出镜像（离线传输）

```bash
# 导出镜像为 tar 文件
docker save --platform linux/amd64 -o wechat-webhook-ems1.1.tar registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1
```

## 三、离线导入镜像

将 tar 文件拷贝到目标机器后执行：

```bash
# 加载镜像

docker stop alertmanager-wechatbot-webhook && docker rm alertmanager-wechatbot-webhook

docker load -i wechat-webhook-ems1.1.tar

```

## 四、启动服务

```bash
# 在目标机器上执行
docker run -d \
  --name alertmanager-wechatbot-webhook \
  --network host \
  registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1 \
  /usr/bin/wechat-webhook \
  -RobotKey 6976d6f8-88db-4ee1-98eb-e238f581e378 \
  -addr :8999 \
  -grafanaUrl 172.16.11.236:3000 \
  -alertDomain 172.16.11.236:9093
  
  
# 查看日志
docker logs -f alertmanager-wechatbot-webhook
```

**参数说明：**
| 参数 | 说明 |
|------|------|
| -RobotKey | 企业微信机器人 Webhook Key |
| -addr | 服务监听地址 |
| -grafanaUrl | Grafana 面板地址（不含 http://） |
| -alertDomain | Alertmanager 地址（不含 http://） |

## 五、Alertmanager 配置

在 Alertmanager 配置文件中添加 webhook receiver：

```yaml
receivers:
  - name: 'wechatbot'
    webhook_configs:
      - url: 'http://目标机器IP:8999/webhook'
        send_resolved: true
```

### 5.1 多机器人配置

#### 方式一：通过告警规则注解指定

在 Prometheus 告警规则的 `annotations` 中添加：

```yaml
annotations:
  summary: "CPU 使用率告警"
  wechatRobot: 另一个机器人的key
```

#### 方式二：通过 URL 参数指定

```
http://service:8999/webhook?key=另一个RobotKey
```

## 六、验证部署

```bash
# 查看容器状态
docker ps | grep wechatbot

# 查看日志
docker logs alertmanager-wechatbot-webhook

# 测试 webhook（单行命令）
curl -X POST http://localhost:8999/webhook -H "Content-Type: application/json" -d '{"version":"4","groupKey":"{}:test","status":"firing","receiver":"webhook","groupLabels":{},"commonLabels":{},"commonAnnotations":{"summary":"测试告警：服务实例宕机"},"externalURL":"http://alertmanager.local","alerts":[{"status":"firing","labels":{"alertname":"ServiceDown","severity":"critical","service":"api","instance":"192.168.1.100:9100"},"annotations":{"description":"服务实例 192.168.1.100:9100 指标采集失败超过 1 分钟"},"startsAt":"2026-05-29T12:00:00Z","endsAt":"0001-01-01T00:00:00Z"}]}'
```

## 七、服务管理

```bash
# 启动服务
docker run -d \
  --name alertmanager-wechatbot-webhook \
  --network host \
  registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1 \
  /usr/bin/wechat-webhook \
  -RobotKey 6976d6f8-88db-4ee1-98eb-e238f581e378 \
  -addr :8999 \
  -grafanaUrl 172.16.11.236:3000 \
  -alertDomain 172.16.11.236:9093

# 停止服务 删除容器
docker stop alertmanager-wechatbot-webhook && docker rm alertmanager-wechatbot-webhook

# 重启服务
docker restart alertmanager-wechatbot-webhook

# 查看日志
docker logs -f alertmanager-wechatbot-webhook
```

## 八、卸载

```bash
docker stop alertmanager-wechatbot-webhook
docker rm alertmanager-wechatbot-webhook
```

## 九、文件清单

离线部署需要拷贝到目标机器的文件：

| 文件 | 说明 |
|------|------|
| `wechat-webhook-ems1.1.tar` | 镜像包 |
| `alertmanager.yaml` | Alertmanager 配置示例 |