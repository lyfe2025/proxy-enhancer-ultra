#!/bin/bash

# ========================================
# 智能反向代理平台 - 统一启动管理脚本
# ========================================

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 项目路径
PROJECT_ROOT=$(pwd)
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
        
        print_message $GREEN "✅ 配置文件加载成功"
        print_message $CYAN "   🔧 API服务器地址：$SERVER_HOST:$SERVER_PORT"
        print_message $CYAN "   🎨 管理后台访问地址：localhost:$VITE_DEV_PORT"
        print_message $CYAN "   🔗 API地址：$VITE_API_BASE_URL"
    else
        print_message $YELLOW "⚠️  警告：未找到 .env 文件，使用默认配置"
        export SERVER_HOST="0.0.0.0"
        export SERVER_PORT="8080"
        export VITE_DEV_PORT="3000"
        export VITE_API_BASE_URL="http://localhost:8080/api"
        export VITE_APP_TITLE="智能反向代理平台"
    fi
}

# 函数：打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 函数：检查端口是否被占用
check_port_usage() {
    local port=$1
    local service_name=$2
    
    if lsof -i ":$port" >/dev/null 2>&1; then
        local pid=$(lsof -ti ":$port" | head -1)
        local process_name=$(ps -p "$pid" -o comm= 2>/dev/null || echo "unknown")
        
        print_message $YELLOW "⚠️  警告：$service_name 端口 $port 被进程 $process_name (PID: $pid) 占用"
        
        # 询问是否强制停止
        read -p "是否强制停止占用端口的进程？(y/N): " force_kill
        
        if [[ "$force_kill" =~ ^[Yy]$ ]]; then
            print_message $BLUE "🛑 正在强制停止占用端口的进程..."
            kill -9 "$pid" 2>/dev/null
            sleep 2
            
            # 再次检查端口
            if lsof -i ":$port" >/dev/null 2>&1; then
                print_message $RED "❌ 无法释放端口 $port，请手动处理"
                return 1
            else
                print_message $GREEN "✅ 端口 $port 已释放"
                return 0
            fi
        else
            print_message $RED "❌ 端口被占用，无法启动 $service_name"
            return 1
        fi
    fi
    
    return 0
}

# 函数：强制清理端口占用
force_cleanup_ports() {
    print_message $BLUE "🧹 正在强制清理端口占用..."
    
    # 清理后端端口
    if lsof -i ":$SERVER_PORT" >/dev/null 2>&1; then
        local backend_pid=$(lsof -ti ":$SERVER_PORT" | head -1)
        print_message $YELLOW "🛑 强制停止占用后端端口 $SERVER_PORT 的进程 (PID: $backend_pid)"
        kill -9 "$backend_pid" 2>/dev/null
    fi
    
    # 清理前端端口
    if lsof -i ":$VITE_DEV_PORT" >/dev/null 2>&1; then
        local frontend_pid=$(lsof -ti ":$VITE_DEV_PORT" | head -1)
        print_message $YELLOW "🛑 强制停止占用前端端口 $VITE_DEV_PORT 的进程 (PID: $frontend_pid)"
        kill -9 "$frontend_pid" 2>/dev/null
    fi
    
    # 等待端口释放
    sleep 3
    
    # 清理 PID 文件
    rm -f "$BACKEND_PID_FILE" "$FRONTEND_PID_FILE"
    
    print_message $GREEN "✅ 端口清理完成"
}

