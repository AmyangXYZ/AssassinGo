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

var html_shadow_url = `
        <tr>
        <td>{url}</td>
        </tr>
`

var html_shadow_email = `
        <tr>
            <td>
                <i class="mail icon"></i>
                {email}
            </td>
        </tr>
`
var html_attack_sqli_url = `
        <tr>
        <td>{url}</td>
        </tr>
`

var html_attack_xss_url = `
        <tr>
        <td>{url}</td>
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

function reset() {
    $("#port-table").html("");
    $("#ip").html("IP Address: ");
    $("#server").html("Web Server: ");
    $("#cms").html("CMS: ");
    $("#url-table").html("");
    $("#email-table").html("");
    $("#sqli-url-table").html("");
    $("#xss-url-table").html("");
}

function portScan() {
    $("#port-table").html("");
    $.ajax({
        url: "/api/info/port",
        type: "GET",
        dataType: "JSON",
        beforesend: $("#port-loading").show(),
    }).done(function (result) {
        $("#port-loading").hide();
        ports = result.data.ports.sort(function sequence(a,b){
            return a - b;
        });
        if (ports.length>0) {
            for (var i=0; i<ports.length;i++) {
                h = html_shadow_port.format({"port":ports[i],"service":"open"})
                $("#port-table").append(h)
            }
        }
    })
}

function basicInfo() {
    $("#ip").html("IP Address: ")
    $("#server").html("Web Server: ")
    $.ajax({
        url: "/api/info/basic",
        type: "GET",
        dataType: "JSON",
        beforesend: $("#ip-loading, #server-loading").show(),
    }).done(function (result) {
        $("#ip-loading, #server-loading").hide();
        ip=result.data.ip; server=result.data.webServer;
        $("#ip").html("IP Address: "+ip);
        $("#server").html("Web Server: "+server);
    })
}

function cmsDetect() {
    $("#cms").html("CMS: ");
    $.ajax({
        url: "/api/info/cms",
        type: "GET",
        dataType: "JSON",
        beforesend: $("#cms-loading").show(),
    }).done(function (result) {
        $("#cms-loading").hide();
        cms=result.data.cms;
        if (cms.length==0) {
            cms = "Unknown";
        }
        $("#cms").html("CMS: "+cms);
    })
}

function crawl() {
    $("#url-table").html("");
    $("#email-table").html("");
    var socket = new WebSocket("ws://localhost:8080/api/ws/crawl")
    // socket.onopen = function() {
    //     container.append("<p>Socket is open</p>");
    // };
    socket.onmessage = function (e) {
        hu = html_shadow_url.format({"url": e.data})
        $("#url-table").append(hu);
    }
    // socket.onclose = function () {
    //     container.append("<p>Socket closed</p>");
    // }
    return socket;
}

function sqliCheck() {
    $("#sqli-url-table").html("")
    $.ajax({
        url: "/api/vul/sqli",
        type: "GET",
        dataType: "JSON",
        beforesend: $("#sqli-url-loading").show(),
    }).done(function (result) {
        $("#sqli-url-loading").hide();
        sqli_urls = result.data.sqli_urls;
        if (sqli_urls.length>0) {
            for (var i=0; i<sqli_urls.length;i++) {
                h = html_attack_sqli_url.format({"url":sqli_urls[i]})
                $("#sqli-url-table").append(h)
            }
        }
    })
}

function xssCheck() {
    $("#xss-url-table").html("")
    $.ajax({
        url: "/api/vul/xss",
        type: "GET",
        dataType: "JSON",
        beforesend: $("#xss-url-loading").show(),
    }).done(function (result) {
        $("#xss-url-loading").hide();
        xss_urls = result.data.xss_urls;
        if (xss_urls.length>0) {
            for (var i=0; i<xss_urls.length;i++) {
                h = html_attack_xss_url.format({"url":xss_urls[i]})
                $("#xss-url-table").append(h)
            }
        }
    })
}

$(document).ready(function(){
    router();
    changeColor();
    $("#bt-set-target").click(function(){
        reset();
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
        basicInfo();
        portScan();
        cmsDetect();
        crawl();
    });

    $("#start-attack").click(function(){
        sqliCheck();
        xssCheck();
    })
});