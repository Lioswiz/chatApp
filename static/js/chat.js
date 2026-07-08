// ======================================
// WebSocket Connection
// ======================================

const protocol = window.location.protocol === "https:" ? "wss://" : "ws://";

const socket = new WebSocket(protocol + window.location.host + "/ws");

socket.onopen = () => {
    console.log("✅ Connected to WebSocket");
};

socket.onclose = () => {
    console.log("❌ WebSocket disconnected");
};

socket.onerror = (err) => {
    console.error("WebSocket error:", err);
};

socket.onmessage = (event) => {
    const data = JSON.parse(event.data);

    console.log("Received:", data);

    appendMessage(data);
};



// ======================================
// DOM Elements
// ======================================

const messageForm = document.getElementById("messageForm");
const messageInput = document.getElementById("messageInput");
const messages = document.getElementById("messages");
const typingIndicator = document.getElementById("typingIndicator");

// Currently selected private chat.
// null = public room
let activeReceiver = null;

// ======================================
// Send Message
// ======================================

messageForm.addEventListener("submit", function (e) {
    e.preventDefault();

    const text = messageInput.value.trim();

    if (text === "") {
        return;
    }

    sendChat(text, activeReceiver);

    messageInput.value = "";
});

// ======================================
// Typing Indicator
// ======================================

let typingTimeout;

messageInput.addEventListener("input", () => {

    sendTyping(activeReceiver, true);

    clearTimeout(typingTimeout);

    typingTimeout = setTimeout(() => {
        sendTyping(activeReceiver, false);
    }, 1000);
});

// ======================================
// Append Message
// ======================================

function appendMessage(message) {

    const wrapper = document.createElement("div");

    wrapper.classList.add("message");

    // If the backend sends current_user,
    // this will automatically style your own messages.
    if (message.current_user === true) {
        wrapper.classList.add("sent");
    } else {
        wrapper.classList.add("received");
    }

    wrapper.innerHTML = `
        <div class="message-author">
            ${message.username || "Unknown User"}
        </div>

        <div class="message-body">
            ${message.message}
        </div>

        <div class="message-time">
            ${formatTime(message.created_at)}
        </div>
    `;

    messages.appendChild(wrapper);

    scrollToBottom();
}
function sendChat(message, receiverID) {

    const payload = {
        message: message,
        message_type: "text"
    };

    if (receiverID !== null) {
        payload.receiver_id = receiverID;
    }

    socket.send(JSON.stringify(payload));
}
// ======================================
// Online Users
// ======================================

function updatePresence(users) {

    const list = document.getElementById("onlineUsers");

    list.innerHTML = "";

    users.forEach(user => {

        const div = document.createElement("div");

        div.className = "user-item";

        div.innerHTML = `
            🟢 ${user.username}
        `;

        div.onclick = () => {

            activeReceiver = user.id;

            document.querySelector(".chat-header h2").textContent =
                "Chat with " + user.username;
        };

        list.appendChild(div);
    });
}

// ======================================
// Typing Display
// ======================================

function showTyping(data) {

    if (data.is_typing) {
        typingIndicator.textContent =
            `${data.username} is typing...`;
    } else {
        typingIndicator.textContent = "";
    }
}

// ======================================
// Read Receipts
// ======================================

function updateReadReceipt(data) {

    console.log("Read receipt:", data);
}

// ======================================
// Helpers
// ======================================

function scrollToBottom() {

    messages.scrollTop = messages.scrollHeight;
}

function formatTime(timestamp) {

    if (!timestamp) return "";

    const date = new Date(timestamp);

    return date.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit"
    });
}
function sendTyping(receiverID, typing) {
    // We'll implement typing events later.
}