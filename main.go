package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type Company struct {
	Name       string
	ListingID  string
	ListingURL string
	WebsiteURL string
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://appexchangejp.salesforce.com/consulting"),
		chromedp.Sleep(5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	// すべてのリストを読み込む
	if err := clickLoadMoreUntilDone(ctx); err != nil {
		log.Println("❌ データ読み込み失敗:", err)
	}

	// ページHTML取得
	var html string
	err = chromedp.Run(ctx, chromedp.OuterHTML("body", &html))
	if err != nil {
		log.Fatal(err)
	}

	// 取得済みIDの読み込み
	fetched := loadFetchedIDs("fetched.txt")

	// HTMLから会社情報抽出（最初の28件スキップ）
	allEntries := strings.Split(html, `appx-tile appx-tile-consultant`)
	companies := []Company{}
	for i, entry := range allEntries {
		if i < 29 { // 先頭28件をスキップ
			continue
		}
		id := extractAttr(entry, `data-listing-id="`)
		if id == "" || fetched[id] {
			continue
		}
		name := htmlUnescape(extractAttr(entry, `data-listing-name="`))
		url := extractAttr(entry, `data-listing-url="`)
		if name != "" && url != "" {
			companies = append(companies, Company{
				Name:       name,
				ListingID:  id,
				ListingURL: url,
			})
		}
	}

	fmt.Printf("🔎 新規取得対象企業数: %d\n", len(companies))

	// 詳細ページへアクセスしてWebsite URL取得
	for i, c := range companies {
		delay := 3 + rand.Intn(3)
		fmt.Printf("[%d/%d] %s にアクセス中...（%d秒待機）\n", i+1, len(companies), c.Name, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		var detailHTML string
		err := chromedp.Run(ctx,
			chromedp.Navigate(c.ListingURL),
			chromedp.Sleep(3*time.Second),
			chromedp.OuterHTML("body", &detailHTML),
		)
		if err == nil {
			website := extractAttr(detailHTML, `data-event="listing-publisher-website" href="`)
			companies[i].WebsiteURL = website
		}

		// fetched に保存
		appendFetchedID("fetched.txt", c.ListingID)
	}

	// CSV出力
	if err := writeCSV("result.csv", companies); err != nil {
		log.Fatal("❌ CSV出力失敗:", err)
	}
	fmt.Println("✅ CSV出力完了: result.csv")
}

// JavaScriptで「もっと見る」を繰り返し呼び出し
func clickLoadMoreUntilDone(ctx context.Context) error {
	for i := 0; i < 20; i++ {
		var initialCount int
		err := chromedp.Run(ctx,
			chromedp.Evaluate(`document.querySelectorAll(".appx-tile.appx-tile-consultant").length`, &initialCount),
		)
		if err != nil {
			return fmt.Errorf("初期件数の取得失敗: %w", err)
		}

		fmt.Printf("📦 [%d] 現在の件数: %d → JSでロード実行\n", i+1, initialCount)

		err = chromedp.Run(ctx,
			chromedp.Evaluate(`loadMoreListingsJS();`, nil),
		)
		if err != nil {
			fmt.Println("✅ JS実行失敗か終了済。")
			break
		}

		success := false
		for w := 0; w < 20; w++ {
			time.Sleep(500 * time.Millisecond)
			var currentCount int
			chromedp.Run(ctx,
				chromedp.Evaluate(`document.querySelectorAll(".appx-tile.appx-tile-consultant").length`, &currentCount),
			)
			if currentCount > initialCount {
				success = true
				break
			}
		}
		if !success {
			fmt.Println("⏱️ データが増加しません。終了。")
			break
		}
	}
	return nil
}

// 属性抽出
func extractAttr(s, prefix string) string {
	idx := strings.Index(s, prefix)
	if idx == -1 {
		return ""
	}
	start := idx + len(prefix)
	end := strings.Index(s[start:], `"`)
	if end == -1 {
		return ""
	}
	return s[start : start+end]
}

// HTMLエスケープ解除
func htmlUnescape(s string) string {
	replacer := strings.NewReplacer(
		"&#x2F;", "/",
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", `"`,
		"&#39;", `'`,
	)
	return replacer.Replace(s)
}

// fetched.txt読み込み
func loadFetchedIDs(path string) map[string]bool {
	result := make(map[string]bool)
	if b, err := os.ReadFile(path); err == nil {
		lines := strings.Split(string(b), "\n")
		for _, line := range lines {
			if line != "" {
				result[line] = true
			}
		}
	}
	return result
}

// fetched.txtに追記
func appendFetchedID(path, id string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		f.WriteString(id + "\n")
	}
}

// CSV出力
func writeCSV(path string, companies []Company) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"会社名", "詳細ページURL", "WebサイトURL"})
	for _, c := range companies {
		writer.Write([]string{c.Name, c.ListingURL, c.WebsiteURL})
	}
	return nil
}
