# memento

✨️ VRChat内で撮影した写真をワールド別に自動整理するやつ

## できること

- タスクバーに常駐して、撮影された写真をワールド別に自動整理

## できないこと

- 過去に撮影した写真の整理
- Windows以外での動作保証

## ビルド

Windows環境を想定しています。

```sh
go build -ldflags "-H=windowsgui"
```

### アプリアイコンの埋め込み

[tc-hib/go-winres: Command line tool for adding Windows resources to executable files](https://github.com/tc-hib/go-winres)
を使っています。

```sh
go install github.com/tc-hib/go-winres@latest
go-winres make
```
