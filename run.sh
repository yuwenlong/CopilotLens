#!/bin/bash
APP_NAME="copilotlens"
BIN_DIR="bin"
PID_FILE="$BIN_DIR/.pid"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 检测操作系统
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
    APP_NAME="$APP_NAME.exe"
fi

start() {
    if [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
        echo "服务已在运行 (PID: $(cat "$PID_FILE"))"
        return 1
    fi
    echo "启动 $APP_NAME ..."
    cd "$SCRIPT_DIR/$BIN_DIR"
    nohup "./$APP_NAME" > /dev/null 2>&1 &
    echo $! > "$SCRIPT_DIR/$PID_FILE"
    cd "$SCRIPT_DIR"
    echo "服务已启动 (PID: $(cat "$PID_FILE"))"
}

stop() {
    if [ ! -f "$PID_FILE" ]; then
        echo "PID 文件不存在，服务未运行"
        return 0
    fi
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "停止服务 (PID: $PID) ..."
        kill "$PID"
        sleep 1
        kill -0 "$PID" 2>/dev/null && kill -9 "$PID"
        echo "服务已停止"
    else
        echo "进程 $PID 不存在"
    fi
    rm -f "$PID_FILE"
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
