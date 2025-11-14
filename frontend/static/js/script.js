
window.onload = () => {
  const chat = document.getElementById("chat");
  const replyBox = document.getElementById("replyBox");
  const replyText = document.getElementById("replyText");
  const messageInput = document.getElementById("messageInput");
  const fileInput = document.getElementById("fileInput");

  if (!chat || !fileInput || !messageInput || !replyBox || !replyText) {
    console.error("DOM missing");
    return;
  }

  // Get WebSocket URL from current location
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const wsUrl = `${protocol}//${window.location.host}/ws`;
  const ws = new WebSocket(wsUrl);

  let replyTo = null;
  let replyToId = null;
  let pendingFile = null;
  const user = /iphone/i.test(navigator.userAgent) ? "J" : "A";
  let inactivityTimer = null;
  const messageMap = {}; // Store messages by ID for replies

  // Inactivity timeout - redirect to / after 30 seconds
  function resetInactivityTimer() {
    clearTimeout(inactivityTimer);
    inactivityTimer = setTimeout(() => {
      window.location.href = "/";
    }, 30000); // 30 seconds
  }

  // Track user activity
  document.addEventListener("keydown", resetInactivityTimer);
  document.addEventListener("click", resetInactivityTimer);
  fileInput.addEventListener("change", resetInactivityTimer);

  // Initialize inactivity timer
  resetInactivityTimer();

  // Load previous messages on page load
  fetch("/messages")
    .then(res => res.json())
    .then(messages => {
      if (messages && Array.isArray(messages)) {
        messages.forEach(msg => render(msg));
      }
    })
    .catch(err => console.error("Failed to load messages:", err));

  ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    render(msg);
  };

  function render(m) {
    const div = document.createElement("div");
    div.className = `msg ${m.from}`;

    // Store message for reply lookups - include fileUrl if present
    messageMap[m.id] = { text: m.text || "(file)", fileUrl: m.fileUrl };

    let html = "";
    if (m.replyToId) {
      const repliedTo = messageMap[m.replyToId];
      if (repliedTo && repliedTo.fileUrl) {
        html += `<small style='opacity:.6;'>Replying to image:</small><br><img src='${repliedTo.fileUrl}' style='max-width:100px;max-height:100px;border-radius:4px;margin-bottom:8px;'>`;
      } else {
        html += `<small style='opacity:.6'>Replying to: ${repliedTo?.text || m.replyToId}</small><br>`;
      }
    }
    if (m.text) html += `<strong style='font-weight:600;'>${m.from}:</strong> ${m.text}`;
    if (m.fileUrl) html += `<br><img src='${m.fileUrl}' style='max-width:200px;border-radius:6px;'>`;

    div.innerHTML = html;
    div.onclick = () => startReply(m.id, m.text, m.fileUrl);

    chat.appendChild(div);
    chat.scrollTop = chat.scrollHeight;
  }

  function startReply(msgId, content, fileUrl) {
    replyToId = msgId;
    replyTo = content || "(file)";
    replyBox.style.display = "block";
    
    if (fileUrl) {
      replyText.innerHTML = `Replying to image:<br><img src='${fileUrl}' style='max-width:100px;max-height:100px;border-radius:4px;'>`;
    } else {
      replyText.innerText = "Replying to: " + replyTo;
    }
  }

  window.sendMessage = () => {
    const text = messageInput.value.trim();
    if (!text && !pendingFile) return;

    const msg = {
      from: user,
      text: text,
      fileUrl: pendingFile,
      replyToId: replyToId,
      tag: user
    };

    ws.send(JSON.stringify(msg));

    messageInput.value = "";
    replyTo = null;
    replyToId = null;
    replyBox.style.display = "none";
    pendingFile = null;
    messagePreview.innerHTML = "";
  };

  document.getElementById("cancelReplyBtn").onclick = () => {
    replyTo = null;
    replyToId = null;
    replyBox.style.display = "none";
  };

  const messagePreview = document.createElement("div");
  messagePreview.id = "messagePreview";
  messagePreview.style.cssText = "margin: 5px 10px; padding: 5px; border-radius: 6px; background: #44475a;";
  document.getElementById("inputBar").insertBefore(messagePreview, document.getElementById("inputBar").firstChild);

  fileInput.onchange = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const form = new FormData();
    form.append("file", file);

    const res = await fetch("/upload", { method: "POST", body: form });
    const data = await res.json();

    pendingFile = data.url;
    
    // Show preview in input box area
    messagePreview.innerHTML = `<img src='${pendingFile}' style='max-width:150px;max-height:150px;border-radius:6px;'>`;
  };
};