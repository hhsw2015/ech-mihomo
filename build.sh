#!/bin/bash

# Mihomo-ECH 跨平台编译脚本

set -e

echo "========================================="
echo "  Mihomo-ECH 跨平台编译脚本"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

appName="ech-mihomo"
# 版本信息
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
GO_VERSION=$(go version | awk '{print $3}')

echo -e "${BLUE}版本信息:${NC}"
echo "  版本: $VERSION"
echo "  构建时间: $BUILD_TIME"
echo "  Go 版本: $GO_VERSION"
echo ""

# 创建输出目录
OUTPUT_DIR="build"
rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR


# 编译单个平台
build_single() {
    local platform=$1
    local goos=$(echo "$platform" | cut -d'/' -f1)
    local goarch_full=$(echo "$platform" | cut -d'/' -f2)
    
    # 解析架构和变体
    local goarch="$goarch_full"
    local goarm=""
    local gomips=""
    local arch_suffix="$goarch_full"
    
    # 处理ARM变体
    if [[ "$goarch_full" =~ ^armv([5-8])$ ]]; then
        goarch="arm"
        goarm="${BASH_REMATCH[1]}"
        arch_suffix="armv${goarm}"
    # 处理MIPS变体
    elif [[ "$goarch_full" =~ ^mips(le)?-(hard|soft)$ ]]; then
        if [[ "$goarch_full" == *"le"* ]]; then
            goarch="mipsle"
        else
            goarch="mips"
        fi
        # 转换MIPS变体名称为Go编译器认可的格式
        if [[ "${BASH_REMATCH[2]}" == "hard" ]]; then
            gomips="hardfloat"
        else
            gomips="softfloat"
        fi
        arch_suffix="$goarch_full"
    fi
    
    echo -e "正在编译 ${platform}..."
    
    # 设置输出文件名
    local output_name="${appName}-${goos}-${arch_suffix}"
    if [[ "$goos" == "windows" ]]; then
        output_name="${appName}-${goos}-${arch_suffix}.exe"
    fi
    
    mkdir -p "$OUTPUT_DIR"
    
    # 准备编译命令
    local env_vars="GOOS=$goos GOARCH=$goarch"
    if [[ -n "$goarm" ]]; then
        env_vars="$env_vars GOARM=$goarm"
    fi
    if [[ -n "$gomips" ]]; then
        env_vars="$env_vars GOMIPS=$gomips"
    fi
    

    local build_cmd="$env_vars go build -trimpath -ldflags=\"-s -w -X 'github.com/metacubex/mihomo/constant.Version=$VERSION' -X 'github.com/metacubex/mihomo/constant.BuildTime=$BUILD_TIME'\" -o '${OUTPUT_DIR}/${output_name}' ."

    # 执行编译
    if eval "$build_cmd" 2>/dev/null; then
        local file_size=$(du -h "${output_dir}/${output_name}" | cut -f1)
        echo -e "✓ ${platform} 编译成功 (${file_size})"
        echo -e "  输出文件: ${output_dir}/${output_name}"
    else
        echo -e "✗ ${platform} 编译失败"
    fi
}

MakeRelease() {
  cd $OUTPUT_DIR
  if [ -d compress ]; then
    rm -rv compress
  fi
  mkdir compress

  
  for i in $(find . -type f -name "$appName-linux-*"); do
    tar -czvf compress/"$i".tar.gz "$i"
  done
  for i in $(find . -type f -name "$appName-android-*"); do
    tar -czvf compress/"$i".tar.gz "$i"
  done
  for i in $(find . -type f -name "$appName-darwin-*"); do
    tar -czvf compress/"$i".tar.gz "$i"
  done
  for i in $(find . -type f -name "$appName-freebsd-*"); do
    tar -czvf compress/"$i".tar.gz "$i"
  done
  for i in $(find . -type f -name "$appName-dragonfly-*"); do
    tar -czvf compress/"$i".tar.gz "$i"
  done
  for i in $(find . -type f -name "$appName-netbsd-*"); do
    tar -czvf compress/"$i".tar.gz "$i"
  done
  for i in $(find . -type f -name "$appName-openbsd-*"); do
    tar -czvf compress/"$i".tar.gz "$i"
  done
  for i in $(find . -type f -name "$appName-plan9-*"); do
    tar -czvf compress/"$i".tar.gz "$i"
  done
  for i in $(find . -type f -name "$appName-solaris-*"); do
    tar -czvf compress/"$i".tar.gz "$i"
  done
  for i in $(find . -type f \( -name "$appName-windows-*" -o -name "$appName-windows7-*" \)); do
    zip compress/$(echo $i | sed 's/\.[^.]*$//').zip "$i"
  done
  
  cd compress
  sha256sum * > SHA256SUMS.txt
  echo "ech-tunnel 构建完成！共 $(ls -1 | grep -E '\.(tar\.gz|zip)$' | wc -l) 个文件"
  
  cd ../..
}

BuildRelease() {
    rm -rf $OUTPUT_DIR
    mkdir -p $OUTPUT_DIR

    # 开始编译
    echo "========================================="
    echo "  开始编译"
    echo "========================================="
    echo ""

    # Linux AMD64
    build_single "linux/amd64"

    # Linux ARM64
    build_single "linux/arm64"

    # Windows AMD64
    build_single "windows/amd64"
    
    build_single "darwin/arm64" 

    # 编译完成
    echo "========================================="
    echo "  编译完成!"
    echo "========================================="
    echo ""
    echo "输出文件:"
    ls -lh $OUTPUT_DIR/
    echo ""
    echo -e "${GREEN}所有文件已保存到 $OUTPUT_DIR/ 目录${NC}"
}

case "$1" in
  release)
    BuildRelease
    MakeRelease
    ;;
  *)
    echo "用法: $0 release"
    ;;
esac
