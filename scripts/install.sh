#!/usr/bin/env bash

# 配置参数
REPO_OWNER="dongfg"
REPO_NAME="dogecli"
VERSION="${1:-latest}"  # 支持版本参数

# 系统信息检测
OS=$(uname -s)
case $(uname -m) in
  "x86_64") ARCH="amd64" ;;
  "aarch64") ARCH="arm64" ;;
  "armv7l") ARCH="armv7" ;;
  "armv6l") ARCH="armv6" ;;
  *) ARCH=$(uname -m) ;;
esac

# 构建 API URL
if [ "$VERSION" = "latest" ]; then
  API_URL="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"
else
  API_URL="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/tags/${VERSION}"
fi

# 获取 Release 信息
RELEASE_JSON=$(curl -fsSL "$API_URL") || { echo "API 请求失败"; exit 1; }

# 构建文件名匹配模式 (基于 GoReleaser 标准命名)
FILE_PATTERN="${REPO_NAME}_${OS}_${ARCH}\."

# 提取下载链接
DOWNLOAD_URL=$(echo "$RELEASE_JSON" | grep -o "\"browser_download_url\": \"[^\"]*${FILE_PATTERN}\(tar\.gz\|zip\)\"" | cut -d'"' -f4)

if [ -z "$DOWNLOAD_URL" ]; then
  echo "错误：未找到 ${OS}/${ARCH} 的发布文件"
  echo "支持架构：amd64, arm64, armv7, armv6"
  exit 1
fi

# 下载文件
FILENAME=$(basename "$DOWNLOAD_URL")
echo "正在下载: $FILENAME"
curl -fsSL -O "$DOWNLOAD_URL" || { echo "下载失败"; exit 1; }

# 校验和验证
CHECKSUM_URL=$(echo "$RELEASE_JSON" | grep -o "\"browser_download_url\": \"[^\"]*checksums.txt\"" | cut -d'"' -f4)

curl -fsSL -o checksums.txt "$CHECKSUM_URL" || { echo "校验文件下载失败"; exit 1; }

ACTUAL_SUM=$(sha256sum "$FILENAME" | cut -d' ' -f1)
grep -q "$ACTUAL_SUM" checksums.txt || {
  echo "校验失败！文件可能损坏"
  rm -f "$FILENAME" checksums.txt
  exit 1
}

echo "正在解压: $FILENAME"

# 解压处理
case "$FILENAME" in
  *.tar.gz) tar xzf "$FILENAME" ;;
  *.zip) unzip -q "$FILENAME" ;;
  *) echo "未知文件格式: $FILENAME" ;;
esac

# 清理文件 (可选)
rm -f "$FILENAME" checksums.txt