var state = {}

console.log('modified time', new Date())

setTimeout(() => {
  state['foo'] = 'bar'
  console.log(state)
  console.log('time', new Date())
}, 1e3)

setTimeout(() => {
}, 100e3)
