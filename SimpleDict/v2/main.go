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
	"sync"
)

var wg sync.WaitGroup

// DictRequestCY
// @Description:用于构建向彩云翻译API发送请求的数据结构
type DictRequestCY struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

// DictResponseCY
// @Description: 彩云翻译API的响应数据结构
type DictResponseCY struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

// DictResponseBD
// @Description: 百度翻译API的响应数据结构
type DictResponseBD struct {
	Errno int `json:"errno"`
	Data  []struct {
		K string `json:"k"`
		V string `json:"v"`
	} `json:"data"`
}

// queryBD
//
//	@Description: 构建 HTTP 请求，使用百度翻译API查询给定单词的中文翻译
//	@param word
func queryBD(word string) {
	client := &http.Client{}
	var data = strings.NewReader("kw=" + word)
	req, err := http.NewRequest("POST", "https://fanyi.baidu.com/sug", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "BIDUPSID=1CDF4C78BCFB3D90B2D6A594920DA6E4; PSTM=1630726121; __yjs_duid=1_02f9d28d226309370f287ee032114e3f1630763189674; REALTIME_TRANS_SWITCH=1; HISTORY_SWITCH=1; FANYI_WORD_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; BAIDUID=396D4FB45E415311D3C31A0DE0D1AF80:FG=1; BDUSS=V6NWkyMzJnaTB-OWplSFdUR25tRkxsTzFnVnBDanN4ZEs2bHFMc2EzaXJOcDVrRVFBQUFBJCQAAAAAAAAAAAEAAACpaydUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKupdmSrqXZkMV; BDUSS_BFESS=V6NWkyMzJnaTB-OWplSFdUR25tRkxsTzFnVnBDanN4ZEs2bHFMc2EzaXJOcDVrRVFBQUFBJCQAAAAAAAAAAAEAAACpaydUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKupdmSrqXZkMV; APPGUIDE_10_6_2=1; MCITY=-75%3A; MAWEBCUID=web_jcjIPqwGUVRAlyJeYpGiILrRIWXcMWNDAjChgRXUZoHLpLyApr; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; delPer=0; PSINO=1; BAIDUID_BFESS=396D4FB45E415311D3C31A0DE0D1AF80:FG=1; BA_HECTOR=2h852l8h2ha5840k8k0k050g1iai2fd1o; ZFY=kYWZ1Y6d6SV8cMcgDFNKAhwEEhlIK7W5gu1C:APN07HM:C; H_PS_PSSID=36546_39026_39022_38858_38957_38954_39009_38918_38972_38815_38637_26350_39041_38948_39046; ab_sr=1.0.1_MjBiYTAxMmUwMWE2MzIyODYxY2JiYzIxNjFjMWY2OTIwOTc5NTA1MTI5MzAzZGI2OWNiMGEwM2I1Y2IwNjJiZWEzNDQwOTg1ZTkyYTJkNmU1OTc3MjY1MTJhZTM1MGEwNWNlM2NkM2VjNWM0NDMwZjY0ZWJlZTEyOTQyZjFkZjE1YzhiYzY4YjBhYzQ4NzE2NWI4MjNkZjA1NTVkMDk2MmRjMjZhYTA4NThmOTYwMmEzOWMxMjAxMGM2OTdjODhm")
	req.Header.Set("Origin", "https://fanyi.baidu.com")
	req.Header.Set("Referer", "https://fanyi.baidu.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.43")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Microsoft Edge";v="114"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
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
	var dictResponse DictResponseBD
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("百度翻译")
	for _, item := range dictResponse.Data {
		fmt.Println(item.V)
	}

	wg.Done()

}

// queryCY
//
//	@Description: 构建 HTTP 请求，使用彩云翻译API查询给定单词的中文翻译和词性解释。 发送 POST 请求到彩云翻译API，将要查询的单词作为参数。 解析 API 响应，提取出中文翻译和词性解释信息，并输出到命令行。
//
//	@param word
func queryCY(word string) {
	client := &http.Client{}
	request := DictRequestCY{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
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
	var dictResponse DictResponseCY
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("彩云翻译")
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}

	wg.Done()
}

func main() {

	fmt.Println("请输入要查询的内容:")
	reader := bufio.NewReader(os.Stdin)
	word, _ := reader.ReadString('\n')
	word = strings.Trim(word, "\r\n")

	wg.Add(2)

	go queryCY(word)

	go queryBD(word)

	wg.Wait()
}
