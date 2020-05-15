// 引入模块
const nconf = require('nconf')
const winston = require('winston')
const path = require('path')
const fs = require('fs')

const pkg = require('../package.json')

function printCopyright () {
  const colors = require('colors/safe')
  const date = new Date()
  console.log(colors.bgBlue(colors.black(' ' + pkg.name + ' v' + pkg.version + ' © ' + date.getFullYear() + ' All Rights Reserved. ')) + '   ' + colors.bgRed(colors.black(' Powered by minded-heart ')))
  console.log('')
  console.log(colors.bgCyan(colors.black(' 我们一路奋战，不是为了改变世界，而是为了不让世界改变我们。 ')))
}

function loadConfig (configFile) {
  winston.verbose('* using configuration stored in: %s', configFile)
  nconf.argv().env()
  // 检测配置文件是否存在
  if (!fs.existsSync(configFile)) {
    winston.error('配置文件不存在，程序初始化失败！')
    process.exit(1)
  }
  nconf.file({
    file: configFile,
    format: require('nconf-yaml')
  })

  nconf.defaults({
    base_dir: path.join(__dirname, '../'),
    version: pkg.version
  })
}

function initWinston () { // 初始化 Winston
  winston.remove(winston.transports.Console)
  winston.add(winston.transports.Console, {
    colorize: true,
    timestamp: function () {
      var date = new Date()
      return date.toISOString() + ' [' + global.process.pid + ']'
    },
    level: nconf.get('log_level') || (global.env === 'production' ? 'info' : 'verbose')
  })
}

function Run (configFile) {
  printCopyright()
  if (!configFile) {
    configFile = path.join(__dirname, '../', 'config.yml')
  }
  initWinston()
  loadConfig(configFile)
  initWinston()
}

module.exports = {
  Run
}
