# Docker ToDoアプリ（Go + MySQL）

Goで作成したToDoアプリとMySQLを、Docker Composeでまとめて起動する開発環境です。タスクの追加、完了状態の切り替え、削除ができます。MySQLのデータはDockerボリュームに保存されるため、コンテナを作り直しても残ります。

## 起動手順

1. `.env.example` を `.env` にコピーします。
2. `.env` の `MYSQL_ROOT_PASSWORD` と `DB_PASSWORD` にパスワードを設定します。
3. 次のコマンドで起動します。

   ```bash
   docker compose up -d --build
   ```

4. ブラウザで http://localhost:8080 を開きます。

## 停止方法

```bash
docker compose down
```

この操作ではDBデータが残ります。DBデータも削除して初期状態に戻す場合のみ、次を実行します。

```bash
docker compose down -v
```
