//
// ----------------------------------------------------------------------------------
//
// 					Answer Codes (Коды ответов функций БД)
//
// ----------------------------------------------------------------------------------
//

package db

const (
	DB_ANSWER_SUCCESS             = iota // Успех
	DB_ANSWER_OBJECT_EXISTS              // Объект существует
	DB_ANSWER_OBJECT_NOT_FOUND           // Объект НЕ существует
	DB_ANSWER_INVALID_CREDENTIALS        // Неверные данные
	DB_ANSWER_PERMISSION_DENIED          // Отказано в доступе
	DB_ANSWER_DELETE_ERROR               // Ошибка удаления объекта
	DB_ANSWER_UNEXPECTED_ERROR           // Незапланированная ошибка
)
