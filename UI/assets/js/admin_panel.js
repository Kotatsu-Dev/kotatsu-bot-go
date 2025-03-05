$(function () {

    if(location.protocol !== "https:") {
        location.protocol = "https:";
    }
    
    // ===================================
    // 
    //                INIT
    // 
    // ===================================
    $("#box-error").hide();
    $("#div-requests-list").hide();

    $(document).on('click', '#box-message-close', function () {
        $("#box-error").hide();
        return false;
    });

    // Нажатие на ссылки в таблице заявок
    $(document).on('click', '.tg_url', function () {
        let tg_url = $(this).text();
        window.open(tg_url, '_blank');
        return false;
    });

    // Кнопка - Выйти из сессии
    $("#button-log-out").on("click", function () {
        let outSessionRequest = ajax_GET(CONFIG_APP_URL_BASE+"logout",{},{});
        handler_getRequest("out_session", outSessionRequest);
        return false;
    });

    // Создание превью файлов для импута картинок в рассылке
    $('#files-message').on('change', function() {
        let previewContainer = $('#image-preview-message-container');
        previewContainer.empty(); // Очищаем контейнер с предыдущими превью

        let files = $(this)[0].files;
        for (let i = 0; i < files.length; i++) {
            let file = files[i];

            // Проверяем, что выбранный файл является изображением
            if (file.type.match('image.*')) {
                let reader = new FileReader();

                // Создаем элемент для превью
                let previewElement = $('<img>');
                previewElement.addClass('image-preview');

                reader.onload = function(e) {
                    // Устанавливаем данные изображения в src атрибут превью
                    previewElement.attr('src', e.target.result);

                    // Добавляем превью в контейнер
                    previewContainer.append(previewElement.clone()); // Используйте clone() для создания копии элемента
                };

                // Читаем файл как Data URL
                reader.readAsDataURL(file);
            }
        }
    });

    // Создание превью файлов для импута картинок создания мероприятия
    $('#files-activity').on('change', function() {
        let previewContainer = $('#image-preview-activity-container');
        previewContainer.empty(); // Очищаем контейнер с предыдущими превью

        let files = $(this)[0].files;
        for (let i = 0; i < files.length; i++) {
            let file = files[i];

            // Проверяем, что выбранный файл является изображением
            if (file.type.match('image.*')) {
                let reader = new FileReader();

                // Создаем элемент для превью
                let previewElement = $('<img>');
                previewElement.addClass('image-preview');

                reader.onload = function(e) {
                    // Устанавливаем данные изображения в src атрибут превью
                    previewElement.attr('src', e.target.result);

                    // Добавляем превью в контейнер
                    previewContainer.append(previewElement.clone()); // Используйте clone() для создания копии элемента
                };

                // Читаем файл как Data URL
                reader.readAsDataURL(file);
            }
        }
    });

    sessionStorage.setItem("list_requests", "hide");

    // Отслеживание изменения типа заявок
    $("#request-type").on("change", function() {
        let state_list_requests = sessionStorage.getItem("list_requests");
        if (state_list_requests == "show") {
            let requestsGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/requests", {}, {});
            handler_getRequest("get_requests_list", requestsGetListRequest);
        }
    });

    // Закрываем модальное окно при клике на кнопку "Закрыть" или за его пределами
    $('#closeModalBtn').click(function() {
        $('#myModal').css('display', 'none');
    });

    $(window).click(function(event) {
        if (event.target === $('#myModal')[0]) {
            $('#myModal').css('display', 'none');
        }
    });

    let is_click_input_activity_date_meeting = false;

    // Ставим сегодняшнюю дату в input
    $("#datetime-picker").val(formatTodayDateTime());
    
    // Регистрируем нажатие на выбор даты проведения мероприятия
    $("#datetime-picker").on("change", function () {
        // Этот код будет выполнен, когда пользователь выберет новую дату и время
        is_click_input_activity_date_meeting = true;
    });

    // Добавляем обработчик события "input" (когда что-то вводится в поле)
    $("#input-user-tg-id").on("input", function() {
        // Получаем текущее значение поля ввода
        let inputValue = $(this).val();

        // Используем регулярное выражение, чтобы удалить все символы, кроме цифр
        let numericValue = inputValue.replace(/\D/g, '');

        // Устанавливаем очищенное значение обратно в поле ввода
        $(this).val(numericValue);
    });

    flatpickr("#datetime-picker", {
        enableTime: true,
        dateFormat: "Y-m-d H:i",
        defaultDate: "2023-09-09 12:00",
        minDate: "today",
        locale: "ru", // Указываем, что используем русскую локализацию
      });

    // ===================================
    // 
    //                КНОПКИ
    // 
    // ===================================

    // Кнопка - Получить JSON всех пользователей
    $("#button-get-json-users").on("click", function () {
        let usersGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/users", {}, {});
        handler_getRequest("get_users_list_json", usersGetListRequest);
        return false;
    });

    // Кнопка - Получить JSON всех мероприятий
    $("#button-get-json-activities").on("click", function () {
        let activitiesGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/activities", {}, {});
        handler_getRequest("get_activities_list_json", activitiesGetListRequest);
        return false;
    });

    // Кнопка - Получить JSON всех заявок
    $("#button-get-json-requests").on("click", function () {
        let requestsGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/requests", {}, {});
        handler_getRequest("get_requests_list_json", requestsGetListRequest);
        return false;
    });

    // Кнопка - Получить JSON всех подписчиков клуба
    $("#button-get-json-users-club-subscribers").on("click", function () {
        let usersGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/users", {}, {});
        handler_getRequest("get_users_club_subscribers_list_json", usersGetListRequest);
        return false;
    });

    // Кнопка - Получить JSON всех, кто не подписан на клуб
    $("#button-get-json-users-no-club-subscribers").on("click", function () {
        let usersGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/users", {}, {});
        handler_getRequest("get_users_no_club_subscribers_list_json", usersGetListRequest);
        return false;
    });

    // Кнопка - Получить подписчиков клуба
    $("#button-get-club-subscribers").on("click", function () {
        let usersGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/users", {}, {});
        handler_getRequest("get_users_list_club_subscribers", usersGetListRequest);
        return false;
    });

    // Кнопка - Получить подписчиков рассылки
    $("#button-get-subscribers-news-letter").on("click", function () {
        let usersGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/users", {}, {});
        handler_getRequest("get_users_list_subscribers_news_letter", usersGetListRequest);
        return false;
    });

    // Кнопка - Получить список мероприятий
    $("#button-get-activities-list").on("click", function () {
        let activitiesGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/activities", {}, {});
        handler_getRequest("get_activities_list", activitiesGetListRequest);
        return false;
    });

    // Кнопка Подписчики в списке мероприятий
    $('#modal-content-data').on('click', '#p-activity-participants', function() {
        let participants_str = this.getAttribute("participants");
        let activity_title = this.getAttribute("activity");
        let participants = JSON.parse(participants_str);
        if (participants.length == 0) {
            printMessage("error","На данное мероприятие никто не подписан");
        } else {
            let list_participants = "";
            participants.forEach(function(participant) {

                if (participant.IsITMO) {
                    list_participants += `<div class="col_data">
                                        <p class="title">Полное имя</p>
                                        <p>${participant.FullName}</p>

                                        <p class="title">Телеграмм ID</p>
                                        <p>${participant.UserTgID}</p>

                                        <p class="title">Из ИТМО?</p>
                                        <p>Да</p>

                                        <p class="title">ИСУ</p>
                                        <p>${participant.ISU}</p>

                                        <p class="title">Сслыка на ТГ</p>
                                        <a href="${participant.TgURL}">${participant.TgURL}</a>

                                     </div>`;
                } else {
                    list_participants += `<div class="col_data">
                                        <p class="title">Полное имя</p>
                                        <p>${participant.FullName}</p>

                                        <p class="title">Телеграмм ID</p>
                                        <p>${participant.UserTgID}</p>

                                        <p class="title">Из ИТМО?</p>
                                        <p>Нет</p>

                                        <p class="title">Номер телефона</p>
                                        <p>${participant.PhoneNumber}</p>

                                        <p class="title">Сслыка на ТГ</p>
                                        <a href="${participant.TgURL}">${participant.TgURL}</a>
                                     </div>`;
                }
            });
        
            $("#modal-title").text("Подписчики мероприятия: " + activity_title);
            $("#modal-content-data").html(list_participants);
        }
        // Добавьте ваш код обработки события здесь
    });

    // Кнопка - Зарегистрировать мероприятие
    $("#button-activity-create").on("click", function () {
        let activity_title = $("#input-activity-title").val();
        let activity_date_meeting = $("#datetime-picker").val();
        if (activity_title == "") {
            printMessage("error","Вы не указали название мероприятия");
            return false;
        }

        // let images = $("#files-activity")[0].files;
        // let send_images = [];

        if (!is_click_input_activity_date_meeting) {
            printMessage("error","Вы не указали дату проведения мероприятия");
            return false;
        }
        
        let activity_location = $("#input-activity-location").val();
        if (activity_location == "") {
            printMessage("error","Вы не указали место проведения мероприятия");
            return false;
        }

        let activity_description = $("#input-activity-description").val();
        if (activity_location == "") {
            printMessage("error","Вы не указали описание мероприятия");
            return false;
        }

        // Создаем объект Date для сегодняшней даты
        let today = new Date();

        // Создаем объект Date на основе значения из input
        let selectedDate = new Date(activity_date_meeting);

        // Сравниваем даты
        if (selectedDate < today) {
            printMessage("error","Вы указали дату проведения мероприятия из прошлого");
            return false;
        }

        // Если были добавлены картинки к мероприятию, то загрузить их
        // if (images.length != 0) {
        //     for (let i = 0; i < images.length; i++) {
        //         send_images.push(images[i]);
        //     }
        // }

        let formData = new FormData();
        let fileInput = document.getElementById('files-activity');
        let files = fileInput.files;
        
        for (let i = 0; i < files.length; i++) {
          formData.append("send_images", files[i]);
        }
        
        // Далее добавьте остальные данные в formData, как вы это сделали в вашем коде.
        formData.append("title", activity_title);
        formData.append("date_meeting", activity_date_meeting);
        formData.append("description", activity_description);
        formData.append("location", activity_location);

        // let credentials = {
        //     "title": activity_title,
        //     "date_meeting": activity_date_meeting,
        //     "description": activity_description,
        //     "location": activity_location,
        // };
        
        let activityCreateRequest = ajax_SendFile(CONFIG_APP_URL_BASE+"api/activities/", formData, {});
        handler_postRequest("create_activity", activityCreateRequest);
        return false;
    });

    // Кнопка - Отправить файл
    $("#button-send-file").on("click", function () {
        
        // Получить выбранный файл
        let fileInput = $("#fileInput")[0];
        let file = fileInput.files[0];

        // Проверка, что файл выбран и это изображение
        if (file && file.type.startsWith("image/")) {
            // Создать объект FormData и добавить файл в него
            let formData = new FormData();
            formData.append("image", file);

            let sendFileRequest = ajax_SendFile(CONFIG_APP_URL_BASE+"upload-file-calendar-activities", formData, {});
            handler_postRequest("send_file", sendFileRequest);
        } else {
            printMessage("error","Вы можете загрузить только изображение в формате PNG");
        }

        
        return false;
    });

    // Кнопка - Исключить пользователя из клуба
    $("#button-user-exclude-club").on("click", function () {

        let user_tg_id = $("#input-user-tg-id").val();
        if (user_tg_id == "") {
            printMessage("error", "Вы не указали Телеграмм ID пользователя");
            return false;
        }

        let credentials = {
            "user_tg_id": user_tg_id,
            "is_club_member": 0,
        };

        let userUpdateRequest = ajax_PUT(CONFIG_APP_URL_BASE+"api/users/club-member", credentials, {});
        handler_updateRequest("user_exclude_club", userUpdateRequest);
        return false;
    });

    // Кнопка - Добавить пользователя в клуб
    $("#button-user-add-to-club").on("click", function () {

        let user_tg_id = $("#input-user-tg-id").val();
        if (user_tg_id == "") {
            printMessage("error", "Вы не указали Телеграмм ID пользователя");
            return false;
        }

        let credentials = {
            "user_tg_id": user_tg_id,
            "is_club_member": 1,
        };

        let userUpdateRequest = ajax_PUT(CONFIG_APP_URL_BASE+"api/users/club-member", credentials, {});
        handler_updateRequest("user_add_club", userUpdateRequest);
        return false;
    });

    // Кнопка - Отправить сообщение
    $("#button-send-message").on("click", function () {
        let message = $("#textarea-message").val();

        if (message == "") {
            printMessage("error", "Вы не ввели сообщение для отправки");
            return false;
        }
  
        let formData = new FormData();
        let fileInput = document.getElementById('files-message');
        let files = fileInput.files;
        
        for (let i = 0; i < files.length; i++) {
          formData.append("send_images", files[i]);
        }

        // Далее добавьте остальные данные в formData, как вы это сделали в вашем коде.
        formData.append("message", message);

        let sendMessageRequest = ajax_SendFile(CONFIG_APP_URL_BASE+"send-message-user", formData, {});
        handler_postRequest("send_message", sendMessageRequest);
        return false;
    });

    // Кнопка - Посмотреть список заявок
    $("#button-get-requests-list").on("click", function () {

        let list_requests = sessionStorage.getItem("list_requests");
        
        if (list_requests == "show") {
            $("#div-requests-list").hide();
            $("#button-get-requests-list").text("Получить список заявок");
            sessionStorage.setItem("list_requests", "hide");
        } else {
            let requestsGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/requests", {}, {});
            handler_getRequest("get_requests_list", requestsGetListRequest);
        }
        return false;
    });

    // Кнопка Одобрить
    $("#button-request-approve").on("click", function () {
        let answer = confirm('Вы уверены что хотите одобрить данную заявку?');
        if(answer) {
            
            let request_id = $("#input-del-request-id").val();
            if (request_id == "") {
                printMessage("error", "Вы не указали ID заявки");
                return false;
            }

            let credentials = {
                "request_id": request_id,
                "status": 1,
            };
    
            let requestUpdateRequest = ajax_PUT(CONFIG_APP_URL_BASE+"api/requests/choice", credentials, {});
            handler_updateRequest("update_request_approve", requestUpdateRequest);
        } else {
            return false;
        }
        return false;
    });

    // Кнопка Отклонить
    $("#button-request-deny").on("click", function () {
        let answer = confirm('Вы уверены что хотите отклонить данную заявку?');
        if(answer) {

            let request_id = $("#input-del-request-id").val();
            if (request_id == "") {
                printMessage("error", "Вы не указали ID заявки");
                return false;
            }

            let credentials = {
                "request_id": request_id,
                "status": 2,
            };

            let requestUpdateRequest = ajax_PUT(CONFIG_APP_URL_BASE+"api/requests/choice", credentials, {});
            handler_updateRequest("update_request_deny", requestUpdateRequest);
        } else {
            return false;
        }
        return false;
    });

    // Кнопка ОЧИСТИТЬ БД
    $("#button-all-delete").on("click", function () {
        let answer = confirm('Вы уверены что хотите удалить всю БД?');
        if(answer) {
            let deleteAllRequest = ajax_DELETE(CONFIG_APP_URL_BASE+`all-db`, {}, {});
            handler_deleteRequest("del_all", deleteAllRequest);
        } else {
            return false;
        }
        return false;
    });

    // Кнопка Удалить пользователей
    $("#button-user-delete").on("click", function () {
        let answer = confirm('Вы уверены что хотите удалить всех зарегистрированных пользователей?');
        if(answer) {
            let deleteUsersRequest = ajax_DELETE(CONFIG_APP_URL_BASE+`api/users`, {}, {});
            handler_deleteRequest("del_users", deleteUsersRequest);
        } else {
            return false;
        }
        return false;
    });

    // Кнопка Удалить мероприятия
    $("#button-activity-delete").on("click", function () {
        let answer = confirm('Вы уверены что хотите удалить все мероприятия?');
        if(answer) {
            let deleteActivitiesRequest = ajax_DELETE(CONFIG_APP_URL_BASE+`api/activities`, {}, {});
            handler_deleteRequest("del_activities", deleteActivitiesRequest);
        } else {
            return false;
        }
        return false;
    });

});

