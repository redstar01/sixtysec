# 60 sec quiz game
It's a variant of a quiz game developed in Golang as Telegram bot

# How to migrate 


```
migrate -source file://database/migrations -database sqlite3://_data/test.db up
migrate -source file://database/migrations -database sqlite3://_data/test.db down

migrate -source file://database/migrations -database sqlite3://_data/test.db force 1
```