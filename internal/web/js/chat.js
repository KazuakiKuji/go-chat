document.addEventListener("DOMContentLoaded", function () {
  const messageForm = document.getElementById("messageForm");
  const messageInput = document.getElementById("js-messageInput");
  const messageArea = document.getElementById("js-messageArea");
  const sendButton = document.getElementById("js-sendButton");
  const buttonText = sendButton.querySelector(".js-buttonText");

  // テキストエリアの高さを自動調整する関数
  function adjustTextareaHeight(textarea) {
    textarea.style.height = "auto";
    textarea.style.height = textarea.scrollHeight + "px";
  }

  // テキストエリアの入力時に高さを自動調整
  messageInput.addEventListener("input", function () {
    adjustTextareaHeight(this);
  });

  // Ctrl + Enter で送信
  messageInput.addEventListener("keydown", function (e) {
    if (e.key === "Enter" && e.ctrlKey) {
      e.preventDefault();
      messageForm.dispatchEvent(new Event("submit"));
    }
  });

  messageForm.addEventListener("submit", async function (e) {
    e.preventDefault();

    if (!messageInput.value.trim()) {
      return;
    }

    const formData = new FormData(messageForm);
    sendButton.disabled = true;
    buttonText.textContent = "送信中";

    try {
      const response = await fetch("/chat", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        throw new Error("メッセージの送信に失敗しました");
      }

      const data = await response.json();

      // メッセージ要素を作成
      const messageDiv = document.createElement("div");
      messageDiv.className = "l-chatMain__message p-message --sent";
      messageDiv.innerHTML = `
        <div class="l-chatMain__content p-message__content">
          <p class="p-message__text c-txt">${escapeHtml(data.content)}</p>
          <time class="p-message__time c-time">${data.created_at}</time>
        </div>
      `;

      // メッセージを追加
      messageArea.appendChild(messageDiv);

      // 最下部にスクロール
      messageArea.scrollTop = messageArea.scrollHeight;

      // 入力欄をクリア
      messageInput.value = "";
      adjustTextareaHeight(messageInput);
    } catch (error) {
      console.error("Error:", error);
      alert("メッセージの送信に失敗しました");
    } finally {
      sendButton.disabled = false;
      buttonText.textContent = "送信";
    }
  });

  // 初期表示時にすべてのメッセージエリアの高さを調整
  messageInput.dispatchEvent(new Event("input"));
});

// HTMLエスケープ
function escapeHtml(unsafe) {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}
