package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/rr_debug"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Executor interface {
	Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
}

type ChainedExecutor[T any] interface {
	Executor
	Then(func(T) Executor) ChainedExecutor[T]
	Otherwise(Executor) ChainedExecutor[T]
}

type ChainedExecutorS[T any, Self ChainedExecutor[T]] struct {
	then      func(T) Executor
	otherwise Executor
	self      Self
}

func (ce *ChainedExecutorS[T, Self]) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	return ce.self.Execute(ctx, b, update)
}

func (ce *ChainedExecutorS[T, Self]) Then(then func(T) Executor) ChainedExecutor[T] {
	ce.then = then
	return ce.self
}

func (ce *ChainedExecutorS[T, Self]) Otherwise(otherwise Executor) ChainedExecutor[T] {
	ce.otherwise = otherwise
	return ce.self
}

type SendMessageS struct {
	chat_id  int64
	text     string
	keyboard models.ReplyMarkup
}

func (msg *SendMessageS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      msg.chat_id,
		ParseMode:   models.ParseModeHTML,
		Text:        msg.text,
		ReplyMarkup: msg.keyboard,
	})

	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "SendMessageRaw", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}

	return true, err
}

func SendMessageRawE(chat_id int64, text string, keyboard models.ReplyMarkup) Executor {
	return &SendMessageS{
		chat_id:  chat_id,
		text:     text,
		keyboard: keyboard,
	}
}

func SendMessageE(chat_id int64, text string, keyboard models.ReplyMarkup) Executor {
	return SendMessageRawE(chat_id, config.T(text), keyboard)
}

func SendMessageTE(chat_id int64, text string, data any, keyboard models.ReplyMarkup) Executor {
	return SendMessageRawE(chat_id, config.TT(text, data), keyboard)
}

type SendMessageMS struct {
	text     string
	keyboard models.ReplyMarkup
}

func (msg *SendMessageMS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chat_id,
		ParseMode:   models.ParseModeHTML,
		Text:        msg.text,
		ReplyMarkup: msg.keyboard,
	})

	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "SendMessageRaw", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}

	return true, err
}

func SendMessageRawME(text string, keyboard models.ReplyMarkup) Executor {
	return &SendMessageMS{
		text:     text,
		keyboard: keyboard,
	}
}

func SendMessageME(text string, keyboard models.ReplyMarkup) Executor {
	return SendMessageRawME(config.T(text), keyboard)
}

func SendMessageMTE(text string, data any, keyboard models.ReplyMarkup) Executor {
	return SendMessageRawME(config.TT(text, data), keyboard)
}

type OneOfS []Executor

func (oos OneOfS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	var (
		err error
		res bool
	)

	for _, x := range oos {
		res, err = x.Execute(ctx, b, update)
		if res {
			return res, err
		}
	}

	return res, err
}

func OneOf(executors ...Executor) Executor {
	return OneOfS(executors)
}

type GuardS struct {
	guard    func(context.Context, *bot.Bot, *models.Update) bool
	executor Executor
}

func (guard *GuardS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if guard.guard(ctx, b, update) {
		return guard.executor.Execute(ctx, b, update)
	}
	return false, nil
}

func Guard(guard func(context.Context, *bot.Bot, *models.Update) bool) func(Executor) *GuardS {
	return func(executor Executor) *GuardS {
		return &GuardS{guard: guard, executor: executor}
	}
}

func TextGuardRaw(text string, executor Executor) *GuardS {
	return Guard(func(ctx context.Context, b *bot.Bot, u *models.Update) bool {
		return u != nil && u.Message != nil && u.Message.Text == text
	})(executor)
}

func TextGuard(text string, executor Executor) *GuardS {
	return TextGuardRaw(config.T(text), executor)
}

func QueryGuard(text string, executor Executor) *GuardS {
	return Guard(func(ctx context.Context, b *bot.Bot, u *models.Update) bool {
		return u != nil && u.CallbackQuery != nil && strings.HasPrefix(u.CallbackQuery.Data, text)
	})(executor)
}

func StepGuard(user *db.User_ReadJSON, step int, executor Executor) *GuardS {
	return Guard(func(ctx context.Context, b *bot.Bot, u *models.Update) bool {
		return user.Step == step
	})(executor)
}

type GetCurrentUser struct {
	ChainedExecutorS[*db.User_ReadJSON, *GetCurrentUser]
}

