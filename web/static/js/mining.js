function checkFieldUpdate(field,newText){
    if($(field).text()!=newText){
        $(field).text(newText);
    }
}

function updateInfo(data) {
    data = JSON.parse(data);
    data = JSON.parse(data);
    for (var i = 0; i < data.length; i++) {
        console.log(data[i],data.length)
        var device = $('[data="'+data[i].IP+'"]');
        if (!device.length) {
            device = $("#card-template").clone();
            $(device).attr('id',"")
            $(device).show();
            $("#card-container").append(device);
            $(device).attr("data",data[i].IP)
            $(device).find(".card-header").text(data[i].Name);
            $(device).find(".startMiner").submit(function (i,device) {
                $(device).find(".startMiner").preventDefault();
                $.ajax({
                    url: '/start/'+i,
                    type: 'post',
                    data: $(device).find('.startMiner').serialize(),
                    success: function () {
                    }
                });
            }(i,device);
            $(device).find('.stopMiner').submit(function (i,device) {
                $(device).find(".stopMiner").preventDefault();
                $.ajax({
                    url: '/stop/'+i,
                    type: 'post',
                    data: $(device).find('.stopMiner').serialize(),
                    success: function () {
                    }
                });
            }(i,device);
        }
        mining = data[i].State.MiningParams
        system = data[i].State.SystemParams
        if (system.CurrentState == false) {
            $(device).find(".start").show();
            $(device).find(".running").hide();
        } else {
            $(device).find(".start").hide();
            $(device).find(".running").show();
            checkFieldUpdate($(device).find(".poolURL"),mining.PoolURL);
            checkFieldUpdate($(device).find(".username"),mining.Username);
            checkFieldUpdate($(device).find(".threads"),mining.Threads);
            checkFieldUpdate($(device).find(".cpuUsage"),mining.CPUUse);
            checkFieldUpdate($(device).find(".cpuLoad"),system.SystemCPU);
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
}
