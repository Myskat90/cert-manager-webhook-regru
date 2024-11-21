package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultBaseURL    = "https://api.reg.ru/api/regru2"
	defaultHTTPClient = 10 * time.Second
)

// RegruClient предоставляет методы для взаимодействия с API Reg.ru.
type RegruClient struct {
	username string
	password string
	zone     string
}

// NewRegruClient создает новый экземпляр клиента для API Reg.ru с использованием логина и пароля.
func NewRegruClient(username string, password string, zone string) *RegruClient {
	return &RegruClient{
		username: username,
		password: password,
		zone:     zone,
	}
}

// getRecords получает записи DNS для указанной зоны.
func (c *RegruClient) getRecords() error {
	fmt.Println("getRecords: Запрашиваем записи DNS для зоны:", c.zone)

	payload := map[string]string{
		"domain": c.zone,
	}

	// Путь для получения записей теперь соответствует документации.
	response, err := c.makePostRequest("zone/get_resource_records", payload)
	if err != nil {
		return fmt.Errorf("ошибка выполнения getRecords: %w", err)
	}
	fmt.Println("getRecords: Ответ от API:", response)
	return nil
}

// createTXT добавляет TXT-запись в DNS-зону.
func (c *RegruClient) createTXT(domain, value string) error {
	fmt.Println("createTXT: Добавляем TXT-запись для домена:", domain)

	payload := map[string]string{
		"domain":    c.zone,
		"subdomain": domain,
		"txt":       value,
	}

	// Путь для добавления TXT записи теперь соответствует документации.
	response, err := c.makePostRequest("zone/add_txt", payload)
	if err != nil {
		return fmt.Errorf("ошибка выполнения createTXT: %w", err)
	}
	fmt.Println("createTXT: Ответ от API:", response)
	return nil
}

// deleteTXT удаляет TXT-запись из DNS-зоны.
func (c *RegruClient) deleteTXT(domain, value string) error {
	fmt.Println("deleteTXT: Удаляем TXT-запись для домена:", domain)

	payload := map[string]string{
		"domain":    c.zone,
		"subdomain": domain,
		"txt":       value,
	}

	// Путь для удаления TXT записи теперь соответствует документации.
	response, err := c.makePostRequest("zone/remove_record", payload)
	if err != nil {
		return fmt.Errorf("ошибка выполнения deleteTXT: %w", err)
	}
	fmt.Println("deleteTXT: Ответ от API:", response)
	return nil
}

// makePostRequest выполняет POST-запрос к API Reg.ru с логином и паролем в заголовке для авторизации.
func (c *RegruClient) makePostRequest(command string, payload map[string]string) (map[string]interface{}, error) {
	payload["input_format"] = "json"
	payload["output_format"] = "json"

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	url := fmt.Sprintf("%s/%s", defaultBaseURL, command)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Добавление логина и пароля в заголовок Authorization (Basic Auth)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: defaultHTTPClient}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// Логируем полный ответ от API
	fmt.Println("Ответ от API:", string(respBody))

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("ошибка десериализации JSON: %w", err)
	}

	// Проверяем результат выполнения API-запроса
	if result["result"] != "success" {
		return nil, fmt.Errorf("API вернул ошибку: %v", result["error_text"])
	}

	return result, nil
}
