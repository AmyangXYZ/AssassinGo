String.prototype.format = function(args) {
    var result = this;
    if (arguments.length < 1) {
        return result;
    }

    var data = arguments; 
    if (arguments.length == 1 && typeof (args) == "object") {
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

var routes = ["#home", "#seek", "#shadow", "#attack", "#assassinate"];

var animates = ["horizontal flip", "vertical flip", "drop", "zoom", 
        "slide down", "slide up", "slide left", "slide right", 
        "fade up", "fade left","fade down", "fade right"];

var iconMap = {'#home':'home', 
            '#seek':'search',
            '#shadow':'spy',
            '#attack':'block layout',
            '#assassinate':'smile'};

var html_shadow_port = `
        <tr>
        <td>{port}</td>
        <td>{service}</td>
        </tr>
        `

function router() {
    hash = window.location.hash;
    if (hash == "") {
        $(location).attr('href', '/#home');
    }
    for (var i=0; i<routes.length; i++) {
        if (routes[i]!=hash) {
            $(routes[i]).transition('hide');
        }
    }
    j = Math.floor(Math.random() * Math.floor(animates.length));
    $(hash).transition(animates[j], '400ms');
}

function changeColor() {
    hash = window.location.hash;
    $(hash+"-sd").attr("class", "orange "+iconMap[hash]+" icon");
    for (var i=0; i<routes.length; i++) {
        if (routes[i]!=hash) {
            $(routes[i]+"-sd").attr("class", "teal "+iconMap[routes[i]]+" icon");
        }
    }
    
}

$(window).bind('hashchange', function() {
    router();
    changeColor();
});

$(document).ready(function(){
    router();
    changeColor();
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
        $(location).attr('href', '/#shadow');
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