# 函数：检查服务是否运行
is_backend_running() {
    if [ -f "$BACKEND_PID_FILE" ]; then
        local pid=$(cat "$BACKEND_PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0
        else
            rm -f "$BACKEND_PID_FILE"
        fi
    fi
    return 1
}

is_frontend_running() {
    if [ -f "$FRONTEND_PID_FILE" ]; then
        local pid=$(cat "$FRONTEND_PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0
        else
            rm -f "$FRONTEND_PID_FILE"
        fi
    fi
    return 1
}

# 函数：启动后端
start_backend() {
    if is_backend_running; then
        print_message $YELLOW "⚠️  后端服务已经在运行中..."
        return
    fi
    
    # 检查端口占用
    if ! check_port_usage "$SERVER_PORT" "后端服务"; then
        return 1
    fi
    
    print_message $BLUE "🚀 正在启动后端服务..."
    
    # 检查 Go 环境
    if ! command -v go &> /dev/null; then
        print_message $RED "❌ 错误：未找到 Go 环境，请先安装 Go"
        return 1
    fi
    
    # 启动后端服务
    cd "$BACKEND_DIR"
    nohup go run cmd/server/main.go > "$BACKEND_LOG_FILE" 2>&1 &
    local backend_pid=$!
    echo $backend_pid > "$BACKEND_PID_FILE"
    
    # 等待服务启动
    sleep 3
    
    if is_backend_running; then
        print_message $GREEN "✅ 后端服务启动成功！PID: $backend_pid"
        print_message $CYAN "📊 后端日志：$BACKEND_LOG_FILE"
        print_message $CYAN "🔧 API服务器地址：http://$SERVER_HOST:$SERVER_PORT"
    else
        print_message $RED "❌ 后端服务启动失败，请检查日志：$BACKEND_LOG_FILE"
        rm -f "$BACKEND_PID_FILE"
        return 1
    fi
}

# 函数：启动前端
start_frontend() {
    if is_frontend_running; then
        print_message $YELLOW "⚠️  前端服务已经在运行中..."
        return
    fi
    
    # 检查端口占用
    if ! check_port_usage "$VITE_DEV_PORT" "前端服务"; then
        return 1
    fi
    
    print_message $BLUE "🚀 正在启动前端服务..."
    
    # 检查 Node.js 环境
    if ! command -v node &> /dev/null; then
        print_message $RED "❌ 错误：未找到 Node.js 环境，请先安装 Node.js"
        return 1
    fi
    
    # 检查 pnpm
    if ! command -v pnpm &> /dev/null; then
        print_message $YELLOW "⚠️  警告：未找到 pnpm，尝试使用 npm..."
        if ! command -v npm &> /dev/null; then
            print_message $RED "❌ 错误：未找到 npm，请先安装 Node.js"
            return 1
        fi
        PACKAGE_MANAGER="npm"
    else
        PACKAGE_MANAGER="pnpm"
    fi
    
    # 进入前端目录
    cd "$FRONTEND_DIR"
    
    # 检查依赖是否安装
    if [ ! -d "node_modules" ]; then
        print_message $YELLOW "📦 正在安装前端依赖..."
        $PACKAGE_MANAGER install
        if [ $? -ne 0 ]; then
            print_message $RED "❌ 依赖安装失败"
            return 1
        fi
    fi
    
    # 启动前端服务
    nohup $PACKAGE_MANAGER dev > "$FRONTEND_LOG_FILE" 2>&1 &
    local frontend_pid=$!
    echo $frontend_pid > "$FRONTEND_PID_FILE"
    
    # 等待服务启动
    sleep 5
    
    if is_frontend_running; then
        print_message $GREEN "✅ 前端服务启动成功！PID: $frontend_pid"
        print_message $CYAN "📊 前端日志：$FRONTEND_LOG_FILE"
        print_message $CYAN "🎨 管理后台访问地址：http://localhost:$VITE_DEV_PORT"
    else
        print_message $RED "❌ 前端服务启动失败，请检查日志：$FRONTEND_LOG_FILE"
        rm -f "$FRONTEND_PID_FILE"
        return 1
    fi
}

# 函数：停止后端
stop_backend() {
    if ! is_backend_running; then
        print_message $YELLOW "⚠️  后端服务未在运行"
        return
    fi
    
    local pid=$(cat "$BACKEND_PID_FILE")
    print_message $BLUE "🛑 正在停止后端服务 (PID: $pid)..."
    
    kill "$pid" 2>/dev/null
    sleep 2
    
    if ps -p "$pid" > /dev/null 2>&1; then
        print_message $YELLOW "⚠️  强制停止后端服务..."
        kill -9 "$pid" 2>/dev/null
    fi
    
    rm -f "$BACKEND_PID_FILE"
    
    # 确保端口完全释放
    if lsof -i ":$SERVER_PORT" >/dev/null 2>&1; then
        local port_pid=$(lsof -ti ":$SERVER_PORT" | head -1)
        print_message $YELLOW "⚠️  端口 $SERVER_PORT 仍被进程 $port_pid 占用，强制释放..."
        kill -9 "$port_pid" 2>/dev/null
        sleep 1
    fi
    
    print_message $GREEN "✅ 后端服务已停止，端口已释放"
}

# 函数：停止前端
stop_frontend() {
    if ! is_frontend_running; then
        print_message $YELLOW "⚠️  前端服务未在运行"
        return
    fi
    
    local pid=$(cat "$FRONTEND_PID_FILE")
    print_message $BLUE "🛑 正在停止前端服务 (PID: $pid)..."
    
    kill "$pid" 2>/dev/null
    sleep 2
    
    if ps -p "$pid" > /dev/null 2>&1; then
        print_message $YELLOW "⚠️  强制停止前端服务..."
        kill -9 "$pid" 2>/dev/null
    fi
    
    rm -f "$FRONTEND_PID_FILE"
    
    # 确保端口完全释放
    if lsof -i ":$VITE_DEV_PORT" >/dev/null 2>&1; then
        local port_pid=$(lsof -ti ":$VITE_DEV_PORT" | head -1)
        print_message $YELLOW "⚠️  端口 $VITE_DEV_PORT 仍被进程 $port_pid 占用，强制释放..."
        kill -9 "$port_pid" 2>/dev/null
        sleep 1
    fi
    
    print_message $GREEN "✅ 前端服务已停止，端口已释放"
}

# 函数：重启后端
restart_backend() {
    print_message $BLUE "🔄 正在重启后端服务..."
    stop_backend
    sleep 2
    start_backend
}

# 函数：重启前端
restart_frontend() {
    print_message $BLUE "🔄 正在重启前端服务..."
    stop_frontend
    sleep 2
    start_frontend
}

# 函数：显示服务状态
show_status() {
    print_message $CYAN "📊 服务状态："
    echo "----------------------------------------"
    
    if is_backend_running; then
        local backend_pid=$(cat "$BACKEND_PID_FILE")
        print_message $GREEN "✅ 后端服务：运行中 (PID: $backend_pid)"
        print_message $CYAN "   🔧 API服务器地址：http://$SERVER_HOST:$SERVER_PORT"
        print_message $CYAN "   📊 日志：$BACKEND_LOG_FILE"
    else
        print_message $RED "❌ 后端服务：未运行"
    fi
    
    echo ""
    
    if is_frontend_running; then
        local frontend_pid=$(cat "$FRONTEND_PID_FILE")
        print_message $GREEN "✅ 前端服务：运行中 (PID: $frontend_pid)"
        print_message $CYAN "   🎨 管理后台访问地址：http://localhost:$VITE_DEV_PORT"
        print_message $CYAN "   📊 日志：$FRONTEND_LOG_FILE"
    else
        print_message $RED "❌ 前端服务：未运行"
    fi
    
    echo ""
    print_message $CYAN "🔗 API地址：$VITE_API_BASE_URL"
    print_message $CYAN "📱 应用标题：$VITE_APP_TITLE"
    echo "----------------------------------------"
}

# 函数：显示日志
show_logs() {
    print_message $CYAN "📋 日志查看选项："
    echo "1. 查看后端日志"
    echo "2. 查看前端日志"
    echo "3. 实时监控后端日志"
    echo "4. 实时监控前端日志"
    echo "5. 返回主菜单"
    
    read -p "请选择 (1-5): " log_choice
    
    case $log_choice in
        1)
            if [ -f "$BACKEND_LOG_FILE" ]; then
                print_message $CYAN "📊 后端日志："
                cat "$BACKEND_LOG_FILE"
            else
                print_message $YELLOW "⚠️  后端日志文件不存在"
            fi
            ;;
        2)
            if [ -f "$FRONTEND_LOG_FILE" ]; then
                print_message $CYAN "📊 前端日志："
                cat "$FRONTEND_LOG_FILE"
            else
                print_message $YELLOW "⚠️  前端日志文件不存在"
            fi
            ;;
        3)
            if [ -f "$BACKEND_LOG_FILE" ]; then
                print_message $CYAN "🔍 实时监控后端日志 (按 Ctrl+C 退出)："
                tail -f "$BACKEND_LOG_FILE"
            else
                print_message $YELLOW "⚠️  后端日志文件不存在"
            fi
            ;;
        4)
            if [ -f "$FRONTEND_LOG_FILE" ]; then
                print_message $CYAN "🔍 实时监控前端日志 (按 Ctrl+C 退出)："
                tail -f "$FRONTEND_LOG_FILE"
            else
                print_message $YELLOW "⚠️  前端日志文件不存在"
            fi
            ;;
        5)
            return
            ;;
        *)
            print_message $RED "❌ 无效选择"
            ;;
    esac
}

