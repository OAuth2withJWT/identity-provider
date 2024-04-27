function addEntry() {
    var scopes = document.getElementById('scope');
    var container = document.getElementById('scope_container');

    var wrapper = document.createElement('div');
    wrapper.classList.add('wrapper');

    var enteredText = document.createElement('div');
    enteredText.textContent = scopes.value;
    enteredText.classList.add('entered_text');

    var deleteButton = document.createElement('button');
    deleteButton.innerHTML = '&times;';
    deleteButton.classList.add('deleteButton');
    deleteButton.onclick = function () {
        container.removeChild(wrapper)
    };

    wrapper.appendChild(enteredText);
    wrapper.appendChild(deleteButton);
    container.appendChild(wrapper);

    scopes.value = '';
}