
Данный файл отвечает за взаимодействие с базой данных с использованием GORM. Ниже описаны все основные методы и их назначение.

## 1. Инициализация

* **`Initdb()`**
  Подключается к локальной базе данных SQLite (`database.db`) и автоматически создает/обновляет таблицы на основе структур данных (миграция).

## 2. Создание (Create)

Все методы создания пользователей используют транзакции (`db.Transaction`) для предотвращения создания неполных записей (например, если аккаунт создался, а профиль — нет).

* **`CreateDoctor(rcd *models.RequestCreateDoctor)`**
  Хэширует пароль, создает учетную запись `Auth` с ролью "Doctor", а затем профиль `Doctor`, привязанный к указанной клинике.

* **`CreatePatient(req *models.RequestCreatePatient)`**
  Хэширует пароль, создает учетную запись `Auth` с ролью "Patient" и профиль `Patient` с указанием диагноза.

* **`CreateClinic(req *models.RequestCreateClinic)`**
  Создает новую клинику, а затем сразу же регистрирует для нее главного администратора (врача с уровнем прав `Root: 1`).

* **`CreateMedicalRecord(req *models.RequestCreateRecord, doctorID uint)`**
  Добавляет медицинскую запись (диагноз и рецепт) для пациента от лица указанного врача, автоматически проставляя текущую дату.

* **`CreateMessage(req models.RequestSendMessage, IDJWT uint)`**
  Находит ID получателя по его Email и сохраняет сообщение в чат от лица отправителя.

* **`CreateAppointment(req models.RequestAppointment)`**
  Создает запись пациента к врачу в таблице очередей (`Quere`).

## 3. Чтение (Read)

* **`ReadDoctor(ID uint)` / `ReadPatient(ID uint)` / `ReadClinic(ID uint)`**
  Возвращают полную структуру сущности по её первичному ключу.

* **`GetDoctorsByClinic(clinicID uint)` / `GetPatientsByClinic(clinicID uint)`**
  Возвращают массивы врачей и пациентов, прикрепленных к конкретной клинике.

* **`GetAllClinics()`**
  Возвращает список всех существующих клиник.

* **`GetAllRecordsByPatient(patientID uint)`**
  Возвращает историю болезни (все медицинские записи) конкретного пациента.

* **`GetMessages(SenderID uint, RecipientID uint)`**
  Возвращает историю переписки между двумя пользователями (в обе стороны).

* **`GetQueryForPatient(ID uint)` / `GetQueryForDoctor(ID uint)`**
  Возвращают списки запланированных приемов для конкретного пациента или врача из таблицы очередей.

## 4. Обновление (Update)

Эти методы также используют транзакции для безопасного обновления связанных таблиц (например, одновременное обновление пароля в `Auth` и личных данных в профиле).

* **`UpdateDoctor(req *models.RequestUpdateDoctor, ID uint)` / `UpdatePatient(...)`**
  Обновляют личные данные профиля, а также перехэшируют и обновляют пароль в таблице аккаунтов.

* **`UpdateClinic(req *models.RequestUpdateClinic, ID uint)`**
  Изменяет название и адрес клиники.

* **`UpdateMedicalRecord(req *models.RequestUpdateRecord)`**
  Обновляет рецепт, дату и ID лечащего врача в существующей медицинской карточке.

* **`UpdateMessage(req *models.RequestUpdateMessage, SenderID uint)`**
  Позволяет отправителю изменить текст своего сообщения.

## 5. Удаление (Delete)

Методы удаляют записи по их первичному ключу `ID`. Благодаря каскадному удалению (`OnDelete:CASCADE`) в `models.go`, удаление клиники автоматически удалит всех ее врачей и пациентов.

* **Базовое удаление:**
  `DeleteDoctor`, `DeletePatient`, `DeleteClinic`, `DeleteMessage`, `DeleteAppointment`.

* **`DeleteMedicalRecord(recordID uint, patientID uint)`**
  Удаляет медицинскую запись. Содержит дополнительную проверку принадлежности к конкретному пациенту для защиты от случайного удаления чужих записей.

## 6. Специальные функции безопасности и доступа (Special func)

* **`FindAuthByEmail(Email string)`**
  Ищет внутренний ID аккаунта по почте. Используется при логине и отправке сообщений.

* **`FindDoctorByAuthID(AuthID uint)` / `FindPatientByAuthID(AuthID uint)`**
  Связывают ID аккаунта (из JWT токена) с ID конкретного профиля. Используются мидлварями и хэндлерами для определения прав пользователя.

* **`MakeClinicIsAvailable(ClinicID uint, AuthIDJWT uint)`**
  Активирует клинику. Транзакция проверяет, что инициатор является суперадмином (`Root == 2`).

* **`MakeDoctorIsAvailable(DoctorID uint, AuthIDJWT uint)`**
  Активирует аккаунт врача. Требует прав локального админа клиники (`Root >= 1`).

* **`PromoteDoctor(DoctorID uint, AuthIDJWT uint)` / `DemoteDoctor(...)`**
  Повышают или понижают уровень прав (`Root`) врача. локальный админ (`Root == 1`) может управлять только своей клиникой и не может выдать права выше своих собственных. Суперадмин (`Root == 2`) имеет глобальный доступ.