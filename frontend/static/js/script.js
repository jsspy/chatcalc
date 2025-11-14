
const chatDisplay = document.getElementById('chatDisplay');
const messageInput = document.getElementById('messageInput');

const scheme = (location.protocol === 'https:') ? 'wss' : 'ws';
const ws = new WebSocket(`${scheme}://${location.host}/ws`);

ws.addEventListener('open', () => {
  console.log('WebSocket connected');
  // tell server we're ready to receive history to avoid a race
  ws.send(JSON.stringify({ type: 'ready' }));
});

ws.addEventListener('message', (ev) => {
  try {
    const data = JSON.parse(ev.data);
    if (data.type === 'history' && Array.isArray(data.posts)) {
      data.posts.forEach(appendMessage);
    } else if (data.type === 'message' && data.post) {
      appendMessage(data.post);
    }
    chatDisplay.scrollTop = chatDisplay.scrollHeight;
  } catch (err) {
    console.error('Invalid ws message', err);
  }
});

function appendMessage(m) {
  const div = document.createElement('div');
  // If message looks like an uploaded image URL, render image preview
  const isImage = typeof m.message === 'string' && (m.message.startsWith('/uploads/') || /\.(png|jpe?g|gif|webp)(\?|$)/i.test(m.message));
  if (m.author === 'J') {
    if (isImage) {
      div.innerHTML = `<b style=\"color: #8A0082\">${m.author}:</b> `;
      const img = document.createElement('img');
      img.src = m.message;
      img.className = 'chat-image';
      div.appendChild(img);
    } else {
      div.innerHTML = `<b style=\"color: #8A0082\">${m.author}:</b> ${m.message}`;
    }
  } else if (m.author === 'A') {
    if (isImage) {
      div.innerHTML = `<b style=\"color: gold\">${m.author}:</b> `;
      const img = document.createElement('img');
      img.src = m.message;
      img.className = 'chat-image';
      div.appendChild(img);
    } else {
      div.innerHTML = `<b style=\"color: gold\">${m.author}:</b> ${m.message}`;
    }
  } else {
    if (isImage) {
      div.innerHTML = `<b>${m.author}:</b> `;
      const img = document.createElement('img');
      img.src = m.message;
      img.className = 'chat-image';
      div.appendChild(img);
    } else {
      div.innerHTML = `<b>${m.author}:</b> ${m.message}`;
    }
  }
  chatDisplay.appendChild(div);
}

function sendMessage() {
  console.log("d")
  const text = messageInput.value.trim();
  const isIphone = /iPhone/i.test(navigator.userAgent);
  const user = isIphone ? 'J' : 'A';

  if (!text || ws.readyState !== WebSocket.OPEN) return;

  ws.send(JSON.stringify({ user, text }));
  messageInput.value = '';
}

let inactivityTimer;

function resetTimer() {
  clearTimeout(inactivityTimer);
  inactivityTimer = setTimeout(() => {
    window.location.href = '/';
  }, 30000); // 30 seconds
}

// reset on any interaction
['mousemove', 'keydown', 'scroll', 'touchstart'].forEach(evt => {
  document.addEventListener(evt, resetTimer);
});

resetTimer(); // start timer

document.getElementById('imageUploader').addEventListener('change', function () {
    uploadImage();
});

function uploadImage() {
    const fileInput = document.getElementById('imageUploader');
    const file = fileInput.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('imageFile', file);

    fetch('/api/upload-image', {
        method: 'POST',
        body: formData
    })
    .then(res => res.json())
    .then(data => {
    console.log('Upload successful:', data);
    if (data && data.status === 'success' && data.file) {
      const isIphone = /iPhone/i.test(navigator.userAgent);
      const user = isIphone ? 'J' : 'A';
      const imageUrl = '/uploads/' + data.file;
      // send as chat message via websocket so it gets saved and broadcast
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ user, text: imageUrl }));
      }
    }
    })
    .catch(err => {
        console.error('Upload failed:', err);
    });
}
