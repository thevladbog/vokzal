import{_ as a,o as i,c as n,ae as t}from"./chunks/framework.D5T0pHrq.js";const r=JSON.parse('{"title":"🔄 Вокзал.ТЕХ — Полный функционал и диаграммы взаимодействия","description":"","frontmatter":{},"headers":[],"relativePath":"initial/12.md","filePath":"initial/12.md"}'),p={name:"initial/12.md"};function l(E,s,e,k,h,c){return i(),n("div",null,[...s[0]||(s[0]=[t(`<h1 id="🔄-вокзал-тех-—-полныи-функционал-и-диаграммы-взаимодеиствия" tabindex="-1">🔄 Вокзал.ТЕХ — Полный функционал и диаграммы взаимодействия <a class="header-anchor" href="#🔄-вокзал-тех-—-полныи-функционал-и-диаграммы-взаимодеиствия" aria-label="Permalink to &quot;🔄 Вокзал.ТЕХ — Полный функционал и диаграммы взаимодействия&quot;">​</a></h1><h2 id="🧩-основные-сценарии-uml-sequence-diagrams" tabindex="-1">🧩 Основные сценарии (UML Sequence Diagrams) <a class="header-anchor" href="#🧩-основные-сценарии-uml-sequence-diagrams" aria-label="Permalink to &quot;🧩 Основные сценарии (UML Sequence Diagrams)&quot;">​</a></h2><h3 id="_1-продажа-билета" tabindex="-1">1. Продажа билета <a class="header-anchor" href="#_1-продажа-билета" aria-label="Permalink to &quot;1. Продажа билета&quot;">​</a></h3><div class="language-mermaid vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">mermaid</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">sequenceDiagram</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Кассир</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant POS as POS (Tauri)</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Ticket as ticket.vokzal.tech</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Fiscal as fiscal.vokzal.tech</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Printer as Принтер</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant ККТ as ККТ (АТОЛ)</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Кассир-&gt;&gt;POS: Выбирает рейс, место</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    POS-&gt;&gt;Ticket: GET /trips/{id}</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket--&gt;&gt;POS: Данные рейса</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Кассир-&gt;&gt;POS: Вводит данные, оплачивает</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    POS-&gt;&gt;Fiscal: POST /fiscal/receipt (приход)</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Fiscal-&gt;&gt;ККТ: Печать чека</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    ККТ--&gt;&gt;Fiscal: Успех</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Fiscal--&gt;&gt;POS: URL чека</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    POS-&gt;&gt;Ticket: POST /tickets/sell</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket--&gt;&gt;POS: ticket_id, QR</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    POS-&gt;&gt;Printer: Печать билета</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Printer--&gt;&gt;Кассир: Билет распечатан</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">2. Возврат билетаsequenceDiagram</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Кассир</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant POS</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Ticket</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Audit as audit.vokzal.tech</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Fiscal</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant ККТ</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Кассир-&gt;&gt;POS: Сканирует QR билета</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    POS-&gt;&gt;Ticket: GET /tickets/{id}</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket--&gt;&gt;POS: Данные билета, статус</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    POS-&gt;&gt;Ticket: Проверка: был ли посажен?</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket--&gt;&gt;POS: Да/Нет</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    alt Если не посажен</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        POS-&gt;&gt;Fiscal: POST /fiscal/refund</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Fiscal-&gt;&gt;ККТ: Печать чека возврата</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        ККТ--&gt;&gt;Fiscal: Успех</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Fiscal--&gt;&gt;POS: URL чека</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        POS-&gt;&gt;Ticket: PATCH /tickets/{id}/return</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Ticket-&gt;&gt;Audit: Логирование возврата</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Audit--&gt;&gt;Ticket: Успех</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        Ticket--&gt;&gt;POS: Возврат подтверждён</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    else Если посажен</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        POS--&gt;&gt;Кассир: Возврат запрещён</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    end</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">3. Начало посадкиsequenceDiagram</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Диспетчер</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Admin as admin.vokzal.tech</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Board as board.vokzal.tech</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant Ticket</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant PA as Голосовой робот</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Диспетчер-&gt;&gt;Admin: Нажимает &quot;Начать посадку&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Admin-&gt;&gt;Ticket: POST /boarding/start</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket-&gt;&gt;Ticket: Блокировка возвратов на рейс</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Ticket--&gt;&gt;Admin: Успех</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Admin-&gt;&gt;Board: Обновить табло</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Board--&gt;&gt;Табло: &quot;Началась посадка&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    Admin-&gt;&gt;PA: POST /announcements</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    PA-&gt;&gt;PA: Генерация голоса (TTS)</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    PA--&gt;&gt;Аудиосистема: &quot;Посадка на рейс в Казань начнётся у перрона 3&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">4. Синхронизация между вокзаламиsequenceDiagram</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant ВокзалА</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant API as api.vokzal.tech</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    participant ВокзалБ</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    ВокзалА-&gt;&gt;API: POST /sync/trip-status</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    API-&gt;&gt;API: Валидация, обновление</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    API-&gt;&gt;ВокзалБ: Вебхук: рейс прибыл</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    ВокзалБ-&gt;&gt;API: GET /trips/{id}</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    API--&gt;&gt;ВокзалБ: Обновлённые данные</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    ВокзалБ-&gt;&gt;Табло: Обновить расписание</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">📋 Полный список функцийМодульФункцияПриоритетРасписаниеСоздание маршрутовВысокийКорректировка рейсовВысокийПроходящие рейсыСреднийПродажаПродажа с местомВысокийПродажа без местаВысокийАнонимная продажаВысокийВозврат с штрафамиВысокийПосадкаРучная фиксацияВысокийСканирование ШКСреднийБлокировка возвратовВысокийТаблоОбщее таблоВысокийПерронное таблоВысокийАвтообновлениеВысокийДиспетчеризацияНазначение перроновВысокийОповещение (TTS)СреднийУправление задержкамиВысокийДокументыГенерация ПД-2ВысокийКастомные шаблоныСреднийИнтеграцииККТ (АТОЛ)ВысокийЭквайрингВысокийTelegram-ботыСреднийYandex MapsВысокийБезопасностьАудит (152-ФЗ)ВысокийФискализация (54-ФЗ)ВысокийRBACВысокий</span></span></code></pre></div>`,4)])])}const d=a(p,[["render",l]]);export{r as __pageData,d as default};
