var state = {}

console.log("unmodified time", new Date())

process.env['FAKETIME'] = '2017-01-01 00:00:00'

console.log("modified time", new Date())

setTimeout(() => {
  state['foo'] = 'bar'
  console.log(state)
}, 1e3)
