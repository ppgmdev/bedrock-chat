console.log("javascript loaded")

const form = document.getElementById("form");
const messageArea = document.getElementById("message-area");

function appendMessage(message, isUser = false) {
    const messageElement = document.createElement("div");
    messageElement.className = isUser ? "message user-message" : "message bot-message";
    messageElement.textContent = message;
    messageArea.appendChild(messageElement);
    messageArea.scrollTop = messageArea.scrollHeight;
}

form.addEventListener("submit", (e) => {
    e.preventDefault()
    const messageInput = document.getElementById("message");
    const message = messageInput.value;
    appendMessage(`user: ${message}`)
    console.log(`sending ${message} to the server`)

    fetch("/send-message", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ message: message, user: "pepe" }),
    })
    .then(response => response.text())
    .then(data => {
        console.log(`Success: ${data}`)
        appendMessage(`bot: ${data}`)
    })
    .catch((error) => {
        console.error(`Error: ${error}`)
    })
})
