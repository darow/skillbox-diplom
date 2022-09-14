<h2>Cетевой многопоточный сервис для Statuspage</h2>

Что было необходимо реализовать? [Текст](ТЗ%20по%20дипломному%20проекту.pdf) задания.

<details>
  <summary>Приложение работает на heroku</summary>

https://skillbox-diplom1.herokuapp.com/status_page.html
</details>

<details>
  <summary>Запустить его и simulator</summary>

В первом терминале: 
```bash
go run ./third_party/simulator
```

Во втором терминале:
```bash
go run ./cmd/statuspage
```

Проверяем ссылки в браузере
http://localhost:8383/mms
http://localhost:8000/status_page.html
</details>

<details>
  <summary>Создать docker контейнер и запустить</summary>

```bash
docker build -t skillbox:v1 .
docker run -p 1000:8383 -p 8000:8000 skillbox:v1
```

Проверяем ссылки в браузере
http://localhost:1000/mms
http://localhost:8000/status_page.html
</details>

