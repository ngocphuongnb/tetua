function approvePost(postID, e) {
  fetch(`/manage/posts/${postID}/approve`, { method: "POST" })
    .then(function (response) {
      if (response.status !== 200) {
        alert(`Error approve post: ${postID}`);
        return;
      }
      var elm = e.target.closest("li");
      elm.querySelector(".status.error").remove();
    })
    .catch(function (err) {
      console.error(err);
      alert(`Error approve post: ${postID}`);
    });
}

window.addEventListener("load", function () {
  var selector = `.approve-post`;
  var nodeElms = Array.from(document.querySelectorAll(selector));

  for (var nodeElm of nodeElms) {
    nodeElm.addEventListener("click", function (e) {
      e.preventDefault();
      e.stopImmediatePropagation();
      var postID = e.target.getAttribute("data-id");

      if (!postID || !confirm(`Are you sure you want to approve this post?`)) {
        return;
      }

      approvePost(postID, e);
    });
  }
});
