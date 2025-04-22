const urlParams = new URLSearchParams(window.location.search);
const successMessage = urlParams.get("success");
if (successMessage) {
  alert(successMessage);
}