# 函数：清理环境
cleanup() {
    print_message $BLUE "🧹 正在清理环境..."
    
    # 停止所有服务
    stop_backend
    stop_frontend
    
    # 清理 PID 文件
    rm -f "$BACKEND_PID_FILE" "$FRONTEND_PID_FILE"
    
    print_message $GREEN "✅ 环境清理完成"
}

# 函数：显示帮助信息
show_help() {
    print_message $CYAN "📖 帮助信息："
    echo "----------------------------------------"
    echo "这是一个智能反向代理平台的统一启动脚本，提供以下功能："
    echo ""
    echo "🚀 启动服务："
    echo "  - 启动前后端：同时启动 Go 后端和 Vue 前端"
    echo "  - 单独启动后端：只启动 Go 后端服务"
    echo "  - 单独启动前端：只启动 Vue 前端服务"
    echo ""
    echo "🛑 停止服务："
    echo "  - 停止所有服务：同时停止前后端"
    echo "  - 单独停止后端：只停止 Go 后端"
    echo "  - 单独停止前端：只停止 Vue 前端"
    echo ""
    echo "🔄 重启服务："
    echo "  - 重启所有服务：重启前后端"
    echo "  - 单独重启后端：重启 Go 后端"
    echo "  - 单独重启前端：重启 Vue 前端"
    echo ""
    echo "📊 其他功能："
    echo "  - 查看服务状态：显示前后端运行状态"
    echo "  - 查看日志：查看或监控服务日志"
    echo "  - 清理环境：停止所有服务并清理"
    echo "  - 强制清理端口：清理端口占用"
    echo ""
    echo "⚙️  配置管理："
    echo "  - 所有配置从 .env 文件读取"
    echo "  - 支持动态端口和地址配置"
    echo "  - 实时显示正确的访问地址"
    echo "  - 自动检测和清理端口占用"
    echo ""
    echo "⚠️  注意事项："
    echo "  - 确保已安装 Go 1.21+ 和 Node.js"
    echo "  - 确保项目根目录有 .env 配置文件"
    echo "  - 前端依赖会自动安装（如果不存在）"
    echo "  - 修改 .env 文件后重启服务生效"
    echo "  - 脚本会自动检测和清理端口占用"
    echo "----------------------------------------"
}

