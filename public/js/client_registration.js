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
    deleteButton.classList.add('delete-button');
    deleteButton.onclick = function () {
        container.removeChild(wrapper)
    };

    wrapper.appendChild(enteredText);
    wrapper.appendChild(deleteButton);
    container.appendChild(wrapper);

    scopes.value = '';
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

function showFieldErrors(data){
    Object.keys(data.formErrors).forEach(fieldName => {
        const errorElement = document.getElementById(`${fieldName}`);
        if (errorElement) {
            errorElement.innerHTML = `<p>${data.formErrors[fieldName]}</p>`;
        }
    });
}

function showOverlay() {
    const overlay = document.createElement('div');
    overlay.classList.add('overlay');
    document.body.appendChild(overlay);
}

function showSuccessPopup(data){
    console.log(data.clientId);
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
    copyText.setSelectionRange(0, 99999);
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