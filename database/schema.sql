-- Tabela de Usuários
CREATE TABLE IF NOT EXISTS usuarios (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    nome VARCHAR(255) NOT NULL,
    senha VARCHAR(255) NOT NULL,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_usuarios_email ON usuarios(email);

-- Tabela de Produtos (Estoque)
CREATE TABLE IF NOT EXISTS produtos (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    preco DECIMAL(10, 2) NOT NULL,
    estoque INT NOT NULL DEFAULT 0,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_produtos_nome ON produtos(nome);

-- Tabela de Clientes
CREATE TABLE IF NOT EXISTS clientes (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    telefone VARCHAR(14) NOT NULL,
    endereco VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clientes_telefone ON clientes(telefone);
CREATE INDEX idx_clientes_email ON clientes(email);

-- Tabela de Vendas
CREATE TABLE IF NOT EXISTS vendas (
    id BIGSERIAL PRIMARY KEY,
    produto_id BIGINT NOT NULL REFERENCES produtos(id),
    quantidade INT NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    tipo_pagamento VARCHAR(50) NOT NULL,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vendas_data ON vendas(criado_em);
CREATE INDEX idx_vendas_produto ON vendas(produto_id);

-- Tabela de Ordens de Serviço
CREATE TABLE IF NOT EXISTS ordens_servico (
    id BIGSERIAL PRIMARY KEY,
    numero_ordem BIGINT NOT NULL UNIQUE,
    cliente_id BIGINT NOT NULL REFERENCES clientes(id),
    defeito TEXT NOT NULL,
    servicos TEXT NOT NULL,
    assinatura TEXT,
    status VARCHAR(50) DEFAULT 'pendente',
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ordens_servico_status ON ordens_servico(status);
CREATE INDEX idx_ordens_servico_data ON ordens_servico(criado_em);
CREATE INDEX idx_ordens_servico_numero ON ordens_servico(numero_ordem);

-- Tabela de Auditoria
CREATE TABLE IF NOT EXISTS auditoria (
    id BIGSERIAL PRIMARY KEY,
    usuario_id BIGINT NOT NULL REFERENCES usuarios(id),
    acao TEXT NOT NULL,
    tabela VARCHAR(50),
    registro_id BIGINT,
    valores_antigos JSONB,
    valores_novos JSONB,
    horario TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_auditoria_usuario ON auditoria(usuario_id);
CREATE INDEX idx_auditoria_horario ON auditoria(horario);
CREATE INDEX idx_auditoria_tabela ON auditoria(tabela);