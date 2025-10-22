package config

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func T(message string) string {
	if bundle == nil {
		bundle = i18n.NewBundle(language.Russian)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		bundle.LoadMessageFile("locales/ru.toml")
	}

	localizer := i18n.NewLocalizer(bundle, "ru")

	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: message,
	})
}

func TT(message string, data any) string {
	if bundle == nil {
		bundle = i18n.NewBundle(language.Russian)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		bundle.LoadMessageFile("locales/ru.toml")
	}

	localizer := i18n.NewLocalizer(bundle, "ru")

	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    message,
		TemplateData: data,
	})
}