func (gcu *GetCurrentUser) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	var (
		code int
		user *db.User_ReadJSON
	)
	if update != nil && update.Message != nil && update.Message.From != nil {
		code, user = db.DB_GET_User_BY_UserTgID(update.Message.From.ID)
	} else if update != nil && update.CallbackQuery != nil {
		code, user = db.DB_GET_User_BY_UserTgID(update.CallbackQuery.From.ID)
	} else {
		// TODO: Log and return meaningful error
		return true, nil
	}
	if code == db.DB_ANSWER_SUCCESS {
		if gcu.then != nil {
			return gcu.then(user).Execute(ctx, b, update)
		}
		return true, nil
	} else {
		if gcu.otherwise != nil {
			return gcu.otherwise.Execute(ctx, b, update)
		}
		return true, nil
	}
}

func GetCurrentUserE() ChainedExecutor[*db.User_ReadJSON] {
	gcu := &GetCurrentUser{}
	gcu.self = gcu
	return gcu
}

type SeqS []Executor

func (seq SeqS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	var (
		res bool
		err error
	)
	for _, executor := range seq {
		res, err = executor.Execute(ctx, b, update)
		if err != nil {
			return res, err
		}
	}
	return res, err
}

func Seq(seq ...Executor) Executor {
	return SeqS(seq)
}

type UpdateUserS struct {
	user   *db.User_ReadJSON
	update map[string]any
}

func UpdateCurrentUser(user *db.User_ReadJSON, update map[string]any) (int, *db.User, bool) {
	update_user := maps.Clone(update)
	update_user["user_tg_id"] = user.UserTgID
	return db.DB_UPDATE_User(update_user)
}

func (ups *UpdateUserS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	UpdateCurrentUser(ups.user, ups.update)
	// TODO: Fail on update
	return true, nil
}

func UpdateUserE(user *db.User_ReadJSON, update map[string]any) Executor {
	return &UpdateUserS{user: user, update: update}
}

func UpdateGenderE(user *db.User_ReadJSON, gender db.Gender) Executor {
	return UpdateUserE(user, map[string]any{
		"gender": gender,
	})
}

func UpdateVisitedE(user *db.User_ReadJSON, is_visited_events bool) Executor {
	return UpdateUserE(user, map[string]any{
		"is_visited_events": is_visited_events,
	})
}

func UpdateStepE(user *db.User_ReadJSON, step int) Executor {
	return UpdateUserE(user, map[string]any{
		"step": step,
	})
}

type GetActiveRouletteS struct {
	ChainedExecutorS[*db.AnimeRoulette_ReadJSON, *GetActiveRouletteS]
}

func (gar *GetActiveRouletteS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	code, roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	if code == db.DB_ANSWER_SUCCESS {
		if gar.then != nil {
			return gar.then(roulette).Execute(ctx, b, update)
		}
		return true, nil
	} else {
		if gar.otherwise != nil {
			return gar.otherwise.Execute(ctx, b, update)
		}
		return true, nil
	}
}

func GetActiveRoulette() ChainedExecutor[*db.AnimeRoulette_ReadJSON] {
	res := &GetActiveRouletteS{}
	res.self = res
	return res
}

type HasActiveRouletteS struct {
	ChainedExecutorS[bool, *HasActiveRouletteS]
}

func (har *HasActiveRouletteS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	return GetActiveRoulette().
		Then(func(ar *db.AnimeRoulette_ReadJSON) Executor {
			if har.then != nil {
				return har.then(true)
			}
			return Empty()
		}).
		Otherwise(Func(func() (bool, error) {
			if har.then != nil {
				return har.then(false).Execute(ctx, b, update)
			}
			return true, nil
		})).
		Execute(ctx, b, update)
}

func HasActiveRoulette() ChainedExecutor[bool] {
	res := &HasActiveRouletteS{}
	res.self = res
	return res
}

type ParseTextIntS struct {
	ChainedExecutorS[int, *ParseTextIntS]
}

func (pti *ParseTextIntS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	res, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		if pti.otherwise != nil {
			return pti.otherwise.Execute(ctx, b, update)
		}
		return true, err
	} else {
		if pti.then != nil {
			return pti.then(res).Execute(ctx, b, update)
		}
		return true, nil
	}
}

func ParseTextInt() ChainedExecutor[int] {
	res := &ParseTextIntS{}
	res.self = res
	return res
}

type MatchTextS struct {
	ChainedExecutorS[string, *MatchTextS]
	r *regexp.Regexp
}

func (mts *MatchTextS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	matched := mts.r.MatchString(update.Message.Text)
	if matched {
		if mts.then != nil {
			return mts.then(update.Message.Text).Execute(ctx, b, update)
		}
		return true, nil
	} else {
		if mts.otherwise != nil {
			return mts.otherwise.Execute(ctx, b, update)
		}
		return true, nil
	}
}

