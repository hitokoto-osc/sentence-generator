'use strict'
const database = require('./database')
const winston = require('winston')
const nconf = require('nconf')
// const diff = require('diff')
const _ = require('lodash')
const semver = require('semver')
const colors = require('colors/safe')

const fs = require('fs')
const path = require('path')

const workdir = path.join(__dirname, '../', nconf.get('workdir'))
console.log(workdir)
const git = require('simple-git/promise')(workdir)
/**
   * .catch(e => {
    console.error(colors.red(e.stack))
    winston.error('无法初始化 GIT 工作目录，程序结束。')
    process.exit(1)
  })
  */

// 缓存记录
let sentenceCount = 0
let currentCategoriesList = []

// 检测版控文件是否存在
let versionData
const versionFile = path.join(workdir, './version.json')
const categoriesData = {}
if (fs.existsSync(versionFile)) {
  versionData = JSON.parse(fs.readFileSync(versionFile, { encoding: 'utf8' }))
  // 读取分类列表
  currentCategoriesList = JSON.parse(fs.readFileSync(path.join(workdir, versionData.categories.path)))
  // 迭代读取句子分类数据
  for (const category of currentCategoriesList) {
    categoriesData[category.key] = JSON.parse(fs.readFileSync(path.join(workdir, category.path)))
    console.log(categoriesData[category.key].length)
  }
} else {
  versionData = {
    protocol_version: nconf.get('version'),
    bundle_version: '1.0.0',
    updated_at: 0,
    categories: {
      path: './categories.json',
      timestamp: 0
    },
    sentences: []
  }
}

async function checkUpdate () {
  try {
    winston.verbose('开始执行同步句子操作...')
    let isUpdate = false
    const conn = await database.getConnection()
    // 获得一言分类列表
    let [rows] = await conn.query('SELECT * FROM `hitokoto_categories`')
    rows.map(v => {
      v.path = `./sentences/${v.key}.json`
      return v
    })
    if (currentCategoriesList.length === 0 || JSON.stringify(currentCategoriesList) !== JSON.stringify(rows)) {
      currentCategoriesList = rows
      isUpdate = true
      // 更新数据
      fs.writeFileSync(path.join(workdir, versionData.categories.path), JSON.stringify(currentCategoriesList, null, 2))
      versionData.categories.timestamp = Date.now()
      fs.writeFileSync(path.join(workdir, './version.json'), JSON.stringify(versionData, null, 2))
    }

    // 获得一言句子集合
    [rows] = await conn.query('SELECT * FROM `hitokoto_sentence`')
    if (rows.length !== sentenceCount) { // TODO:  寻找更经济、全面的比对算法。尽管如果碰巧删除的句子数目与新增句子数目相等时此规则将不合预期忽略掉变化，但目前这是最经济的做法。
    // 建立数据集缓存
      const tmp = {}
      for (const sentence of rows) {
        sentence.length = sentence.hitokoto.length
        if (!tmp[sentence.type]) {
          tmp[sentence.type] = [sentence]
        } else {
          tmp[sentence.type].push(sentence)
        }
      }
      // 按分类比对更新
      for (const category of currentCategoriesList) {
        const {key} = category
        if (!categoriesData[key] || categoriesData[key].length === 0 || JSON.stringify(categoriesData[key]) !== JSON.stringify(tmp[key])) {
          isUpdate = true
          categoriesData[key] = tmp[key]
          let categoryVersion = _.find(versionData.sentences, { key })
          if (!categoryVersion) {
            categoryVersion = {
              name: category.name,
              key: category.key,
              path: `./sentences/${key}.json`,
              timestamp: Date.now()
            }
            versionData.sentences.push(categoryVersion)
          } else {
            categoryVersion.timestamp = Date.now()
          }
          if (!fs.existsSync(path.join(workdir, categoryVersion.path, '../'))) {
            fs.mkdirSync(path.join(workdir, categoryVersion.path, '../'))
          }
          fs.writeFileSync(path.join(workdir, categoryVersion.path), JSON.stringify(tmp[key], null, 2))
          fs.writeFileSync(path.join(versionFile), JSON.stringify(versionData, null, 2))
        }
      }
    }
    // 调用 Git，生成新的版本号
    if (isUpdate) {
      winston.verbose('文件生成完毕，开始发布 GIT 版本。')
      versionData.bundle_version = semver.inc(versionData.bundle_version, 'patch')
      versionData.updated_at = Date.now()
      fs.writeFileSync(path.join(versionFile), JSON.stringify(versionData, null, 2))

      // GIT 操作
      await git.add('*')
      await git.commit(`build(auto): v${versionData.bundle_version}`)
      await git.addAnnotatedTag(`v${versionData.bundle_version}`, `v${versionData.bundle_version}. Auto released by hitokoto-sentence-generator.`)
      await Promise.all([
        git.push(),
        git.pushTags()
      ])
    } else {
      winston.verbose('文件内容无需更新。')
    }
    winston.verbose('更新操作执行完毕。')
  } catch (e) {
    console.error(colors.red(e.stack))
    winston.error('执行同步操作时发生错误！')
  }
}

module.exports = {
  checkUpdate
}
