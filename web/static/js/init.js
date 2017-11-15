$(document).ready(function () {
    setInterval(function () {
        $.ajax({
            url: '/status',
            type: 'get',
            success: function (data) {
                updateInfo(data);
            }
        })
    }, 1000)
    $("#card-template").hide()
})