func MatchText(r *regexp.Regexp) ChainedExecutor[string] {
	res := &MatchTextS{r: r}
	res.self = res
	return res
}

type GetContactS struct {
	ChainedExecutorS[string, *GetContactS]
}

func (gc *GetContactS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message != nil && update.Message.Contact != nil {
		if gc.then != nil {
			return gc.then(update.Message.Contact.PhoneNumber).Execute(ctx, b, update)
		}
		return true, nil
	}
	if gc.otherwise != nil {
		return gc.otherwise.Execute(ctx, b, update)
	}
	return true, nil
}

func GetContact() ChainedExecutor[string] {
	res := &GetContactS{}
	res.self = res
	return res
}

type EmptyS struct{}

func (*EmptyS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	return true, nil
}

func Empty() Executor {
	return &EmptyS{}
}

type FuncS struct {
	f func() (bool, error)
}

func (wc *FuncS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	return wc.f()
}

func Func(f func() (bool, error)) Executor {
	return &FuncS{f: f}
}

type WithCtxS struct {
	f func(ctx context.Context, b *bot.Bot, update *models.Update) Executor
}

func (wc *WithCtxS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	return wc.f(ctx, b, update).Execute(ctx, b, update)
}

func WithCtx(f func(ctx context.Context, b *bot.Bot, update *models.Update) Executor) Executor {
	return &WithCtxS{f: f}
}

type GetActiveActivitiesS struct {
	ChainedExecutorS[[]db.Activity_ReadJSON, *GetActiveActivitiesS]
}

func (gaa *GetActiveActivitiesS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	activities := db.DB_GET_Active_Activities()
	if gaa.then != nil {
		return gaa.then(activities).Execute(ctx, b, update)
	}
	return true, nil
}

func GetActiveActivities() ChainedExecutor[[]db.Activity_ReadJSON] {
	res := &GetActiveActivitiesS{}
	res.self = res
	return res
}

type GetActivityByIDS struct {
	ChainedExecutorS[*db.Activity_ReadJSON, *GetActivityByIDS]
	id uint
}

func (gai *GetActivityByIDS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	db_answer_code, activity := db.DB_GET_Activity_BY_ID(gai.id)
	if db_answer_code == db.DB_ANSWER_SUCCESS {
		if gai.then != nil {
			return gai.then(activity).Execute(ctx, b, update)
		}
		return true, nil
	}
	if gai.otherwise != nil {
		return gai.otherwise.Execute(ctx, b, update)
	}
	return true, nil
}

func GetActivityByID(id uint) ChainedExecutor[*db.Activity_ReadJSON] {
	res := &GetActivityByIDS{id: id}
	res.self = res
	return res
}

type CreateRequestS struct {
	ChainedExecutorS[*db.User_ReadJSON, *CreateRequestS]
	user *db.User_ReadJSON
}

func (cr *CreateRequestS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	db_answer_code := db.DB_CREATE_Request(cr.user.ID)
	if db_answer_code == db.DB_ANSWER_SUCCESS {
		cont, err := UpdateUserE(cr.user, map[string]any{
			"is_sent_request": true,
		}).Execute(ctx, b, update)

		if err != nil {
			if cr.otherwise != nil {
				return cr.otherwise.Execute(ctx, b, update)
			}
			return cont, err
		}

		if cr.then != nil {
			return cr.then(cr.user).Execute(ctx, b, update)
		}
		return true, nil
	}
	if cr.otherwise != nil {
		return cr.otherwise.Execute(ctx, b, update)
	}
	return true, nil
}

type AddParticipantS struct {
	ChainedExecutorS[*db.Activity_ReadJSON, *AddParticipantS]
	activity *db.Activity_ReadJSON
	user     *db.User_ReadJSON
}

func (ap *AddParticipantS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	code := db.DB_UPDATE_Activity_ADD_Participants(ap.activity.ID, ap.user.ID)
	if code == db.DB_ANSWER_SUCCESS {
		if ap.then != nil {
			return ap.then(ap.activity).Execute(ctx, b, update)
		}
		return true, nil
	}
	if ap.otherwise != nil {
		return ap.otherwise.Execute(ctx, b, update)
	}
	return true, nil
}

func AddParticipant(activity *db.Activity_ReadJSON, user *db.User_ReadJSON) ChainedExecutor[*db.Activity_ReadJSON] {
	res := &AddParticipantS{activity: activity, user: user}
	res.self = res
	return res
}

func CreateRequest(user *db.User_ReadJSON) ChainedExecutor[*db.User_ReadJSON] {
	res := &CreateRequestS{user: user}
	res.self = res
	return res
}

