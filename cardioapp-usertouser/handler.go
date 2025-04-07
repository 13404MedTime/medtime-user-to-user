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