// -----------------------------------
//        Misc(прочие функции)
// -----------------------------------
function formatTodayDateTime() {
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0'); // Месяц начинается с 0
    const day = String(now.getDate()).padStart(2, '0');
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    
    const formattedDateTime = `${year}-${month}-${day} ${hours}:${minutes}`;
    return formattedDateTime;
}
// Форматированная сегодняшняя дата для input - date
function formatTodayDateTimeOld() {
    // Получаем текущую дату и время
    let currentDate = new Date();
    let year = currentDate.getFullYear();
    let month = ('0' + (currentDate.getMonth() + 1)).slice(-2); // Добавляем 1, так как месяцы считаются с 0
    let day = ('0' + currentDate.getDate()).slice(-2);
    let hours = ('0' + currentDate.getHours()).slice(-2);
    let minutes = ('0' + currentDate.getMinutes()).slice(-2);

    // Форматируем дату и время в формат, который поддерживает input type="datetime-local"
    let formattedDateTime = year + '-' + month + '-' + day + 'T' + hours + ':' + minutes;
    return formattedDateTime
}

// -----------------------------------
//      Views(представление данных)
// -----------------------------------
// Показать подписчиков клуба в модальном окне
function view_ClubSubscribers(users) {
    let club_subscribers = [];
    let list_club_subscribers = "";

    users.forEach(function(element) {
        if (element.IsClubMember) {
            club_subscribers.push(element);
        }
    });

    if (club_subscribers.length != 0) {

        club_subscribers.forEach(function(club_subscriber) {
            list_club_subscribers += `<div class="row_data">
                                        <p class="title-name">Имя:</p>
                                        <p>${club_subscriber.UserName}</p>
                                        <p class="title-tg-id">Телеграмм ID:</p>
                                        <p>${club_subscriber.UserTgID}</p>
                                     </div>`;
        });

        $('#myModal').css('display', 'block');
        $("#modal-title").text("Подписчики клуба");
        $("#modal-content-data").html(list_club_subscribers);
    } else {
        printMessage("error","У клуба нет ни одного подписчика");
    }
    
}