type AnswerQueryS struct{}

func (*AnswerQueryS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})
	}
	return true, nil
}

func AnswerQueryE() Executor {
	return &AnswerQueryS{}
}

type QueryDataS struct {
	ChainedExecutorS[string, *QueryDataS]
}

func (qd *QueryDataS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.CallbackQuery == nil {
		if qd.otherwise != nil {
			return qd.otherwise.Execute(ctx, b, update)
		}
		return true, nil
	}
	parts := strings.Split(update.CallbackQuery.Data, "::")
	data := parts[1]

	if qd.then != nil {
		return qd.then(data).Execute(ctx, b, update)
	}
	return true, nil
}

func QueryDataE() ChainedExecutor[string] {
	res := &QueryDataS{}
	res.self = res
	return res
}

type QueryDataUintS struct {
	ChainedExecutorS[uint64, *QueryDataUintS]
}

func (qdu *QueryDataUintS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.CallbackQuery == nil {
		if qdu.otherwise != nil {
			return qdu.otherwise.Execute(ctx, b, update)
		}
		return true, nil
	}
	parts := strings.Split(update.CallbackQuery.Data, "::")
	data := parts[1]
	res, err := strconv.ParseUint(data, 10, 64)

	if err != nil {
		if qdu.otherwise != nil {
			return qdu.otherwise.Execute(ctx, b, update)
		}
		return true, nil
	}

	if qdu.then != nil {
		return qdu.then(res).Execute(ctx, b, update)
	}
	return true, nil
}

func QueryDataUintE() ChainedExecutor[uint64] {
	res := &QueryDataUintS{}
	res.self = res
	return res
}

type OpenCalendarS struct {
	ChainedExecutorS[io.Reader, *OpenCalendarS]
}

func OpenCalendar() (*os.File, string) {
	directory := config.ByUI("./img/calendar_activities")
	files, err := os.ReadDir(directory)
	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "OpenCalendar", "os.ReadDir", "Ошибка поиска файла календаря", err.Error())
	}

	if len(files) <= 0 {
		return nil, ""
	}

	fileInfo := files[0]
	filePath := filepath.Join(directory, fileInfo.Name())

	file, err := os.Open(filePath)
	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "OpenCalendar", "os.Open", "Ошибка чтения файла календаря", err.Error())
		return nil, ""
	}

	return file, fileInfo.Name()
}

func (oc *OpenCalendarS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	file, _ := OpenCalendar()
	if file == nil {
		if oc.otherwise != nil {
			return oc.otherwise.Execute(ctx, b, update)
		}
		return true, nil
	}

	data, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "OpenCalendarE", "io.ReadAll", "Ошибка чтения файла календаря", err.Error())
		if oc.otherwise != nil {
			return oc.otherwise.Execute(ctx, b, update)
		}
		return true, nil
	}

	if oc.then != nil {
		return oc.then(bytes.NewReader(data)).Execute(ctx, b, update)
	}
	return true, nil
}

func OpenCalendarE() ChainedExecutor[io.Reader] {
	res := &OpenCalendarS{}
	res.self = res
	return res
}

type SendPhotoMS struct {
	file     io.Reader
	text     string
	keyboard models.ReplyMarkup
}

func (msg *SendPhotoMS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	_, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:    chat_id,
		ParseMode: models.ParseModeHTML,
		Photo: &models.InputFileUpload{
			Filename: fmt.Sprintf("photo_%d.jpg", time.Now().UnixNano()),
			Data:     msg.file,
		},
		Caption:     msg.text,
		ReplyMarkup: msg.keyboard,
	})

	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "SendPhotoMS", "bot.SendPhoto", "Ошибка отправки фотографии", err.Error())
	}

	return true, err
}

func SendPhotoRawME(file io.Reader, text string, keyboard models.ReplyMarkup) Executor {
	return &SendPhotoMS{file: file, text: text, keyboard: keyboard}
}

type SendPhotosMS struct {
	files []io.Reader
	text  string
}

func (msg *SendPhotosMS) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	media_group := make([]models.InputMedia, len(msg.files))
	for i, file := range msg.files {
		if i == 0 {
			media_group[i] = &models.InputMediaPhoto{
				Media:           fmt.Sprintf("attach://photo_%d", i),
				MediaAttachment: file,
				ParseMode:       models.ParseModeHTML,
				Caption:         msg.text,
			}
		} else {
			media_group[i] = &models.InputMediaPhoto{
				Media:           fmt.Sprintf("attach://photo_%d", i),
				MediaAttachment: file,
			}
		}
	}
	_, err := b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
		ChatID: chat_id,
		Media:  media_group,
	})

	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "SendPhotosRaw", "bot.SendMediaGroup", "Ошибка отправки фотографии", err.Error())
	}

	return true, err
}

