String.prototype.format = function(args) {
    var result = this;
    if (arguments.length < 1) {
        return result;
    }

    var data = arguments;       //如果模板参数是数组
    if (arguments.length == 1 && typeof (args) == "object") {
        //如果模板参数是对象
        data = args;
    }
    for (var key in data) {
        var value = data[key];
        if (undefined != value) {
            result = result.replace("{" + key + "}", value);
        }
    }
    return result;
}

html_shadow_port = `
    <tr>
    <td>{port}</td>
    <td>{service}</td>
    </tr>
`

function router() {
    routes = ["#/seek", "#/shadow", "#/attack", "#/assassinate"]
    hash = window.location.hash;
    if (hash == "#") {
        for (var i=0; i<routes.length; i++) {
            $(routes[i].replace("/", "")).transition('hide');
        }
    }
    for (var i=0; i<routes.length; i++) {
        if (routes[i]!=hash) {
            $(routes[i].replace("/", "")).transition('hide');
        }
    }
    animates = ["vertical flip", "drop'", "slide down", "slide left", "slide right", "fade up", "fade right", "zoom"]
    j = Math.floor(Math.random() * Math.floor(animates.length));
    $(hash.replace("/","")).transition(animates[j], '400ms');
    
}

$(window).bind('hashchange', function() {
    router();
});

$(document).ready(function(){
    router();
    $("#bt-set-target").click(function(){
        $.ajax({
            url:" /api/target",
            type: "POST",
            data: "target="+$("#target").val(),
            dataType: "JSON",
        }).done(function (result) {
            if (result.status == "success") {
                
            } else {
                alert(result.message)
            }
        }).fail(function(result){
            alert("some thing error")
        });
        $(location).attr('href', '/#/shadow');
    });

    $("#start-shadow").click(function(){
        $.ajax({
            url: "/api/info/port",
            type: "GET",
            dataType: "JSON",
        }).done(function (result) {
            ports = result.data.sort(function sequence(a,b){
                return a - b;
            });
            if (ports.length>0) {
                for (var i=0; i<ports.length;i++) {
                    h = html_shadow_port.format({"port":ports[i],"service":"open"})
                    $("#port-table").append(h)
                }
            }
        })
    })

});