// Показать подписчиков рассылки в модальном окне
function view_SubscribersNewsLetter(users) {
    let subscribers_news_letter = [];
    let list_subscribers_news_letter = "";

    users.forEach(function(element) {
        if (element.IsSubscribeNewsletter) {
            subscribers_news_letter.push(element);
        }
    });

    if (subscribers_news_letter.length != 0) {

        subscribers_news_letter.forEach(function(club_subscriber) {
            list_subscribers_news_letter += `<div class="row_data">
                                        <p class="title-name">Имя:</p>
                                        <p>${club_subscriber.UserName}</p>
                                        <p class="title-tg-id">Телеграмм ID:</p>
                                        <p>${club_subscriber.UserTgID}</p>
                                     </div>`;
        });

        $('#myModal').css('display', 'block');
        $("#modal-title").text("Подписчики рассылки");
        $("#modal-content-data").html(list_subscribers_news_letter);
    } else {
        printMessage("error","Никто не подписался на рассылку");
    }
}



// Показать список мероприятий в модальном окне
function view_ActivitiesList(activities) {

    let list_activities = "";

    let activity_date_meeting;
    let formatted_activity_date_meeting;
    let participants_str;
    let participants_array = [];
    let x = 0;
    let list_participants;
    let status = "";
    let status_class = "";
    let element_button_del = "";

    activities.forEach(function(element) {

        activity_date_meeting = moment(element.DateMeeting);
        formatted_activity_date_meeting = activity_date_meeting.zone('+0000').format("DD.MM HH:mm");

        list_participants = {
            "activity_title" : element.Title,
            "participants" : element.Participants,
        };

        if (element.Status) {
            status = "Активно"
            status_class = "status-active"
            element_button_del = `<button class="status-button" onclick="inactiveStatusActivity('${element.ID}','${element.Title}')">Удалить</button>`
        } else {
            status = "Не активно"
            status_class = "status-inactive"
            element_button_del = ""
        }

        participants_array.push(list_participants);

        list_activities += `<div class="col_data">
                                    <p class="title">Название</p>
                                    <p>${element.Title}</p>
                                    
                                    <p class="title">Статус</p>
                                    <p class="${status_class}">${status}</p>

                                    <p class="title">Описание</p>
                                    <p>${element.Description}</p>
                                    
                                    <p class="title">Дата проведения</p>
                                    <p>${formatted_activity_date_meeting}</p>

                                    <p class="title">Место проведения</p>
                                    <p>${element.Location}</p>

                                    <p class="title-participants" onclick="getExcelTableParticipants(${x})">Подписчики</p>
                                    ${element_button_del}
                                 </div>`;
                                //  <p class="title-participants" id="p-activity-participants" activity="${element.Title}" participants='${participants_str}'>Подписчики</p>
                                x = x + 1;
    });

    $('#myModal').css('display', 'block');
    $("#modal-title").text("Список мероприятий");
    $("#modal-content-data").html(list_activities);
    
    participants_str = JSON.stringify(participants_array);
    sessionStorage.setItem("participants", participants_str);
}

