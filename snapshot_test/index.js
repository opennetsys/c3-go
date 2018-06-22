var state = {}

console.log('modified time', new Date())

setTimeout(() => {
  state['foo'] = 'bar'
  console.log(state)
}, 1e3)
