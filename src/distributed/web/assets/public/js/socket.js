(function() {
  var socket = new WebSocket("ws://localhost:3000/ws");
  var sources = [];

  socket.addEventListener("message", function(e) {
    var msg = JSON.parse(e.data);
    switch (msg.type) {
      case "source":
        if (
          !sources.some(function(source) {
            return source == msg.data.name;
          })
        ) {
          sources.push(msg.data.name);
          createChart($("#chartContainer"), msg.data);
        }

        break;
      case "reading":
        updateChart(msg.data);
        break;
    }
  });

  socket.addEventListener("open", function(e) {
    socket.send(
      JSON.stringify({
        type: "discover"
      })
    );
  });
})();