// Получение Excel таблицы участников мероприятия
function getExcelTableParticipants(index) {
    let participants_str;
    let participants_array;
    let participants;

    participants_str = sessionStorage.getItem("participants");
    participants_array = JSON.parse(participants_str);
    
    participants = participants_array[index].participants;
    if (participants == null) {
        printMessage("error","Никто не является участником данного мероприятия");
    } else {
        let list_participants = [];
        let current_participant;
        let x = 1;
        let isu_str = "";
        let full_name;
        // let first_name;
        // let last_name;
        let phone_number = "";

        participants.forEach(participant => {
            
            if (participant.ISU == "") {
                isu_str = "подписчик не из ИТМО";
            } else {
                isu_str = participant.ISU;
            }

            if (participant.PhoneNumber == "") {
                phone_number = "-";
            } else {
                phone_number = participant.PhoneNumber;
            }

            full_name = participant.FullName;
            // first_name = full_name[1]; // Имя
            // last_name = full_name[0]; // Фамилия

            // if (first_name == undefined) {
            //     first_name = last_name;
            //     last_name = "Не указана";
            // }

            // if (last_name == undefined) {
                // last_name = "Не указана";
            // }

            current_participant = {
                "number": x,
                "ISU": isu_str,
                "full_name": full_name,
                "phone_number": phone_number,
                "tg_url": participant.TgURL,
            }

            list_participants.push(current_participant);

            x = x + 1;
        });

        csv_ParticipantsList(list_participants, participants_array[index].activity_title);
    }
}

