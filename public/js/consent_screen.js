document.addEventListener("DOMContentLoaded", function () {
    var denyButton = document.getElementById('deny-button');

    denyButton.addEventListener("click", function (event) {
        var form = document.getElementById('consent-form');
        var checkboxes = form.querySelectorAll('input[type="checkbox"]');
        checkboxes.forEach(function (checkbox) {
            if (checkbox.checked) {
                checkbox.checked = false;
            }
        });
        form.submit();
    });
});

