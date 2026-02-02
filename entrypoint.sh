#!/bin/sh

# 执行数据库迁移
echo "正在执行数据库迁移..."
/www/main artisan migrate

# 启动主程序
echo "正在启动主程序..."
exec /www/main