// Установка статуса мероприятия как - Неактивно
function inactiveStatusActivity(activity_id, activity_title) {
    let answer = confirm('Вы уверены что хотите удалить мероприятие - ' + activity_title + "?");
    if (answer) {
        
        let credentials = {
            "activity_id": activity_id,
            "status": 1,
        };

        let userUpdateRequest = ajax_PUT(CONFIG_APP_URL_BASE+"api/activities", credentials, {});
        handler_updateRequest("activity_inactive", userUpdateRequest);
        return false;
    }
}

// -----------------------------------
//   Requests(обработчики ответов)
// -----------------------------------
// Формирование JSON списка юзеров
function getJSONList(data) {
    
    let jsonContent = JSON.stringify(data, null, 4); // Преобразование в красивый JSON
      
    // Создание Blob объекта
    let blob = new Blob([jsonContent], { type: 'application/json' });

    // Создаем Object URL
    const url = URL.createObjectURL(blob);

    // Открываем URL в новой вкладке
    window.open(url, '_blank');

    // Очищаем URL после использования (не обязательно, но рекомендуется)
    URL.revokeObjectURL(url);
}

// Формирование списка подписчиков клуба
function formatClubSubscribers(users) {
    
    let is_club_subscribers = false;
    let club_subscribers = [];

    users.forEach(element => {
        if (element.IsClubMember) {
            club_subscribers.push(element);
            is_club_subscribers = true;
        }
    });

    if (is_club_subscribers) {
        getJSONList(club_subscribers);
    } else {
        printMessage("error", "Ни один пользователь не подписан на клуб");
    }
}

