const c3 = require('./c3')

var items = {}

c3.register('setItem', (key, value) => {
  items[key] = value
})

c3.register('getItem', (key) => {
  return items[key]
})
