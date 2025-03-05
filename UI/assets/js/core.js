//      RZHEVSKY ROBOTICS
//  Web-development Division
// 
//  core.js - Базоваые функции и конфигурация проекта
//

const CONFIG_APP_URL_BASE = "https://bot.kotatsu.spb.ru/";


const CONFIG_DATE_TIME_FORMAT = new Intl.DateTimeFormat('ru', {
    year: 'numeric',
    month: 'numeric',
    day: 'numeric',
    hour: 'numeric',
    minute: 'numeric',
	second: 'numeric',
});

const CONFIG_DATE_TIME_FORMAT_SHORT = new Intl.DateTimeFormat('ru', {
	year: 'numeric',
	month: 'numeric',
	day: 'numeric',
});

//AJAX
function ajax_JSON(url,request_type, data,custom_headers){
	return $.ajax({
		url : url,
		type : request_type,
		headers: custom_headers,
		data:	JSON.stringify(data),
		contentType: 'application/json; charset=utf-8',
		dataType: 'json',
	});
}

function ajax_PUT(url, data,custom_headers){
	return $.ajax({
		url : url,
		type : "PUT",
		headers: custom_headers,
		data:	JSON.stringify(data),
		contentType: 'application/json; charset=utf-8',
		dataType: 'json',
	});
}

function ajax_GET(url,data,custom_headers){
	return $.ajax({
		url : url,
		type : "GET",
		headers: custom_headers,
		data:	data,
	});
}

function ajax_POST(url,data,custom_headers){
	return $.ajax({
		url : url,
		type : "POST",
		headers: custom_headers,
		data:	data,
	});
}

function ajax_DELETE(url,data,custom_headers){
	return $.ajax({
		url : url,
		type : "DELETE",
		headers: custom_headers,
		data:	data,
	});
}

function ajax_SendFile(url,formData,custom_headers){
	return $.ajax({
		url : url,
		type : "POST",
		headers: custom_headers,
		timeout: 60000,
		contentType: false,
		processData: false,
		data: formData,
	});
}

// Получение GET параметров
function getUrlParameter(sParam) {
    let sPageURL = window.location.search.substring(1),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : decodeURIComponent(sParameterName[1]);
        }
    }
};

// Печать сообщений об успехе или ошибке
function printMessage(type,text,html){

	switch (type) {
		case "success":
			$("#box-error").addClass("success");
			break;
	
		case "warning":
			$("#box-error").removeClass("success");
			$("#box-error").addClass("warning");
			break;

		default:
			$("#box-error").removeClass("success");
			$("#box-error").removeClass("warning");
			break;
	}

	if (html) {
		$("#box-error").show();
		$("#box-error").html(text);
	} else {
		$("#box-error").show();
		$("#box-error p").text(text);
	}

}

//Обработчик ошибки запроса
function console_RequestError(message, request){
    console.error("------\n"
    +"Error! Message: "+message +"\n"
    +"Status code: "+request.status +"\n"
    +"Answer: "+ request.responseText +"\n"
    +"------");
}
//Показать объект запроса
function console_RequestShowObject(object){
    console.log("%c ------",'color: green');
    console.log("%c Ok! Answer: ",'color: green');
    console.log(object);
    console.log("%c ------",'color: green');
}