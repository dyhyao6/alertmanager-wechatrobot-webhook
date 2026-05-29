#!/usr/bin/env bash

# 默认配置
WECOM_KEY=${WECOM_KEY:-6976d6f8-88db-4ee1-98eb-e238f581e378}
FEISHU_KEY=${FEISHU_KEY:-0a3f8d87-9fd3-4d14-a58d-25ff7dfd9ef5}
DINGTALK_KEY=${DINGTALK_KEY:-your-dingtalk-key}
ADDR=${ADDR:-:8999}
LOG_DIR=${LOG_DIR:-./logs}
LOG_FILE=${LOG_FILE:-webhook.log}

# 创建日志目录
mkdir -p "$LOG_DIR"

./wechat-webhook \
  -wecom.key=$WECOM_KEY \
  -feishu.key=$FEISHU_KEY \
  -dingtalk.key=$DINGTALK_KEY \
  -addr=$ADDR \
  -log.dir=$LOG_DIR \
  -log.file=$LOG_FILE