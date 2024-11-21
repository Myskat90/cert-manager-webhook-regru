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
	defaultBaseURL = "https://api.reg.ru/api/regru2"
	httpTimeout    = 10 * time.Second
)

var client = &http.Client{Timeout: httpTimeout}

type RegruClient struct {
	username string
	password string
	zone     string
}

// NewRegruClient создает новый экземпляр клиента для работы с API
func NewRegruClient(username string, password string, zone string) *RegruClient {
	return &RegruClient{
		username: username,
		password: password,
		zone:     zone,
	}
}

// sendRequest отправляет POST-запрос с параметрами и возвращает ответ
func (c *RegruClient) sendRequest(method string, params map[string]interface{}) ([]byte, error) {
	url := defaultBaseURL + "/" + method
	params["username"] = c.username
	params["password"] = c.password

	// Преобразуем параметры в JSON
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать данные в JSON: %v", err)
	}

	// Создаем POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ: %v", err)
	}

	// Логируем ответ API
	fmt.Println("Ответ API:", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка API: %s", string(respBody))
	}

	return respBody, nil
}

// CreateTXT создает TXT запись для домена
func (c *RegruClient) createTXT(domain string, value string) error {
	params := map[string]interface{}{
		"zone":   c.zone,
		"domain": domain,
		"txt":    value,
	}

	_, err := c.sendRequest("zone/add_txt", params)
	if err != nil {
		return fmt.Errorf("createTXT: API вернуло ошибку: %v", err)
	}

	return nil
}

// DeleteTXT удаляет TXT запись для домена
func (c *RegruClient) deleteTXT(domain string, value string) error {
	params := map[string]interface{}{
		"zone":   c.zone,
		"domain": domain,
		"txt":    value,
	}

	_, err := c.sendRequest("zone/remove_record", params)
	if err != nil {
		return fmt.Errorf("deleteTXT: API вернуло ошибку: %v", err)
	}

	return nil
}

// getRecords получает TXT записи для зоны
func (c *RegruClient) getRecords() ([]byte, error) {
	params := map[string]interface{}{
		"zone": c.zone,
	}

	// Запрашиваем записи для зоны
	respBody, err := c.sendRequest("zone/get_resource_records", params)
	if err != nil {
		return nil, fmt.Errorf("getRecords: ошибка при получении записей: %v", err)
	}

	return respBody, nil
}
