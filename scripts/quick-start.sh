#!/bin/bash

# ========================================
# 智能反向代理平台 - 快速启动脚本
# ========================================

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# 项目路径
PROJECT_ROOT=$(dirname "$0")/..
SCRIPTS_DIR="$PROJECT_ROOT/scripts"
BACKEND_DIR="$PROJECT_ROOT"
FRONTEND_DIR="$PROJECT_ROOT/web"

# 配置文件
ENV_FILE="$PROJECT_ROOT/.env"

# PID 文件路径
BACKEND_PID_FILE="$PROJECT_ROOT/.backend.pid"
FRONTEND_PID_FILE="$PROJECT_ROOT/.frontend.pid"

# 日志文件路径
BACKEND_LOG_FILE="$PROJECT_ROOT/logs/backend.log"
FRONTEND_LOG_FILE="$PROJECT_ROOT/logs/frontend.log"

# 创建日志目录
mkdir -p "$PROJECT_ROOT/logs"

# 函数：从.env文件读取配置
load_env_config() {
    if [ -f "$ENV_FILE" ]; then
        # 加载环境变量
        export $(grep -v '^#' "$ENV_FILE" | xargs)
        
        # 设置默认值
        export SERVER_HOST=${SERVER_HOST:-"0.0.0.0"}
        export SERVER_PORT=${SERVER_PORT:-"8080"}
        export VITE_DEV_PORT=${VITE_DEV_PORT:-"3000"}
        export VITE_API_BASE_URL=${VITE_API_BASE_URL:-"http://localhost:8080/api"}
        export VITE_APP_TITLE=${VITE_APP_TITLE:-"智能反向代理平台"}
        
        echo -e "${GREEN}✅ 配置文件加载成功${NC}"
        echo -e "${CYAN}   🌐 后端地址：$SERVER_HOST:$SERVER_PORT${NC}"
        echo -e "${CYAN}   🎨 前端地址：localhost:$VITE_DEV_PORT${NC}"
        echo -e "${CYAN}   🔗 API地址：$VITE_API_BASE_URL${NC}"
    else
        echo -e "${YELLOW}⚠️  警告：未找到 .env 文件，使用默认配置${NC}"
        export SERVER_HOST="0.0.0.0"
        export SERVER_PORT="8080"
        export VITE_DEV_PORT="3000"
        export VITE_API_BASE_URL="http://localhost:8080/api"
        export VITE_APP_TITLE="智能反向代理平台"
    fi
}

# 函数：检查端口是否被占用
check_port_usage() {
    local port=$1
    local service_name=$2
    
    if lsof -i ":$port" >/dev/null 2>&1; then
        local pid=$(lsof -ti ":$port" | head -1)
        local process_name=$(ps -p "$pid" -o comm= 2>/dev/null || echo "unknown")
        
        echo -e "${YELLOW}⚠️  警告：$service_name 端口 $port 被进程 $process_name (PID: $pid) 占用${NC}"
        echo -e "${BLUE}🛑 正在强制停止占用端口的进程...${NC}"
        
        kill -9 "$pid" 2>/dev/null
        sleep 2
        
        # 再次检查端口
        if lsof -i ":$port" >/dev/null 2>&1; then
            echo -e "${RED}❌ 无法释放端口 $port，请手动处理${NC}"
            return 1
        else
            echo -e "${GREEN}✅ 端口 $port 已释放${NC}"
            return 0
        fi
    fi
    
    return 0
}

# 函数：强制清理端口占用
force_cleanup_ports() {
    echo -e "${BLUE}🧹 正在强制清理端口占用...${NC}"
    
    # 清理后端端口
    if lsof -i ":$SERVER_PORT" >/dev/null 2>&1; then
        local backend_pid=$(lsof -ti ":$SERVER_PORT" | head -1)
        echo -e "${YELLOW}🛑 强制停止占用后端端口 $SERVER_PORT 的进程 (PID: $backend_pid)${NC}"
        kill -9 "$backend_pid" 2>/dev/null
    fi
    
    # 清理前端端口
    if lsof -i ":$VITE_DEV_PORT" >/dev/null 2>&1; then
        local frontend_pid=$(lsof -ti ":$VITE_DEV_PORT" | head -1)
        echo -e "${YELLOW}🛑 强制停止占用前端端口 $VITE_DEV_PORT 的进程 (PID: $frontend_pid)${NC}"
        kill -9 "$frontend_pid" 2>/dev/null
    fi
    
    # 等待端口释放
    sleep 3
    
    # 清理 PID 文件
    rm -f "$BACKEND_PID_FILE" "$FRONTEND_PID_FILE"
    
    echo -e "${GREEN}✅ 端口清理完成${NC}"
}

