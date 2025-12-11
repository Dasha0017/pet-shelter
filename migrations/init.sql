-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу животных
CREATE TABLE IF NOT EXISTS animals (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    species VARCHAR(50) NOT NULL,
    breed VARCHAR(100),
    age INTEGER,
    gender VARCHAR(10),
    health_status VARCHAR(50),
    description TEXT,
    adopted BOOLEAN DEFAULT FALSE,
    arrival_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу каталога товаров
CREATE TABLE IF NOT EXISTS catalog_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    category VARCHAR(50),
    price DECIMAL(10, 2),
    quantity INTEGER DEFAULT 0,
    description TEXT,
    supplier VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_animals_updated_at BEFORE UPDATE ON animals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_catalog_items_updated_at BEFORE UPDATE ON catalog_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Вставляем тестового администратора (пароль: admin123)
-- В реальном приложении используйте хеширование пароля
INSERT INTO users (username, email, password, role) 
VALUES ('admin', 'admin@petshelter.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.vs7QpJ.9Z5pF3dJd8cU1bHpB6bG1zO', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Вставляем тестовые данные животных
INSERT INTO animals (name, species, breed, age, gender, health_status, description) VALUES
('Барсик', 'Кошка', 'Персидская', 3, 'male', 'Здоров', 'Ласковый и игривый кот'),
('Шарик', 'Собака', 'Лабрадор', 2, 'male', 'Здоров', 'Дружелюбная и активная собака'),
('Мурка', 'Кошка', 'Дворовая', 1, 'female', 'Здорова', 'Скромная и нежная кошечка')
ON CONFLICT DO NOTHING;

-- Вставляем тестовые данные каталога
INSERT INTO catalog_items (name, type, category, price, quantity, description, supplier) VALUES
('Сухой корм для собак', 'Корм', 'Для собак', 25.99, 50, 'Премиум корм для взрослых собак', 'Royal Canin'),
('Игрушка для кошек', 'Игрушка', 'Для кошек', 5.99, 100, 'Мягкая игрушка с кошачьей мятой', 'PetFun'),
('Наполнитель для туалета', 'Аксессуар', 'Для кошек', 15.50, 30, 'Древесный наполнитель, 5 кг', 'CleanPets')
ON CONFLICT DO NOTHING;