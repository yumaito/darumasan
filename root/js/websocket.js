var uri = "ws://" + location.host + "/monitor";
var webSocket = null;

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
    var json_before = JSON.parse(event.data)
    var json = JSON.stringify(json_before, null, "    ");
    $("#Json pre").text(json);
}

function onError(event) {
}

function onClose(event) {
}

// 初期処理
$(init);
