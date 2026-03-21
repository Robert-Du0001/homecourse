#!/bin/sh

/www/main artisan migrate
/www/main artisan app:init

# 启动主程序
exec /www/main