package tests

import (
	"fmt"
	"strings"
	"testing"

	"resty.dev/v3"
)

func TestName(t *testing.T) {

	// 将字段名改为大写，并添加 json 标签
	resp := struct {
		Message string `json:"message"`
	}{}

	client := resty.New().SetBaseURL("https://example.k8s.shub.us/")
	var one, two int
	for i := 0; i < 100; i++ {
		client.R().SetResult(&resp).Get("")
		message := resp.Message
		if strings.Contains(message, "01") {
			one++
		} else if strings.Contains(message, "02") {
			two++
		}
	}
	fmt.Println("example-01 is ", one, ", example-02 is ", two, ". Ratio[01/02] is ", one, "/", two)
}

func Test(t *testing.T) {
	client := resty.New().SetBaseURL("https://www.storehubhq.com")

	for i := 0; i < 100; i++ {
		get, err := client.R().
			SetContentType("content-type:application/json;charset=utf-8").
			SetHeaders(map[string]string{
				"accept-language":     "en",
				"accept-encoding":     "gzip,deflate",
				"storehub-version":    "2.51.3.0",
				"storehub-business":   "kafe123",
				"storehub-registerid": "66c473291eac1300079fbf24",
				"storehub-token":      "694b3ed083df11f0b972dda5953608f2",
			}).
			SetQueryParams(map[string]string{
				"bn":               "kafe123",
				"registerId":       "66c473291eac1300079fbf24",
				"includeAllStores": "true",
			}).
			Get("/api/syncStoreInfo")

		if err != nil {
			t.Error(err)
		}
		fmt.Println(get)
	}
}
