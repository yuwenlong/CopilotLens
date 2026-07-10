#!/bin/bash
APP_NAME="copilotlens"
BIN_DIR="bin"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 检测操作系统
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
    APP_NAME="$APP_NAME.exe"
fi

start() {
    if pgrep -f "$APP_NAME" > /dev/null 2>&1; then
        echo "服务已在运行 (PID: $(pgrep -f "$APP_NAME"))"
        return 1
    fi
    echo "启动 $APP_NAME ..."
    cd "$SCRIPT_DIR/$BIN_DIR"
    nohup "./$APP_NAME" > /dev/null 2>&1 &
    cd "$SCRIPT_DIR"
    echo "服务已启动"
}

stop() {
    PIDS=$(pgrep -f "$APP_NAME")
    if [ -n "$PIDS" ]; then
        echo "停止服务 (PID: $PIDS) ..."
        echo "$PIDS" | xargs kill
        sleep 1
        PIDS=$(pgrep -f "$APP_NAME")
        if [ -n "$PIDS" ]; then
            echo "$PIDS" | xargs kill -9
        fi
        echo "服务已停止"
    else
        echo "服务未运行"
    fi
}

case "${1:-}" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        stop
        sleep 1
        start
        ;;
    reload)
        echo "重新构建..."
        cd "$SCRIPT_DIR"
        bash build.sh
        stop
        sleep 1
        start
        ;;
    *)
        echo "用法: $0 {start|stop|restart|reload}"
        exit 1
        ;;
esac
