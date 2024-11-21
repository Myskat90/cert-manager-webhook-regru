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

// RegruClient представляет клиента для работы с API Reg.ru
type RegruClient struct {
	username string
	password string
	zone     string
}

// NewRegruClient создает новый клиент для работы с API Reg.ru
func NewRegruClient(username, password, zone string) *RegruClient {
	return &RegruClient{
		username: username,
		password: password,
		zone:     zone,
	}
}

// GetRecords получает список записей для указанной зоны
func (c *RegruClient) getRecords() ([]map[string]interface{}, error) {
	fmt.Println("getRecords: запрашиваем записи для зоны:", c.zone)
	endpoint := fmt.Sprintf("%s/zone/get_resource_records", defaultBaseURL)
	payload := map[string]interface{}{
		"username": c.username,
		"password": c.password,
		"dname":    c.zone,
	}

	response, err := c.makePostRequest(endpoint, payload)
	if err != nil {
		return nil, err
	}

	records, ok := response["answer"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("поле 'answer' отсутствует в ответе")
	}

	var result []map[string]interface{}
	for _, record := range records {
		if r, ok := record.(map[string]interface{}); ok {
			result = append(result, r)
		}
	}

	return result, nil
}

// createTXT создает TXT-запись в указанной зоне
func (c *RegruClient) createTXT(domain, value string) error {
	fmt.Println("createTXT: создаем TXT-запись для зоны:", c.zone, "домен:", domain, "значение:", value)
	endpoint := fmt.Sprintf("%s/zone/add_txt", defaultBaseURL)
	payload := map[string]interface{}{
		"username": c.username,
		"password": c.password,
		"dname":    c.zone,
		"name":     domain,
		"text":     value,
	}

	_, err := c.makePostRequest(endpoint, payload)
	return err
}

// deleteTXT удаляет TXT-запись в указанной зоне
func (c *RegruClient) deleteTXT(domain, value string) error {
	fmt.Println("deleteTXT: удаляем TXT-запись для зоны:", c.zone, "домен:", domain, "значение:", value)
	endpoint := fmt.Sprintf("%s/zone/remove_record", defaultBaseURL)
	payload := map[string]interface{}{
		"username": c.username,
		"password": c.password,
		"dname":    c.zone,
		"name":     domain,
		"text":     value,
		"type":     "TXT",
	}

	_, err := c.makePostRequest(endpoint, payload)
	return err
}

// makePostRequest выполняет POST-запрос с указанными параметрами
func (c *RegruClient) makePostRequest(endpoint string, payload map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println("makePostRequest: выполняем POST-запрос на:", endpoint)
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	client := &http.Client{Timeout: defaultHTTPClient}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("ошибка обработки JSON: %w", err)
	}

	if result["result"] != "success" {
		return nil, fmt.Errorf("API вернул ошибку: %v", result["error_text"])
	}

	return result, nil
}
