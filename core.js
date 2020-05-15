const preStart = require('./src/prestart')

// 初始化程序
preStart.Run()

// 预先执行一次检测
require('./src/check').checkUpdate()

// 启动定时任务
require('./src/cron')
