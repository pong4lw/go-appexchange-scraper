# go-appexchange-scraper

Salesforce AppExchange（https://appexchange.salesforce.com）からアプリ情報をスクレイピングする Go 製のツールです。`chromedp` を使用してブラウザ操作を自動化し、CSVファイルとしてデータを出力します。

## 🧩 概要

このツールは、AppExchange に掲載されているアプリ情報（名前、ID、URLなど）を企業名ベースで検索し、自動的に収集します。マーケティングや競合調査、リード獲得などに活用できます。

## 🚀 主な機能

- AppExchange上での企業・アプリ情報の自動検索と取得
- 以下の情報をCSV形式で出力
  - アプリ名
  - リスティングID
  - リスティングURL
  - ウェブサイトURL
- クロール間にランダムディレイを挿入し、Bot検出を回避
- 単体で動作するシンプルなGoプログラム

## 📦 インストール方法

```bash
git clone https://github.com/pong4lw/go-appexchange-scraper.git
cd go-appexchange-scraper
go mod tidy
```

🔧 使い方
```
go run main.go
```
実行後、result.csv に以下の形式で結果が出力されます：

fetched.txt
```
ListingId
```

csv
```
Name,ListingURL,WebsiteURL
```

🛠 必要な環境
```
Go 1.18 以上
Google Chrome（インストール済であること）
インターネット接続
```

📁 プロジェクト構成

```
.
├── main.go         // メイン処理
├── go.mod          // Go モジュール設定
├── fetched.txt     // 実行済みの企業ID（取得し直したい場合は対象のIDを削除）
├── result.csv      // 出力ファイル（スクレイピング結果）
```

⚠️ 注意事項
このツールは研究・学習目的で作成されたものです。
実際にウェブサイトをスクレイピングする際は、対象サイトの利用規約やロボッツ.txt を必ずご確認ください。
利用は自己責任でお願いいたします。

