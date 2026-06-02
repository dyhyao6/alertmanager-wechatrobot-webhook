#!/bin/bash

# Test Payload - WeChat/Feishu/Dingtalk 通用告警测试数据
ALERT_PAYLOAD='{
   "version": "4",
   "groupKey": "{}:{alertname=\"node节点cpu使用率过高\"}",
   "status": "firing",
   "Receiver": "webhook",
   "GroupLabels": {
           "alertname": "node节点cpu使用率过高"
   },
   "CommonLabels": {
           "alertname": "node节点cpu使用率过高",
           "container": "kube-rbac-proxy",
           "endpoint": "https",
           "instance": "172.32.2.191",
           "job": "node-exporter",
           "namespace": "monitoring",
           "pod": "node-exporter-8xpz6",
           "prometheus": "monitoring/k8s",
           "service": "node-exporter",
           "severity": "warning"
   },
   "CommonAnnotations": {
           "description": "集群名称:储能-ems-cn  node名称:172.32.2.191  cpu使用率超过85%,当前值:4%",
           "summary": "node节点cpu使用率过高"
   },
   "ExternalURL": "http://alertmanager-main-0:9093",
   "Alerts": [
           {
                   "labels": {
                           "alertname": "node节点cpu使用率过高",
                           "container": "kube-rbac-proxy",
                           "endpoint": "https",
                           "instance": "172.32.2.191",
                           "job": "node-exporter",
                           "namespace": "monitoring",
                           "pod": "node-exporter-8xpz6",
                           "prometheus": "monitoring/k8s",
                           "service": "node-exporter",
                           "severity": "warning"
                   },
                   "Annotations": {
                           "description": "集群名称:储能-ems-cn  node名称:172.32.2.191  cpu使用率超过85%,当前值:4%",
                           "summary": "node节点cpu使用率过高"
                   },
                   "startsAt": "2024-09-10T06:47:48.741Z",
                   "endsAt": "0001-01-01T00:00:00Z"
           }
   ]
}'

# 服务地址
SERVER="http://localhost:8999"

# ========== 微信测试 ==========
echo "========== 微信 WeChat 测试 =========="
curl -X POST \
     -H "Content-Type: application/json" \
     -d "${ALERT_PAYLOAD}" \
     "${SERVER}/webhook/wecom?key=6976d6f8-88db-4ee1-98eb-e238f581e378"

echo -e "\n"

# ========== 飞书测试 ==========
echo "========== 飞书 Feishu 测试 =========="
curl -X POST \
     -H "Content-Type: application/json" \
     -d "${ALERT_PAYLOAD}" \
     "${SERVER}/webhook/feishu?key=0a3f8d87-9fd3-4d14-a58d-25ff7dfd9ef5"

echo -e "\n"

# ========== 钉钉测试 ==========
echo "========== 钉钉 DingTalk 测试 =========="
curl -X POST \
     -H "Content-Type: application/json" \
     -d "${ALERT_PAYLOAD}" \
     "${SERVER}/webhook/dingtalk?key=c46edefd4b38546527c7e21bb88e0b11a7f9cfa87d00c8e2162f6a82890fca2b"

echo -e "\n"