package cb

import (
	"bytes"
	"fmt"
	"io"
	"os"
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"
	"time"
)

func check_is_participant(user *db.User_ReadJSON, activity *db.Activity_ReadJSON) bool {
	for _, participant := range activity.Participants {
		if user.UserTgID == participant.UserTgID {
			return true
		}
	}
	return false
}

func JoinClubQuery(user *db.User_ReadJSON) Executor {
	return Do(
		AnswerQuery(),
		QueryData().
			Then(func(data string) Executor {
				switch data {
				case "from_ITMO_student":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_ITMO_ENTER_ISU,
							"itmo_status": db.Student,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)

				case "from_ITMO_graduate":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_ITMO_ENTER_ISU,
							"itmo_status": db.Graduate,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)
				case "from_ITMO_employee":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_ITMO_ENTER_ISU,
							"itmo_status": db.Employee,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)
				case "from_ITMO_student_employee":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_ITMO_ENTER_ISU,
							"itmo_status": db.StudentEmployee,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)
				case "from_ITMO_graduate_employee":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_ITMO_ENTER_ISU,
							"itmo_status": db.GraduateEmployee,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)
				default:
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_NOITMO_ENTER_FULLNAME,
							"itmo_status": db.Guest,
						}),
						SendMessageM(
							"request.enter_full_name", nil,
						),
					)
				}
			}),
	)
}

func RelevancePhoneQuery(current_user *db.User_ReadJSON) Executor {
	return Do(
		AnswerQuery(),
		GetActivityByID(uint(current_user.TempActivityID)).
			Then(func(activity *db.Activity_ReadJSON) Executor {
				return QueryData().
					Then(func(data string) Executor {
						if data == "yes" {
							if activity.Status {
								if activity.GuestRegistrationUntil != nil &&
									!current_user.IsITMO &&
									activity.GuestRegistrationUntil.Before(time.Now()) {
									return Do(
										UpdateStep(current_user, config.STEP_DEFAULT),
										SendMessageM(
											"events.registration_closed",
											keyboards.ListEvents,
										),
									)
								} else {
									return Do(
										UpdateStep(current_user, config.STEP_DEFAULT),
										AddParticipant(activity, current_user),
										SendMessageMT(
											"events.registered", activity,
											keyboards.ListEvents,
										),
									)
								}
							} else {
								return SendMessageM(
									"events.non_existent",
									keyboards.ListEvents,
								)
							}

						} else {
							return Do(
								UpdateStep(current_user, config.STEP_CHANGING_PHONE),
								SendMessageM(
									"request.send_phone", keyboards.Keyboard_RequestContact,
								),
							)
						}
					})
			}),
	)
}

func UnsubscribeQuery(user *db.User_ReadJSON) Executor {
	return Do(
		AnswerQuery(),
		QueryDataUint().
			Then(func(activity_id uint64) Executor {
				return GetActivityByID(uint(activity_id)).
					Then(func(activity *db.Activity_ReadJSON) Executor {
						return RemoveParticipant(activity, user).
							Then(func(activity *db.Activity_ReadJSON) Executor {
								return SendMessageMT(
									"events.unregistered", activity,
									keyboards.ListEvents,
								)
							}).
							Otherwise(SendMessageM(
								"events.not_registered",
								keyboards.ListEvents,
							))
					})
			}),
	)
}

func SubscribeQuery(user *db.User_ReadJSON) Executor {
	return Do(
		AnswerQuery(),
		QueryDataUint().
			Then(func(activity_id uint64) Executor {
				return If(
					user.IsFilledData,
					If(user.IsITMO,
						GetActivityByID(uint(activity_id)).
							Then(func(activity *db.Activity_ReadJSON) Executor {
								if activity.Status {
									// Is ITMO = true
									return Do(
										AddParticipant(activity, user),
										UpdateStep(user, config.STEP_DEFAULT),
										SendMessageMT(
											"events.registered", activity,
											keyboards.ListEvents,
										),
									)
								} else {
									return Do(
										UpdateStep(user, config.STEP_DEFAULT),
										SendMessageM(
											"events.non_existent",
											keyboards.ListEvents,
										),
									)
								}
							}),
						Do(
							UpdateUser(user, map[string]any{
								"step":             config.STEP_DEFAULT,
								"temp_activity_id": int(activity_id),
							}),
							SendMessageMT(
								"events.phone_number", user.PhoneNumber,
								keyboards.InlineKbd_RelevancePhoneNumber,
							),
						),
					),
					Do(
						UpdateUser(user, map[string]any{
							"temp_activity_id": int(activity_id),
						}),
						SendMessageM(
							"request.unknown", keyboards.InlineKbd_Appointment,
						),
					),
				)
			}),
	)
}

