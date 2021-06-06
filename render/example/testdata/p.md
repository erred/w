---
title: p
description: simple paste service
---

<form id="form" autocomplete="off" action="/api/v0/form" enctype="multipart/form-data" method="POST"></form>

<textarea
  id="paste"
  name="paste"
  form="form"
  autofocus
  placeholder="paste something here..."
  rows="10"
  cols="40"
></textarea>

<div>
  <label for="upload" id="uploadlabel"><em>Or upload:</em></label>
  <input type="file" id="upload" name="upload" form="form" />
  <input type="submit" id="submit" value="Send" form="form" />
</div>

<script>
  const file = document.querySelector("#upload");
  file.addEventListener("change", (e) => {
    // Get the selected file
    const [file] = e.target.files;
    document.querySelector("#uploadlabel").innerHTML = `<em>Or upload:</em> ${file.name}`;
  });
</script>
