package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

type DictRequestHS struct {
	Source         string   `json:"source"`
	Words          []string `json:"words"`
	SourceLanguage string   `json:"source_language"`
	TargetLanguage string   `json:"target_language"`
}

type DictResponseHS struct {
	Details []struct {
		Detail string `json:"detail"`
		Extra  string `json:"extra"`
	} `json:"details"`
}

type DictResponseHSData struct {
	Result []struct {
		Ec struct {
			Basic struct {
				Explains []struct {
					Pos   string `json:"pos"`
					Trans string `json:"trans"`
				} `json:"explains"`
			} `json:"basic"`
		} `json:"ec"`
	} `json:"result"`
}

func query() {
	for {
		fmt.Println("请输入要查询的内容:")
		reader := bufio.NewReader(os.Stdin)
		word, _ := reader.ReadString('\n')
		word = strings.Trim(word, "\r\n")
		if IsEnglishString(word) {
			queryHS(word)
			break
		} else {
			fmt.Println("请输入英语")
		}
	}

}
func queryHS(word string) {

	client := &http.Client{}
	request := DictRequestHS{"youdao", []string{word}, "en", "zh"}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/detail/v1/?msToken=&X-Bogus=DFSzswVOQDaibrQ3tJHN7cppgiFh&_signature=_02B4Z6wo00001g0lO6gAAIDD-FrRNX0w-.4NJT8AAOfuf7", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16888064002651706; i18next=zh-CN; s_v_web_id=verify_ljtrq6kx_UW3ieIzP_8gQX_4abc_B8D8_AoHwuLysn026; ttcid=db98bce9149b4f09b905a71503d9331e36")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/?category=&home_language=zh&source_language=detect&target_language=zh&text=bad")
	req.Header.Set("sec-ch-ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Microsoft Edge";v="114"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.43")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}

	var dictResponse DictResponseHS

	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("火山翻译")
	item := dictResponse.Details[0]
	jsonStr := item.Detail

	var HSData DictResponseHSData
	err = json.Unmarshal([]byte(jsonStr), &HSData)
	if err != nil {
		panic(err)
	}

	for _, item := range HSData.Result[0].Ec.Basic.Explains {
		fmt.Println(item.Pos, item.Trans)
	}
}

func IsEnglishString(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Scripts["Han"], r) {
			continue
		} else {
			return false
		}
	}
	return true
}

func main() {
	defer func() {
		if msg := recover(); msg != nil {
			fmt.Println(msg, "输入不合法,请重新输入")
		}
	}()
	query()
}
