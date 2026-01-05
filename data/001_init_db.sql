\set ON_ERROR_STOP on

CREATE TABLE IF NOT EXISTS public.notes
(
    id          uuid                                          NOT NULL PRIMARY KEY,
    title       text                                          NOT NULL,
    path        text                                          NOT NULL,
    class       text                                          NOT NULL,
    tags        text[]                   DEFAULT '{}'::text[] NOT NULL,
    links       text[]                   DEFAULT '{}'::text[] NOT NULL,
    content     text[]                   DEFAULT '{}'::text[] NOT NULL,
    create_time timestamp with time zone DEFAULT now()        NOT NULL,
    update_time timestamp with time zone DEFAULT now()        NOT NULL
);

-- В init-скриптах официального образа Postgres psql запускается от имени $POSTGRES_USER,
-- поэтому владелец = текущий пользователь (без подстановок переменных).
ALTER TABLE public.notes OWNER TO CURRENT_USER;

-- Загрузка CSV (файл должен существовать в контейнере по пути /data/notes.csv).
\copy public.notes (id,title,path,class,tags,links,content,create_time,update_time) FROM '/docker-entrypoint-initdb.d/notes.csv' WITH (FORMAT csv, HEADER true);

