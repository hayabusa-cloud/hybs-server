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
|サーバー|URL|APIドキュメント
|----|----|----
|プラットフォーム|http://localhost:8086/v1/system/|http://localhost:8086/docs/
|ゲーム01|http://localhost:8087/sample-game/|http://localhost:8087/sample-game/docs/
|ゲーム02|http://localhost:8089/sample-game/|http://localhost:8089/sample-game/docs/
|セントラル|http://localhost:8085/v1/|http://localhost:8085/docs/
   
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
#### /application/controller
コントローラーを定義する所
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
- 起動コンフィグファイル

### /csv
マスタデータの置く場所

### /log
出力されたログファイルの場所   

---
## シンプルなテストWeb APIを実装してみよう
具体例を挙げてWeb APIの実装する方法を解説します。

- ステップ1：DB設計とデータモデル定義  
データベースを設計して、そして"/application/**model**/{server-name}/"フォルダに"{table-name}.go"ファイルを新規作成します。  
新規作成したファイルにデータテーブルの構造体を定義します（インデックスなど定義も含む）。
例えば、
```go
type PlayerExampleData struct {
	HayabusaID  string `sql:"hayabusa_id" json:"hayabusaId"`
	IntField    int    `sql:"int_field" json:"intField"`
	StringField string `sql:"string_field" json:"stringField"`
}
```   
*ApplicationUpメソッドに索引の管理、オートマイグレーションなどを入れるのはおすすめします。*

- ステップ2：コントローラーを定義  
コントローラーはリクエストパスとリクエストメソッドによって、サーバーはどのような挙動をするのか、を定義するコンフィグファイルです。例を挙げて説明します。  
```yaml
  # まず、リクエストパスを指定します
- location: /sample-game/v1/
  # リクエストパスにルートパラメータが入れられます、例えば：
  # - location: /sample-game/v1/:user-id
  # "description"はAPIドキュメントの生成に使います
  description: サンプルゲーム仮想サーバーV1：ルートパス
  # 接続できるIPアドレスを制限できる
  # "allow"がオミットされた場合はIPアドレスを制限しません
  allow:
    - 0.0.0.0/0
  # IPアドレスのブラックリストも設定できます。オミットの場合は無制限です
  deny:
    - 100.100.100.100/32 # 記入例
    - 99.99.99.0/24      # 記入例
  # リクエストのパラメータに入力ルールが設定できます
  # リクエストのパラメータはパスパラメータ、クエリパラメータ、ポストパラメータ、定数パラメータ4種類に分けられています。
  # パスパラメータの入力ルール
  # パスパラメータはURLに埋め込まれる特定のリソースを識別するためのパラメータです
  # 例えば：/sample-game/v1/foo/:id　というパスでしたら、idがパスパラメータになります
  # [METHOD] /sample-game/v1/foo/100 をリクエストするときに、"foo"が100となります。
  path_params: # （記入例）
    - name: id # パラメータ名
      description: サンプルフィールドです
      # 入力ルールを決めます
      allow: // ホワイトリスト。オミットの場合は無制限
        - ^[a-zA-Z0-9]{4,10}$ # 正規表現でルールを書く
      deny: // ブラックリスト。オミットの場合は無制限
        - ^root$ # 記入例
        - ^sys$  # 記入例
      example: hayabusa00  # 例
  # クエリパラメータはURLの最後に「？」が付いたパラメータです。検索やフィルタなどに関する条件がクエリパラメータとして扱われます。
  # 例えば、リクエストが/sample-game/v1/foo/?name=hayabusa
  query_args:    # （記入例）
    - name: name # パラメータ名
      description: サンプルフィールドです
      example: hayabusa
  # フォームパラメーターはリクエストボディです。通常に更新や追加する時に入れます。
  # 例えば：
  # [POST] /sample-game/v1/foo/100
  # m=1000&n=2000&gender=male
  form_args:     # （記入例）
    - name: name # パラメータ名
      description: サンプルフィールドです
      example: hayabusa
  # 定数パラメータはコントローラーで指定する定数です
  # 例えば、"server=hayabusa"と指定すれば、ソースコード上"server"というキーで"hayabusa"という値が取得できます
  const_params:
    - name: server
      value: hayabusa
  # "middlewares"はそのリクエストパスをプレフィックスとする全てのAPIが実行する共通ロジックを定義します
  middlewares:
    - HybsLog             # ログ出力（内装ミドルウェア）
    - Authentication      # ベーシックユーザー認証（内装ミドルウェア）
    - ResponseJSON        # JSON形式のレスポンスを生成（内装ミドルウェア）
    - UseCache            # 名前が"Cache"のプラグインを使用
    - UseRedis            # 名前が"Redis"のプラグインを使用
    - UseMongoSampleGame  # 名前が"MongoSampleGame"のプラグインを使用
    - UseMySQLSampleGame  # 名前が"MySQLSampleGame"のプラグインを使用
    - UseSqliteSampleGame # 名前が"SqliteSampleGame"のプラグインを使用
    - AuthOnetimeToken    # ワンタイムトークンによるアクセス権限の承認
    - CheckPlayerStatus   # アカウントが凍結中かどうかをチェック
  slow_query_warn: 80ms   # API実行時間が80msを超えると、Warnレベルのスロークエリログが記録されます
  slow_query_error: 200ms # API実行時間が200msを超えると、Errorレベルのスロークエリログが記録されます
```  
では、/sample-game/v1の下に、APIを定義しましょう
```yaml
- location: /sample-game/v1/foo/
  # "location"は"/sample-game/v1/"プレフィックスを含む為、"/sample-game/v1/"のミドルウェアやルールなどは適用されます
  description: テストロケーション
  # 実行するAPIを定義する
  services: 
    - method: GET  # メソッド
      description: テストAPI
      # ここでもパラメーターにルールが設定できます
      query_args:
        - name: x
          description: テストパラメーター
          allow:
            - ^1|2|3|4|5$
          # 最も重要な項目です。GETメソッドで/sample-game/v1/foo/にアクセスしたら、どのAPI処理関数を呼び出すのを指定します。
      service_id: SampleGameTestAPI # 登録した名前が"SampleGameTestAPI"の処理関数を呼び出す
      response: # ドキュメント生成の為の項目です
        - status_code: 200 # 成功した場合
          description: 成功
          fields: # 返すフィールド
            - name: x
              description: 説明文
            - name: y
              description: 説明文
        - status_code: 400 
          description: 失敗
```
"Use〇〇〇"のようなミドルウェア名は"Use"+プラグイン名の組み合わせです。例えば、起動コンフィグファイルで名前が"Redis001"のプラグインを定義して、「UseRedis001」というミドルウェア名をコントローラーの"middlewares"の下で書けば、そのプラグインがそのパスの下のAPIに使えるになります。  

