const getBeanText = async () => {
  console.log("call text api");
  // I can't figure out how to catch errors on fecth()
  const request = await fetch("/api/beans/text");
  if (!request.ok) {
    const message = `An error has occured: ${request.status}`;
    throw new Error(message);
  }
  const data = await request.json();
  console.log(data);
  document.getElementById("ai-text-container").innerText = data;
};
const getNewBeanImage = async () => {
  console.log("call image api");
  const request = await fetch("/api/beans/image");
  if (!request.ok) {
    const message = `An error has occured: ${request.status}`;
    throw new Error(message);
  }
  const data = await request.text();
  console.log(data);
  document.getElementById("ai-image-container").src = data;
};

window.onload = function () {
  document
    .getElementById("magic-beans")
    .addEventListener("click", getObjectsFromAI);
};

function getObjectsFromAI() {
  const element = document.getElementById("hidden-col");
  element.style.display = null;
  element.scrollIntoView();
  getBeanText();
  getNewBeanImage();
}
