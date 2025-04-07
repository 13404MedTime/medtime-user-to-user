package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cast"
)

const (
	apiKey = "P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS"

	createUrl = "https://api.auth.u-code.io/v2/user"
	deleteUrl = "https://api.auth.u-code.io/v2/user"
	checkUrl  = "https://api.auth.u-code.io/v2/user/check"
)

// func main() {
// 	str := `{"data":{"additional_parameters":[],"app_id":"P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS","method":"CREATE","object_data":{"category_id":"aa06cc06-d0b7-4afb-8124-154455c6f863","desc":"\u003cp\u003e\u003cspan style=\"color: rgb(55, 55, 55);\"\u003eХамраев Аброр Асрорович- гастроэнтеролог, врач высшей категории, доктор медицинских наук, профессор. Проходил стажировку в России. Разработано и опубликовано более 100 научных работ. Является постоянным участником международных конференций. Руководитель отделения Гастроэнтерологии.\u003c/span\u003e\u003c/p\u003e\u003cp\u003e\u003cspan style=\"color: rgb(55, 55, 55);\"\u003eСпециализируется на диагностике и лечении заболеваний желудочно-кишечного тракта и гепатобилиарной системы.\u003c/span\u003e\u003c/p\u003e","doctor_name":"Хамраев Аброр Асрорович","experience":12,"guid":"c650875d-28f2-4211-927b-cc25b2170f62","hospital":"City Med","login":"test_doctor","medic_photo":"https://cdn.u-code.io/ucode/775cb374-dc49-4399-b3de-a697c88d531c_entrant__career-guidance__about-the-profession-1.jpg","multi":[],"password":"test_doctor","phone_number":"+998933332323","telegram_nick":"test_doctor"},"object_ids":["c650875d-28f2-4211-927b-cc25b2170f62"],"table_slug":"doctor","user_id":"c384727f-408a-4dd0-b242-db3cc9355edd"}}`
// 	fmt.Println(Handle([]byte(str)))
// }

// Handle a serverless request
func Handle(req []byte) string {
	// Send2(string(req))
	var (
		request NewRequestBody
	)

	if err := json.Unmarshal(req, &request); err != nil {
		return Handler("error", "error 1")
	}

	if cast.ToString(request.Data["method"]) == "CREATE" {
		if err := CreateUser(request); err != nil {
			return Handler("error", "eror from CREATE")
		}
	}

	return Handler("OK", "OK")
}

func CreateUser(request NewRequestBody) error {

	// ! firsst we must check user: is it existed in auth of the ucode

	var (
		objectData        = cast.ToStringMap(request.Data["object_data"])
		userGuid          = cast.ToString(objectData["guid"])
		responseCheckUser ResponseUserModel
		checkRequestBody  = map[string]interface{}{
			"email":         "",
			"login":         objectData["login"],
			"phone":         "",
			"resource_type": 1,
		}
	)
	responseBodyCheck, err := DoRequest(checkUrl, "POST", checkRequestBody)
	if err != nil {
		Handler("error", "error 2")
		return err
	}

	if err := json.Unmarshal(responseBodyCheck, &responseCheckUser); err != nil {
		// Handler("error", " error 3"+string(responseBodyCheck))
		// return err

		var (
			responseUser ResponseUserModel

			requestBody = map[string]interface{}{
				"active":                  1,
				"client_type_id":          "1cbb3a57-2b66-4fa2-8466-4882b2ef9b2e",
				"role_id":                 "48f1fe05-7db7-4fa3-9319-4988c2adfe57",
				"login":                   objectData["login"],
				"name":                    objectData["doctor_name"],
				"password":                objectData["password"],
				"phone":                   objectData["phone_number"],
				"project_id":              "a4dc1f1c-d20f-4c1a-abf5-b819076604bc",
				"resource_type":           0,
				"year_of_birth":           "",
				"base_url":                "",
				"client_platform_id":      "",
				"company_id":              "",
				"email":                   "",
				"expires_at":              "",
				"photo_url":               "",
				"resource_environment_id": "",
			}
		)

		responseBody, err := DoRequest(createUrl, "POST", requestBody)
		if err != nil {
			Handler("error", "error 4")
			return err
		}

		if err := json.Unmarshal(responseBody, &responseUser); err != nil {
			Handler("error", " error 5"+string(responseBody))
			return err
		}

		userDeleteUrl := fmt.Sprintf("https://api.admin.u-code.io/v1/object/doctor/%s", responseUser.Data.ID)
		_, err = DoRequest(userDeleteUrl, "DELETE", Request{Data: map[string]interface{}{}})
		if err != nil {
			Handler("error", "error 6")
			return err
		}

		var (
			userUpdateBody = Request{
				Data: map[string]interface{}{
					"guid":      userGuid,
					"auth_guid": responseUser.Data.ID,
				},
			}
		)

		_, err = DoRequest("https://api.admin.u-code.io/v1/object/doctor", "PUT", userUpdateBody)
		if err != nil {
			Handler("error", "error 7")
			return err
		}
	} else {

		var (
			userUpdateBody = Request{
				Data: map[string]interface{}{
					"guid":      userGuid,
					"auth_guid": responseCheckUser.Data.ID,
				},
			}
		)

		_, err = DoRequest("https://api.admin.u-code.io/v1/object/doctor", "PUT", userUpdateBody)
		if err != nil {
			Handler("error", "error 7")
			return err
		}
	}

	return nil
}

// ! MAKE MESSAGE FOR SENDING
func Handler(status, message string) string {

	var (
		response Response
		Message  = make(map[string]interface{})
	)

	// sendMessage("cardio user-to-user", status, message)
	response.Status = status
	Message["message"] = message
	respByte, _ := json.Marshal(response)
	return string(respByte)
}

func DoRequest(url string, method string, body interface{}) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if method == "PUT" || method == "DELETE" {
		request.Header.Add("authorization", "API-KEY")
		request.Header.Add("X-API-KEY", apiKey)
	}
	if method == "POST" {
		request.Header.Add("Resource-Id", "a97e8954-5d8e-4469-a241-9a9af2ea2978")
		request.Header.Add("Environment-Id", "dcd76a3d-c71b-4998-9e5c-ab1e783264d0")
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}

func Send(text string) {
	bot, _ := tgbotapi.NewBotAPI("6241555505:AAHPpkXj-oHBGblWd_7O9kxc9a05tJUIFRw")

	msg := tgbotapi.NewMessage(1194897882, text)

	bot.Send(msg)
}

func Send2(text string) {
	bot, _ := tgbotapi.NewBotAPI("6519849383:AAHw5BnPuFvtER6MNW6cNgcrVG6bMvElgac")

	msg := tgbotapi.NewMessage(1546926238, text)

	bot.Send(msg)
}
