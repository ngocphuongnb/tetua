function fileInputPreviewer(element) {
  element.addEventListener("click", function () {
    var fileInputId = element.getAttribute("for");
    var previewerElm = element.querySelector("img");
    var fileInput = document.getElementById(fileInputId);
    fileInput.click();

    fileInput.onchange = function () {
      var file = fileInput.files[0];
      var reader = new FileReader();
      reader.onload = function (e) {
        previewerElm.src = e.target.result;
      };
      reader.readAsDataURL(file);
    };
  });
}

function deleteNode(nodeType, url, nodeID, callback, e) {
  fetch(url + `/${nodeID}`, { method: "DELETE" })
    .then(function (response) {
      if (response.status !== 200) {
        alert(`Error deleting ${nodeType}`);
        return;
      }
      if (typeof callback === "string") {
        window.location.href = callback;
      }
      if (typeof callback === "function") {
        callback(e, response);
      }
    })
    .catch(function (err) {
      console.error(err);
      alert(`Error deleting ${nodeType}`);
    });
}

function listenDeleteNodeEvents(nodeType, url, callback) {
  var selector = `.delete-${nodeType}`;
  var nodeElms = Array.from(document.querySelectorAll(selector));

  for (var nodeElm of nodeElms) {
    nodeElm.addEventListener("click", function (e) {
      e.preventDefault();
      e.stopImmediatePropagation();
      var nodeID = e.target.getAttribute("data-id");

      if (
        !nodeID ||
        !confirm(`Are you sure you want to delete this ${nodeType}?`)
      ) {
        return;
      }

      deleteNode(nodeType, url, nodeID, callback, e);
    });
  }
}

function uploadHandler(file, callback) {
  const formData = new FormData();
  formData.append("file", file);
  fetch("/files/upload", { method: "POST", body: formData })
    .then((res) => {
      if (!res.ok) {
        throw new Error("File upload failed");
      }

      return res.json()
    })
    .then((res) => callback(res.url))
    .catch((e) => callback(null, e));
}

window.addEventListener('load', function () {
  var imagePreviewers = Array.from(
    document.querySelectorAll(".image-upload-previewer")
  );

  for (var imagePreviewer of imagePreviewers) {
    fileInputPreviewer(imagePreviewer);
  }

  if (window.hljs) {
    var hlNodes = document.querySelectorAll("pre code");
    var languages = hljs.listLanguages();

    for (var hlNode of Array.from(hlNodes)) {
      var language = (hlNode.className || "").substring("language-".length) || "auto";
      hlNode.setAttribute("data-language", language);
      hlNode.parentNode.setAttribute("data-language", language);

      if (languages.includes(language)) {
        hljs.highlightElement(hlNode, { language });
      } else {
        var rs = hljs.highlightAuto(hlNode.textContent);
        hlNode.innerHTML = rs.value;
      }
    }
  }

  var commentInputs = Array.from(
    document.querySelectorAll(".comments textarea")
  );
  for (var commentInput of commentInputs) {
    commentInput.addEventListener("keyup", function (e) {
      e.target.style.height = e.target.scrollHeight + "px";
    });
    commentInput.addEventListener("mouseup", function (e) {
      e.target.style.height = e.target.scrollHeight + "px";
    });
  }

  var commentEditBtns = Array.from(document.querySelectorAll(".edit-comment"));
  for (var commentEditBtn of commentEditBtns) {
    commentEditBtn.addEventListener("click", function (e) {
      e.preventDefault();
      e.stopImmediatePropagation();
      var commentElm = e.target.closest(".comment");
      var textareaElm = commentElm.querySelector("textarea");
      commentElm.classList.add("editing");
      textareaElm.style.height = textareaElm.scrollHeight + "px";
    });
  }

  listenDeleteNodeEvents("comment", "/comments", function (ev) {
    ev.target.closest(".comment").remove();
  });
});