# 函数：快速启动
quick_start() {
    print_message $BLUE "🚀 快速启动所有服务..."
    
    # 启动后端
    start_backend
    if [ $? -ne 0 ]; then
        print_message $RED "❌ 后端启动失败，停止操作"
        return 1
    fi
    
    # 启动前端
    start_frontend
    if [ $? -ne 0 ]; then
        print_message $RED "❌ 前端启动失败，停止后端服务"
        stop_backend
        return 1
    fi
    
    # 显示成功信息
    echo ""
    print_message $GREEN "🎉 所有服务启动成功！"
    echo "========================================"
    print_message $CYAN "🎨 管理后台界面：http://localhost:$VITE_DEV_PORT"
    print_message $CYAN "🔧 API服务器地址：http://$SERVER_HOST:$SERVER_PORT"
    print_message $CYAN "🔗 API地址：$VITE_API_BASE_URL"
    print_message $CYAN "📱 应用标题：$VITE_APP_TITLE"
    print_message $CYAN "📊 后端日志：$BACKEND_LOG_FILE"
    print_message $CYAN "📊 前端日志：$FRONTEND_LOG_FILE"
    echo ""
    print_message $BLUE "💡 提示："
    echo "- 使用菜单选项 10 查看服务状态"
    echo "- 使用菜单选项 11 查看服务日志"
    echo "- 使用菜单选项 0 退出管理"
    echo ""
    print_message $GREEN "🚀 开始享受你的智能反向代理平台吧！"
}

