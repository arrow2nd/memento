<div align="center">

![Frame 11](https://github.com/user-attachments/assets/514d65e2-5926-42c8-a7a7-6c2771f55af0)
✨️ VRChat内で撮影した写真をワールド毎に自動整理するやつ

<a href="https://github.com/arrow2nd/memento/wiki">取扱説明書</a> |
<a href="https://github.com/arrow2nd/memento/releases/latest">ダウンロード</a>

</div>

## これなに

タスクバーに常駐して、VRChat 内で撮影した写真をリアルタイムに JPEG
へ変換・ワールド毎に整理するアプリケーションです。

### できること

- タスクバーに常駐してリアルタイムに処理
- `YYYY-MM/[撮影ワールド名]/` のディレクトリを作成して自動で整理
- JPEGへ変換し、撮影日時・ワールド名をEXIF情報として書込む

### できないこと

- 過去に撮影した写真の整理
- Windows以外での動作 (ビルドはできるはず)

## ビルド

> [!NOTE]
> Windows環境を想定しています。

```sh
make build
```

もしくは

```sh
go install github.com/tc-hib/go-winres@latest
go generate
go build -tags prod -ldflags="-H=windowsgui -s -w -X github.com/arrow2nd/memento/app.appVersion=v.x.x.x" -o "dist/memento_v.x.x.x.exe"
```

### アプリアイコンの埋め込み

[tc-hib/go-winres: Command line tool for adding Windows resources to executable files](https://github.com/tc-hib/go-winres)
を使っています。
