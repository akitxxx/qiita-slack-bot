# 元となるイメージを指定
FROM golang:latest

# 作業ディレクトリを指定
WORKDIR /go/src/qiita-bot
# 左側(ローカル)のディレクトリをイメージの作業ディレクトリにコピー
COPY . .

# イメージのビルド時に実行するコマンド
# depのインストールとdep ensureで依存関係のパッケージをインストール
RUN go get -u github.com/golang/dep/cmd/dep \
  && dep ensure

# freshのインストール
RUN go get github.com/pilu/fresh
# イメージからコンテナを作成する際に実行
CMD ["fresh"]
