# kz-domain-monitor

Утилита для мониторинга истечения срока регистрации доменных имён в зоне .kz.

- Проверка срока оплаты доменного имени
- Уведомления в Telegram
- Поддерживает Windows/Linux/Docker/Kubernetes
- Написана на Go

Основана на API [ps.kz](https://ps.kz/). 
Для работы необходима регистрация и получение токена в личном кабинете (бесплатно).

## Быстрый старт (Docker)
```shell
docker run --rm -e DOMAIN_LIST=example.kz -e PS_GRAPHQL_TOKEN='****' kz-domain-monitor
```

## Установка

Для начала загрузите файл .env для настройки утилиты:
[Ссылка]()
```shell
wget https://github.com/kz-domain-monitor/.env.example .env
```

### Загрузка pre-built binary
#### Linux
```shell
wget https://github.com/kz-domain-monitor

chmod +x kz-domain-monitor
```
#### Windows
[Страница Release в Github]()

### Docker
```shell
docker run --rm -v $(pwd)/.env:/app/.env kz-domain-monitor
```

### Kubernetes
```shell
kubectl create namespace kz-domain-monitor
kubectl apply -n kz-domain-monitor -f k8s/configmap.yml
kubectl apply -n kz-domain-monitor -f k8s/cronjob.yml
```

### Source
```shell
git clone https://github.com/kz-domain-monitor
cd kz-domain-monitor
go build -o kz-domain-monitor
```

## Настройка

Отредактируйте файл .env (или ConfigMap при установке в Kubernetes) 
и установите значения для обязательных переменных.

### Получение и настройка доступа к API ps.kz
1. Создайте токен в кабинете ps.kz. https://console.ps.kz/account/iam/tokens?tab=my
2. Укажите роль "Только чтение".
3. Скопируйте сгенерированный токен в переменную `PS_GRAPHQL_TOKEN`

### Настройка уведомлений в Telegram
1. Создайте Telegram-бота с помощью [BotFather](https://telegram.me/BotFather).
2. Создайте и скопируйте токен нового бота.
3. Установите токен в переменную `TELEGRAM_BOT_TOKEN`.
4. Напишите боту любое сообщение
5. Перейдите по ссылке (замените <BOT_TOKEN> на токен бота, полученный на шаге 2) https://api.telegram.org/bot<BOT_TOKEN>/getUpdates
6. Найдите ID чата в полученном JSON `"chat":{"id":123456789}`
7. Установите полученный ID в переменную `TELEGRAM_CHAT_ID`

### Доменные имена
Перечислите доменные имена, которые вы хотите отслеживать в переменной `DOMAIN_LIST`.
Формат: `example.kz,example2.kz,example3.kz`

## Использование
Запуск проверки доменов:

Linux:
```shell
./kz-domain-monitor
```

Docker:
```shell
docker run --rm -v $(pwd)/.env:/app/.env kz-domain-monitor
```

### Планировщик
Для периодической проверки доменов необходимо добавить запуск команды в планировщик системы.

Рекомендуется устанавливать проверку не чаще 1 раза в сутки, чтобы не столкнуться с Rate Limit ps.kz.

#### Linux
Добавьте новую строку в файл /etc/crontab
```shell
0 14 * * * <user> /path/to/kz-domain-monitor
```
(замените `<user>` на название пользователя в вашей системе и `/path/to` на путь куда загружен kz-domain-monitor)

#### Windows
Создание периодической задачи в планировщике Windows:

```shell
schtasks /create ^
  /tn "KZDomainMonitor" ^
  /tr "C:\domain-monitor\domain-monitor.exe" ^
  /sc daily ^
  /st 12:00 ^
  /f
```

## Разработка
```shell
go mod vendor
go run main.go
```

#### Сборка Docker образа
```shell
docker build -t kz-domain-monitor .
docker run --rm -v $(pwd)/.env:/app/.env kz-domain-monitor
```

Полезные ссылки:
- [Инструкция по API GraphQL](https://console.ps.kz/docs/faq/pscloud-api/ps-cloud-api/instrukciya-po-api-graphql)
- [GraphQL Playground](https://console.ps.kz/domains/graphql)

## Лицензия

Проект лицензирован под лицензией [MIT](LICENSE)
