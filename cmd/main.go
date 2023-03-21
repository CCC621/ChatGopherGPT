package main

import (
	"fmt"
    "log"
	"strings"
    "time"
    "os"
    "os/signal"
    "io"
	"io/ioutil"
	"net/http"
	"syscall"
	"bytes"
	"encoding/json"

    "github.com/bwmarrin/discordgo"
)

// 各種構造体の定義
type OpenAiRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAiResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
	Usages  Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Messages     Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// APIのURLを設定
const openaiURL = "https://api.openai.com/v1/chat/completions"

// 環境変数の読み込み
var (
    apiKey = os.Getenv("OPENAI_API_KEY")
    token = os.Getenv("DISCORD_TOKEN")
)

// chatGPTとのやりとりを保持する変数
var messages []Message

func main() {
    
    // Discord Botの設定
    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        fmt.Println("Error creating Discord session: ", err)
    }

    // メッセージ受信時のイベントハンドラを登録
    dg.AddHandler(messageCreate)

    // Discord botの起動
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error starting Discord bot: %v", err)
	}
	defer dg.Close()

    // 終了シグナルを受け取るためのチャネルを作成
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

    // 終了シグナルが送信されるまで待機
    fmt.Println("Bot is running...") 
    <-quit

}

// メッセージ受信時のイベントハンドラ
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

    // メッセージの送信者がBotだった場合は無視する
    if m.Author.ID == s.State.User.ID {
        return
    }

    // ログファイルをオープン（存在しない場合は作成）
    file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // メッセージ処理
    messages = append(messages, Message{
        Role:    "user",
        Content: m.Content,
    })

    // API Call
    response := getOpenAIResponse()

    // レスポンスを扱いやすい値に変換
    msg := strings.TrimSpace(response.Choices[0].Messages.Content)

    // ログファイルに受信したメッセージを追記
    if _, err := file.WriteString(time.Now().Format(time.Stamp) + "\t" + m.Author.Username + "  >  " + m.Content + "\n"); err != nil {
        log.Fatal(err)
    }

    // ログファイルに返信されるメッセージを追記
    if _, err := file.WriteString(time.Now().Format(time.Stamp) + "\t ChatGPT >  " + msg + "\n"); err != nil {
        log.Fatal(err)
    }

    // ChatGPTからのメッセージを返信する
    s.ChannelMessageSend(m.ChannelID, msg)
}

func getOpenAIResponse() OpenAiResponse {

	requestBody := OpenAiRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	requestJSON, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response OpenAiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		println("Error: ", err.Error())
		return OpenAiResponse{}
	}

	messages = append(messages, Message{
		Role:    "assistant",
		Content: response.Choices[0].Messages.Content,
	})

	return response
}