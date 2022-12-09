package main

import (
	// "errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func fecha(texto string, formato string) string {
	patron := `\/fecha\s\d{4}-(0[0-9]|1[0-2])-([0-2][0-9]|3[01])\s[+-]\s\d{1,6}`
	bool, _ := regexp.MatchString(patron, texto)

	if bool {
		mensaje := strings.Split(texto, " ")
		fecha := mensaje[1]
		signo := mensaje[2]
		dias, _ := strconv.Atoi(mensaje[3])
		t, err := time.Parse(formato, fecha)
		if err != nil {
			return "🗓 Revisa la fecha..."
		}

		if len(mensaje[3]) > 5 {
			return "🤔 Contemos menos de cien mil de días mejor 😉."
		}
		if signo == "+" {
			newT := "La nueva fecha es " + t.AddDate(0, 0, dias).Format(formato)
			return newT
		}
		if signo == "-" {
			newT := "La nueva fecha es " + t.AddDate(0, 0, dias*(-1)).Format(formato)
			return newT
		}
	}

	return "🤔 Algo salió mal. 👨🏽‍🏫 Recuerda que es importante colocar bien los espacios y respetar el formato de la fecha YYYY-MM-dd."
}

func dias(texto string, formato string) string {
	patron := `\/dias\s\d{4}-(0[0-9]|1[0-2])-([0-2][0-9]|3[1])\s:\s\d{4}-(0[0-9]|1[0-2])-([0-2][0-9]|3[01])`
	bool, _ := regexp.MatchString(patron, texto)

	if bool {
		mensaje := strings.Split(texto, " ")
		fecha1 := mensaje[1]
		fecha2 := mensaje[3]

		f1, err1 := time.Parse(formato, fecha1)
		f2, err2 := time.Parse(formato, fecha2)
		if err1 != nil || err2 != nil {
			return "🗓 Revisa las fechas..."
		}
		dias := math.Ceil(f2.Sub(f1).Hours()) / 24

		if dias > -106500 && dias < 106500 {
			flo := fmt.Sprintf("%g", math.Ceil(f2.Sub(f1).Hours()/24))
			return "Hay una diferencia de " + flo + " días."
		}
		if dias < -106500 || dias > 106500 {
			return "🤯 Hay más de 106500 días de diferencia... No puedo contar tanto. 😥"
		}
	}

	return "🤔 Algo salió mal. 👨🏽‍🏫 Recuerda que es importante colocar bien los espacios y respetar el formato de la fecha YYYY-MM-dd."
}

func main() {

	TOKEN := "5653661791:AAFuSHacevEgPMhhqsayDCj4Yt1rNOWJL9U"

	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	formato := "2006-01-02"

	help := "Hola 🤖 te digo la fecha después de ❔ días.\n\nPor ejemplo, si hoy fuese 2022-12-07 y quisieras saber la fecha después de 180 días.\nEscribe lo siguiente👇🏽:\n\n/fecha 2022-12-07 + 180\nTambién puedes restar los días:\n/fecha 2022-12-07 - 180\n\n🤖 Y te diré:\n\nLa nueva fecha es 2023-06-05\nO si has restado...\nLa nueva fecha es 2022-06-10\n\nTambién te digo cuántos días han pasado entre dos fechas, escribe el comando /dias y luego dos fechas con el siguiente formato YYYY-MM-dd separadas por :\t\t\t.\nPor ejemplo:\n\n/dias 2023-06-05 : 1969-07-21\n\n🤖 Y te responderé:\n\nHay un diferencia de 19677 días\n\n❗❗Es importante colocar bien los espacios y respetar el formato de la fecha YYYY-MM-dd.❗❗"

	// error := "🤔 Algo salió mal. 👨🏽‍🏫 Recuerda que es importante colocar bien los espacios y respetar el formato de la fecha YYYY-MM-dd. "

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = help
		case "help":
			msg.Text = help
		case "fecha":
			msg.Text = fecha(update.Message.Text, formato)
		case "dias":
			msg.Text = dias(update.Message.Text, formato)
		default:
			msg.Text = "🥴No conozco ese comando...\n👨🏽‍🏫 Los comandos son:\n/help\n/fecha\n/dias"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

	}

	bot.Debug = true
}
