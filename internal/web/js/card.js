document.addEventListener("DOMContentLoaded", () => {
  const iconWraps = document.querySelectorAll("#js-iconWrap");

  iconWraps.forEach((wrap) => {
    wrap.addEventListener("click", (e) => {
      e.stopPropagation();
      e.preventDefault();
      const userId = wrap.getAttribute("data-user-id");
      if (userId) {
        window.location.href = `/profile/${userId}`;
      }
    });
  });
});
