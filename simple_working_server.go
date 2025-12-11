package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Определяем порт
	port := "8080"

	// Создаем маршруты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Pet Shelter</title>
            <style>
                body { font-family: Arial; padding: 40px; }
                h1 { color: #667eea; }
                button { padding: 10px 20px; background: #667eea; color: white; border: none; cursor: pointer; }
            </style>
        </head>
        <body>
            <h1>🐾 Pet Shelter API - РАБОТАЕТ!</h1>
            <p>Сервер запущен на порту %s</p>
            <button onclick="testAPI()">Тест API</button>
            <div id="result" style="margin-top: 20px;"></div>
            <script>
                function testAPI() {
                    fetch('/api/animals')
                        .then(r => r.json())
                        .then(data => {
                            document.getElementById('result').innerHTML = 
                                '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
                        });
                }
            </script>
			<button onclick="testAPI1()">Тест API1</button>
            <div id="result" style="margin-top: 20px;"></div>
            <script>
                function testAPI1() {
                    fetch('/api/catalog')
                        .then(r => r.json())
                        .then(data => {
                            document.getElementById('result').innerHTML = 
                                '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
                        });
                }
            </script>
			<button onclick="testAPI2()">Тест API2</button>
            <div id="result" style="margin-top: 20px;"></div>
            <script>
                function testAPI2() {
                    fetch('/api/login')
                        .then(r => r.json())
                        .then(data => {
                            document.getElementById('result').innerHTML = 
                                '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
                        });
                }
            </script>
        </body>
        </html>
        `, port)
	})

	http.HandleFunc("/api/animals", func(w http.ResponseWriter, r *http.Request) {
		animals := []map[string]interface{}{
			{"id": 1, "name": "Барсик", "species": "Кошка", "age": 3},
			{"id": 2, "name": "Шарик", "species": "Собака", "age": 2},
			{"id": 3, "name": "Мурка", "species": "Кошка", "age": 1},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(animals)
	})

	http.HandleFunc("/api/catalog", func(w http.ResponseWriter, r *http.Request) {
		catalog := []map[string]interface{}{
			{"id": 1, "name": "Корм для собак", "price": 25.99},
			{"id": 2, "name": "Игрушка для кошек", "price": 5.99},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(catalog)
	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token":  "test-token-123",
			"status": "success",
		})
	})

	// Запускаем сервер
	log.Printf("✅ Сервер запущен!")
	log.Printf("🌐 Откройте: http://localhost:%s", port)

	// Автоматически открываем браузер через 1 секунду
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Printf("\n✨ Откройте браузер по адресу: http://localhost:%s\n", port)
		fmt.Println("   Или нажмите Ctrl+клик по ссылке выше")
	}()

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("❌ Ошибка запуска сервера:", err)
	}
}