echo -e "${GREEN}🚀 智能反向代理平台 - 快速启动${NC}"
echo "========================================"

# 加载配置文件
load_env_config

# 强制清理端口占用
force_cleanup_ports

# 检查环境
echo -e "${BLUE}🔍 检查环境...${NC}"

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo -e "❌ 错误：未找到 Go 环境，请先安装 Go"
    exit 1
fi

# 检查 Node.js 环境
if ! command -v node &> /dev/null; then
    echo -e "❌ 错误：未找到 Node.js 环境，请先安装 Node.js"
    exit 1
fi

echo -e "${GREEN}✅ 环境检查通过${NC}"

# 启动后端
echo -e "${BLUE}🚀 启动后端服务...${NC}"
cd "$BACKEND_DIR"
nohup go run cmd/server/main.go > "$BACKEND_LOG_FILE" 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID > "$BACKEND_PID_FILE"

# 等待后端启动
echo "等待后端服务启动..."
sleep 5

# 检查后端状态
if ps -p $BACKEND_PID > /dev/null; then
    echo -e "${GREEN}✅ 后端服务启动成功！PID: $BACKEND_PID${NC}"
    echo -e "${CYAN}🌐 后端地址：http://$SERVER_HOST:$SERVER_PORT${NC}"
else
    echo -e "❌ 后端服务启动失败，请检查日志：$BACKEND_LOG_FILE"
    rm -f "$BACKEND_PID_FILE"
    exit 1
fi

# 启动前端
echo -e "${BLUE}🚀 启动前端服务...${NC}"
cd "$FRONTEND_DIR"

# 检查依赖
if [ ! -d "node_modules" ]; then
    echo -e "${BLUE}📦 安装前端依赖...${NC}"
    if command -v pnpm &> /dev/null; then
        pnpm install
    else
        npm install
    fi
fi

# 启动前端
if command -v pnpm &> /dev/null; then
    nohup pnpm dev > "$FRONTEND_LOG_FILE" 2>&1 &
else
    nohup npm run dev > "$FRONTEND_LOG_FILE" 2>&1 &
fi
FRONTEND_PID=$!
echo $FRONTEND_PID > "$FRONTEND_PID_FILE"

# 等待前端启动
echo "等待前端服务启动..."
sleep 8

# 检查前端状态
if ps -p $FRONTEND_PID > /dev/null; then
    echo -e "${GREEN}✅ 前端服务启动成功！PID: $FRONTEND_PID${NC}"
    echo -e "${CYAN}🌐 前端地址：http://localhost:$VITE_DEV_PORT${NC}"
else
    echo -e "❌ 前端服务启动失败，请检查日志：$FRONTEND_LOG_FILE"
    rm -f "$FRONTEND_PID_FILE"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 所有服务启动成功！${NC}"
echo "========================================"
echo -e "${CYAN}🌐 前端界面：http://localhost:$VITE_DEV_PORT${NC}"
echo -e "${CYAN}🔧 后端API：http://$SERVER_HOST:$SERVER_PORT${NC}"
echo -e "${CYAN}🔗 API地址：$VITE_API_BASE_URL${NC}"
echo -e "${CYAN}📱 应用标题：$VITE_APP_TITLE${NC}"
echo -e "${CYAN}📊 后端日志：$BACKEND_LOG_FILE${NC}"
echo -e "${CYAN}📊 前端日志：$FRONTEND_LOG_FILE${NC}"
echo ""
echo -e "${BLUE}💡 提示：${NC}"
echo "- 使用 ../manage.sh 进行交互式管理"
echo "- 使用 ../scripts/stop.sh 快速停止"
echo "- 查看日志了解服务运行状态"
echo ""
echo -e "${GREEN}🚀 开始享受你的智能反向代理平台吧！${NC}"
