# kz-domain-monitor

[Орысша](README.md) | **Қазақша**

.kz доменді аттарының тіркеу мерзімін бақылауға арналған утилита.

- Домендік атаудың төлем мерзімін тексеру
- Telegram хабарландырулары
- Windows/Linux/Docker/Kubernetes қолдайды
- Go тілінде жазылған

Әдепкі бойынша жалпыға қол жетімді RDAP-сервисті [rdap.nic.kz](https://rdap.nic.kz) пайдаланады — тіркелу және токен қажет емес.
Балама драйвер ретінде [ps.kz](https://ps.kz/) API да қолдауланады.

![ps.png](.github/ps.png)

## Бұл утилита кімге арналған?
Домен төлемінде мәселелер туындағанда алдын ала хабардар болғысы келетін жүйелік әкімшілерге, сайт иелеріне, әзірлеу командаларына — домен өшіп, сайт жұмыс жасамай қалғанға дейін.

## Жылдам бастау (Docker)
```shell
docker run --rm -e DOMAIN_LIST=example.kz kravets1996/kz-domain-monitor
```

## Демо
[Мемлекеттік сервистер домендерін бақылайтын жалпыға қол жетімді Telegram-канал](https://t.me/kz_gov_domain_monitor)

## Орнату

Алдымен утилитаны баптауға арналған .env файлын жүктеп алыңыз:
[Сілтеме](https://raw.githubusercontent.com/Kravets1996/kz-domain-monitor/refs/heads/main/.env.example)
```shell
wget -O .env https://raw.githubusercontent.com/Kravets1996/kz-domain-monitor/refs/heads/main/.env.example
```

### Дайын binary файлды жүктеп алу
#### Linux
```shell
wget https://github.com/Kravets1996/kz-domain-monitor/releases/latest/download/kz-domain-monitor

chmod +x kz-domain-monitor
```
#### Windows / Windows Server
[Жүктеп алу](https://github.com/Kravets1996/kz-domain-monitor/releases/latest/download/kz-domain-monitor.exe)
[Github Release беті](https://github.com/Kravets1996/kz-domain-monitor/releases)

### Docker
```shell
docker run --rm -v $(pwd)/.env:/app/.env kravets1996/kz-domain-monitor
```

### Kubernetes
```shell
kubectl create namespace kz-domain-monitor
kubectl apply -n kz-domain-monitor -f https://raw.githubusercontent.com/Kravets1996/kz-domain-monitor/refs/heads/main/k8s/configmap.yml
kubectl apply -n kz-domain-monitor -f https://raw.githubusercontent.com/Kravets1996/kz-domain-monitor/refs/heads/main/k8s/cronjob.yml
```

### Source
```shell
git clone https://github.com/kravets1996/kz-domain-monitor
cd kz-domain-monitor
go build -o kz-domain-monitor
```

## Баптау

.env файлын (немесе Kubernetes-те орнатқанда ConfigMap) өңдеп, міндетті айнымалылардың мәндерін орнатыңыз.

### Драйверді (деректер провайдерін) таңдау
Драйвер `DOMAIN_PROVIDER` айнымалысымен таңдалады:
- `rdap` — әдепкі, жалпыға қол жетімді RDAP-сервисі rdap.nic.kz, тіркелу және токен қажет емес.
- `pskz` — ps.kz API, кіру токені қажет.

### ps.kz API-на қол жеткізуді алу және баптау
1. ps.kz кабинетінде токен жасаңыз. https://console.ps.kz/account/iam/tokens?tab=my
2. "Тек оқу" рөлін көрсетіңіз.
3. Жасалған токенді `PS_GRAPHQL_TOKEN` айнымалысына көшіріңіз

### Хабарландыруларды баптау
#### Telegram
1. [BotFather](https://telegram.me/BotFather) арқылы Telegram-бот жасаңыз.
2. Жаңа боттың токенін жасап, көшіріңіз.
3. Токенді `TELEGRAM_BOT_TOKEN` айнымалысына орнатыңыз.
4. Ботқа кез келген хабарлама жіберіңіз
5. Мына сілтемеге өтіңіз (<BOT_TOKEN> орнына 2-қадамда алған бот токенін қойыңыз): `https://api.telegram.org/bot<BOT_TOKEN>/getUpdates`
6. Алынған JSON-нан чат ID-ін табыңыз: `"chat":{"id":123456789}`
7. Алынған ID-ді `TELEGRAM_CHAT_ID` айнымалысына орнатыңыз
8. `TELEGRAM_ENABLED` айнымалысы арқылы хабарландыруларды қосыңыз

#### Slack
1. Кіріс webhook жасаңыз.
2. URL-ді `SLACK_WEBHOOK_URL` айнымалысына көшіріңіз
3. `SLACK_ENABLED` айнымалысы арқылы хабарландыруларды қосыңыз

#### Email
1. .env.example-да көрсетілген қажетті айнымалыларды толтырыңыз
2. `EMAIL_ENABLED` айнымалысы арқылы хабарландыруларды қосыңыз

#### Webhook
1. `WEBHOOK_URL` айнымалысына Webhook URL-ін көрсетіңіз
2. `WEBHOOK_ENABLED` айнымалысы арқылы хабарландыруларды қосыңыз

### Домендік атаулар
Бақылағыңыз келетін домендік атауларды `DOMAIN_LIST` айнымалысында тізіңіз.
Формат: `example.kz,example2.kz,example3.kz`

## Пайдалану
Домендерді тексеруді іске қосу:

Linux:
```shell
./kz-domain-monitor
```

Docker:
```shell
docker run --rm -v $(pwd)/.env:/app/.env kravets1996/kz-domain-monitor
```

### Жоспарлаушы
Домендерді мерзімді тексеру үшін командар жұмысын жүйе жоспарлаушысына қосу қажет.

rdap.nic.kz немесе ps.kz Rate Limit-ке тап болмау үшін тексеруді күніне 1 реттен жиі орнатпау ұсынылады.

#### Linux
/etc/crontab файлына жаңа жол қосыңыз
```shell
0 14 * * * <user> /path/to/kz-domain-monitor
```
(`<user>` орнына жүйедегі пайдаланушы атын және `/path/to` орнына kz-domain-monitor жүктелген жолды қойыңыз)

#### Windows
Windows жоспарлаушысында мерзімді тапсырма жасау:

```shell
schtasks /create ^
  /tn "KZDomainMonitor" ^
  /tr "C:\domain-monitor\kz-domain-monitor.exe" ^
  /sc daily ^
  /st 12:00 ^
  /f
```

## Әзірлеу
```shell
go mod vendor
go run main.go
```

Жинау:
```shell
go build -o kz-domain-monitor
GOOS=windows GOARCH=amd64 go build -o kz-domain-monitor.exe
```

#### Docker образын жинау
```shell
docker build -t kz-domain-monitor .
docker run --rm -v $(pwd)/.env:/app/.env kz-domain-monitor
```

Пайдалы сілтемелер:
- [RDAP KazNIC](https://nic.kz/docs/announc_20_01_2026.jsp)
- [GraphQL API нұсқаулығы](https://console.ps.kz/docs/faq/pscloud-api/ps-cloud-api/instrukciya-po-api-graphql)
- [GraphQL Playground](https://console.ps.kz/domains/graphql)

## Үлес қосу нұсқаулығы
Кодта қате тапсаңыз, жақсарту ұсынғыңыз немесе сұрақ қойғыңыз келсе —
[issue](https://github.com/Kravets1996/kz-domain-monitor/issues) және [pull request](https://github.com/Kravets1996/kz-domain-monitor/pulls) бөлімдерін пайдаланыңыз.

## Лицензия

Жоба [MIT](LICENSE) лицензиясы бойынша лицензияланған
