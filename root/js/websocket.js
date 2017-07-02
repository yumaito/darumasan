var uri = "ws://" + location.host + "/monitor";
var webSocket = null;
var data;

function init() {
    open()
}

// 接続
function open() {
    if (webSocket == null) {
        webSocket = new WebSocket(uri);
        webSocket.onopen = onOpen;
        webSocket.onmessage = onMessage;
        webSocket.onclose = onClose;
        webSocket.onerror = onError;
    }
}

// 接続イベント
function onOpen(event) {
}

function onMessage(event) {
    data = JSON.parse(event.data)
    data.clients.sort(function(a,b) {
        if (a < b) return -1;
        if (a > b) return 1;
        return 0;
    })
    data.dead_clients.sort(function(a,b) {
        if (a < b) return -1;
        if (a > b) return 1;
        return 0;
    })
    var disp_json = JSON.stringify(data, null, "    ");
    $("#Json pre").text(disp_json);
}

function onError(event) {
}

function onClose(event) {
}

// 初期処理
$(init);