// Формирование списка подписчиков клуба
function formatNoClubSubscribers(users) {
    
    let is_no_club_subscribers = false;
    let no_club_subscribers = [];

    users.forEach(element => {
        if (!element.IsClubMember) {
            no_club_subscribers.push(element);
            is_no_club_subscribers = true;
        }
    });

    if (is_no_club_subscribers) {
        getJSONList(no_club_subscribers);
    } else {
        printMessage("success", "Все зарегистрированные пользователи подписаны на клуб");
    }
}

// Формирование списка заявок
function formatRequests(list_requests) {
    
    sessionStorage.setItem("list_requests", "show");
    let request_type = $("#request-type").val();
    let list_element = "";
    let main_list = "";
    let created_at_date;
    let created_at_form;

    let is_itmo = "";
    let secret_code = "";

    let header_table = `<table class="request-table">
                    <thead>
                    <tr>
                        <th>ID заявки</th>
                        <th>Дата подачи заявки</th>
                        <th>ID телеграмма пользователя</th>
                        <th>Пользователь из ИТМО</th>
                        <th>ИСУ пользователя</th>
                        <th>ФИО пользователя</th>
                        <th>Секретный код</th>
                        <th>Номер телефона пользователя</th>
                        <th>Ссылка на телеграмм пользователя</th>
                    </tr>
                    </thead>
                    <tbody>`;

    let footer_table = `</tbody>
                        </table>`;

    for (const element of list_requests) {
        
        created_at_date = Date.parse(element.CreatedAt);
        created_at_form = CONFIG_DATE_TIME_FORMAT.format(created_at_date);

        if (element.UserInfo.IsITMO) {
            is_itmo = "Да";
        } else {
            is_itmo = "Нет";
        }

        if (element.UserInfo.SecretCode == "0") {
            secret_code = "Не имеется";
        } else {
            secret_code = element.UserInfo.SecretCode;
        }

        switch (request_type) {
            case "all_users":
                list_element += `<tr>
                                    <td>${element.ID}</td>
                                    <td>${created_at_form}</td>
                                    <td>${element.UserInfo.UserTgID}</td>
                                    <td>${is_itmo}</td>
                                    <td>${element.UserInfo.ISU}</td>
                                    <td>${element.UserInfo.FullName}</td>
                                    <td>${secret_code}</td>
                                    <td>${element.UserInfo.PhoneNumber}</td>
                                    <td class="tg_url">${element.UserInfo.TgURL}</td>
                                </tr>`
                break;
        
            case "itmo_users":
                if (element.UserInfo.IsITMO) {
                    list_element += `<tr>
                                    <td>${element.ID}</td>
                                    <td>${created_at_form}</td>
                                    <td>${element.UserInfo.UserTgID}</td>
                                    <td>${is_itmo}</td>
                                    <td>${element.UserInfo.ISU}</td>
                                    <td>${element.UserInfo.FullName}</td>
                                    <td>${secret_code}</td>
                                    <td>${element.UserInfo.PhoneNumber}</td>
                                    <td>${element.UserInfo.TgURL}</td>
                                </tr>`
                }
                break;

            case "no_itmo_users":
                if (!element.UserInfo.IsITMO) {
                    list_element += `<tr>
                                    <td>${element.ID}</td>
                                    <td>${created_at_form}</td>
                                    <td>${element.UserInfo.UserTgID}</td>
                                    <td>${is_itmo}</td>
                                    <td>${element.UserInfo.ISU}</td>
                                    <td>${element.UserInfo.FullName}</td>
                                    <td>${secret_code}</td>
                                    <td>${element.UserInfo.PhoneNumber}</td>
                                    <td>${element.UserInfo.TgURL}</td>
                                </tr>`
                }
                break;

            default:
                break;
        }

    }

    main_list = header_table + list_element + footer_table;
    $("#div-requests-list").html(main_list);
    $("#div-requests-list").show();
    $("#button-get-requests-list").text("Закрыть список заявок");
}

