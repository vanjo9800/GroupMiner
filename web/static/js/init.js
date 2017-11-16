$(document).ready(function () {
	var webSocketServerLocation='ws://'+window.location.host+"/ws";
	var ws = new ReconnectingWebSocket(webSocketServerLocation);
	ws.addEventListener('message', function (e) {
		updateInfo(e.data);
	});
	$("#card-template").hide()
})
