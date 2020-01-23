import { writable } from 'svelte/store';

const serverAddress = 'ws://localhost:5003/echo';

export const connected = writable(false);

export const messages = writable([])

export const ConnectionHandler = (function () {

    let ws = null;
    const messageCallbacks = [];

    let userdata = {
        id: null,
        username: null
    };

    function parseMessage(msgData) {
        const message = JSON.parse(msgData);
        return {
            own: message.uuid === userdata.id,
            user: message.user,
            text: message.text
        };
    }

    function connect(cb) {
        ws = new WebSocket(serverAddress);
        ws.addEventListener('open', function () {
            if (typeof cb === 'function') {
                cb();
            }

            connected.set(true);
            messages.set([]);
        });
    }

    function disconnect(cb) {
        ws.addEventListener('close', function () {
            if (typeof cb === 'function') {
                cb();
            }

            connected.set(false);
        });
        ws.close();

        messageCallbacks.length = 0;
        userdata.id = null;
    }

    function handshake(username) {
        userdata.username = username;

        ws.send(username);

        ws.addEventListener('message', function (event) {
            userdata.id = +event.data;

            ws.addEventListener('message', function (event) {
                messages.update(value => [...value, parseMessage(event.data)]);

                messageCallbacks.forEach(cb => {
                    cb(event.data);
                });
            });
        }, {once: true});
    }

    function sendMessage(msg) {
        ws.send(msg);
    }

    function addMessageCallback(cb) {
        if (typeof cb === 'function') {
            messageCallbacks.push(cb);
        }
    }

    function getId() {
        return ws.readyState === WebSocket.OPEN ? userdata.id : null;
    }

    function getUsername() {
        return ws.readyState === WebSocket.OPEN ? userdata.username : null;
    }

    return {
        connect,
        disconnect,
        handshake,
        sendMessage,
        addMessageCallback,
        getId,
        getUsername
    }
}) ();