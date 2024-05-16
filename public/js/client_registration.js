function addEntry() {
    var scopes = document.getElementById('scope_input');
    var container = document.getElementById('scope_container');
    var errorElement = document.getElementById('Scope');

    if (!scopes.value.trim()) {
        container.style.marginTop = '8';
        errorElement.textContent = "Scope cannot be empty";
        return;
    }
    
    if (scopes.value.includes(' ')) {
        container.style.marginTop = '8';
        errorElement.textContent = "Scope cannot contain spaces";
        return;
    }

    container.style.marginTop = '0';

    var wrapper = document.createElement('div');
    wrapper.classList.add('wrapper');

    var enteredText = document.createElement('div');
    enteredText.textContent = scopes.value;
    enteredText.classList.add('entered_text');

    var deleteButton = document.createElement('button');
    deleteButton.innerHTML = '&times;';
    deleteButton.classList.add('delete-button');
    deleteButton.onclick = function () {
        container.removeChild(wrapper);
        updateScopeValues();
        checkScopeError();
    };

    wrapper.appendChild(enteredText);
    wrapper.appendChild(deleteButton);
    container.appendChild(wrapper);

    scopes.value = '';
    updateScopeValues();
    checkScopeError();
}

function updateScopeValues() {
    var container = document.getElementById('scope_container');
    var scopeValuesInput = document.getElementById('scope');

    var scopeValues = Array.from(container.querySelectorAll('.entered_text')).map(function (element) {
        return element.textContent;
    }).join(', ');

    scopeValuesInput.value = scopeValues;
}

function checkScopeError() {
    const scopeErrorElement = document.getElementById("Scope");
    if (scopeErrorElement && scopeErrorElement.classList.contains("error-message") && scopeErrorElement.textContent.includes("empty")) {
        scopeErrorElement.textContent = "";
    }
}

document.addEventListener("DOMContentLoaded", function () {
    const form = document.querySelector("form");
    form.addEventListener("submit", function (event) {
        event.preventDefault();

        const formData = new FormData(form);
        fetch("/client-registration", {
            method: "POST",
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                if (data.formErrors) {
                    clearErrorDivs()
                    showFieldErrors(data)
                }
                else {
                    showOverlay();
                    showSuccessPopup(data);
                }
            })
            .catch(error => {
                console.error('An error occurred:', error);
            });
    });
});

function clearErrorDivs() {
    const errorElements = document.querySelectorAll(".error-message");
    errorElements.forEach(element => {
        element.textContent = "";
    });
}

function showFieldErrors(data) {
    Object.keys(data.formErrors).forEach(fieldName => {
        const errorElement = document.getElementById(`${fieldName}`);
        if (errorElement) {
            errorElement.innerHTML = `<p>${data.formErrors[fieldName]}</p>`;
        }
    });
}

function showOverlay() {
    const errorElements = document.querySelectorAll(".error-message");
    errorElements.forEach(element => {
        element.textContent = "";
    });
    const overlay = document.createElement('div');
    overlay.classList.add('overlay');
    document.body.appendChild(overlay);
}

function showSuccessPopup(data) {
    const form = document.querySelector("form");
    form.reset();
    const popup = document.createElement("div");
    popup.classList.add("popup");
    popup.innerHTML = `
        <h3>OAuth client created</h3>
        <div class="popup-content">
            <span class="close-btn">&times;</span>
            <div class="popup-explanation">Your client access has been successfully created. Copy your client id and secret now. You will not be able to access them again.</div>
            <p>Your Client ID </p>
            <div class="popup-wrapper">
                <input type="text" id="clientID" value="${data.clientId}" readonly />
                <button class="copy-button" onclick="copyText('clientID')">COPY</button>
            </div>
            <p>Your Client Secret </p>
            <div class="popup-wrapper">    
                <input type="text" id="clientSecret" value= ${data.clientSecret} readonly />
                <button class="copy-button" onclick="copyText('clientSecret')">COPY</button>
            </div>
        </div>
`;
    document.body.appendChild(popup);

    localStorage.setItem("popupShown", true);

    const closeBtn = popup.querySelector(".close-btn");
    closeBtn.addEventListener("click", function () {
        hideOverlay();
        document.body.removeChild(popup);
    });

    var container = document.getElementById('scope_container');
    container.innerHTML = '';
    updateScopeValues();
    checkScopeError();
}

function hideOverlay() {
    const overlay = document.querySelector('.overlay');
    if (overlay) {
        document.body.removeChild(overlay);
    }
}

function copyText(id) {
    var copyText = document.getElementById(id);
    copyText.select();
    navigator.clipboard.writeText(copyText.value);
}

document.addEventListener("DOMContentLoaded", function () {
    const form = document.querySelector("form");

    form.addEventListener("input", function (event) {
        const target = event.target;
        const errorElement = target.nextElementSibling;
        const pElement = errorElement.querySelector("p");
        if (errorElement && errorElement.classList.contains("error-message") && pElement.textContent.includes("empty")) {
            errorElement.textContent = "";
        }
    });
}); 