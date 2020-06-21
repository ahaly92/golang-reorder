# golang-reorder

# Goal 

Updating a backend/database display order with drag and drop in SQL

# Models - pkg/models/
User 

Todo

TodoList

# Running Server
To run the server run:
```
go run cmd/main.go   
```

# DB Migrations - Goose 
1. Navigate to pkg/migrations/
```
cd pkg/repository/migrations/
```

2. To run the migrations and create tables run:
```
goose postgres "host=localhost port=5432 user=postgres dbname=reorder sslmode=disable password=postgres" up
```
