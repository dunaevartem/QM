-- Таблица для хранения сообщений чата
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для быстрой выборки последних сообщений (понадобится при загрузке истории)
CREATE INDEX IF NOT EXISTS idx_created_at ON messages (created_at DESC);

-- Тестовое сообщение (чтобы сразу видеть, что база работает)
INSERT INTO messages (username, content) 
VALUES ('System', 'Добро пожаловать в быстрый Go-чат на K8s!');
