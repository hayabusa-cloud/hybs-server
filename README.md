┏┓╋╋╋╋╋╋╋╋╋╋╋┏┓  
┃┃╋╋╋╋╋╋╋╋╋╋╋┃┃  
┃┗━┳━━┳┓╋┏┳━━┫┗━┳┓┏┳━━┳━━┓  
┃┏┓┃┏┓┃┃╋┃┃┏┓┃┏┓┃┃┃┃━━┫┏┓┃  
┃┃┃┃┏┓┃┗━┛┃┏┓┃┗┛┃┗┛┣━━┃┏┓┃  
┗┛┗┻┛┗┻━┓┏┻┛┗┻━━┻━━┻━━┻┛┗┛  
╋╋╋╋╋╋┏━┛┃  
╋╋╋╋╋╋┗━━┛  
   
---
## 「hybs-server」とは？
- 「hybs-server」（以下はやぶさサーバー）は「hayabusa」フレームワークを使った、Web APIとリアルタイム通信を統合したゲームサーバサイドアーキテクチャです。
- 設計プロセスと実装プロセスが手軽に分けられます。
- 使いやすさとハイパフォーマンスを両立しています。
- HTTP/1.1、HTTP/3、RUDPプロトコルが対応されてます。
- 対応されているデータベース等：MySQL、Mongodb、Redis。
- MITライセンスのもとで無料配布しています。
---
## サーバー構成
|サーバー|役割  
|----|----  
|ゲート|簡単なリバースプロキシとロードバランシングを行う
|プラットフォーム|ユーザーIDの発行・認証などを行う
|ゲーム|ゲームAPIリクエストを実行する
|セントラル|リアルタイム通信サーバーに使うIDの発行・認証などを行う
|リアルタイム|RUDPソケット通信サーバー

---
## インストール
- ステップ1：はやぶさフレームワークをインストールします
```bash
go get -u github.com/hayabusa-cloud/hayabusa 
```
- ステップ2：はやぶさサーバーをインストールします
```bash
go get -u github.com/hayabusa-cloud/hybs-server
```

### 環境の準備
- Docker Composeを用意しておきます：
https://docs.docker.jp/compose/install.html

### ローカルサーバーを起動
- ステップ1：（略）/github.com/hayabusa-cloud/hybs-server/dockerに移動します
- ステップ2：Dockerを起動します   
次回以降の起動は"--build"不要
```bash
docker-compose up --build
```
- ステップ3：（略）/github.com/hayabusa-cloud/hybs-serverに移動します
- ステップ4：複数のサービスを一気に起動
```bash
go run main.go -f config/localhost-all-compose.yml
```
サービスを一つずつ起動したい場合は：
```bash
go run main.go -f config/{service-configfile}
```

### 疎通テストしてみよう
|サーバー|URL
|----|----
|プラットフォーム|http://localhost:8086/v1/system/
|ゲーム01|http://localhost:8087/sample-game/
|ゲーム02|http://localhost:8089/sample-game/
|セントラル|http://localhost:8085/v1
レスポンス例：
```json
{"env":"local","serverId":"platform-local","serverTime":1600000000}
```
---
## プロジェクトの構造
### /application
最も重要な部分でアプリケーションを実装するところです。
#### /application/batch
バッチタスクを実装する所
#### /application/middleware
ミドルウェア（共通ロジック）を実装する所
#### /application/model
データテーブルなどの構造体を定義する所
#### /application/service
ビジネスロジックを実装する所。実装例をご参考してください：   
モデル：/application/model/samplegame/example.go   
サービス：/application/service/samplegame/example.go   
コントローラー：/config/service-sample-game/example.yml

### /config
- 起動コンフィグファイルと
- コントローラー定義

### /csv
マスタデータの置く場所

### /log
出力されたログファイルの場所   

---
## JMeterのシナリオ
/hybs-server/test/テスト計画.jmxを使ってJMeterで性能テストしましょう。

---
## ライセンス
MITライセンスのもとで配布されています   
https://github.com/hayabusa-cloud/hybs-server/LICENSE   
© 2021 hayabusa-cloud   

## 問い合わせ
Eメール：hayabusa-cloud@outlook.jp   