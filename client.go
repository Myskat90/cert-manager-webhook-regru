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
	defaultBaseURL = "https://api.reg.ru/api/regru2" // Базовый URL API
)

// RegruClient представляет клиента для работы с API Reg.ru
type RegruClient struct {
	Username   string
	Password   string
	Zone       string
	HTTPClient *http.Client
}

// NewRegruClient создает новый клиент для работы с API Reg.ru
func NewRegruClient(username, password, zone string) *RegruClient {
	return &RegruClient{
		Username:   username,
		Password:   password,
		Zone:       zone,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// getRecords получает все записи для указанной зоны
func (c *RegruClient) getRecords() error {
	fmt.Println("getRecords: отправляем запрос на получение записей для зоны", c.Zone)
	endpoint := fmt.Sprintf("%s/zone/get_resource_records", defaultBaseURL)
	payload := map[string]interface{}{
		"username": c.Username,
		"password": c.Password,
		"dname":    c.Zone,
	}

	fmt.Println("getRecords: запрос с параметрами:", payload)

	response, err := c.makePostRequest(endpoint, payload)
	if err != nil {
		fmt.Println("getRecords: ошибка выполнения запроса:", err)
		return err
	}

	// Проверяем наличие записей в ответе
	if answer, ok := response["answer"].([]interface{}); ok {
		fmt.Println("getRecords: получены записи", answer)
	} else {
		fmt.Println("getRecords: ошибка получения записей")
		return fmt.Errorf("ошибка получения записей: нет поля 'answer' в ответе")
	}

	return nil
}

// createTXT создает TXT-запись для указанной зоны
func (c *RegruClient) createTXT(name, content string) error {
	fmt.Println("createTXT: создаем TXT-запись с именем", name, "и содержимым", content)
	endpoint := fmt.Sprintf("%s/zone/add_txt", defaultBaseURL)
	payload := map[string]interface{}{
		"username": c.Username,
		"password": c.Password,
		"dname":    c.Zone,
		"name":     name,
		"text":     content,
	}

	fmt.Println("createTXT: отправляем запрос с параметрами", payload)

	_, err := c.makePostRequest(endpoint, payload)
	if err != nil {
		fmt.Println("createTXT: ошибка выполнения запроса:", err)
		return err
	}

	fmt.Println("createTXT: запись успешно создана")
	return nil
}

// deleteTXT удаляет TXT-запись для указанной зоны
func (c *RegruClient) deleteTXT(name string) error {
	fmt.Println("deleteTXT: удаляем TXT-запись с именем", name)
	endpoint := fmt.Sprintf("%s/zone/remove_record", defaultBaseURL)
	payload := map[string]interface{}{
		"username": c.Username,
		"password": c.Password,
		"dname":    c.Zone,
		"name":     name,
		"type":     "TXT",
	}

	fmt.Println("deleteTXT: отправляем запрос с параметрами", payload)

	_, err := c.makePostRequest(endpoint, payload)
	if err != nil {
		fmt.Println("deleteTXT: ошибка выполнения запроса:", err)
		return err
	}

	fmt.Println("deleteTXT: запись успешно удалена")
	return nil
}

// makePostRequest выполняет POST-запрос с параметрами к API Reg.ru
func (c *RegruClient) makePostRequest(endpoint string, payload map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println("makePostRequest: выполняется POST запрос на", endpoint, "с телом", payload)
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("makePostRequest: ошибка сериализации данных:", err)
		return nil, fmt.Errorf("ошибка сериализации данных: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("makePostRequest: ошибка создания запроса:", err)
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("makePostRequest: ошибка выполнения запроса:", err)
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("makePostRequest: ошибка, статус код:", resp.StatusCode)
		return nil, fmt.Errorf("ошибка, статус код: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("makePostRequest: ошибка чтения ответа:", err)
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("makePostRequest: ошибка обработки JSON:", err)
		return nil, fmt.Errorf("ошибка обработки JSON: %w", err)
	}

	if result["result"] != "success" {
		fmt.Println("makePostRequest: ошибка API:", result["error_text"])
		return nil, fmt.Errorf("ошибка API: %s", result["error_text"])
	}

	fmt.Println("makePostRequest: запрос успешно выполнен, результат:", result)
	return result, nil
}
