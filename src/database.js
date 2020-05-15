const mysql = require('mysql2/promise')
const nconf = require('nconf')
const winston = require('winston')
const colors = require('colors/safe')

let pool = null

async function createConnectionsPool () {
  if (!nconf.get('mysql')) {
    winston.error('数据库配置不存在，程序结束。')
    process.exit(1)
  }
  try {
    pool = await mysql.createPool({
      host: nconf.get('mysql:host') || 'localhost',
      user: nconf.get('mysql:user') || 'root',
      password: nconf.get('mysql:password') || null,
      port: nconf.get('mysql:port') || 3306,
      charset: nconf.get('mysql:charset') || 'utf8mb4',
      database: nconf.get('mysql:database') || null
    })
    await pool.query('SELECT 1')
  } catch (e) {
    console.error(colors.red(e.stack))
    winston.error('与数据库建立连接时发生错误，程序结束。')
    process.exit(1)
  }
}

async function getConnection () {
  if (!pool) {
    await createConnectionsPool()
  }
  return pool
}

module.exports = {
  getConnection
}
