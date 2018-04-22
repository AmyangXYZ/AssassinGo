// AJAX API

function TestNewAssassinHandler() {
    $.ajax({
        url: "/api/assassin",
        type: "POST",
        data: "target=207.246.110.215",
        dataType: "JSON",
    }).done(function (result) {
        console.log(result.data)
    })
}

function TestNewAssassinDadHandler() {
    $.ajax({
        url: "/api/assassin-dad",
        type: "POST",
        data: "targets=207.246.110.215,47.94.136.141",
        dataType: "JSON",
    }).done(function (result) {
        console.log(result.data)
    })
}

function TestBasicInfoHandler() {
    $.ajax({
        url: "/api/info/basic",
        type: "GET",
        dataType: "JSON",
    }).done(function (result) {
        console.log(result.data)
    })
}

function TestCMSDetectHandler() {
    $.ajax({
        url: "/api/info/cms",
        type: "GET",
        dataType: "JSON",
    }).done(function (result) {
        console.log(result.data)
    })
}

function TestHoneypotHandler() {
    $.ajax({
        url: "/api/info/honeypot",
        type: "GET",
        dataType: "JSON",
    }).done(function (result) {
        console.log(result.data)
    })
}

function TestWhoisHandler() {
    $.ajax({
        url: "/api/info/whois",
        type: "GET",
        dataType: "JSON",
    }).done(function (result) {
        console.log(result.data)
    })
}

function TestGetPoCListHandler() {
    $.ajax({
        url: "/api/poc",
        type: "GET",
        dataType: "JSON",
    }).done(function (result) {
        console.log(result.data)
    })
}

function TestRunPoCHandler() {
    $.ajax({
        url: "/api/poc/drupal-rce",
        type: "GET",
        dataType: "JSON",
    }).done(function (result) {
        console.log(result.data)
    })
}

// WebSocket API


function TestPortHandler() {
    var socket = new WebSocket("ws://207.246.110.215:8000/ws/info/port");
    socket.onopen = function(e) {
        var msg = {
            method: "tcp",
        }
        socket.send(JSON.stringify(msg))
    }
    socket.onmessage = function (e) {
        console.log(JSON.parse(e.data))
    }
    socket.onclose = function () {
        console.log("finished")
    }
}

function TestTracertHandler() {
    var socket = new WebSocket("ws://207.246.110.215:8000/ws/info/tracert")
    socket.onmessage = function (e) {
        console.log(JSON.parse(e.data))
    }
    socket.onclose = function () {
        console.log("finished")
    }
}

function TestDirbHandler() {
    var socket = new WebSocket("ws://207.246.110.215:8000/ws/info/dirb")
    socket.onopen = function(e) {
        var msg = {
            gort_count: 20,
            dict: "php",
        }
        socket.send(JSON.stringify(msg))
    }
    socket.onmessage = function (e) {
        console.log(JSON.parse(e.data))
    }
    socket.onclose = function () {
        console.log("finished")
    }
}

function TestCrawlHandler() {
    var socket = new WebSocket("ws://207.246.110.215:8000/ws/attack/crawl")
    socket.onmessage = function (e) {
        console.log(JSON.parse(e.data))
    }
    socket.onclose = function () {
        console.log("finished")
    }
}

function TestSQLiHandler() {
    var socket = new WebSocket("ws://207.246.110.215:8000/ws/attack/sqli")
    socket.onmessage = function (e) {
        console.log(JSON.parse(e.data))
    }
    socket.onclose = function () {
        console.log("finished")
    }
}

function TestXSSHandler() {
    var socket = new WebSocket("ws://207.246.110.215:8000/ws/attack/xss")
    socket.onmessage = function (e) {
        console.log(JSON.parse(e.data))
    }
    socket.onclose = function () {
        console.log("finished")
    }
}

function TestIntruderHandler() {
    var socket = new WebSocket("ws://207.246.110.215:8000/ws/attack/intrude")
    socket.onopen = function(e) {
        var msg = {
            header: `GET /$$1$$ HTTP/1.1
Host: 47.94.136.141`,
            payload: "1,2,3",
            gort_count: 5,
        }
        socket.send(JSON.stringify(msg))
    }
    socket.onmessage = function (e) {
        console.log(JSON.parse(e.data))
    }
    socket.onclose = function () {
        console.log("finished")
    }
}

function TestSeekHandler() {
    var socket = new WebSocket("ws://207.246.110.215:8000/ws/seek")
    socket.onopen = function(e) {
        var msg = {
            query: "information security",
            se: "bing",
            max_page: 1,
        }
        socket.send(JSON.stringify(msg))
    }
    socket.onmessage = function (e) {
        console.log(JSON.parse(e.data))
    }
    socket.onclose = function () {
        console.log("finished")
    }
}

function TestRunDadPoCHandler() {
    var socket = new WebSocket("ws://207.246.110.215:8000/ws/poc/drupal-rce")
    socket.onopen = function(e) {
        var msg = {
            gort_count: 10,
        }
        socket.send(JSON.stringify(msg))
    }
    socket.onmessage = function (e) {
        console.log(JSON.parse(e.data))
    }
    socket.onclose = function () {
        console.log("finished")
    }
}