# 函数：快速停止
quick_stop() {
    print_message $BLUE "🛑 快速停止所有服务..."
    stop_backend
    stop_frontend
    print_message $GREEN "🎉 所有服务已停止！"
}

# 函数：主菜单
main_menu() {
    while true; do
        echo ""
        print_message $PURPLE "🎯 智能反向代理平台 - 统一启动管理"
        echo "========================================"
        echo "1. 🚀 启动前后端"
        echo "2. 🚀 单独启动后端"
        echo "3. 🚀 单独启动前端"
        echo "4. 🛑 停止所有服务"
        echo "5. 🛑 单独停止后端"
        echo "6. 🛑 单独停止前端"
        echo "7. 🔄 重启所有服务"
        echo "8. 🔄 单独重启后端"
        echo "9. 🔄 单独重启前端"
        echo "10. 📊 查看服务状态"
        echo "11. 📋 查看日志"
        echo "12. 🧹 清理环境"
        echo "13. 📖 帮助信息"
        echo "14. ⚡ 快速启动"
        echo "15. 🛑 快速停止"
        echo "16. 🧹 强制清理端口"
        echo "0. 🚪 退出"
        echo "========================================"
        
        read -p "请选择操作 (0-16): " choice
        
        case $choice in
            1)
                print_message $BLUE "🚀 正在启动前后端服务..."
                start_backend
                if [ $? -eq 0 ]; then
                    start_frontend
                fi
                ;;
            2)
                start_backend
                ;;
            3)
                start_frontend
                ;;
            4)
                print_message $BLUE "🛑 正在停止所有服务..."
                stop_backend
                stop_frontend
                ;;
            5)
                stop_backend
                ;;
            6)
                stop_frontend
                ;;
            7)
                print_message $BLUE "🔄 正在重启所有服务..."
                stop_backend
                stop_frontend
                sleep 2
                start_backend
                if [ $? -eq 0 ]; then
                    start_frontend
                fi
                ;;
            8)
                restart_backend
                ;;
            9)
                restart_frontend
                ;;
            10)
                show_status
                ;;
            11)
                show_logs
                ;;
            12)
                cleanup
                ;;
            13)
                show_help
                ;;
            14)
                quick_start
                ;;
            15)
                quick_stop
                ;;
            16)
                force_cleanup_ports
                ;;
            0)
                print_message $GREEN "👋 再见！"
                exit 0
                ;;
            *)
                print_message $RED "❌ 无效选择，请重新输入"
                ;;
        esac
        
        echo ""
        read -p "按回车键继续..."
    done
}

# 主程序入口
main() {
    # 检查是否在项目根目录
    if [ ! -f "go.mod" ] || [ ! -d "web" ]; then
        print_message $RED "❌ 错误：请在项目根目录运行此脚本"
        exit 1
    fi
    
    # 设置信号处理
    trap cleanup EXIT
    trap 'echo ""; print_message $YELLOW "⚠️  收到中断信号，正在清理..."; cleanup; exit 1' INT TERM
    
    # 加载配置文件
    load_env_config
    
    # 显示欢迎信息
    print_message $GREEN "🎉 欢迎使用智能反向代理平台统一启动脚本！"
    print_message $CYAN "📁 项目路径：$PROJECT_ROOT"
            print_message $CYAN "🔧 API服务器路径：$BACKEND_DIR"
        print_message $CYAN "🎨 管理后台路径：$FRONTEND_DIR"
    print_message $CYAN "⚙️  配置文件：$ENV_FILE"
    echo ""
    
    # 启动主菜单
    main_menu
}

# 运行主程序
main
