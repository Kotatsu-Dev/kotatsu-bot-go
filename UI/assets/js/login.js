$(function () {

    if(location.protocol !== "https:") {
        location.protocol = "https:";
    }

    $("#box-error").hide();

    // Нажатие на кнопку Логин при помощи клавиши Enter
    $(document).keypress(function(event) {
        // проверяем, нажата ли клавиша Enter (код 13)
        if (event.which === 13) {
            $("#button-login").trigger("click");
        }
    });

    // Клик по кнопке вход
    $("#button-login").on("click", function () {
 
        //Собираем информацию с полей email и пароль   
        let login = $("#input-login").val();
        let password = $("#input-password").val();

        //если одно из полей пустое
        if(login == "" || password == "") {
            printMessage("error","Все поля должны быть заполнены!");
            return false;
        };

        let credentials = {
            "login": login,
            "password": password,
        };

        loginRequest = ajax_JSON(CONFIG_APP_URL_BASE + "login", "POST", credentials, {});
        handler_sendLoginRequest(loginRequest);
        return false;
    });

    $("#box-message-close").on("click", function () {
        $("#box-error").hide();
        return false;
    });
});



// -----------------------------------
// 
//          AJAX HANDLERS
// 
// -----------------------------------
function handler_sendLoginRequest(request) {
    request.always(function () {
        switch (request.status) {
            //Успех
            case 200:
                $("#box-error").hide();

                //Перебрасываем на adm-panel
                window.location.replace("/admin-panel");
                break;

                //Неверный логин-пароль
            case 401:
                printMessage("error","Неверный логин или пароль!");
                console_RequestError("Invalid auth!", request);
                break;

                //В ином случае
            default:
                printMessage("error","Неизвестная ошибка!");
                console_RequestError("Error!", request);
                break;
        }
    });
}
