#!/bin/bash
set -e

echo "=== 构建 copilotlens ==="
go build -o bin/copilotlens .

rm -rf bin/web
cp -r web bin/
cp -r data bin/
mkdir -p bin/conf
[ -f bin/conf/config.toml ] || cp toml/config_template.toml bin/conf/config.toml

echo "构建完成: bin/copilotlens"
