# testtaskmyoffice
test task for myoffice

Строки в консоль выводятся не в том же порядке что строки в файле с URL т.к. этого условия небыло в задании

Запуск с дефолтными настройками:
  make run path=./source/testurl.txt   

Тесты:
  make test

Запуск с полными настройками:
  go run cmd\cli\main.go -path=testurl.txt -cpucount=4 -countWorker=10 -maxprocessurldurationmsec=1000 -maxtotaldurationsecond=600 

  -path путь к файлу с URL
  -cpucount количество используемых процессоров 
  -countWorker количество запускаемых воркеров
  -maxprocessurldurationmsec время таймаута при вызове URL,мсек
  -maxtotaldurationsecond максимальное время работы программы,секунд 
  
source:
Тестовое задание Go
Необходимо реализовать CLI-утилиту, которая реализует асинхронную обработку входящих URL из файла, переданного в качестве аргумента данной утилите.
Формат входного файла: на каждой строке – один URL. URL может быть очень много! Но могут быть и невалидные URL.

Пример входного файла:
https://myoffice.ru
https://yandex.ru

По каждому URL получить контент и вывести в консоль его размер и время обработки. Предусмотреть обработку ошибок.

Код должен быть размещён на Gitlab/Github (в форме заявки необходимо дать на него ссылку).