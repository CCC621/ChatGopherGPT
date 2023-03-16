package main

import (
    "context"
	"fmt"
	"log"
	"os"
	"strings"
    "time"
    "os/signal"
    "syscall"
    "net/http"

    "github.com/bwmarrin/discordgo"
	"github.com/PullRequestInc/go-gpt3"    
)

// 環境変数の読み込み
var (
    apiKey = os.Getenv("OPENAI_API_KEY")
    token = os.Getenv("DISCORD_TOKEN")
)

// chatGPTの呼び出し時間を延長
var httpClient = &http.Client{
    Timeout: time.Duration(600 * time.Second),
}

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

    // OpenAI APIの設定
	client := gpt3.NewClient(apiKey,gpt3.WithHTTPClient(httpClient))
    ctx := context.Background()
    
    // メッセージをChatGPTで処理する
    resp, err := client.Completion(ctx, gpt3.CompletionRequest{
        Prompt:    []string{m.Content},

        // MaxTokens: gpt3.IntPtr(30),
        // Stop:      []string{"."},
        // Echo:      true,
    })
    if err != nil {
        log.Fatalln(err)
        return
    }

    // レスポンスを扱いやすい値に変換
    msg := strings.TrimSpace(resp.Choices[0].Text)

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
