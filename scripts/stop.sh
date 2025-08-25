#!/bin/bash

# ========================================
# 智能反向代理平台 - 停止脚本
# ========================================

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 项目路径
PROJECT_ROOT=$(dirname "$0")/..
SCRIPTS_DIR="$PROJECT_ROOT/scripts"
BACKEND_PID_FILE="$PROJECT_ROOT/.backend.pid"
FRONTEND_PID_FILE="$PROJECT_ROOT/.frontend.pid"

# 配置文件
ENV_FILE="$PROJECT_ROOT/.env"

# 加载配置
if [ -f "$ENV_FILE" ]; then
    export $(grep -v '^#' "$ENV_FILE" | xargs)
    export SERVER_PORT=${SERVER_PORT:-"8080"}
    export VITE_DEV_PORT=${VITE_DEV_PORT:-"3000"}
else
    export SERVER_PORT="8080"
    export VITE_DEV_PORT="3000"
fi

echo -e "${BLUE}🛑 智能反向代理平台 - 停止服务${NC}"
echo "========================================"

# 停止后端服务
if [ -f "$BACKEND_PID_FILE" ]; then
    BACKEND_PID=$(cat "$BACKEND_PID_FILE")
    if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
        echo -e "${BLUE}🛑 正在停止后端服务 (PID: $BACKEND_PID)...${NC}"
        kill "$BACKEND_PID" 2>/dev/null
        sleep 2
        
        if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
            echo -e "${YELLOW}⚠️  强制停止后端服务...${NC}"
            kill -9 "$BACKEND_PID" 2>/dev/null
        fi
        
        rm -f "$BACKEND_PID_FILE"
        echo -e "${GREEN}✅ 后端服务已停止${NC}"
    else
        echo -e "${YELLOW}⚠️  后端服务未在运行${NC}"
        rm -f "$BACKEND_PID_FILE"
    fi
else
    echo -e "${YELLOW}⚠️  后端服务未在运行${NC}"
fi

# 停止前端服务
if [ -f "$FRONTEND_PID_FILE" ]; then
    FRONTEND_PID=$(cat "$FRONTEND_PID_FILE")
    if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
        echo -e "${BLUE}🛑 正在停止前端服务 (PID: $FRONTEND_PID)...${NC}"
        kill "$FRONTEND_PID" 2>/dev/null
        sleep 2
        
        if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
            echo -e "${YELLOW}⚠️  强制停止前端服务...${NC}"
            kill -9 "$FRONTEND_PID" 2>/dev/null
        fi
        
        rm -f "$FRONTEND_PID_FILE"
        echo -e "${GREEN}✅ 前端服务已停止${NC}"
    else
        echo -e "${YELLOW}⚠️  前端服务未在运行${NC}"
        rm -f "$FRONTEND_PID_FILE"
    fi
else
    echo -e "${YELLOW}⚠️  前端服务未在运行${NC}"
fi

# 强制清理端口占用
echo -e "${BLUE}🧹 正在清理端口占用...${NC}"

# 清理后端端口
if lsof -i ":$SERVER_PORT" >/dev/null 2>&1; then
    local backend_pid=$(lsof -ti ":$SERVER_PORT" | head -1)
    echo -e "${YELLOW}🛑 强制停止占用后端端口 $SERVER_PORT 的进程 (PID: $backend_pid)${NC}"
    kill -9 "$backend_pid" 2>/dev/null
    sleep 1
fi

# 清理前端端口
if lsof -i ":$VITE_DEV_PORT" >/dev/null 2>&1; then
    local frontend_pid=$(lsof -ti ":$VITE_DEV_PORT" | head -1)
    echo -e "${YELLOW}🛑 强制停止占用前端端口 $VITE_DEV_PORT 的进程 (PID: $frontend_pid)${NC}"
    kill -9 "$frontend_pid" 2>/dev/null
    sleep 1
fi

echo ""
echo -e "${GREEN}🎉 所有服务已停止，端口已释放！${NC}"
echo "========================================"
echo -e "${BLUE}💡 提示：${NC}"
echo "- 使用 ../manage.sh 进行交互式管理"
echo "- 使用 ../scripts/quick-start.sh 快速启动所有服务"
echo "- 使用 ../scripts/stop.sh 快速停止所有服务"
