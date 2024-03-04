# Отчет о выполнении ДЗ по MongoDB

1) Сначала скачиваем докер образ и запускаем контейнер *MongoDB*

    `docker run -d -p 27017:27017 --name mongodb mongo`

![Снимок экрана 2024-03-04 в 12.00.51.png](..%2F..%2F..%2Fvar%2Ffolders%2Fcd%2Fdh9jtjl112ggmxcf8wnkdsgc0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_fraAnN%2F%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD%D0%B0%202024-03-04%20%D0%B2%2012.00.51.png)

2) Ставим mongo клиент, чтобы из командной строки выполнять команды

    ```
   brew tap mongodb/brew
   brew update
   brew install mongodb-community@7.0`
   ```
   
3) Проверяем, что клиент скачан и может заходить в контейнер

`mongosh --host localhost --port 27017`

![Снимок экрана 2024-03-04 в 12.32.50.png](..%2F..%2F..%2Fvar%2Ffolders%2Fcd%2Fdh9jtjl112ggmxcf8wnkdsgc0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_Mod7Zq%2F%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD%D0%B0%202024-03-04%20%D0%B2%2012.32.50.png)

4) Скачиваем датасет из [хабра](https://habr.com/ru/companies/edison/articles/480408/)
Я выбрал датасет расшифровки радиограмм NPR:
   
   https://www.kaggle.com/datasets/shuyangli94/interview-npr-media-dialog-transcripts?resource=download

5) Далее распаковываем zip архив, достаем какой-нибудь `csv` формата файл и загружаем этот датасет в MongoDB:

 ```
 mongoimport --host localhost --port 27017 --db mydatabase \

 --collection mycollection --type csv --file utterances-2sp.csv --headerline
 ```

![Снимок экрана 2024-03-04 в 12.33.56.png](..%2F..%2F..%2Fvar%2Ffolders%2Fcd%2Fdh9jtjl112ggmxcf8wnkdsgc0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_aFdDRd%2F%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD%D0%B0%202024-03-04%20%D0%B2%2012.33.56.png)

6) Используем опять клиент, чтоб выполнить CRUD операци в рамках импортированного датасета. Для начала прочитаем данные:

    ```
    mongosh --host localhost --port 27017
    test> use mydatabase
    mydatabase> db.mycollection.find()
    ```
   
![Снимок экрана 2024-03-04 в 12.40.00.png](..%2F..%2F..%2Fvar%2Ffolders%2Fcd%2Fdh9jtjl112ggmxcf8wnkdsgc0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_0UFD7A%2F%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD%D0%B0%202024-03-04%20%D0%B2%2012.40.00.png)

7) Операция update:

    ```
   mydatabase> db.mycollection.updateOne(
   ... { "_id": ObjectId("65e594be66bebabc5113b108") }, // Условие для поиска записи
   ... { $set: { "utterance": "Hellow world from Mongo!" } } // Новые значения полей
   ... )
   ```
   
![Снимок экрана 2024-03-04 в 12.52.30.png](..%2F..%2F..%2Fvar%2Ffolders%2Fcd%2Fdh9jtjl112ggmxcf8wnkdsgc0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_RPeZw6%2F%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD%D0%B0%202024-03-04%20%D0%B2%2012.52.30.png)

8) Операция delete:

`mydatabase> db.mycollection.deleteOne({ "_id": ObjectId("65e594be66bebabc5113b108") })`

![Снимок экрана 2024-03-04 в 12.55.04.png](..%2F..%2F..%2Fvar%2Ffolders%2Fcd%2Fdh9jtjl112ggmxcf8wnkdsgc0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_NlHQTi%2F%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD%D0%B0%202024-03-04%20%D0%B2%2012.55.04.png)

9) Создание индекса и сравнение производительности запросов:

```
// Сначала выполним поиск без индекса замерив время выполнения
mydatabase> db.mycollection.find({ "episode": 1, "speaker_order": 0 }).explain("executionStats")
// Получаем 44188 ms = 44 s - довольно долго!
// Далее создаём составной индекс
mydatabase> db.mycollection.createIndex({ "episode": 1, "speaker_order": 1 })
// И повторяем операцию поиска
mydatabase> db.mycollection.find({ "episode": 1, "speaker_order": 0 }).explain("executionStats")
// Получаем executionTimeMillis: 309, то есть 309 ms
```

Получили прирост во времени выполнения в 140 раз!
И научились работать с mongo)

![Снимок экрана 2024-03-04 в 13.12.08.png](..%2F..%2F..%2Fvar%2Ffolders%2Fcd%2Fdh9jtjl112ggmxcf8wnkdsgc0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_KsCjnA%2F%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD%D0%B0%202024-03-04%20%D0%B2%2013.12.08.png)

![Снимок экрана 2024-03-04 в 13.13.35.png](..%2F..%2F..%2Fvar%2Ffolders%2Fcd%2Fdh9jtjl112ggmxcf8wnkdsgc0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_BU0ylm%2F%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD%D0%B0%202024-03-04%20%D0%B2%2013.13.35.png)

------------
# Вывод
Самое главное что могу сказать: у mongoDB очень приятная и понятная документация из-за этого работать с ней довольно просто.
Ну и к тому же это одна из самых популярных NoSQL баз во многом из засвоей простоты.
