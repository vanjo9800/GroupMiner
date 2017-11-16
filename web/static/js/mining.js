function checkFieldUpdate(field, newText) {
    if ($(field).text() != newText) {
        $(field).text(newText);
    }
}

function updateInfo(data) {
    data = JSON.parse(data);
    currentIPs = []
    for (var i = 0; i < data.length; i++) {
        var device = $('[data="' + data[i].IP + '"]');
        currentIPs.push(data[i].IP)
        if (!device.length) {
            device = $("#card-template").clone();
            $(device).attr('id', "")
            $(device).show();
            $("#card-container").append(device);
            $(device).attr("data", data[i].IP)
            $(device).find(".card-header").text(data[i].Name);
            var startForm = $(device).find(".startMiner")
            $(startForm).attr("action", "/start/" + i)
            $(startForm).submit(function (e) {
                e.preventDefault();
                $.ajax({
                    url: $(this).attr("action"),
                    type: 'post',
                    data: $(this).serialize(),
                    success: function () {
                    }
                });
            });
            var stopForm = $(device).find(".stopMiner")
            $(stopForm).attr("action", "/stop/" + i)
            $(stopForm).submit(function (e) {
                e.preventDefault();
                $.ajax({
                    url: $(this).attr("action"),
                    type: 'post',
                    data: $(this).serialize(),
                    success: function () {
                    }
                });
            });
        }
        mining = data[i].State.MiningParams
        system = data[i].State.SystemParams
        if (system.CurrentState == false) {
            $(device).find(".start").show();
            $(device).find(".running").hide();
        } else {
            $(device).find(".start").hide();
            $(device).find(".running").show();
            checkFieldUpdate($(device).find(".poolURL"), mining.PoolURL);
            checkFieldUpdate($(device).find(".username"), mining.Username);
            checkFieldUpdate($(device).find(".threads"), mining.Threads);
            checkFieldUpdate($(device).find(".cpuUsage"), mining.CPUUse);
            checkFieldUpdate($(device).find(".cpuLoad"), system.SystemCPU);
            if (system.SystemCPU < 50.0) {
                $(device).find(".cpuLoad").css("color", "green");
            } else {
                if (system.SystemCPU < 85.0) {
                    $(device).find(".cpuLoad").css("color", "orange");
                } else {
                    $(device).find(".cpuLoad").css("color", "red");
                }
            }
        }
    }
    $('#card-container').children().each(function() { 
        var divIP=$(this).attr("data");
        if(currentIPs.indexOf(divIP)==-1){
            $(this).remove();
        }
    });
}
