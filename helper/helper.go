package helper

import (
	"encoding/base64"
	"encoding/json"
	"github.com/simabdi/vodka-authservice/models"
	"log"
	"os"
	"time"
)

const (
	LayoutISO          = "2006-01-02"
	LayoutISO24Hour    = "2006-01-02 15:04:05"
	LayoutISO12Hour    = "2006-01-02 03:04:05"
	LayoutISO24HourLog = "2006-01-02_150405"
	LayoutISO12HourLog = "2006-01-02_030405"
	LayoutID           = "01-02-2006"
)

func JsonResponse(code int, message string, success bool, error string, data interface{}) models.Response {
	meta := models.Meta{
		Code:    code,
		Status:  success,
		Message: message,
		Error:   error,
	}

	response := models.Response{
		Meta: meta,
		Data: data,
	}

	return response
}

func Logger(requestType string, request interface{}, bodyBytes []byte) {
	fo, err := os.OpenFile("storage/logs/log-"+time.Now().Format("2006-01-02")+".txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	req, err := json.Marshal(&request)
	if err != nil {
		log.Println("[Error Log Marshal] : ", err.Error())
	}

	var resJson map[string]interface{}
	err = json.Unmarshal(bodyBytes, &resJson)
	if err != nil {
		log.Println("[Error Log Unmarshal] : ", err.Error())
	}

	resDecode, err := json.Marshal(resJson)
	if err != nil {
		log.Println("[Error Log res Marshal] : ", err.Error())
	}

	text := []byte(
		"===========================================================================================\n" +
			requestType + " " + time.Now().Format("2006-01-02 15:04:05") +
			"\n===========================================================================================\n" +
			"=================================REQUEST===================================\n" +
			string(req) + "\n" +
			"=================================RESPONSE==================================\n" +
			string(resDecode) + "\n\n\n")

	_, err = fo.WriteString(string(text))
	if err != nil {
		log.Println("[Error Log WriteString] : ", err.Error())
	}

	defer fo.Close()
}

func GetFormattedDate(date time.Time, format string) string {
	return date.Format(format)
}

func GetDate(format string) string {
	date := time.Now()
	return date.Format(format)
}

func ParseDate(s string, format string) time.Time {
	date, _ := time.Parse(format, s)
	return date
}

func Std64Encode(plainText string) string {
	return base64.StdEncoding.EncodeToString([]byte(plainText))
}

func Std64Decode(encoded string) string {
	decodedByte, _ := base64.StdEncoding.DecodeString(encoded)
	return string(decodedByte)
}