func SendPhotosRawME(files []io.Reader, text string) Executor {
	return &SendPhotosMS{files: files, text: text}
}

func SendPhotoME(file io.Reader, text string, keyboard models.ReplyMarkup) Executor {
	return SendPhotoRawME(file, config.T(text), keyboard)
}

func SendPhotosME(files []io.Reader, text string) Executor {
	return SendPhotosRawME(files, config.T(text))
}

func SendMessageRaw(ctx context.Context, b *bot.Bot, chat_id int64, text string, keyboard models.ReplyMarkup) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chat_id,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})

	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "SendMessageRaw", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}

	return err
}

func SendMessage(ctx context.Context, b *bot.Bot, chat_id int64, text string, keyboard models.ReplyMarkup) error {
	return SendMessageRaw(ctx, b, chat_id, config.T(text), keyboard)
}

func SendMessageT(ctx context.Context, b *bot.Bot, chat_id int64, text string, data any, keyboard models.ReplyMarkup) error {
	return SendMessageRaw(ctx, b, chat_id, config.TT(text, data), keyboard)
}

func SendMessageM(ctx context.Context, b *bot.Bot, update *models.Update, text string, keyboard models.ReplyMarkup) error {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	return SendMessage(ctx, b, chat_id, text, keyboard)
}

func SendMessageMT(ctx context.Context, b *bot.Bot, update *models.Update, text string, data any, keyboard models.ReplyMarkup) error {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	return SendMessageT(ctx, b, chat_id, text, data, keyboard)
}

func SendPhotoRaw(ctx context.Context, b *bot.Bot, chat_id int64, filename string, file io.Reader, text string, keyboard models.ReplyMarkup) error {
	_, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:    chat_id,
		ParseMode: models.ParseModeHTML,
		Photo: &models.InputFileUpload{
			Filename: filename,
			Data:     file,
		},
		Caption:     text,
		ReplyMarkup: keyboard,
	})

	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "SendPhotoRaw", "bot.SendPhoto", "Ошибка отправки фотографии", err.Error())
	}

	return err
}

func SendPhoto(ctx context.Context, b *bot.Bot, chat_id int64, filename string, file io.Reader, text string, keyboard models.ReplyMarkup) error {
	return SendPhotoRaw(ctx, b, chat_id, filename, file, config.T(text), keyboard)
}

func SendPhotoM(ctx context.Context, b *bot.Bot, update *models.Update, filename string, file io.Reader, text string, keyboard models.ReplyMarkup) error {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	return SendPhoto(ctx, b, chat_id, filename, file, text, keyboard)
}

func SendPhotosRaw(ctx context.Context, b *bot.Bot, chat_id int64, files []io.Reader, text string) error {
	media_group := make([]models.InputMedia, len(files))
	for i, file := range files {
		if i == 0 {
			media_group[i] = &models.InputMediaPhoto{
				Media:           fmt.Sprintf("attach://photo_%d", i),
				MediaAttachment: file,
				ParseMode:       models.ParseModeHTML,
				Caption:         text,
			}
		} else {
			media_group[i] = &models.InputMediaPhoto{
				Media:           fmt.Sprintf("attach://photo_%d", i),
				MediaAttachment: file,
			}
		}
	}
	_, err := b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
		ChatID: chat_id,
		Media:  media_group,
	})

	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "SendPhotosRaw", "bot.SendMediaGroup", "Ошибка отправки фотографии", err.Error())
	}

	return err
}

func SendPhotos(ctx context.Context, b *bot.Bot, chat_id int64, files []io.Reader, text string) error {
	return SendPhotosRaw(ctx, b, chat_id, files, config.T(text))
}

func SendPhotosM(ctx context.Context, b *bot.Bot, update *models.Update, files []io.Reader, text string) error {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	return SendPhotos(ctx, b, chat_id, files, text)
}

func AnswerQuery(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	return b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
}

func UpdateGender(user *db.User_ReadJSON, gender db.Gender) (int, *db.User, bool) {
	return UpdateCurrentUser(user, map[string]any{
		"gender": gender,
	})
}

func UpdateVisited(user *db.User_ReadJSON, is_visited_events bool) (int, *db.User, bool) {
	return UpdateCurrentUser(user, map[string]any{
		"is_visited_events": is_visited_events,
	})
}

func UpdateStep(user *db.User_ReadJSON, step int) (int, *db.User, bool) {
	return UpdateCurrentUser(user, map[string]any{
		"step": step,
	})
}

func ITE[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
