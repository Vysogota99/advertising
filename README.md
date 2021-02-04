# advertising

<h2>Инструкция по запуску</h2>
<h3>Запуск контейнеров</h3>
<p>В директории ./deployments прописать:
    <br>
    <code>docker-compose up -d</code>
    <br>
    После этого соберутся и запустятся контейнеры для работы с go и postgres
</p>
<p>
 Посмотреть их названия можно командой:
    <br>
    <code>
        docker ps
    </code>
</p>

<h3>Создание пользователя для работы с базой данных </h3>
<ul>
    <li>Войти в контейнере c postgres в учетную запись postgres (пароль - "qwerty"):
    <br>
    <code>docker exec -it deployments_store_1 /bin/bash</code>
    <br>
        <code>psql -U postgres -p 5432 -h store</code>
    <br>
    </li>
    <li>Создать нового пользователя:
    <br>
        <code>CREATE ROLE user1 WITH PASSWORD 'password' LOGIN CREATEDB;</code>
    <br>
    </li>
    <li>Создать базы данных:
    <br>
        <code>CREATE DATABASE user1;
        <br>
        CREATE DATABASE app
        </code>
    <br>
    </li>
</ul>

<h3>Запуск миграций </h3>
<ul>
    <li>
        Зайти в контейнер с golang:
        <br>
        <code>docker exec -it deployments_backend_1 /bin/bash</code>
    </li>
    <li>В директории ./build запустить миграции командой:
    <br>
    <code>migrate -database ${POSTGRESQL_URL} -path ./migrations up</code>
</ul>


<h3>Запуск сервера</h3>
<ul>
    <li>
        В директории ./build необходимо создать .env файл и заполнить его по примеру .env_example(скопировать все из .env_example в .env)
        <br>
        <code>
            cp .env_example .env
        </code>
    </li>
    <li>
        В директории ./build необходимо прописать команду для сборки сервера:
        <br>
            <code>go build ../cmd/app/main.go</code>
        <br>
    </li>
    <li>
        Теперь его можно запустить командой:
        <br>
            <code>./main</code>
        <br>
    </li>
</ul>
<h3>Структура проекта</h3>
<a href="https://github.com/golang-standards/project-layout">https://github.com/golang-standards/project-layout</a>
<br>
<ul>
    <li>build - содержит все необходимое для запуска и работы сервер: .env, Dockerfile для golang, миграции и файлы базы данных. После сборки сервера здесь появится скомпилированный файл для запуска;</li>
    <li>cmd - содержит main package;</li>
    <li>deployments - содержит файл для docker-compose.yml. Контейнеры запускаются здесь;</li>
    <li>internal - содержит пакеты для работы сервера.
        <ul>
            <li>models - модели сущностей;</li>
            <li>server - конфиг, сервер, роутер и обработчики;</li>
            <li>store - включает в себя две реализации интерфейса для работы с б.д: mock для тестирования http и postgres для работы сервера;</li>
        </ul>
    </li>
</ul>
<h3>Архитектура</h3>
<img src="./app_arch.png">
<p>
    Пользователь взаимодействием с приложением при помощи rest API. Модуль server обрабатывает входящий запрос на 3000 порту, после чего отправляет его в соответствующий обработчик, расположенный в модуле Store, который в зависимости от бизнес логики записывает данные в базу данных или забирает их из нее. 
</p>

<h3>Описание методов сервиса</h3>
<ul>
    <li>
        <h4>Метод создания объявления</h4>
        <p>
        POST /ad
        </p>
        <p>
        пример тела запроса в формате JSON:
        </p>
        <pre>
{
    "name": "iphones",
    "description": "Продам айфон 2",
    "photos": [
        "https://test/image1",
        "https://test/image2",
        "https://test/image3"
    ],
    "price": 124
}
        </pre>
        <p>
        пример отета:
        </p>
        <pre>
в случае успеха
{
    "result": {
        "id": 25,
        "status": "success"
    }
}
в случае неудачи, например если отправлено 4 ссылки на фото
{
    "result": {
        "error": "Key: 'Ad.Links' Error:Field validation for 'Links' failed on the 'max' tag",
        "message": "Некорректный запрос",
        "status": "fail"
    }
}
        </pre>
    </li>
        <li>
        <h4>Метод получения конкретного объявления</h4>
        <p>
        GET /ad/id 
        <br>
        id - id объявления
        <br>
        Чтобы получить дополнительное описание, необходимо добавить в запрос description=true
        <br>
        Чтобы получить все фотографии, необходимо добавить в запрос photos=true. 
        </p>
        <p>
        пример запроса:
        </p>
        <pre>
http://127.0.0.1:3000/ad/3?description=true&photos=true
        </pre>
        <p>
        пример отета:
        </p>
        <pre>
в случае успеха
{
    "result": {
        "name": "iphones",
        "description": "Продам айфон 2",
        "photos": [
            "https://test/image1",
            "https://test/image2",
            "https://test/image3"
        ],
        "price": 100
    }
}
в случае неудачи
{
    "result": {
        "error": "strconv.ParseBool: parsing \"true1\": invalid syntax",
        "message": "Некорректный запрос",
        "status": "fail"
    }
}
        </pre>
    </li>
            <h4>Метод получения всех объявлений</h4>
        <p>
        GET /ads
        <br>
        Чтобы указать номер нужной страницы, необходимо добавить в запрос p=n, где n - положительное число
        <br>
        Чтобы указать по какому полю сортировать, необходимо добавить в запрос sort_by=field, где field может принимать значение price или created_at
        <br>
        Чтобы указать направление сортировки, необходимо добавить в запрос sort_direction=type, где type может принимать значение asc для возрастания или desc для убывания 
        </p>
        <p>
        пример запроса:
        </p>
        <pre>
http://127.0.0.1:3000/ads?p=1&sort_by=price&sort_direction=desc
        </pre>
        <p>
        В ответе содержится массив с объявлениями и общее количество страниц 
        <br>
        пример отета:
        </p>
        <pre>
{
    "result": {
        "data": [
            {
                "name": "iphones",
                "photos": [
                    "https://test/image1"
                ],
                "price": 124
            },
            {
                "name": "iphones",
                "photos": [
                    "https://test/image1"
                ],
                "price": 123
            },
...
        "n_pages": 3
    }
}
        </pre>
    </li>
</ul>
