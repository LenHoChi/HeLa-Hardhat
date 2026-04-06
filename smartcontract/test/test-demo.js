function addNumbers(a, b) {
  return a + b;
}

console.log("addNumbers(2, 3) =", addNumbers(2, 3));

async function getDemoMessage() {
  return "hello from async function";
}

getDemoMessage().then((message) => {
  console.log("async message =", message);
});
// node test/test-demo.js