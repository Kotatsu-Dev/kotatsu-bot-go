package config

import (
	"path"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func Localizer() *i18n.Localizer {
	if bundle == nil {
		bundle = i18n.NewBundle(language.Russian)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		bundle.LoadMessageFile(path.Join(GetConfig().LOCALES_DIR, "ru.toml"))
	}
	return i18n.NewLocalizer(bundle, "ru")
}

func T(message string) string {
	return Localizer().MustLocalize(&i18n.LocalizeConfig{
		MessageID: message,
	})
}

func TT(message string, data any) string {
	return Localizer().MustLocalize(&i18n.LocalizeConfig{
		MessageID:    message,
		TemplateData: data,
	})
}
