function getAddressInputValue(inputElement) {
    return inputElement.value.trim();
}

function getNumberInputValue(inputElement) {
  return Number(inputElement.value);
}

function showError(message) {
  alert(message);
}

function showSuccess(message) {
  alert(message);
}

export { getAddressInputValue, getNumberInputValue, showError, showSuccess };
