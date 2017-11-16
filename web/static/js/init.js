$(document).ready(function () {
	var ws = new WebSocket('ws://' + window.location.host + '/ws');
	ws.addEventListener('message', function (e) {
		updateInfo(e.data);
	});
	$("#card-template").hide()
})
