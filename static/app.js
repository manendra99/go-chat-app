let ws;
let username;

let messageInput = document.getElementById("message-input");
let messageDiv = document.getElementById("messages");
let chatWindow = document.getElementById("chat-window");
let usernameSection = document.getElementById("username-section");
let leaveButton = document.getElementById("leave-btn");

// to join the chat
function joinChat() {
  username = document.getElementById("username").value;

  // to check if there is any username or not
  if (!username) {
    alert(" Enter a user name ");
    return;
  }

  // to connect to web-socket

  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = function () {
    console.log("connect to the chat");

    // to remove the username input
    usernameSection.style.display = "none";

    // to show the chat window
    chatWindow.style.display = "block";

    ws.send(JSON.stringify({ action: "join", username: username }));

    ws.onmessage = function (event) {
      let message = event.data;

      let messageElement = document.createElement("div");

      messageElement.textContent = message;

      messageDiv.appendChild(messageElement);

      messageDiv.scrollTop = messageDiv.scrollHeight;
    };

    ws.onclose = function () {
      console.log("Dissconnect from the chat");

      resetChatUI();
    };

    ws.onerror = function () {
      console.error("error: ", error);
    };
  };
}

//to send the message

function sendMessage() {
    let message = messageInput.value.trim();
    if (message) {
      ws.send(JSON.stringify({ action: "message", message: message }));
      messageInput.value = "";                                                                                                                                                                                                                                                                                                                                                                                                               
    }
  }
  
// to leave the chat
function leaveChat() {
  ws.send(JSON.stringify({ action: "leave", username: username }));

  chatWindow.style.display = "none";

  usernameSection.style.display = "block";
}

function resetChatUI() {
    usernameSection.style.display = "block"; 
    chatWindow.style.display = "none";
    messageDiv.innerHTML = "";
    messageInput.value = "";
  }
  