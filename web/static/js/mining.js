$(document).ready(function () {
    $('#minerStart').submit(function (e) {
        e.preventDefault();
        $.ajax({
            url: '/start',
            type: 'post',
            data: $('#minerStart').serialize(),
            success: function () {
            }
        });
    });
    $('#minerStop').submit(function (e) {
        e.preventDefault();
        $.ajax({
            url: '/stop',
            type: 'post',
            data: $('#minerStop').serialize(),
            success: function () {
            }
        });
    });
});
function updateInfo(data) {
    data = JSON.parse(data);
        $("#card-container").empty();
    for (var i = 0; i < data.length; i++) {
        var device = $("#card-template").clone();
        $(device).attr('id',"")
        $(device).show();
        $("#card-container").append(device);
        $(device).find(".card-header").text(data[i].Name);
        $(device).find("#start").attr('id', 'start' + i);
        $(device).find("#running").attr('id', 'running' + i);
        $(device).find("#startMiner").attr('id', 'startMiner' + i);
        $(device).find("#stopMiner").attr('id', 'stopMiner' + i);
        $(device).find("#poolURL").attr('id',"poolURL"+i)
        $(device).find("#username").attr('id',"username"+i)
        $(device).find("#password").attr('id',"password"+i)
        $(device).find("#threads").attr('id',"threads"+i)
        $(device).find("#cpuUsage").attr('id',"cpuUsage"+i)
        $(device).find("#cpuLoad").attr('id',"cpuLoad"+i)
        mining = data[i].State.MiningParams
        system = data[i].State.SystemParams
        if (system.CurrentState == false) {
            $("#start" + i).show();
            $("#running" + i).hide();
        } else {
            $("#start" + i).hide();
            $("#running" + i).show();
            $("#poolURL" + i).text(mining.PoolURL);
            $("#username" + i).text(mining.Username);
            $("#threads" + i).text(mining.Threads);
            $("#cpuUsage" + i).text(mining.CPUUse);
            $("#cpuLoad" + i).text(system.SystemCPU);
            if (system.SystemCPU < 50.0) {
                $("#cpuLoad" + i).css("color", "green");
            } else {
                if (system.SystemCPU < 85.0) {
                    $("#cpuLoad" + i).css("color", "orange");
                } else {
                    $("#cpuLoad" + i).css("color", "red");
                }
            }
        }
    }
}