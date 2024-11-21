package main

import (
	"bytes"
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

// NewRegruClient создает новый экземпляр клиента для API Reg.ru.
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

	payload := map[string]interface{}{
		"auth": map[string]string{
			"username": c.username,
			"password": c.password,
		},
		"zone":    c.zone,
		"command": "get_records",
	}

	response, err := c.makePostRequest(payload)
	if err != nil {
		return fmt.Errorf("ошибка выполнения getRecords: %w", err)
	}
	fmt.Println("getRecords: Ответ от API:", response)
	return nil
}

// createTXT добавляет TXT-запись в DNS-зону.
func (c *RegruClient) createTXT(domain, value string) error {
	fmt.Println("createTXT: Добавляем TXT-запись для домена:", domain)

	payload := map[string]interface{}{
		"auth": map[string]string{
			"username": c.username,
			"password": c.password,
		},
		"zone":    c.zone,
		"command": "add_record",
		"record": map[string]interface{}{
			"type":  "TXT",
			"name":  domain,
			"value": value,
		},
	}

	response, err := c.makePostRequest(payload)
	if err != nil {
		return fmt.Errorf("ошибка выполнения createTXT: %w", err)
	}
	fmt.Println("createTXT: Ответ от API:", response)
	return nil
}

// deleteTXT удаляет TXT-запись из DNS-зоны.
func (c *RegruClient) deleteTXT(domain, value string) error {
	fmt.Println("deleteTXT: Удаляем TXT-запись для домена:", domain)

	payload := map[string]interface{}{
		"auth": map[string]string{
			"username": c.username,
			"password": c.password,
		},
		"zone":    c.zone,
		"command": "del_record",
		"record": map[string]interface{}{
			"type":  "TXT",
			"name":  domain,
			"value": value,
		},
	}

	response, err := c.makePostRequest(payload)
	if err != nil {
		return fmt.Errorf("ошибка выполнения deleteTXT: %w", err)
	}
	fmt.Println("deleteTXT: Ответ от API:", response)
	return nil
}

// makePostRequest выполняет POST-запрос к API Reg.ru.
func (c *RegruClient) makePostRequest(payload map[string]interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	req, err := http.NewRequest("POST", defaultBaseURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.password)

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
