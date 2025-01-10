console.log("javascript loaded")

const form = document.getElementById("form");
const messageArea = document.getElementById("message-area");
const selectModel = document.getElementById("select-model");

function appendMessage(message, isUser) {
    const messageElement = document.createElement("div");
    console.log(`isuser:${isUser}`)
    messageElement.className = isUser ? "message user-message" : "message bot-message";
    messageElement.textContent = message;
    messageArea.appendChild(messageElement);
    messageArea.scrollTop = messageArea.scrollHeight;
}

function getModelText() {
    const modelText = selectModel.options[selectModel.selectedIndex].value;
    console.log("model text: ", modelText)
    return modelText
}

function markdownMessageArea() {
    var messageAreas = document.getElementById('message-area');
    
    for (var i = 0; i < messageAreas.length; i++) {
        var content = messageAreas[i].innerHTML;
        var html = marked.parse(content);
        messageAreas[i].innerHTML = html;
    }
};

form.addEventListener("submit", (e) => {
    e.preventDefault()
    const messageInput = document.getElementById("message");
    const message = messageInput.value;
    const model = getModelText()

    appendMessage(`user: ${message}`, true)
    console.log("hello")
    console.log(`sending ${message} to the server`)
    messageInput.value = ""

    fetch("/send-message", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ message: message, model: model, user: "pepe" }),
    })
    .then(response => response.json())
    .then(data => {
        console.log(data.output)
        console.log(data.latency)
        appendMessage(`bot: ${data.output}`, false);
        appendMessage(`latency: ${data.latency} ms`, false);
        markdownMessageArea()
    })
    .catch((error) => {
        console.error(`Error: ${error}`)
    })
})
