$(function () {

    if(location.protocol !== "https:") {
        location.protocol = "https:";
    }

    $("#box-error").hide();

    let user_tg_id = getUrlParameter("user_tg_id");
    let reference_number = getUrlParameter("reference_number");

    let send_user_tg_id = 0;

    if (user_tg_id == undefined || user_tg_id == "" || reference_number == undefined || reference_number == "") {
        printMessage("error","Ссылка некорректная или неверная. В ссылке должно быть 2 GET параметра");
    } else {
        send_user_tg_id = Number(user_tg_id);
    }

    $("#h2-answer-user").text("Ответ пользователю на обращение №" + reference_number);
    

    // Клик по кнопке Отправить ответ
    $("#button-send-reply").on("click", function () {
        
        if (send_user_tg_id == 0) {
            printMessage("error","ID телеграмма пользователя не был передан в ссылке");
            return false;
        }

        let message = $("#textarea-support-message").val();
        if (message == "") {
            printMessage("error","Вы не ввели сообщение для ответа пользователю");
            return false;
        }

        let credentials = {
            "user_tg_id": send_user_tg_id,
            "reference_number": reference_number,
            "message": message,
        };

        let sendSupportMessageToUserRequest = ajax_JSON(CONFIG_APP_URL_BASE+"send-message-user-from-support", "POST", credentials, {});
        handler_postRequest("send_message_user", sendSupportMessageToUserRequest);
    });

    $("#box-message-close").on("click", function () {
        $("#box-error").hide();
        return false;
    });
    
});

// POST
function handler_postRequest(request_type, request){
    request.always(function(){
    
        switch(request.status){
            //Успех
            case 200:
                switch(request_type){
                    case "send_message_user":
                        printMessage("success","Ответ от тех. поддержки успешно был отправлен пользователю!");
                        break;

                }     
                break;

            case 401:
                printMessage("error","Администратор не авторизован!");
                console_RequestError("Invalid auth!",request);
                break;

            case 404:
                switch(request_type){
                    case "send_message_user":
                        printMessage("error","Пользатель по такому Телеграмм id не найден!");
                        break;
                }
                break;
            
            case 400:
                switch (request.responseJSON.status.code) {
                    case 601:
                        printMessage("error","Ошибка коннекта с БД");
                        console_RequestError("Error message: ",request.responseJSON.status.message);
                        break;
                
                    case 602:
                        printMessage("error","Общая ошибка с БД");
                        console_RequestError("Error message: ",request.responseJSON.status.message);
                        break;

                    default:
                        printMessage("error","Неизвестная ошибка");
                        console_RequestError("Error!", request);
                        break;
                }
            break;

            //В ином случае
            default:
                printMessage("error","Неизвестная ошибка!");
                console_RequestError("Error!", request);
                break;
        }
    });
}