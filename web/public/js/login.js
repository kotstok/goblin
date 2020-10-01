document.addEventListener('DOMContentLoaded', (event) => {

    document.getElementById("login-form").addEventListener("submit", function (e) {
        e.preventDefault();

        if (e.target.nodeName !== "FORM") {
            return false;
        }

        let data = new FormData(e.target);

        fetch('/auth', {
            method: 'POST',
            body: data,
        }).then(function (response) {
            if (response.ok) {
                return response.json();
            }
            return Promise.reject(response);

        }).then(function (data) {
            console.log(data);

        }).catch(function (error) {
            console.warn(error);
        });
    });
});