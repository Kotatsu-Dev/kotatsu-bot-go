package config

const (
	FILE_PHOTO_GEMFEST_PATH = "./img/templates/gemfest.jpg"
)

const (
	// Steps
	STEP_DEFAULT = iota // Меню 1
	STEP_MESSAGE_SUPPORT
	STEP_ITMO_ENTER_ISU
	STEP_ITMO_ENTER_FULLNAME
	STEP_ITMO_ENTER_SECRET_CODE
	STEP_NOITMO_ENTER_FULLNAME
	STEP_NOITMO_ENTER_PHONE
	STEP_NOITMO_ENTER_SECRET_CODE
	STEP_APPOINTMENT_ITMO_ENTER_ISU
	STEP_APPOINTMENT_ITMO_ENTER_FULLNAME
	STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME
	STEP_APPOINTMENT_NOITMO_ENTER_PHONE
	STEP_USER_LEAVES_CLUB
	STEP_CHANGING_PHONE

	// Аниме рулетка
	STEP_ANIME_RUOLETTE_ENTER_ENIGMATIC_TITLE
	STEP_ANIME_RUOLETTE_ENTER_LINK_MY_ANIME_LIST
)

const (
	// stages anime roulette
	ANIME_RUOLETTE_STAGE_START_REGISTRATION = iota // Начало регистрации
	ANIME_RUOLETTE_STAGE_ANIME_GATHERING           // Объявление темы и сбор аниме
	ANIME_RUOLETTE_STAGE_DATA_PROCESSING           // Обработка собранных аниме
	ANIME_RUOLETTE_STAGE_END                       // Окончание проведения рулетки
)
