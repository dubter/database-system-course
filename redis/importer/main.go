package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

type Data struct {
	Id         string   `json:"id"`
	Index      int      `json:"index"`
	Guid       string   `json:"guid"`
	IsActive   bool     `json:"isActive"`
	Balance    string   `json:"balance"`
	Picture    string   `json:"picture"`
	Age        int      `json:"age"`
	EyeColor   string   `json:"eyeColor"`
	Name       string   `json:"name"`
	Gender     string   `json:"gender"`
	Company    string   `json:"company"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Address    string   `json:"address"`
	About      string   `json:"about"`
	Registered string   `json:"registered"`
	Latitude   float64  `json:"latitude"`
	Longitude  float64  `json:"longitude"`
	Tags       []string `json:"tags"`
	Friends    []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"friends"`
	Greeting      string `json:"greeting"`
	FavoriteFruit string `json:"favoriteFruit"`
}

var (
	password = os.Getenv("REDIS_PASSWORD")
	path     = os.Getenv("JSON_PATH")
)

func main() {
	// Подключение к кластеру Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Адреса узлов кластера
		Password: password,         // Пароль для аутентификации
	})

	// Проверка подключения
	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	// Чтение данных из JSON-файла
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var dataset []Data
	err = json.Unmarshal(byteValue, &dataset)
	if err != nil {
		log.Fatalf("Ошибка разбора JSON: %v", err)
	}

	// Сохранение каждого элемента в Redis
	for _, data := range dataset {
		// Преобразование структуры в JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Ошибка при преобразовании в JSON: %v", err)
			continue
		}

		// Сохранение в виде строки
		err = redisClient.Set(context.Background(), data.Id, jsonData, 0).Err()
		if err != nil {
			log.Printf("Ошибка при сохранении строки в Redis: %v", err)
		}

		maped := map[string]string{
			"id":            data.Id,
			"guid":          data.Guid,
			"balance":       data.Balance,
			"picture":       data.Picture,
			"eyeColor":      data.EyeColor,
			"name":          data.Name,
			"gender":        data.Gender,
			"company":       data.Company,
			"email":         data.Email,
			"phone":         data.Phone,
			"address":       data.Address,
			"about":         data.About,
			"registered":    data.Registered,
			"greeting":      data.Greeting,
			"favoriteFruit": data.FavoriteFruit,
		}

		// Сохранение в виде хэша
		err = redisClient.HSet(context.Background(), "hash:"+data.Id, maped).Err()
		if err != nil {
			log.Printf("Ошибка при сохранении хэша в Redis: %v", err)
		}

		// Сохранение в виде сортированного множества (ZSET)
		err = redisClient.ZAdd(context.Background(), "zset", &redis.Z{
			Score:  float64(data.Index),
			Member: maped,
		}).Err()

		// Сохранение в виде списка
		err = redisClient.LPush(context.Background(), "list", maped).Err()
		if err != nil {
			log.Printf("Ошибка при сохранении списка в Redis: %v", err)
		}
	}

	fmt.Println("Данные успешно сохранены в Redis")
}
