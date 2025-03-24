-- Tabela de usuários
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    firebase_uid VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    picture_url TEXT,
    plan_type VARCHAR(50),
    premium_since TIMESTAMP,
    plan_expiry TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Índices para facilitar buscas por email e firebase_uid
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_firebase_uid ON users (firebase_uid);

-- Tabela de gatos
CREATE TABLE IF NOT EXISTS cats (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    breed VARCHAR(255),
    birth_date DATE,
    gender VARCHAR(50),
    weight_kg NUMERIC(5,2),
    picture_url TEXT,
    is_neutered BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Índice para busca rápida dos gatos por usuário
CREATE INDEX IF NOT EXISTS idx_cats_user_id ON cats (user_id);

-- Tabela de planos
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    max_cats INT NOT NULL,
    has_image_analysis BOOLEAN DEFAULT FALSE,
    has_consultation BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Índices para otimizar consultas por preço e quantidade de gatos permitidos
CREATE INDEX IF NOT EXISTS idx_plans_price ON plans (price);
CREATE INDEX IF NOT EXISTS idx_plans_max_cats ON plans (max_cats);
