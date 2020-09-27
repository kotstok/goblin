window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    document.getElementById("form").addEventListener("submit", function (e) {
        e.preventDefault();

        if (!conn || !msg.value) {
            return false;
        }

        conn.send(JSON.stringify({
            "message": msg.value,
            "username": "mycustomtoken"
        }));

        msg.value = "";

        return false;
    });

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");

        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };

        conn.onmessage = function (evt) {
            let data = JSON.parse(evt.data);

            let item = document.createElement("div");
            item.innerHTML = "<b>"+data.username+"</b>: " + data.message;
            appendLog(item);
        };

    } else {
        let item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};
