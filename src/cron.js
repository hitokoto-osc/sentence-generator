'use strict'
const schedule = require('node-schedule')
const {checkUpdate} = require('./check')

schedule.scheduleJob('checkSentencesUpdates', '1 * * * *', checkUpdate)