// -----------------------------------
//              EXCEL
// -----------------------------------
function csv_ParticipantsList(object_array, activity_title){

    //Массив данных для таблицы, который мы передали только что
    let data = object_array;
    
    //Создаем таблицу Exel
	let wb = new ExcelJS.Workbook();
	
    // Пишем тут Название таблицы
    let workbookName = "Участники мероприятия - " + activity_title + ".xlsx";

    //Название вкладки (ну видел внизу такие вкладки есть у Exel).
    let worksheetName = "Участники мероприятия";

    let ws = wb.addWorksheet(worksheetName, 
        {
        properties: {
            tabColor: {argb:'FFFF0000'}
        }
        }
    );

    ws.columns = [
        { 
            key: "number",
            header: "№",
            width: 10
        },
        { 
            key: "ISU", 
            header: "ИСУ", 
            width: 40
        },
        { 
            key: "full_name", 
            header: "ФИО", 
            width: 20 
        },
        {
            key: "phone_number", 
            header: "Номер телефона", 
            width: 20
        },
        {
            key: "tg_url", 
            header: "Ссылка на ТГ", 
            width: 20
        },
        
    ];

    data.forEach((participant, index) => {
        
        ws.addRow({
            number: participant.number,
            ISU: participant.ISU,
            full_name: participant.full_name,
            phone_number: participant.phone_number,
            tg_url: { text: 'Телеграм', hyperlink: participant.tg_url},
        });
        const numberCell = ws.getCell(index + 2, ws.getColumnKey('number').number);
        const tgUrlCell = ws.getCell(index + 2, ws.getColumnKey('tg_url').number);
        tgUrlCell.font = {
            color: {argb: '0000FF'},    // Синий цвет
            underline: true             // Подчеркивание
        };

        tgUrlCell.alignment = {
            vertical: 'middle',         // Вертикальное центрирование
            horizontal: 'center'        // Горизонтальное центрирование
        };

        numberCell.alignment = {
            vertical: 'middle',         // Вертикальное центрирование
            horizontal: 'center'        // Горизонтальное центрирование
        };
    });

    //Делаем полужирный шрифт первой строки
    ws.getRow(1).font = { bold: true };
    ws.getRow(1).alignment = {
        vertical: 'middle',             // Вертикальное центрирование
        horizontal: 'center'            // Горизонтальное центрирование}
    }
    
    // Записываем в файл
    wb.xlsx.writeBuffer()
        .then(function(buffer) {
        saveAs(
            new Blob([buffer], { type: "application/octet-stream" }),
            workbookName
        );
    });
}

