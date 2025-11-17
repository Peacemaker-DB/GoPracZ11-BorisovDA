# Практическое задание № 11 Борисов Денис Александрович ЭФМО-01-25

Цели:

1.	Освоить принципы проектирования REST API.
2.	Научиться разрабатывать структуру проекта backend-приложения на Go.
3.	Спроектировать и реализовать CRUD-интерфейс (Create, Read, Update, Delete) для сущности «Заметка».
4.	Освоить применение слоистой архитектуры (handler → service → repository).
5.	Подготовить основу для интеграции с базой данных и JWT-аутентификацией в следующих занятиях.

# Выполнение практического задания.

1. Структура проекта

  Для выполнения практической работы была сделана следующая структура проекта

<img width="174" height="337" alt="image" src="https://github.com/user-attachments/assets/fd1fd290-60ee-4d61-a66e-eb7bb61010d4" />


   Так же были установлены все необходимые расширения для выполнения практики

<img width="979" height="147" alt="Снимок экрана 2025-11-17 033223" src="https://github.com/user-attachments/assets/9ea8d6af-1a6e-497f-ae17-d90972f5fa5e" />



2.	Модель данных.
  Для реализации практики была написана модель в файл note.go

<img width="244" height="212" alt="image" src="https://github.com/user-attachments/assets/f4f1107a-35f1-4915-a4e6-c9d7ae02f64c" />


3. In-memory репозиторий.
 Для выполнения практики была написана файл note_mem.go. В котором реализован CRUD для обработки заметок.

<img width="450" height="1008" alt="image" src="https://github.com/user-attachments/assets/df4b0477-84dc-43c6-a7e3-2e9c8df198f1" />


4. HTTP-обработчик.

   Затем был создан файл notes.go, в котором происходит обработка запросов

<img width="435" height="994" alt="image" src="https://github.com/user-attachments/assets/2aa91438-cc07-45ef-bb61-e42ebd9e27a5" />


  
 
5. Маршрутизация

  После был напиисан router.go, где осуществляется мушрутизация

<img width="368" height="279" alt="image" src="https://github.com/user-attachments/assets/74ed1a48-1aa1-4a30-87f4-2acb5c3b84a8" />

6. Точка входа

   Для запуска сервера был написан main.go

<img width="364" height="289" alt="image" src="https://github.com/user-attachments/assets/01f5cbba-3639-412d-9cb7-cccfc9600401" />


# Проверка работоспособности

  Для проверки работоспособности был запущен сервер, после в Postman были проверено:

  Создание заметки

<img width="691" height="412" alt="image" src="https://github.com/user-attachments/assets/43696791-2d50-4a7d-be25-89a40aa3b95e" />


  Вывод всех заметок 

<img width="697" height="634" alt="image" src="https://github.com/user-attachments/assets/874f2056-fc1f-4e49-8349-ca7042eadbeb" />

  Вывод определенной заметки

<img width="525" height="414" alt="image" src="https://github.com/user-attachments/assets/fef10779-b872-4cfb-97be-3083eece7ad5" />


  Обновить заметку

<img width="721" height="422" alt="image" src="https://github.com/user-attachments/assets/48257205-b95b-4f17-8613-02d49551addf" />

  Удалить заметку

<img width="711" height="348" alt="image" src="https://github.com/user-attachments/assets/31957b1f-cd5e-488d-964a-bc2bc57c2bc4" />


 После удаления

<img width="700" height="547" alt="image" src="https://github.com/user-attachments/assets/c7163699-671b-44ae-8ecc-148ba3884997" />
