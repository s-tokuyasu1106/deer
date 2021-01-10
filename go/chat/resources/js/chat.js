$(function(){
    var socket = null;
    var msgBox = $("#chatbox textarea");
    var message = $("#message");
    var host = location.hostname;

    $("#chatbox").submit(function(){
        if(!msgBox.val()) return false;
        if(!socket){
            alert("エラー：websocketに接続できませんでした");
            return false;
        }
        socket.send(msgBox.val());
        msgBox.val("");
        return false;
    });
    if(!window["WebSocket"]){
        alert("エラー：websocketに対応していないブラウザです")
    }else{
        socket = new WebSocket("ws://"+host+":8080/room");
        socket.onclose = function() {
            alert("接続が終了しました")
        }
        socket.onmessage = function(e) {
            message.append($("<li>").text(e.data));
        }
    }
});