// -----------------------------------
//   Handlers(обработчики запросов)
// -----------------------------------
// GET
function handler_getRequest(request_type, request) {
    request.always(function () {

        switch (request.status) {
            //Успех
            case 200:
                switch (request_type) {
                    case "get_users_list_json":
                        if (request.responseJSON.data.list_users == null) {
                            printMessage("error","Ни один пользователь не зарегистрирован");
                        } else {
                            getJSONList(request.responseJSON.data.list_users);
                        }
                        break;

                    case "get_requests_list_json":
                        if (request.responseJSON.data.list_requests == null) {
                            printMessage("error","Ни одна заявка не была отправлена");
                        } else {
                            getJSONList(request.responseJSON.data.list_requests);
                        }
                        break;

                    case "get_requests_list":
                        if (request.responseJSON.data.list_requests == null) {
                            printMessage("error","Ни одна заявка не была отправлена");
                        } else {
                            formatRequests(request.responseJSON.data.list_requests);
                        }
                        break;

                    case "get_users_club_subscribers_list_json":
                        if (request.responseJSON.data.list_users == null) {
                            printMessage("error","Ни один пользователь не зарегистрирован");
                        } else {
                            formatClubSubscribers(request.responseJSON.data.list_users);
                        }
                        break;

                    case "get_users_no_club_subscribers_list_json":
                        if (request.responseJSON.data.list_users == null) {
                            printMessage("error","Ни один пользователь не зарегистрирован");
                        } else {
                            formatNoClubSubscribers(request.responseJSON.data.list_users);
                        }
                        break;

                    case "get_activities_list_json":
                        if (request.responseJSON.data.list_activities == null) {
                            printMessage("error","Ни одно мероприятие не зарегистрировано");
                        } else {
                            getJSONList(request.responseJSON.data.list_activities);
                        }
                        break;

                    case "get_users_list_club_subscribers":
                        if (request.responseJSON.data.list_users == null) {
                            printMessage("error","Ни один пользователь не зарегистрирован");
                        } else {
                            view_ClubSubscribers(request.responseJSON.data.list_users);
                        }
                        break;

                    case "get_users_list_subscribers_news_letter":
                        if (request.responseJSON.data.list_users == null) {
                            printMessage("error","Ни один пользователь не зарегистрирован");
                        } else {
                            view_SubscribersNewsLetter(request.responseJSON.data.list_users);
                        }
                        break;

                    case "get_activities_list":
                        if (request.responseJSON.data.list_activities == null) {
                            printMessage("error","Ни одно мероприятие не зарегистрировано");
                        } else {
                            view_ActivitiesList(request.responseJSON.data.list_activities);
                        }
                        break;

                    case "out_session":
                        window.location.replace("/login");
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

// POST
function handler_postRequest(request_type, request){
    request.always(function(){
    
        switch(request.status){
            //Успех
            case 200:
                switch(request_type){
                    case "create_activity":
                        printMessage("success","Мероприятие успешно зарегистрировано!");
                        break;

                    case "send_message":
                        printMessage("success","Сообщение на рассылку было успешно отправлено!");
                        break;

                    case "send_file":
                        printMessage("success","Изображение для списка мероприятий было успешно загружено!");
                        break;
                }     
                break;

            case 401:
                printMessage("error","Администратор не авторизован!");
                console_RequestError("Invalid auth!",request);
                break;

            case 404:
                switch(request_type){
                    case "send_message":
                        printMessage("error","Никто из пользователей не подписался на рассылку!");
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

// PUT
function handler_updateRequest(request_type, request){
    request.always(function(){
    
        switch(request.status){
            //Успех
            case 200:
                switch(request_type){
                    case "user_exclude_club":
                        printMessage("success","Пользователь был успешно исключён из клуба!");
                        break;

                    case "user_add_club":
                        printMessage("success","Пользователь был успешно добавлен в клуб!");
                        break;

                    case "update_request_approve":
                        $("#div-requests-list").hide();
                        $("#button-get-requests-list").text("Получить список заявок");
                        sessionStorage.setItem("list_requests", "hide");
                        printMessage("success","Вы одобрили заявку. Пользователь получил уведомление.");
                        break;
    
                    case "update_request_deny":
                        $("#div-requests-list").hide();
                        $("#button-get-requests-list").text("Получить список заявок");
                        sessionStorage.setItem("list_requests", "hide");
                        printMessage("success","Вы отклонили заявку. Пользователь получил уведомление.");
                        break;

                    case "activity_inactive":
                        let activitiesGetListRequest = ajax_GET(CONFIG_APP_URL_BASE+"api/activities", {}, {});
                        handler_getRequest("get_activities_list", activitiesGetListRequest);
                        printMessage("success","Мероприятие успешно удалено.");
                        break;

                }     
                break;

            case 401:
                printMessage("error","Администратор не авторизован!");
                console_RequestError("Invalid auth!",request);
                break;

            case 404:
                switch(request_type){
                    case "user_exclude_club":
                    case "user_add_club":
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

// DELETE
function handler_deleteRequest(request_type, request){
    request.always(function(){
    
        switch(request.status){
            //Успех
            case 200:
                switch(request_type){
                    case "del_users":
                        printMessage("success","Все пользователи успешно удалены");
                        break;

                    case "del_activities":
                        printMessage("success","Все мероприятия успешно удалены");
                        break;

                    case "del_all":
                        printMessage("success","Все сущности успешно удалены");
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