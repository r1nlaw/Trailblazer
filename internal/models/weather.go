package models

import "time"

// WeatherForecast представляет полный ответ от API OpenWeatherMap (эндпоинт /forecast).
type WeatherForecast struct {
	Cod     string     `json:"cod"`     // Код ответа API (например, "200")
	Message int        `json:"message"` // Сообщение об успехе/ошибке (обычно 0)
	Cnt     int        `json:"cnt"`     // Количество временных интервалов в прогнозе
	List    []Forecast `json:"list"`    // Список прогнозов на 3-часовые интервалы
	City    City       `json:"city"`    // Информация о городе
}
type WeatherResponse struct {
	LandmarkID  int       `json:"landmark_id"`
	Date        time.Time `json:"date"`
	Temperature float64   `json:"temperature"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Rain        float64   `json:"rain"`
	WindSpeed   float64   `json:"wind_speed"`
	WindDegree  float64   `json:"wind_degree"`
}

// Forecast описывает прогноз для одного 3-часового интервала.
type Forecast struct {
	Dt         int64     `json:"dt"`             // Unix timestamp начала интервала
	Main       Main      `json:"main"`           // Основные погодные параметры
	Weather    []Weather `json:"weather"`        // Описание погодных условий
	Clouds     Clouds    `json:"clouds"`         // Облачность
	Wind       Wind      `json:"wind"`           // Параметры ветра
	Visibility int       `json:"visibility"`     // Видимость (в метрах)
	Pop        float64   `json:"pop"`            // Вероятность осадков (0–1)
	Rain       *Rain     `json:"rain,omitempty"` // Объем осадков за 3 часа (может отсутствовать)
	Sys        Sys       `json:"sys"`            // Системная информация (день/ночь)
	DtTxt      string    `json:"dt_txt"`         // Время интервала в текстовом формате
}

// Main содержит основные погодные параметры.
type Main struct {
	Temp      float64 `json:"temp"`       // Температура (°C)
	FeelsLike float64 `json:"feels_like"` // Ощущаемая температура (°C)
	TempMin   float64 `json:"temp_min"`   // Минимальная температура (°C)
	TempMax   float64 `json:"temp_max"`   // Максимальная температура (°C)
	Pressure  int     `json:"pressure"`   // Давление (гПа)
	SeaLevel  int     `json:"sea_level"`  // Давление на уровне моря (гПа)
	GrndLevel int     `json:"grnd_level"` // Давление на уровне земли (гПа)
	Humidity  int     `json:"humidity"`   // Влажность (%)
	TempKf    float64 `json:"temp_kf"`    // Внутренний параметр точности прогноза
}

// Weather описывает погодные условия.
type Weather struct {
	ID          int    `json:"id"`          // Код погодного условия
	Main        string `json:"main"`        // Основной тип погоды (например, "Rain")
	Description string `json:"description"` // Описание погоды (например, "небольшой дождь")
	Icon        string `json:"icon"`        // Код иконки погоды
}

// Clouds содержит информацию об облачности.
type Clouds struct {
	All int `json:"all"` // Процент облачного покрова (0–100)
}

// Wind содержит параметры ветра.
type Wind struct {
	Speed float64 `json:"speed"` // Скорость ветра (м/с)
	Deg   int     `json:"deg"`   // Направление ветра (градусы)
	Gust  float64 `json:"gust"`  // Скорость порывов ветра (м/с)
}

// Rain содержит информацию об осадках.
type Rain struct {
	ThreeHour float64 `json:"3h"` // Объем осадков за 3 часа (мм)
}

// Sys содержит системную информацию.
type Sys struct {
	Pod string `json:"pod"` // Часть дня: "d" (день) или "n" (ночь)
}

// City содержит информацию о городе.
type City struct {
	ID         int    `json:"id"`         // ID города
	Name       string `json:"name"`       // Название города
	Coord      Coord  `json:"coord"`      // Координаты города
	Country    string `json:"country"`    // Код страны
	Population int    `json:"population"` // Население города
	Timezone   int    `json:"timezone"`   // Смещение времени (в секундах)
	Sunrise    int64  `json:"sunrise"`    // Время восхода солнца (Unix timestamp)
	Sunset     int64  `json:"sunset"`     // Время заката солнца (Unix timestamp)
}

// Coord содержит координаты города.
type Coord struct {
	Lat float64 `json:"lat"` // Широта
	Lon float64 `json:"lon"` // Долгота
}
