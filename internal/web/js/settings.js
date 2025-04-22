// ページ読み込み時にフォームの表示状態を設定
document.addEventListener("DOMContentLoaded", function () {
  // ユーザー名変更フォームの表示状態を設定
  const usernameForm = document.querySelector(".l-settings__usernameForm");
  if (usernameForm && usernameForm.classList.contains("is-active")) {
    usernameForm.style.display = "block";
  }

  // パスワード変更フォームの表示状態を設定
  const passwordForm = document.querySelector(".l-settings__passwordForm");
  if (passwordForm && passwordForm.classList.contains("is-active")) {
    passwordForm.style.display = "block";
  }
});

// パスワード変更フォームの表示/非表示を切り替える
function togglePasswordForm() {
  const form = document.querySelector(".l-settings__passwordForm");
  if (form) {
    // フォームの表示状態を切り替え
    if (form.style.display === "none" || form.style.display === "") {
      form.style.display = "block";
      form.classList.add("is-active");
    } else {
      form.style.display = "none";
      form.classList.remove("is-active");

      // フォームが非表示になった場合、入力をクリア
      const inputs = form.querySelectorAll("input[type='password']");
      inputs.forEach((input) => (input.value = ""));
    }
  } else {
    // フォームが存在しない場合は、サーバーにリクエストを送信
    window.location.href = "/settings?show_password_form=true";
  }
}

// ユーザー名変更フォームの表示/非表示を切り替える
function toggleUsernameForm() {
  const form = document.querySelector(".l-settings__usernameForm");
  if (form) {
    // フォームの表示状態を切り替え
    if (form.style.display === "none" || form.style.display === "") {
      form.style.display = "block";
      form.classList.add("is-active");
    } else {
      form.style.display = "none";
      form.classList.remove("is-active");

      // フォームが非表示になった場合、入力をクリア
      const inputs = form.querySelectorAll("input[type='text']");
      inputs.forEach((input) => (input.value = ""));
    }
  } else {
    // フォームが存在しない場合は、サーバーにリクエストを送信
    window.location.href = "/settings?show_username_form=true";
  }
}
