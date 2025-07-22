package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Конфигурация для OpenWeatherMap
const (
	weatherAPIURL = "http://api.openweathermap.org/data/2.5/weather"
)

// Функция для получения погоды
func getWeather(city string) string {
	//Получаем API-ключ из переменных окружения
	weatherAPIKey := os.Getenv("WAK")
	if weatherAPIKey == "" {
		log.Println("API-ключ для OpenWeatherMap отсутствует. Убедитесь, что он задан в файле .env")
		return "Ошибка: API-ключ для OpenWeatherMap отсутствует."
	}
	// Создаем URL для запроса
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric&lang=ru", weatherAPIURL, city, weatherAPIKey)

	// Делаем HTTP-запрос
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка при запросе погоды: %v", err)
		return "Не удалось получить данные о погоде."
	}
	defer resp.Body.Close()

	// Распарсим JSON-ответ
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Ошибка при парсинге ответа: %v", err)
		return "Ошибка обработки данных о погоде."
	}

	// Достаем информацию о погоде
	main, ok := data["main"].(map[string]interface{})
	if !ok {
		return "Не удалось найти информацию о погоде."
	}

	weather, ok := data["weather"].([]interface{})
	if !ok || len(weather) == 0 {
		return "Не удалось найти описание погоды."
	}

	description := weather[0].(map[string]interface{})["description"].(string)
	temp := main["temp"].(float64)

	return fmt.Sprintf("Погода в городе %s: %s, %.1f°C", city, description, temp)
}

func main() {
	// Токен Telegram-бота
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error to loading .env fail")
	}
	botToken := os.Getenv("TBT") // или замените на ваш токен в виде строки

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Включаем режим отладки
	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Настройка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Города
	cityA := "Москва"
	cityB := "Пушкино"
	cityC := "Донецк"

	for update := range updates {
		// Игнорируем обновления без сообщений
		if update.Message == nil {
			continue
		}

		// Логируем сообщения
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Пользовательская клавиатура с кнопками
		buttons := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(cityA),
				tgbotapi.NewKeyboardButton(cityB),
				tgbotapi.NewKeyboardButton(cityC),
			),
		)

		// Если пользователь отправил команду /start
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите город для получения погоды:")
			msg.ReplyMarkup = buttons
			bot.Send(msg)
			continue
		}

		// Обработка выбора города
		if update.Message.Text == cityA {
			customMessage := "Москва ждёт вас с распростёртыми объятиями! 🏙️"
			weather := getWeather(cityA)
			finalText := fmt.Sprintf("%s\n\n%s", weather, customMessage)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, finalText)
			bot.Send(msg)
		} else if update.Message.Text == cityB {
			customMessage := "Пушкино - прекрасный уголок, наполненный уютом ❤️"
			weather := getWeather(cityB)
			finalText := fmt.Sprintf("%s\n\n%s", weather, customMessage)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, finalText)
			bot.Send(msg)
		} else if update.Message.Text == cityC {
			customMessage := "В Донецке солнечно и тепло в сердце 🌞"
			weather := getWeather(cityC)
			finalText := fmt.Sprintf("%s\n\n%s", weather, customMessage)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, finalText)
			bot.Send(msg)
		}

	}
}
