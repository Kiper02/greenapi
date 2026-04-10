package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	greenapi "github.com/green-api/max-api-client-golang"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/api/getSettings", getSettingsHandler)
	http.HandleFunc("/api/getStateInstance", getStateInstanceHandler)
	http.HandleFunc("/api/sendMessage", sendMessageHandler)
	http.HandleFunc("/api/sendFileByUrl", sendFileByUrlHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер запущен на порту %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

func getGreenAPIClient(r *http.Request) (*greenapi.GreenAPI, error) {
	idInstance := r.URL.Query().Get("idInstance")
	apiTokenInstance := r.URL.Query().Get("apiTokenInstance")

	if idInstance == "" || apiTokenInstance == "" {
		return nil, &appError{"Не указаны idInstance или apiTokenInstance", http.StatusBadRequest}
	}

	client := greenapi.GreenAPI{
		APIURL:           "https://api.green-api.com",
		MediaURL:         "https://api.green-api.com",
		IDInstance:       idInstance,
		APITokenInstance: apiTokenInstance,
	}
	return &client, nil
}

func getSettingsHandler(w http.ResponseWriter, r *http.Request) {
	client, err := getGreenAPIClient(r)
	if err != nil {
		http.Error(w, err.Error(), err.(*appError).Code)
		return
	}

	resp, err := client.Account().GetSettings()
	if err != nil {
		http.Error(w, "Ошибка вызова API: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, resp)
}

func getStateInstanceHandler(w http.ResponseWriter, r *http.Request) {
	client, err := getGreenAPIClient(r)
	if err != nil {
		http.Error(w, err.Error(), err.(*appError).Code)
		return
	}

	resp, err := client.Account().GetStateInstance()
	if err != nil {
		http.Error(w, "Ошибка вызова API: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, resp)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	client, err := getGreenAPIClient(r)
	if err != nil {
		http.Error(w, err.Error(), err.(*appError).Code)
		return
	}

	var reqBody struct {
		ChatId  string `json:"chatId"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if reqBody.ChatId == "" || reqBody.Message == "" {
		http.Error(w, "Не указаны chatId или message", http.StatusBadRequest)
		return
	}

	resp, err := client.Sending().SendMessage(reqBody.ChatId, reqBody.Message)
	if err != nil {
		http.Error(w, "Ошибка вызова API: "+err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, resp)
}

func sendFileByUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	client, err := getGreenAPIClient(r)
	if err != nil {
		http.Error(w, err.Error(), err.(*appError).Code)
		return
	}

	var reqBody struct {
		ChatId   string `json:"chatId"`
		UrlFile  string `json:"urlFile"`
		FileName string `json:"fileName"`
		Caption  string `json:"caption"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if reqBody.ChatId == "" || reqBody.UrlFile == "" {
		http.Error(w, "Не указаны chatId или urlFile", http.StatusBadRequest)
		return
	}

	var captionOpt greenapi.SendFileByUrlOption
	if reqBody.Caption != "" {
		captionOpt = greenapi.OptionalCaptionSendUrl(reqBody.Caption)
	}
	resp, err := client.Sending().SendFileByUrl(reqBody.ChatId, reqBody.UrlFile, reqBody.FileName, captionOpt)

	if err != nil {
		http.Error(w, "Ошибка вызова API: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, resp)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

type appError struct {
	Message string
	Code    int
}

func (e *appError) Error() string {
	return e.Message
}
