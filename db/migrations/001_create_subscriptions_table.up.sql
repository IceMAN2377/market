CREATE TABLE subscriptions (
                               id SERIAL PRIMARY KEY,
                               service_name VARCHAR(255) NOT NULL,
                               price INTEGER NOT NULL CHECK (price > 0),
                               user_id UUID NOT NULL,
                               start_date VARCHAR(7) NOT NULL,
                               end_date VARCHAR(7),
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для оптимизации запросов
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);


-- Составной индекс для фильтрации по пользователю и сервису
CREATE INDEX idx_subscriptions_user_service ON subscriptions(user_id, service_name);

-- Комментарии к таблице и полям
COMMENT ON TABLE subscriptions IS 'Таблица для хранения информации о подписках пользователей';
COMMENT ON COLUMN subscriptions.id IS 'Уникальный идентификатор подписки';
COMMENT ON COLUMN subscriptions.service_name IS 'Название сервиса, предоставляющего подписку';
COMMENT ON COLUMN subscriptions.price IS 'Стоимость месячной подписки в рублях';
COMMENT ON COLUMN subscriptions.user_id IS 'Идентификатор пользователя в формате UUID';
COMMENT ON COLUMN subscriptions.start_date IS 'Дата начала подписки в формате MM-YYYY';
COMMENT ON COLUMN subscriptions.end_date IS 'Дата окончания подписки в формате MM-YYYY';