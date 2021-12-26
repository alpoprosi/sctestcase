# GET Request Counter Server

## Задача

Сервер считает количество обработанных запросов за последние 60 секунд.

Инкрементирование счетчика происходит при отправке GET запроса на любой метод. 

В случае остановки приложения структура с временем полученных актуальных 
    запросов сохраняется на диск и считывается при следующем запуске.

## Запуск и использование

Для запуска достаточно запустить `main.go`:

```sh
  go run main.go
```

Сервер запускается на `8888` порту и обрабатывает все запросы:

```sh
    curl 'http://127.0.0.1:8888/anyUrls'
```

Чтобы получить количество запросов на данный момент, не учитывая текущий,
необходимо выполнить следующий запрос:

```sh
    curl 'http://127.0.0.1:8888/count'
```

В ответ придет `json` вида:

```json
    {"count":62}
```

## Хранение данных

Данные сохраняются в директории `./data`.

Учитывается также `/favicon.ico`.

Файлы сохраненные с последнего завершения имеют расширение `.cache`.

Файлы, которые были прочитаны при запуске имеют расширение `.used`.