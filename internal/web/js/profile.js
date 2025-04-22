// アイコン画像の変更を処理する関数
function handleIconChange(event) {
  const file = event.target.files[0];
  if (file) {
    // 画像のプレビューを表示
    const reader = new FileReader();
    reader.onload = function (e) {
      const img = document.getElementById("profile-icon");
      img.src = e.target.result;

      // フォームを自動送信
      event.target.form.submit();
    };
    reader.readAsDataURL(file);
  }
}

// ページ読み込み時の処理
document.addEventListener("DOMContentLoaded", function () {
  // 保存された画像を復元
  const savedIcon = localStorage.getItem("selectedIcon");
  if (savedIcon) {
    const img = document.getElementById("profile-icon");
    img.src = savedIcon;
  }
});
