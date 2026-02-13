import{_ as a,o as i,c as n,ae as t}from"./chunks/framework.D5T0pHrq.js";const E=JSON.parse('{"title":"Архитектура Вокзал.ТЕХ","description":"","frontmatter":{},"headers":[],"relativePath":"architecture/README.md","filePath":"architecture/README.md"}'),e={name:"architecture/README.md"};function l(p,s,r,h,c,d){return i(),n("div",null,[...s[0]||(s[0]=[t(`<h1 id="архитектура-вокзал-тех" tabindex="-1">Архитектура Вокзал.ТЕХ <a class="header-anchor" href="#архитектура-вокзал-тех" aria-label="Permalink to &quot;Архитектура Вокзал.ТЕХ&quot;">​</a></h1><h2 id="обзор" tabindex="-1">Обзор <a class="header-anchor" href="#обзор" aria-label="Permalink to &quot;Обзор&quot;">​</a></h2><p>Вокзал.ТЕХ построен на современной микросервисной архитектуре с использованием event-driven подхода.</p><h2 id="высокоуровневая-архитектура" tabindex="-1">Высокоуровневая архитектура <a class="header-anchor" href="#высокоуровневая-архитектура" aria-label="Permalink to &quot;Высокоуровневая архитектура&quot;">​</a></h2><div class="language-mermaid vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">mermaid</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">graph TB</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    subgraph &quot;Frontend Applications&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        AdminPanel[Admin Panel&lt;br/&gt;React]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        PassengerPortal[Passenger Portal&lt;br/&gt;React]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        BoardDisplay[Board Display&lt;br/&gt;React + WebSocket]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        POS[POS App&lt;br/&gt;Tauri]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Controller[Controller App&lt;br/&gt;React PWA]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    end</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    subgraph &quot;API Gateway&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Traefik[Traefik&lt;br/&gt;API Gateway]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    end</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    subgraph &quot;Microservices&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Auth[Auth Service&lt;br/&gt;:8081]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Schedule[Schedule Service&lt;br/&gt;:8082]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Ticket[Ticket Service&lt;br/&gt;:8083]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Payment[Payment Service&lt;br/&gt;:8085]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Fiscal[Fiscal Service&lt;br/&gt;:8084]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Board[Board Service&lt;br/&gt;:8086]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Notify[Notify Service&lt;br/&gt;:8087]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Audit[Audit Service&lt;br/&gt;:8088]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Document[Document Service&lt;br/&gt;:8089]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Geo[Geo Service&lt;br/&gt;:8090]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    end</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    subgraph &quot;Data Layer&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Postgres[(PostgreSQL)]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Redis[(Redis)]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        MinIO[(MinIO)]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    end</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    subgraph &quot;Message Bus&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        NATS[NATS]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    end</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    subgraph &quot;External Services&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        ATOL[АТОЛ ККТ]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Tinkoff[Tinkoff Acquiring]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        SBP[СБП]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        SMS[SMS.ru]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Yandex[Yandex Maps]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    end</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    AdminPanel --&gt; Traefik</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    PassengerPortal --&gt; Traefik</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    BoardDisplay --&gt; Traefik</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    POS --&gt; Traefik</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Controller --&gt; Traefik</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik --&gt; Auth</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik --&gt; Schedule</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik --&gt; Ticket</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik --&gt; Payment</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik --&gt; Board</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik --&gt; Notify</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik --&gt; Document</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik --&gt; Geo</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Auth --&gt; Postgres</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Schedule --&gt; Postgres</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket --&gt; Postgres</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Payment --&gt; Postgres</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Fiscal --&gt; Postgres</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Audit --&gt; Postgres</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Document --&gt; Postgres</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Notify --&gt; Postgres</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Auth --&gt; Redis</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Board --&gt; Redis</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Payment --&gt; Redis</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Document --&gt; MinIO</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket --&gt; NATS</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Payment --&gt; NATS</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Fiscal --&gt; NATS</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Notify --&gt; NATS</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Audit --&gt; NATS</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Fiscal --&gt; ATOL</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Payment --&gt; Tinkoff</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Payment --&gt; SBP</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Notify --&gt; SMS</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Geo --&gt; Yandex</span></span></code></pre></div><h2 id="микросервисы" tabindex="-1">Микросервисы <a class="header-anchor" href="#микросервисы" aria-label="Permalink to &quot;Микросервисы&quot;">​</a></h2><h3 id="auth-service-port-8081" tabindex="-1">Auth Service (Port 8081) <a class="header-anchor" href="#auth-service-port-8081" aria-label="Permalink to &quot;Auth Service (Port 8081)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>Аутентификация пользователей</li><li>Выдача и обновление JWT токенов</li><li>Управление ролями (RBAC)</li><li>Session management</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>PostgreSQL</li><li>Redis (sessions)</li><li>JWT-Go</li></ul><p><strong>API Endpoints:</strong></p><ul><li><code>POST /v1/auth/login</code></li><li><code>POST /v1/auth/refresh</code></li><li><code>POST /v1/auth/logout</code></li><li><code>GET /v1/auth/me</code></li></ul><h3 id="schedule-service-port-8082" tabindex="-1">Schedule Service (Port 8082) <a class="header-anchor" href="#schedule-service-port-8082" aria-label="Permalink to &quot;Schedule Service (Port 8082)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>CRUD маршрутов</li><li>Управление расписанием</li><li>Создание и отмена рейсов</li><li>Поиск рейсов</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>PostgreSQL</li><li>NATS (events)</li></ul><p><strong>API Endpoints:</strong></p><ul><li><code>GET /v1/routes</code></li><li><code>POST /v1/routes</code></li><li><code>GET /v1/schedules</code></li><li><code>GET /v1/trips/search</code></li></ul><h3 id="ticket-service-port-8083" tabindex="-1">Ticket Service (Port 8083) <a class="header-anchor" href="#ticket-service-port-8083" aria-label="Permalink to &quot;Ticket Service (Port 8083)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>Продажа билетов</li><li>Возврат билетов</li><li>Блокировка возвратов</li><li>Отметка посадки</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>PostgreSQL</li><li>Redis (locks)</li><li>NATS (events)</li></ul><p><strong>API Endpoints:</strong></p><ul><li><code>POST /v1/tickets/sell</code></li><li><code>POST /v1/tickets/:id/refund</code></li><li><code>POST /v1/tickets/:id/board</code></li></ul><h3 id="payment-service-port-8085" tabindex="-1">Payment Service (Port 8085) <a class="header-anchor" href="#payment-service-port-8085" aria-label="Permalink to &quot;Payment Service (Port 8085)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>Инициализация платежей</li><li>Tinkoff Acquiring</li><li>СБП (QR коды)</li><li>Webhooks от провайдеров</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>PostgreSQL</li><li>NATS (events)</li></ul><p><strong>Интеграции:</strong></p><ul><li>Tinkoff Acquiring API v2</li><li>СБП (Система быстрых платежей)</li></ul><h3 id="fiscal-service-port-8084" tabindex="-1">Fiscal Service (Port 8084) <a class="header-anchor" href="#fiscal-service-port-8084" aria-label="Permalink to &quot;Fiscal Service (Port 8084)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>Фискализация продаж</li><li>Z-отчёты</li><li>Интеграция с АТОЛ ККТ</li><li>История фискальных операций</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>PostgreSQL</li><li>АТОЛ Fiscal Driver</li></ul><h3 id="board-service-port-8086" tabindex="-1">Board Service (Port 8086) <a class="header-anchor" href="#board-service-port-8086" aria-label="Permalink to &quot;Board Service (Port 8086)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>WebSocket сервер для табло</li><li>Real-time обновления рейсов</li><li>Кэширование данных</li><li>Голосовые объявления</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>Redis (кэш)</li><li>WebSocket</li><li>TTS (eSpeak)</li></ul><h3 id="notify-service-port-8087" tabindex="-1">Notify Service (Port 8087) <a class="header-anchor" href="#notify-service-port-8087" aria-label="Permalink to &quot;Notify Service (Port 8087)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>SMS уведомления</li><li>Email рассылка</li><li>Telegram боты</li><li>TTS объявления</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>SMTP</li><li>SMS.ru API</li><li>Telegram Bot API</li></ul><h3 id="audit-service-port-8088" tabindex="-1">Audit Service (Port 8088) <a class="header-anchor" href="#audit-service-port-8088" aria-label="Permalink to &quot;Audit Service (Port 8088)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>Логирование всех операций</li><li>Соответствие 152-ФЗ</li><li>Immutable audit log</li><li>Отчёты по операциям</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>PostgreSQL (append-only)</li></ul><h3 id="document-service-port-8089" tabindex="-1">Document Service (Port 8089) <a class="header-anchor" href="#document-service-port-8089" aria-label="Permalink to &quot;Document Service (Port 8089)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>Генерация ПД-2</li><li>Кастомные шаблоны PDF</li><li>Хранение документов</li><li>Электронные билеты</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>MinIO (хранение)</li><li>PDF generation</li></ul><h3 id="geo-service-port-8090" tabindex="-1">Geo Service (Port 8090) <a class="header-anchor" href="#geo-service-port-8090" aria-label="Permalink to &quot;Geo Service (Port 8090)&quot;">​</a></h3><p><strong>Ответственность:</strong></p><ul><li>Геокодирование адресов</li><li>Расчёт расстояний</li><li>Время в пути</li><li>Интеграция с Yandex Maps</li></ul><p><strong>Технологии:</strong></p><ul><li>Go + Gin</li><li>Yandex Maps API</li></ul><h2 id="event-driven-communication" tabindex="-1">Event-Driven Communication <a class="header-anchor" href="#event-driven-communication" aria-label="Permalink to &quot;Event-Driven Communication&quot;">​</a></h2><h3 id="nats-subjects" tabindex="-1">NATS Subjects <a class="header-anchor" href="#nats-subjects" aria-label="Permalink to &quot;NATS Subjects&quot;">​</a></h3><div class="language- vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang"></span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>ticket.sold          — Билет продан</span></span>
<span class="line"><span>ticket.refunded      — Билет возвращён</span></span>
<span class="line"><span>ticket.boarded       — Пассажир сел в автобус</span></span>
<span class="line"><span></span></span>
<span class="line"><span>payment.confirmed    — Платёж подтверждён</span></span>
<span class="line"><span>payment.failed       — Платёж не прошёл</span></span>
<span class="line"><span></span></span>
<span class="line"><span>fiscal.registered    — Чек пробит</span></span>
<span class="line"><span>fiscal.report        — Z-отчёт сформирован</span></span>
<span class="line"><span></span></span>
<span class="line"><span>trip.created         — Рейс создан</span></span>
<span class="line"><span>trip.cancelled       — Рейс отменён</span></span>
<span class="line"><span>trip.departed        — Автобус отправился</span></span></code></pre></div><h3 id="event-flow-example-продажа-билета" tabindex="-1">Event Flow Example: Продажа билета <a class="header-anchor" href="#event-flow-example-продажа-билета" aria-label="Permalink to &quot;Event Flow Example: Продажа билета&quot;">​</a></h3><div class="language-mermaid vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">mermaid</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">sequenceDiagram</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant POS</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Ticket</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Payment</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Fiscal</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Audit</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant NATS</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    POS-&gt;&gt;Ticket: POST /tickets/sell</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket-&gt;&gt;Payment: Инициализация платежа</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Payment-&gt;&gt;Tinkoff: Создать платёж</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Tinkoff--&gt;&gt;Payment: PaymentURL</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Payment--&gt;&gt;POS: PaymentURL</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    POS-&gt;&gt;POS: Покупатель оплачивает</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Tinkoff-&gt;&gt;Payment: Webhook (CONFIRMED)</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Payment-&gt;&gt;NATS: Publish payment.confirmed</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket-&gt;&gt;Ticket: Subscribe payment.confirmed</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket-&gt;&gt;Fiscal: Фискализировать</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Fiscal-&gt;&gt;ATOL: Пробить чек</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    ATOL--&gt;&gt;Fiscal: Фискальный документ</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Fiscal-&gt;&gt;NATS: Publish fiscal.registered</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket-&gt;&gt;NATS: Publish ticket.sold</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Audit-&gt;&gt;Audit: Subscribe all events</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Audit-&gt;&gt;Audit: Записать в audit_log</span></span></code></pre></div><h2 id="data-layer" tabindex="-1">Data Layer <a class="header-anchor" href="#data-layer" aria-label="Permalink to &quot;Data Layer&quot;">​</a></h2><h3 id="postgresql-schema" tabindex="-1">PostgreSQL Schema <a class="header-anchor" href="#postgresql-schema" aria-label="Permalink to &quot;PostgreSQL Schema&quot;">​</a></h3><p><strong>stations</strong></p><ul><li>id, name, city, address, coordinates, timezone</li></ul><p><strong>routes</strong></p><ul><li>id, from_station_id, to_station_id, distance_km, duration_minutes</li></ul><p><strong>schedules</strong></p><ul><li>id, route_id, departure_time, arrival_time, days_of_week</li></ul><p><strong>trips</strong></p><ul><li>id, schedule_id, bus_id, driver_id, departure_date, status</li></ul><p><strong>tickets</strong></p><ul><li>id, trip_id, seat_id, price, passenger_name, status, sold_at</li></ul><p><strong>payments</strong></p><ul><li>id, ticket_id, amount, method, provider, status, confirmed_at</li></ul><p><strong>users</strong></p><ul><li>id, username, password_hash, role, station_id</li></ul><p><strong>audit_logs</strong> (append-only)</p><ul><li>id, user_id, action, resource, details, timestamp</li></ul><h3 id="redis-usage" tabindex="-1">Redis Usage <a class="header-anchor" href="#redis-usage" aria-label="Permalink to &quot;Redis Usage&quot;">​</a></h3><div class="language- vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang"></span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>sessions:{user_id}        — Сессии пользователей</span></span>
<span class="line"><span>locks:seat:{trip_id}:{id} — Блокировки мест</span></span>
<span class="line"><span>cache:trips:{date}        — Кэш рейсов</span></span>
<span class="line"><span>ws:clients:{board_id}     — WebSocket клиенты</span></span></code></pre></div><h3 id="minio-buckets" tabindex="-1">MinIO Buckets <a class="header-anchor" href="#minio-buckets" aria-label="Permalink to &quot;MinIO Buckets&quot;">​</a></h3><div class="language- vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang"></span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>documents    — PDF документы (билеты, отчёты)</span></span>
<span class="line"><span>templates    — Шаблоны документов</span></span>
<span class="line"><span>invoices     — Счета и акты</span></span>
<span class="line"><span>reports      — Z-отчёты, сводки</span></span></code></pre></div><h2 id="security-architecture" tabindex="-1">Security Architecture <a class="header-anchor" href="#security-architecture" aria-label="Permalink to &quot;Security Architecture&quot;">​</a></h2><h3 id="authentication-flow" tabindex="-1">Authentication Flow <a class="header-anchor" href="#authentication-flow" aria-label="Permalink to &quot;Authentication Flow&quot;">​</a></h3><div class="language-mermaid vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">mermaid</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">sequenceDiagram</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Client</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Traefik</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Auth</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Service</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Redis</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Client-&gt;&gt;Traefik: POST /auth/login</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik-&gt;&gt;Auth: Forward request</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Auth-&gt;&gt;Auth: Validate credentials</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Auth-&gt;&gt;Redis: Create session</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Auth-&gt;&gt;Auth: Generate JWT</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Auth--&gt;&gt;Client: access_token + refresh_token</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    </span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Client-&gt;&gt;Traefik: GET /tickets (with token)</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik-&gt;&gt;Traefik: JWT Middleware</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Traefik-&gt;&gt;Service: Forward (с user_id)</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Service--&gt;&gt;Client: Response</span></span></code></pre></div><h3 id="rbac-matrix" tabindex="-1">RBAC Matrix <a class="header-anchor" href="#rbac-matrix" aria-label="Permalink to &quot;RBAC Matrix&quot;">​</a></h3><table tabindex="0"><thead><tr><th>Endpoint</th><th>SuperAdmin</th><th>Admin</th><th>Cashier</th><th>Dispatcher</th><th>Controller</th><th>Viewer</th></tr></thead><tbody><tr><td>POST /tickets/sell</td><td>✅</td><td>✅</td><td>✅</td><td>❌</td><td>❌</td><td>❌</td></tr><tr><td>POST /tickets/refund</td><td>✅</td><td>✅</td><td>✅</td><td>❌</td><td>❌</td><td>❌</td></tr><tr><td>POST /tickets/board</td><td>✅</td><td>✅</td><td>❌</td><td>❌</td><td>✅</td><td>❌</td></tr><tr><td>POST /trips</td><td>✅</td><td>✅</td><td>❌</td><td>✅</td><td>❌</td><td>❌</td></tr><tr><td>POST /routes</td><td>✅</td><td>✅</td><td>❌</td><td>✅</td><td>❌</td><td>❌</td></tr><tr><td>GET /reports</td><td>✅</td><td>✅</td><td>❌</td><td>✅</td><td>❌</td><td>❌</td></tr><tr><td>GET /*</td><td>✅</td><td>✅</td><td>✅</td><td>✅</td><td>✅</td><td>✅</td></tr></tbody></table><h2 id="deployment-architecture" tabindex="-1">Deployment Architecture <a class="header-anchor" href="#deployment-architecture" aria-label="Permalink to &quot;Deployment Architecture&quot;">​</a></h2><h3 id="production-setup" tabindex="-1">Production Setup <a class="header-anchor" href="#production-setup" aria-label="Permalink to &quot;Production Setup&quot;">​</a></h3><div class="language- vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang"></span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>┌─────────────────┐</span></span>
<span class="line"><span>│   CloudFlare    │  ← DNS, CDN, DDoS Protection</span></span>
<span class="line"><span>└────────┬────────┘</span></span>
<span class="line"><span>         │</span></span>
<span class="line"><span>┌────────▼────────┐</span></span>
<span class="line"><span>│  Load Balancer  │  ← AWS ALB / nginx</span></span>
<span class="line"><span>└────────┬────────┘</span></span>
<span class="line"><span>         │</span></span>
<span class="line"><span>┌────────▼────────────────────────────┐</span></span>
<span class="line"><span>│        Kubernetes Cluster           │</span></span>
<span class="line"><span>│  ┌──────────────────────────────┐   │</span></span>
<span class="line"><span>│  │      Ingress (Traefik)       │   │</span></span>
<span class="line"><span>│  └──────────────┬───────────────┘   │</span></span>
<span class="line"><span>│                 │                    │</span></span>
<span class="line"><span>│  ┌──────────────▼───────────────┐   │</span></span>
<span class="line"><span>│  │   Services (Deployments)     │   │</span></span>
<span class="line"><span>│  │  • auth-service (3 replicas) │   │</span></span>
<span class="line"><span>│  │  • ticket-service (5)        │   │</span></span>
<span class="line"><span>│  │  • payment-service (3)       │   │</span></span>
<span class="line"><span>│  │  • etc...                    │   │</span></span>
<span class="line"><span>│  └──────────────────────────────┘   │</span></span>
<span class="line"><span>│                                      │</span></span>
<span class="line"><span>│  ┌──────────────────────────────┐   │</span></span>
<span class="line"><span>│  │     Stateful Sets            │   │</span></span>
<span class="line"><span>│  │  • PostgreSQL (Primary+2Rep) │   │</span></span>
<span class="line"><span>│  │  • Redis Cluster (6 nodes)   │   │</span></span>
<span class="line"><span>│  │  • NATS Cluster (3 nodes)    │   │</span></span>
<span class="line"><span>│  └──────────────────────────────┘   │</span></span>
<span class="line"><span>└──────────────────────────────────────┘</span></span>
<span class="line"><span>         │</span></span>
<span class="line"><span>┌────────▼────────┐</span></span>
<span class="line"><span>│   External      │</span></span>
<span class="line"><span>│   Services      │</span></span>
<span class="line"><span>│  • АТОЛ ККТ     │</span></span>
<span class="line"><span>│  • Tinkoff API  │</span></span>
<span class="line"><span>│  • SMS.ru       │</span></span>
<span class="line"><span>└─────────────────┘</span></span></code></pre></div><h2 id="scaling-strategy" tabindex="-1">Scaling Strategy <a class="header-anchor" href="#scaling-strategy" aria-label="Permalink to &quot;Scaling Strategy&quot;">​</a></h2><h3 id="horizontal-scaling" tabindex="-1">Horizontal Scaling <a class="header-anchor" href="#horizontal-scaling" aria-label="Permalink to &quot;Horizontal Scaling&quot;">​</a></h3><ul><li>Stateless сервисы масштабируются легко (3-10 реплик)</li><li>Используем Kubernetes HPA (CPU/Memory based)</li></ul><h3 id="vertical-scaling" tabindex="-1">Vertical Scaling <a class="header-anchor" href="#vertical-scaling" aria-label="Permalink to &quot;Vertical Scaling&quot;">​</a></h3><ul><li>PostgreSQL: read replicas для тяжёлых запросов</li><li>Redis: cluster mode для высоких RPS</li></ul><h3 id="caching-strategy" tabindex="-1">Caching Strategy <a class="header-anchor" href="#caching-strategy" aria-label="Permalink to &quot;Caching Strategy&quot;">​</a></h3><ul><li>Redis для hot data (30 минут)</li><li>CDN для статики (1 день)</li><li>Browser cache для UI (1 час)</li></ul><h2 id="monitoring-observability" tabindex="-1">Monitoring &amp; Observability <a class="header-anchor" href="#monitoring-observability" aria-label="Permalink to &quot;Monitoring &amp; Observability&quot;">​</a></h2><h3 id="prometheus-metrics" tabindex="-1">Prometheus Metrics <a class="header-anchor" href="#prometheus-metrics" aria-label="Permalink to &quot;Prometheus Metrics&quot;">​</a></h3><div class="language- vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang"></span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>http_requests_total</span></span>
<span class="line"><span>http_request_duration_seconds</span></span>
<span class="line"><span>database_queries_total</span></span>
<span class="line"><span>nats_messages_published</span></span></code></pre></div><h3 id="distributed-tracing" tabindex="-1">Distributed Tracing <a class="header-anchor" href="#distributed-tracing" aria-label="Permalink to &quot;Distributed Tracing&quot;">​</a></h3><ul><li>Jaeger для trace запросов</li><li>Correlation IDs через все сервисы</li></ul><h3 id="logging" tabindex="-1">Logging <a class="header-anchor" href="#logging" aria-label="Permalink to &quot;Logging&quot;">​</a></h3><ul><li>Structured JSON logs</li><li>Centralized в Loki</li><li>Retention: 30 дней</li></ul><hr><p>© 2026 Вокзал.ТЕХ</p>`,115)])])}const k=a(e,[["render",l]]);export{E as __pageData,k as default};
