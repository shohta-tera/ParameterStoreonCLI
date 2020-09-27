# AWSのリソース(ParametrStore)をCLIから触る

AWSのリソースをCLIから触るためのレポジトリ  

カテゴリ: AWS, GO言語  

## 実行環境

WSL: Ubuntu 18.04

### aws

AWSクレデンシャルの設定はしているものとする(.awsディレクトリ)

```
aws-cli: 1.16.279
```

### GO

```
go1.13.4
GOPATHなどの初期設定は済ませて下さい
```

# レポジトリの中身

## ディレクトリ構成

```
HOME
└ handlers
  └ ParameterStoreUpdate
     └ main.go

```

### ParameterStoreUpdate

パラメータストアに格納したパラメータを一括で更新する対話型のCLI

# 実行方法

## 注意事項

1. リソースの値は決め打ちなので、改変する必要がある  
2. コンパイルして実行する手順については記載しない

## ParameterStoreUpdate

### 前提条件
SystemsManagerにパラメータが格納している

### 実行手順

必要なモジュールをインストールする

```
go get -v github.com/aws/aws-sdk-go/...
```

その後コンパイルもしくはコマンドを実行  

```
go run src/handlers/ParamterStoreUpdate/main.go
```
あとは指示に従って入力していくだけ
