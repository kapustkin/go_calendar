# file: features/edit.feature

Feature: Редактирование события в календаре
        Когда у пользователя есть событие и он хочет его удалить
        Тогда пользовател может удалить событие используя UUID
        Если он сделает запрос на удалить событие
        Тогда событие с переданным UUID будет удалено

        Scenario: Сервис отправки сообщений доступен
                When посылаю "GET" запрос к "http://localhost:5000/ping"
                Then ожидаю что код ответа будет 200
                And тело ответа будет равно "OK"

        Scenario: Создать событие в календаре
                When посылаю "POST" запрос к "http://localhost:5000/calendar/1/add" c "application/json" содержимым:
		        """
                { 
                "Message": "встреча с основателем мак дональдс",
                "EventDate": "2019-10-20T12:00"
                }
		        """
                Then ожидаю что код ответа будет 200                

        Scenario: Найти событие в календаре у пользователя
                When посылаю "GET" запрос к "http://localhost:5000/calendar/1" 
                Then ожидаю что код ответа будет 200
                And в ответе будет событие с Message:
                """
                { 
                    "Message": "встреча с основателем мак дональдс",
                    "EventDate": "2019-10-12T12:00"
                }
		        """

        Scenario: Удалить событие из календаря
                When посылаю "POST" запрос к "http://localhost:5000/calendar/1/remove" c "application/json" и новым содержимым:
		        """
                { 
                    "UUID":"{REPLACE_UUID}"
                }
		        """
                Then ожидаю что код ответа будет 200
                
        @negative
        Scenario: Проверить отсутствие события в календаре у пользователя
                When посылаю "GET" запрос к "http://localhost:5000/calendar/1" 
                Then ожидаю что код ответа будет 200
                And в ответе не будет события Message:
                """
                { 
                    "UUID":"{REPLACE_UUID}"
                }
	        """