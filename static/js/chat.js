// ======================================
// Current User Details
// ======================================

const currentUserEl = document.getElementById("currentUser");
const currentUserID = currentUserEl ? parseInt(currentUserEl.getAttribute("data-id")) : 0;

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

// Map of user_id -> username for online users
const onlineUsers = new Map();

const publicLobbyBtn = document.getElementById("publicLobbyBtn");
if (publicLobbyBtn) {
    publicLobbyBtn.addEventListener("click", () => {
        activeReceiver = null;
        document.querySelector(".chat-header h2").textContent = "Public Chat Room";
    });
}

// ======================================
// Send Message
// ======================================

messageForm.addEventListener("submit", function (e) {
    e.preventDefault();

    const text = messageInput.value.trim();

    if (text === "") {
        return;
    }

    const sent = sendChat(text, activeReceiver);
    if (sent) {
        messageInput.value = "";
    } else {
        alert("Cannot send message: Connection is offline. Please wait for reconnection.");
    }
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

    if (message.sender_id === currentUserID) {
        wrapper.classList.add("sent");
    } else {
        wrapper.classList.add("received");
    }

    wrapper.innerHTML = `
        <div class="message-author">
            ${message.sender_username || "Unknown User"}
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
// ======================================
// Online Users
// ======================================

function updatePresence(presence) {

    if (presence.online) {
        // Exclude current user from their own online users list
        if (presence.user_id !== currentUserID) {
            onlineUsers.set(presence.user_id, presence.username);
        }
    } else {
        onlineUsers.delete(presence.user_id);
    }

    const list = document.getElementById("onlineUsers");
    if (!list) return;

    list.innerHTML = "";

    onlineUsers.forEach((username, id) => {
        const div = document.createElement("div");
        div.className = "user-item";
        div.innerHTML = `🟢 ${username}`;

        div.onclick = () => {
            activeReceiver = id;
            document.querySelector(".chat-header h2").textContent =
                "Chat with " + username;
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

window.addEventListener("load", scrollToBottom);