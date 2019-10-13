# file: features/create.feature

Feature: Создание события в календаре
        Когда пользователь обратится к сервису чтобы создать событие в календере
        Если будет получен POST запрос по адресу http://localhost:5000/calendar/1/add
        Тогда пользователь создаст событие в календаре

        Scenario: Сервис отправки сообщений доступен
                When посылаю "GET" запрос к "http://localhost:5000/ping"
                Then ожидаю что код ответа будет 200
                And тело ответа будет равно "OK"

        Scenario: Создать событие в календаре
                When посылаю "POST" запрос к "http://localhost:5000/calendar/1/add" c "application/json" содержимым:
		"""
                { 
                "Message": "купить продукты",
                "EventDate": "2019-10-15T12:00"
                }
		"""
                Then ожидаю что код ответа будет 200

        Scenario: Проверить что событие существует в календаре у пользователя
                When посылаю "GET" запрос к "http://localhost:5000/calendar/1" 
                Then ожидаю что код ответа будет 200
                And в ответе будет событие с Message "купить продукты"
                And дождаться оповещения о событии с сообщением "купить продукты"