import{_ as t,o as e,c as d,j as n}from"./chunks/framework.D5T0pHrq.js";const T=JSON.parse('{"title":"","description":"","frontmatter":{},"headers":[],"relativePath":"initial/14.md","filePath":"initial/14.md"}'),u={name:"initial/14.md"};function a(s,i,r,_,o,l){return e(),d("div",null,[...i[0]||(i[0]=[n("p",null,'erDiagram STATION ||--o{ BUS : "имеет" STATION ||--o{ DRIVER : "имеет" STATION ||--o{ TICKET : "продаёт" STATION ||--o{ BLOCKING_RULE : "управляет"',-1),n("pre",null,[n("code",null,`ROUTE ||--o{ SCHEDULE : "имеет"
SCHEDULE ||--o{ TRIP : "создаёт"
TRIP }|--|{ TICKET : "содержит"
TRIP }|--|| BUS : "использует"
TRIP }|--|| DRIVER : "назначает"
TRIP }|--o{ BOARDING_EVENT : "имеет"

BUS ||--o{ SEAT : "имеет"

TICKET }|--|| SEAT : "занимает"
TICKET }|--|| USER : "продано"
TICKET }|--o{ RETURN_EVENT : "возврат"

BLOCKING_RULE }|--|| ROUTE : "применяется"
BLOCKING_RULE }|--|| STATION : "ограничивает"

USER }|--o{ AUDIT_LOG : "выполняет"
TRIP }|--o{ DOCUMENT : "генерирует"

STATION {
    uuid station_id PK
    string name
    string code
    text address
    string timezone
}

BUS {
    uuid bus_id PK
    string plate_number
    string model
    int capacity
    string status
    uuid station_id FK
}

DRIVER {
    uuid driver_id PK
    string full_name
    string license_number
    int experience_years
    string phone
    uuid station_id FK
}

ROUTE {
    uuid route_id PK
    string name
    jsonb stops
    decimal distance_km
    int duration_min
    bool is_active
}

SCHEDULE {
    uuid schedule_id PK
    uuid route_id FK
    time departure_time
    jsonb days_of_week
    string platform
    bool is_active
}

TRIP {
    uuid trip_id PK
    uuid schedule_id FK
    date date
    string status
    int delay_minutes
    string platform
    timestamp departure_actual
    timestamp arrival_actual
    uuid bus_id FK
    uuid driver_id FK
}

SEAT {
    uuid seat_id PK
    uuid bus_id FK
    int number
    string type
    bool is_available
}

TICKET {
    uuid ticket_id PK
    uuid trip_id FK
    uuid seat_id FK
    string passenger_name
    string passport
    decimal price
    string status
    decimal return_penalty
    timestamp sold_at
    uuid sold_by FK
    uuid station_id FK
}

BLOCKING_RULE {
    uuid rule_id PK
    uuid station_id FK
    uuid route_id FK
    string seat_range
    string reason
    date valid_from
    date valid_to
}

USER {
    uuid user_id PK
    string username
    string full_name
    string role
    uuid station_id FK
}

AUDIT_LOG {
    uuid log_id PK
    string entity_type
    uuid entity_id
    string action
    uuid user_id FK
    jsonb old_values
    jsonb new_values
    timestamp created_at
}

RETURN_EVENT {
    uuid return_id PK
    uuid ticket_id FK
    string reason
    decimal refunded_amount
    timestamp returned_at
    uuid processed_by FK
}

BOARDING_EVENT {
    uuid boarding_id PK
    uuid ticket_id FK
    bool confirmed
    timestamp marked_at
    uuid marked_by FK
}

DOCUMENT {
    uuid document_id PK
    string type
    uuid trip_id FK
    text content
    string format
    timestamp generated_at
    uuid generated_by FK
}
`)],-1)])])}const K=t(u,[["render",a]]);export{T as __pageData,K as default};
