package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sqweek/dialog"
)

func main() {
	// 1. エクスプローラーでフォルダを選択
	fmt.Println("使用するHTMLファイルが含まれるディレクトリを選択してください...")
	directory, err := dialog.Directory().Title("Webビルドの場所を選択").Browse()
	if err != nil {
		log.Fatalf("フォルダ選択がキャンセルされたか、エラーが発生しました: %v", err)
	}

	// 2. サーバー設定
	port := ":8060"
	fs := http.FileServer(http.Dir(directory))

	// 3. Godotに必要なヘッダーを付与するミドルウェア
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// これらが無いとGodot 4以降のWeb版は起動しません
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
		
		// キャッシュ防止（開発用）
		w.Header().Set("Cache-Control", "no-store")
		
		fs.ServeHTTP(w, r)
	})

	fmt.Printf("\n--- Server Started ---\n")
	fmt.Printf("選択されたパス: %s\n", directory)
	fmt.Printf("URL: http://localhost%s\n", port)
	fmt.Printf("終了するにはこのウィンドウを閉じるか Ctrl+C を押してください\n")

	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}