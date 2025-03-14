-- Создание таблицы song
CREATE TABLE song (
    id              SERIAL PRIMARY KEY,      -- Идентификатор песни
    group_name      VARCHAR(255) NOT NULL,   -- Название группы или исполнителя
    song_name       VARCHAR(255) NOT NULL,   -- Название песни
    release_date    DATE,                    -- Дата выпуска песни
    text            TEXT,                    -- Текст песни
    link            VARCHAR(255)             -- Ссылка на песню
);
