#!/bin/sh

/www/main artisan key:generate
/www/main artisan jwt:secret
/www/main artisan migrate
/www/main artisan app:init

# 启动主程序
exec /www/main