// ======================================
// WebSocket Connection
// ======================================

let socket = null;

// Connect to the websocket server.
function connectWebSocket() {

    const protocol =
        window.location.protocol === "https:"
            ? "wss://"
            : "ws://";

    socket = new WebSocket(
        protocol + window.location.host + "/ws"
    );

    socket.onopen = () => {
        console.log("✅ Connected to WebSocket");
    };

    socket.onclose = () => {
        console.log("❌ WebSocket disconnected");

        // Reconnect after 2 seconds.
        setTimeout(connectWebSocket, 2000);
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    socket.onmessage = (event) => {

        const wsMessage = JSON.parse(event.data);

        switch (wsMessage.type) {

            case "chat":
                if (typeof appendMessage === "function") {
                    appendMessage(wsMessage.data);
                }
                break;

            case "presence":
                if (typeof updatePresence === "function") {
                    updatePresence(wsMessage.data);
                }
                break;

            case "typing":
                if (typeof showTyping === "function") {
                    showTyping(wsMessage.data);
                }
                break;

            case "read_receipt":
                if (typeof updateReadReceipt === "function") {
                    updateReadReceipt(wsMessage.data);
                }
                break;

            default:
                console.log("Unknown message:", wsMessage);
        }
    };
}

// ======================================
// Send a chat message
// ======================================

function sendChat(message, receiverID = null) {

    if (!socket || socket.readyState !== WebSocket.OPEN) {
        return;
    }

    socket.send(
        JSON.stringify({
            type: "chat",
            data: {
                receiver_id: receiverID,
                message: message,
                message_type: "text"
            }
        })
    );
}

// ======================================
// Send typing notification
// ======================================

function sendTyping(receiverID = null, isTyping = true) {

    if (!socket || socket.readyState !== WebSocket.OPEN) {
        return;
    }

    socket.send(
        JSON.stringify({
            type: "typing",
            data: {
                receiver_id: receiverID,
                is_typing: isTyping
            }
        })
    );
}

// ======================================
// Read receipt
// ======================================

function sendReadReceipt(messageID) {

    if (!socket || socket.readyState !== WebSocket.OPEN) {
        return;
    }

    socket.send(
        JSON.stringify({
            type: "read_receipt",
            data: {
                message_id: messageID
            }
        })
    );
}

// ======================================
// Connect automatically
// ======================================

window.addEventListener("load", () => {
    connectWebSocket();
});