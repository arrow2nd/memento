<div align="center">

![Frame 11](https://github.com/user-attachments/assets/514d65e2-5926-42c8-a7a7-6c2771f55af0)
✨️ VRChat内で撮影した写真をワールド別に自動整理するやつ

</div>



## できること

- タスクバーに常駐して、撮影された写真をワールド別に自動整理

## できないこと

- 過去に撮影した写真の整理
- Windows以外での動作保証

## ビルド

Windows環境を想定しています。

```sh
go build -tags prod -ldflags="-H=windowsgui -s -w -X github.com/arrow2nd/memento/app.appVersion=v.x.x.x" -o "dist/memento_v.x.x.x.exe"
```

### アプリアイコンの埋め込み

[tc-hib/go-winres: Command line tool for adding Windows resources to executable files](https://github.com/tc-hib/go-winres)
を使っています。

```sh
go install github.com/tc-hib/go-winres@latest
go-winres make
```
