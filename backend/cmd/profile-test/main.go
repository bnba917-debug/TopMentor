package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	post := func(path string, body interface{}) []byte {
		b, _ := json.Marshal(body)
		resp, err := http.Post("http://localhost:8080/api/v1"+path, "application/json", bytes.NewReader(b))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		out, _ := io.ReadAll(resp.Body)
		fmt.Printf("POST %s -> %d %s\n", path, resp.StatusCode, string(out))
		return out
	}

	req := func(method, path, token string, body interface{}) {
		var r io.Reader
		if body != nil {
			b, _ := json.Marshal(body)
			r = bytes.NewReader(b)
		}
		httpReq, _ := http.NewRequest(method, "http://localhost:8080/api/v1"+path, r)
		if body != nil {
			httpReq.Header.Set("Content-Type", "application/json")
		}
		httpReq.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		out, _ := io.ReadAll(resp.Body)
		fmt.Printf("%s %s -> %d %s\n", method, path, resp.StatusCode, string(out))
	}

	post("/auth/sms/send", map[string]string{"phone": "13800000001"})
	raw := post("/auth/sms/login", map[string]string{"phone": "13800000001", "code": "123456"})

	var login struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &login); err != nil || login.Data.Token == "" {
		panic("no token")
	}
	token := login.Data.Token

	req("GET", "/mentor/profile", token, nil)
	req("PUT", "/mentor/profile", token, map[string]interface{}{
		"real_name": "张明", "school_name": "清华大学", "major": "计算机科学",
		"gender": "male", "english_score": "高考148", "bio": "脚本测试", "tags": []string{"阳光幽默"},
		"avatar_url": "", "intro_video_url": "https://example.com/videos/mentor1.mp4",
	})
}
