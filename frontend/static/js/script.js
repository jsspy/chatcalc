
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

  // Send message on Enter key
  messageInput.addEventListener('keydown', function (e) {
    if (e.key === 'Enter') {
      e.preventDefault();
      // call the global send function
      if (typeof window.sendMessage === 'function') {
        window.sendMessage();
      }
    }
  });

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

  // Prevent double-tap zoom on iOS as a fallback: block quick successive taps
  // We use a non-passive listener so we can call preventDefault()
  let lastTouch = 0;
  document.addEventListener('touchend', function (e) {
    const now = Date.now();
    if (now - lastTouch <= 300) {
      // Prevent second tap from triggering native zoom
      e.preventDefault();
    }
    lastTouch = now;
  }, { passive: false });

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

  function formatDate(iso) {
    if (!iso) return "";
    try {
      const d = new Date(iso);
      const day = String(d.getDate()).padStart(2, '0');
      const month = String(d.getMonth() + 1).padStart(2, '0');
      const hours = String(d.getHours()).padStart(2, '0');
      const minutes = String(d.getMinutes()).padStart(2, '0');
      return `${day}/${month} - ${hours}:${minutes}`;
    } catch (e) {
      return iso;
    }
  }

  // Determine file type by URL/extension
  function getFileType(url) {
    if (!url) return 'other';
    const u = url.split('?')[0].toLowerCase();
    if (u.match(/\.(png|jpe?g|gif|webp|bmp)$/)) return 'image';
    if (u.match(/\.(mp4|webm|ogg|mov|mkv)$/)) return 'video';
    if (u.match(/\.(mp3|wav|m4a|aac|ogg)$/)) return 'audio';
    return 'other';
  }

  function render(m) {
    const div = document.createElement("div");
    div.className = `msg ${m.from}`;

    // Store message for reply lookups - include fileUrl if present
    messageMap[m.id] = { text: m.text || "(file)", fileUrl: m.fileUrl };

    let html = "";
    const ts = m.createdAt ? formatDate(m.createdAt) : "";
    if (m.replyToId) {
      const repliedTo = messageMap[m.replyToId];
      if (repliedTo && repliedTo.fileUrl) {
        const t = getFileType(repliedTo.fileUrl);
        if (t === 'image') {
          html += `<small style='opacity:.6;'>Replying to image:</small><br><img src='${repliedTo.fileUrl}' style='max-width:100px;max-height:100px;border-radius:4px;margin-bottom:8px;'>`;
        } else if (t === 'video') {
          html += `<small style='opacity:.6;'>Replying to video:</small><br><video muted playsinline preload='metadata' style='max-width:120px;max-height:80px;border-radius:4px;margin-bottom:8px;' src='${repliedTo.fileUrl}'></video>`;
        } else if (t === 'audio') {
          html += `<small style='opacity:.6;'>Replying to audio:</small><br><audio controls src='${repliedTo.fileUrl}' style='max-width:150px;margin-bottom:8px;'></audio>`;
        } else {
          html += `<small style='opacity:.6'>Replying to file</small><br>`;
        }
      } else {
        html += `<small style='opacity:.6'>Replying to: ${repliedTo?.text || m.replyToId}</small><br>`;
      }
    }
    if (m.text) html += `<div class="msg-header"><strong style='font-weight:600;'>${m.from}:</strong> <span class="time">${ts}</span></div> ${m.text}`;
    if (m.fileUrl) {
      const ft = getFileType(m.fileUrl);
      if (ft === 'image') {
        html += `<br><img src='${m.fileUrl}' style='max-width:200px;border-radius:6px;'>`;
      } else if (ft === 'video') {
        html += `<br><video controls playsinline preload='metadata' style='max-width:320px;border-radius:6px;' src='${m.fileUrl}'></video>`;
      } else if (ft === 'audio') {
        html += `<br><audio controls src='${m.fileUrl}' style='width:100%;max-width:320px;'></audio>`;
      } else {
        html += `<br><a href='${m.fileUrl}' target='_blank' rel='noopener'>Download file</a>`;
      }
    }

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
      const t = getFileType(fileUrl);
      if (t === 'image') {
        replyText.innerHTML = `Replying to image:<br><img src='${fileUrl}' style='max-width:100px;max-height:100px;border-radius:4px;'>`;
      } else if (t === 'video') {
        replyText.innerHTML = `Replying to video:<br><video muted playsinline preload='metadata' style='max-width:180px;max-height:120px;border-radius:4px;' src='${fileUrl}'></video>`;
      } else if (t === 'audio') {
        replyText.innerHTML = `Replying to audio:<br><audio controls src='${fileUrl}' style='max-width:180px;'></audio>`;
      } else {
        replyText.innerText = "Replying to: " + replyTo;
      }
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
    
    // Show preview in input box area (respect file type)
    const ptype = getFileType(pendingFile);
    if (ptype === 'image') {
      messagePreview.innerHTML = `<img src='${pendingFile}' style='max-width:150px;max-height:150px;border-radius:6px;'>`;
    } else if (ptype === 'video') {
      messagePreview.innerHTML = `<video controls playsinline preload='metadata' style='max-width:200px;max-height:150px;border-radius:6px;' src='${pendingFile}'></video>`;
    } else if (ptype === 'audio') {
      messagePreview.innerHTML = `<audio controls src='${pendingFile}' style='max-width:200px;'></audio>`;
    } else {
      messagePreview.innerHTML = `<a href='${pendingFile}' target='_blank' rel='noopener'>Uploaded file</a>`;
    }
  };
};