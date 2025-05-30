# Bank API Go

Микросервис для обработки банковских операций, написанный на Go.

## 📌 Основные функции
- Создание транзакций между счетами
- Просмотр истории операций
- Управление балансом пользователей
- Автоматическая документация API через Swaggo

## 🚀 Технологии
```plaintext
Go 1.24.2+
Docker 23.0+
PostgreSQL
Gin + Swaggo
```

## ⚙️ Установка
### Через Go:
```bash
git clone https://github.com/Alexandrjob/bank_api_go
cd bank_api_go
go mod tidy
go run main.go
```

### Через Docker:
```bash
docker build -t bank-api .
docker run -p 8080:8080 bank-api
```

### Генерация Swagger:
```bash
swag init
```

### Доступ к документации:
После запуска откройте:  
http://localhost:8080/swagger/index.html
