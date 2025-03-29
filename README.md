# Система голосования в опросах

предоставляет REST API для создания опросов, голосования и получения статистики.

### Инстукция

* Склонируйте репозиторий ```git@github.com:grozaqueen/poll.git```
* Создайте файл .env в корне проекта на основе примера и отредактируйте его:
```  
DB_ADDRESS=tarantool_db:3301
DB_USERNAME=adminka
DB_PASSWORD=12345
APP_ENV=local
```
* Подключитесь к tarantool
* Выполните в терминале ```docker-compose up --build```
* Можно проверить доступность по адресу http://localhost:8080
