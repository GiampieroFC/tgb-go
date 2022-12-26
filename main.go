package main

import (
	// "errors"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Token         string `json:"token"`
	Formato_fecha string `json:"formato_fecha"`
	Msg_start     string `json:"msg_start"`
	Msg_error     string `json:"msg_error"`
	Msg_help      string `json:"msg_help"`
	Msg_default   string `json:"msg_default"`
	Desarrollador string `json:"desarrollador"`
}

var config Config

func loadConfig() {
	fmt.Println("leyendo el archivo de configuración...")
	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("hubo un problema al leer el archivo", err)
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal("hubo un problema al convertir el archivo", err)
	}
	fmt.Println("archivo de configuración leído...")
}

func fecha(texto string, formato string, wrong string) string {
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
	return wrong
}

func dias(texto string, formato string, wrong string) string {
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
	return wrong
}

func main() {
	loadConfig()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)
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
			msg.Text = config.Msg_start + config.Desarrollador
		case "help":
			msg.Text = config.Msg_help + config.Desarrollador
		case "fecha":
			msg.Text = fecha(update.Message.Text, config.Formato_fecha, config.Msg_error)
		case "dias":
			msg.Text = dias(update.Message.Text, config.Formato_fecha, config.Msg_error)
		default:
			msg.Text = config.Msg_default
		}
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
	bot.Debug = true
}
