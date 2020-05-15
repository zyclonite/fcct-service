let transpileButton = document.getElementById("transpileButton");
let validateButton = document.getElementById("validateButton");
let clearButton = document.getElementById("clearButton");
let yamlInput = document.getElementById("yamlInput");
let jsonOutput = document.getElementById("jsonOutput");
let strictCheckbox = document.getElementById("strictCheckbox");
let prettyCheckbox = document.getElementById("prettyCheckbox");
let toast = document.getElementById("toast");
let showNotification = (text) => {
  toast.innerHTML = text;
  toast.classList.remove("hidden");
  setTimeout(() => {
    toast.classList.add("hidden");
  }, 1000);
};
clearButton.addEventListener("click", (e) => {
  yamlInput.value = "";
}, false);
validateButton.addEventListener("click", (e) => {
  showNotification("not implemented yet");
}, false);
transpileButton.addEventListener("click", (e) => {
  let input = yamlInput.value;
  if (input.length <= 0) {
    showNotification("no input");
    return;
  }
  fetch('api/v1/transpile?strict='+strictCheckbox.checked+'&pretty='+prettyCheckbox.checked, {
      body: input,
      headers: {
        'content-type': 'text/x-yaml'
      },
      method: 'POST'
    })
    .then((response) => {
      if (!response.ok) {
        response.json().then((json) => {
          showNotification(json.error);
        });
        return response.statusText;
      }
      return response.text();
    })
    .then((output) => {
      showNotification("success");
      jsonOutput.value = output;
    })
    .catch((error) => {
      showNotification(error);
    });
}, false);
