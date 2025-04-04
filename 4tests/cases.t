docker exec -it frappuccino_db_1 psql -U latte -d frappuccino
docker-compose down && docker-compose up --build
python -m SimpleHTTPServer 8000

http://localhost:9090/inventory/{id} PUT, GET, DELETE

1

{"name":"dd","stock":1000,"unit":"g","reorder_threshold":200,"price":0}


-----------------------------------------------------------------

http://localhost:9090/reports/orderedItemsByPeriod?period=day&month=april

http://localhost:9090/reports/orderedItemsByPeriod?period=month&year=2025

http://localhost:9090/reports/orderedItemsByPeriod?period=month&year=2025




http://localhost:9090/menu POST
{
    "id": 2,
    "name": "Latte",
    "description": "Espresso with steamed milk",
    "categories": ["coffee", "hot"],
    "allergens": ["milk"],
    "price": 4.5,
    "available": true,
    "size": "medium",
    "ingredients": [
        {"ingredient_id": 1, "quantity": 30.0, "unit": "g"},
        {"ingredient_id": 2, "quantity": 200.0, "unit": "ml"}
    ]
}



1. Запрос на получение остатков инвентаря с сортировкой по цене и пагинацией

    Метод: GET

    URL: http://localhost:9090/inventory/getLeftOvers?sortBy=price&page=1&pageSize=10

    Описание: Этот запрос получает остатки инвентаря, отсортированные по цене, на первой странице с 10 товарами на странице.

2. Запрос на получение остатков инвентаря с сортировкой по количеству и пагинацией

    Метод: GET

    URL: http://localhost:9090/inventory/getLeftOvers?sortBy=quantity&page=2&pageSize=5

    Описание: Этот запрос получает остатки инвентаря, отсортированные по количеству товара, на второй странице с 5 товарами на странице.

3. Запрос на получение остатков инвентаря с пагинацией по умолчанию (сортировка по имени)

    Метод: GET

    URL: http://localhost:9090/inventory/getLeftOvers?page=1&pageSize=10

    Описание: Этот запрос получает остатки инвентаря, отсортированные по имени (значение по умолчанию для sortBy), на первой странице с 10 товарами на странице.

4. Запрос на получение остатков инвентаря с сортировкой по цене без пагинации (по умолчанию 10 товаров)

    Метод: GET

    URL: http://localhost:9090/inventory/getLeftOvers?sortBy=price

    Описание: Этот запрос получает остатки инвентаря, отсортированные по цене, с дефолтными значениями для страницы (1) и количества товаров на странице (10).

5. Запрос на получение остатков инвентаря с некорректными параметрами

    Метод: GET

    URL: http://localhost:9090/inventory/getLeftOvers?sortBy=price&page=abc&pageSize=10

    Описание: Этот запрос должен вернуть ошибку, так как значение параметра page не является числом.

    Ожидаемый ответ:

        Статус: 400 Bad Request

        Тело: "Invalid page parameter"

6. Запрос на получение остатков инвентаря с некорректным значением параметра sortBy

    Метод: GET

    URL: http://localhost:9090/inventory/getLeftOvers?sortBy=invalid&page=1&pageSize=10

    Описание: Этот запрос должен вернуть ошибку, так как параметр sortBy имеет недопустимое значение.

    Ожидаемый ответ:

        Статус: 400 Bad Request

        Тело: "Invalid sortBy parameter"

7. Запрос на получение остатков инвентаря с параметрами по умолчанию (пагинация 1 страница и 10 элементов на странице)

    Метод: GET

    URL: http://localhost:9090/inventory/getLeftOvers

    Описание: Этот запрос получит остатки инвентаря с параметрами по умолчанию (страница 1, 10 товаров на странице и сортировка по имени).