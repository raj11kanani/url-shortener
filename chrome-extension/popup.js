const shortenBtn = document.getElementById("shorten-btn");
const originalUrlDiv = document.getElementById("original-url");
const resultDiv = document.getElementById("result");
const shortUrlDiv = document.getElementById("short-url");
const copyMsg = document.getElementById("copy-msg");
const errorDiv = document.getElementById("error");

let currentUrl = "";

// When popup opens, get the current tab's URL
chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
  currentUrl = tabs[0].url;
  originalUrlDiv.textContent = currentUrl;
});

// When user clicks Shorten button
shortenBtn.addEventListener("click", async () => {
  errorDiv.style.display = "none";
  resultDiv.style.display = "none";
  shortenBtn.textContent = "Shortening...";
  shortenBtn.disabled = true;

  try {
    const response = await fetch("http://localhost:8080/shorten", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url: currentUrl }),
    });

    if (!response.ok) throw new Error("Server error");

    const data = await response.json();

    // Show the short URL
    shortUrlDiv.textContent = data.short_url;
    resultDiv.style.display = "block";

    // Click to copy
    shortUrlDiv.addEventListener("click", () => {
      navigator.clipboard.writeText(data.short_url);
      copyMsg.style.display = "block";
      setTimeout(() => (copyMsg.style.display = "none"), 2000);
    });
  } catch (err) {
    errorDiv.textContent = "❌ Could not connect to server. Is it running?";
    errorDiv.style.display = "block";
  } finally {
    shortenBtn.textContent = "Shorten This URL";
    shortenBtn.disabled = false;
  }
});