func ShowActivity(user *db.User_ReadJSON, activity_id uint) Executor {
	return GetActivityByID(activity_id).
		Then(func(activity *db.Activity_ReadJSON) Executor {
			var formattedTime, formattedDate string
			is_participant := check_is_participant(user, activity)

			loc, _ := time.LoadLocation("Europe/Moscow")
			formattedTime = activity.DateMeeting.In(loc).Format("15:04")
			formattedDate = FormatDate(activity.DateMeeting.In(loc))

			var files []io.Reader
			if len(activity.PathsImages) != 0 {
				files = make([]io.Reader, len(activity.PathsImages))
				for i, path := range activity.PathsImages {
					fileData, err := os.ReadFile(path)
					if err != nil {
						rr_debug.PrintLOG("Query.go", "ShowActivity", "os.ReadFile(path)", "Ошибка открытия файла", err.Error())
						return Empty()
					}

					files[i] = bytes.NewReader(fileData)
				}
			}
			return Do(
				If(len(files) > 0,
					SendPhotosM(
						files, "",
					),
					Empty(),
				),
				SendMessageMT(
					"events.format", &map[string]any{
						"activity":      activity,
						"formattedDate": formattedDate,
						"formattedTime": formattedTime,
					},
					ITE(is_participant,
						keyboards.CreateInlineKbd_UnsubscribeActivity(int(activity.ID)),
						keyboards.CreateInlineKbd_SubscribeActivity(int(activity.ID)),
					),
				),
				UpdateStep(user, config.STEP_ACTIVITY),
			)

		})
}

func ActivitiesQuery(user *db.User_ReadJSON) Executor {
	return Do(
		AnswerQuery(),
		QueryDataUint().
			Then(func(activity_id uint64) Executor {
				return ShowActivity(user, uint(activity_id))
			}),
	)
}

func AppointQuery(user *db.User_ReadJSON) Executor {
	return Do(
		AnswerQuery(),
		QueryData().
			Then(func(data string) Executor {
				switch data {
				case "from_ITMO_student":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
							"itmo_status": db.Student,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)

				case "from_ITMO_graduate":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
							"itmo_status": db.Graduate,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)
				case "from_ITMO_employee":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
							"itmo_status": db.Employee,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)
				case "from_ITMO_student_employee":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
							"itmo_status": db.StudentEmployee,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)
				case "from_ITMO_graduate_employee":
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
							"itmo_status": db.GraduateEmployee,
						}),
						SendMessageM(
							"request.enter_isu_number", nil,
						),
					)
				default:
					return Do(
						UpdateUser(user, map[string]any{
							"step":        config.STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME,
							"itmo_status": db.Guest,
						}),
						SendMessageM(
							"request.enter_full_name", nil,
						),
					)
				}
			}),
	)
}

func FormatDate(t time.Time) string {
	var weekday, month string
	switch t.Weekday() {
	case time.Monday:
		weekday = "понедельник"
	case time.Tuesday:
		weekday = "вторник"
	case time.Wednesday:
		weekday = "среда"
	case time.Thursday:
		weekday = "четверг"
	case time.Friday:
		weekday = "пятница"
	case time.Saturday:
		weekday = "суббота"
	case time.Sunday:
		weekday = "воскресенье"
	}

	switch t.Month() {
	case time.January:
		month = "января"
	case time.February:
		month = "февраля"
	case time.April:
		month = "апреля"
	case time.March:
		month = "марта"
	case time.May:
		month = "мая"
	case time.June:
		month = "июня"
	case time.July:
		month = "июля"
	case time.August:
		month = "августа"
	case time.September:
		month = "сентября"
	case time.October:
		month = "октября"
	case time.November:
		month = "ноября"
	case time.December:
		month = "декабря"
	}

	return fmt.Sprintf("%d %s (%s)", t.Day(), month, weekday)
}
