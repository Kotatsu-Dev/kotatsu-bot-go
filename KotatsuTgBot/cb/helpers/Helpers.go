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

type SourceS[T any] struct {
	ChainedExecutorS[T, *SourceS[T]]
	get func(ctx context.Context, b *bot.Bot, update *models.Update) (T, bool)
}

func (s *SourceS[T]) Execute(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	value, ok := s.get(ctx, b, update)
	if ok {
		if s.then != nil {
			return s.then(value).Execute(ctx, b, update)
		}
		return true, nil
	}
	if s.otherwise != nil {
		return s.otherwise.Execute(ctx, b, update)
	}
	return true, nil
}

func Source[T any](get func(ctx context.Context, b *bot.Bot, update *models.Update) (T, bool)) ChainedExecutor[T] {
	res := &SourceS[T]{get: get}
	res.self = res
	return res
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

func PrivateMessagesGuard(exectutor Executor) *GuardS {
	return Guard(func(ctx context.Context, b *bot.Bot, u *models.Update) bool {
		return u != nil && u.Message != nil && u.Message.Chat.Type == models.ChatTypePrivate
	})(exectutor)
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

func GetCurrentUserE() ChainedExecutor[*db.User_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (*db.User_ReadJSON, bool) {
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
			return nil, false
		}
		return user, code == db.DB_ANSWER_SUCCESS
	})
}

func CreateOrGetUserE() ChainedExecutor[*db.User_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (*db.User_ReadJSON, bool) {
		full_tg_name := update.Message.From.FirstName + " " + update.Message.From.LastName
		user_to_add := db.User_CreateJSON{
			UserTgID:   update.Message.From.ID,
			UserName:   update.Message.From.Username,
			FullTgName: full_tg_name,
		}

		code, user := db.DB_CREATE_User(&user_to_add)
		return user, code == db.DB_ANSWER_SUCCESS || code == db.DB_ANSWER_OBJECT_EXISTS
	})
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

func GetActiveRoulette() ChainedExecutor[*db.AnimeRoulette_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (*db.AnimeRoulette_ReadJSON, bool) {
		code, roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
		return roulette, code == db.DB_ANSWER_SUCCESS
	})
}

func HasActiveRoulette() ChainedExecutor[bool] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (bool, bool) {
		code, _ := db.DB_GET_AnimeRoulette_BY_Status(true)
		return code == db.DB_ANSWER_SUCCESS, true
	})
}

func ParseTextInt() ChainedExecutor[int] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (int, bool) {
		res, err := strconv.Atoi(update.Message.Text)
		return res, err == nil
	})
}

func MatchText(r *regexp.Regexp) ChainedExecutor[string] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (string, bool) {
		return update.Message.Text, r.MatchString(update.Message.Text)
	})
}

func GetContact() ChainedExecutor[string] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (string, bool) {
		if update.Message != nil && update.Message.Contact != nil {
			return update.Message.Contact.PhoneNumber, true
		}
		return "", false
	})
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

func GetActiveActivities() ChainedExecutor[[]db.Activity_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) ([]db.Activity_ReadJSON, bool) {
		return db.DB_GET_Active_Activities(), true
	})
}

func GetUserActiveActivities(user *db.User_ReadJSON) ChainedExecutor[[]db.Activity_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) ([]db.Activity_ReadJSON, bool) {
		return db.DB_GET_User_Active_Activities(user.ID), true
	})
}

func GetActivityByID(id uint) ChainedExecutor[*db.Activity_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (*db.Activity_ReadJSON, bool) {
		code, activity := db.DB_GET_Activity_BY_ID(id)
		return activity, code == db.DB_ANSWER_SUCCESS
	})
}

func AddParticipant(activity *db.Activity_ReadJSON, user *db.User_ReadJSON) ChainedExecutor[*db.Activity_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (*db.Activity_ReadJSON, bool) {
		code := db.DB_UPDATE_Activity_ADD_Participants(activity.ID, user.ID)
		return activity, code == db.DB_ANSWER_SUCCESS
	})
}

func AddRouletteParticipant(user *db.User_ReadJSON) ChainedExecutor[*db.User_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (*db.User_ReadJSON, bool) {
		code := db.DB_UPDATE_AnimeRoulette_ADD_Participants(user.ID)
		return user, code == db.DB_ANSWER_SUCCESS
	})
}

func RemoveRouletteParticipant(user *db.User_ReadJSON) ChainedExecutor[*db.User_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (*db.User_ReadJSON, bool) {
		code := db.DB_UPDATE_AnimeRoulette_REMOVE_Participants(user.ID)
		return user, code == db.DB_ANSWER_SUCCESS
	})
}

func CreateRequest(user *db.User_ReadJSON) ChainedExecutor[*db.User_ReadJSON] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (*db.User_ReadJSON, bool) {
		if db.DB_CREATE_Request(user.ID) != db.DB_ANSWER_SUCCESS {
			return nil, false
		}
		UpdateUserE(user, map[string]any{
			"is_sent_request": true,
		}).Execute(ctx, b, update)
		return user, true
	})
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

func queryData(update *models.Update) (string, bool) {
	if update.CallbackQuery == nil {
		return "", false
	}
	parts := strings.SplitN(update.CallbackQuery.Data, "::", 2)
	if len(parts) < 2 {
		return "", false
	}
	return parts[1], true
}

func QueryDataE() ChainedExecutor[string] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (string, bool) {
		return queryData(update)
	})
}

func QueryDataUintE() ChainedExecutor[uint64] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (uint64, bool) {
		data, ok := queryData(update)
		if !ok {
			return 0, false
		}
		res, err := strconv.ParseUint(data, 10, 64)
		return res, err == nil
	})
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

func OpenCalendarE() ChainedExecutor[io.Reader] {
	return Source(func(ctx context.Context, b *bot.Bot, update *models.Update) (io.Reader, bool) {
		file, _ := OpenCalendar()
		if file == nil {
			return nil, false
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			rr_debug.PrintLOG("CommandHandlers.go", "OpenCalendarE", "io.ReadAll", "Ошибка чтения файла календаря", err.Error())
			return nil, false
		}

		return bytes.NewReader(data), true
	})
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