- ステップ3：ビジネスロジックを実装  
では、APIを処理する関数を実装しましょう。下記のようにGo言語のソースコードを新規作成します。  
/application/service/{server-name}/{module-name}.go  
実装例は/application/service/sample-game/example.goを参考してください。  
関数の実装が終わったら、"init"関数に"hybs.RegisterService({server-id}, メソッド名)"を書いてAPI処理メソッドを登録し忘れないで下さい。

- ステップ4：単体テスト
Go言語の内装テストツールやPostman、JMeterなど外部ツールでテストしましょう。

- 色々実装例：
コントローラー定義：/config/service-sample-game/example.yml  
実装：/application/service/samplegame/example.go  

---
## JMeterのシナリオ
/hybs-server/test/テスト計画.jmxを使ってJMeterでHTTPサーバーの性能テストをやりましょう。

---
## リアルタイム通信機能をテストしましょう   
### サーバーを起動   
- 起動中のComposeサーバー（統合サーバー）を停止します
```bash
kill {pid}
```
- リアルタイム通信専用サーバーを起動します
```bash
go run main.go -f config/localhost-all-compose.yml
```

### エコーテスト  
まず、エコーテストをしましょう。クライアント1000個から秒間20回頻度でサーバーにサイズが64バイトのメッセージを送ります。  
サーバーは受けたメッセージを変えずにそのままクライアントに返送します。  
./test/realtime-echo/に移動します。  
```bash
cd ./test/realtime-echo/
```
そして、テストプログラムを実行します。  
```bash
go run main.go
```
テストプログラムはコンソールにQPS数を出力します。   
楽に20k qpsに達成しましたか。
4コア以上のPCでテストするのをおすすめです。  
Windows PCをお持ちの方、WSL2環境で実行するのをおすすめです。

### ルームブロードキャストテスト  
./test/realtime-roombroadcast/に移動します。
そして、テストプログラムを実行します。
```bash
go run ./
```
なお、負荷パラメータは指定できます。
```go
const (
  // ルーム数
  roomNum = 10
  // 1ルーム内ユーザー数
  roomSize = 10
  // 通信頻度
  ioInterval = time.Second / 10
)
```
秒間20k回の受信に達成しましたか。

### パフォーマンスにつきまして
本プロジェクトはDEMO版として、秒間数万（4コアCPU）回のデータ受送信が耐えられています。   
更にパフォーマンスを求めるなら、I/O部分ないしアプリケーション層の設計に最適化する必要があります。  
且つ、Linux OS側の設定と運用を正しく行わなければなりません。  
高負荷に耐えられるサーバエンジンの設計・開発・運用についてご要望がございましたら、お問い合わせメールまでご連絡をお願いいたします。

---
## ライセンス
MITライセンスのもとで配布されています   
https://github.com/hayabusa-cloud/hybs-server/LICENSE   
© 2021 hayabusa-cloud   

## お問い合わせ
Eメール：hayabusa-cloud@outlook.jp   
