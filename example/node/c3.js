var net = require('net');

const methods = {}

var server = net.createServer(function(socket) {
  socket.on('data', x => {
    try {
      const json = JSON.parse(x.toString('utf8'))
      const result = methods[json.method](...json.args)
      console.log(result)
      const res = {
        result
      }
      socket.write(`${JSON.stringify(res)}\r\n`)
      //socket.pipe(socket);
    } catch(err) {
      console.log(err)
    }
  })
});

server.listen(9999, '0.0.0.0');

module.exports = {
  register: function(methodName, fn) {
    methods[methodName] = fn
  }
}


/*
var client = new net.Socket();
client.connect(9999, '127.0.0.1', function() {
	console.log('Connected');
	client.write('Hello, server! Love, Client.');
});

client.on('data', function(data) {
	console.log('Received: ' + data);
	client.destroy(); // kill client after server's response
});

client.on('close', function() {
	console.log('Connection closed');
})
*/
