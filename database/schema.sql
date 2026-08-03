CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================
-- USUÁRIOS
-- =========================

CREATE TABLE empresas (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(150) NOT NULL,
    cnpj VARCHAR(150) NOT NULL UNIQUE,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE usuarios (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    nome VARCHAR(150) NOT NULL,
    senha TEXT NOT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_usuario_empresa
        FOREIGN KEY(empresa_id)
        REFERENCES empresas(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_usuarios_email
ON usuarios(email);

CREATE INDEX idx_usuarios_empresa
ON usuarios(empresa_id);

CREATE TABLE assinaturas (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL UNIQUE,
    valor NUMERIC(12,2) NOT NULL DEFAULT 140.00,
    status VARCHAR(30) NOT NULL DEFAULT 'ATIVA',
    inicio_em TIMESTAMP NOT NULL DEFAULT NOW(),
    proxima_cobranca_em TIMESTAMP NOT NULL,
    cancelada_em TIMESTAMP,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_assinatura_empresa
        FOREIGN KEY(empresa_id)
        REFERENCES empresas(id)
        ON DELETE RESTRICT
);

-- =========================
-- PRODUTOS
-- =========================

CREATE TABLE produtos (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(150) NOT NULL,
    preco NUMERIC(12,2) NOT NULL DEFAULT 0,
    estoque INTEGER NOT NULL DEFAULT 0,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_produtos_nome
ON produtos(nome);


CREATE INDEX idx_produtos_estoque
ON produtos(estoque);


-- =========================
-- CLIENTES
-- =========================

CREATE TABLE clientes (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(150) NOT NULL,
    telefone VARCHAR(30),
    endereco TEXT,
    email VARCHAR(255),
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_clientes_nome
ON clientes(nome);


CREATE INDEX idx_clientes_email
ON clientes(email);


CREATE INDEX idx_clientes_telefone
ON clientes(telefone);


-- =========================
-- VENDAS
-- =========================

CREATE TABLE vendas (
    id BIGSERIAL PRIMARY KEY,
    produto_id BIGINT NOT NULL,
    quantidade INTEGER NOT NULL DEFAULT 1,
    total NUMERIC(12,2) NOT NULL,
    tipo_pagamento VARCHAR(50) NOT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_venda_produto
        FOREIGN KEY(produto_id)
        REFERENCES produtos(id)
        ON DELETE RESTRICT
);


CREATE INDEX idx_vendas_produto_id
ON vendas(produto_id);


CREATE INDEX idx_vendas_criado_em
ON vendas(criado_em);


CREATE INDEX idx_vendas_tipo_pagamento
ON vendas(tipo_pagamento);



-- =========================
-- ORDEM DE SERVIÇO
-- =========================

CREATE TABLE ordens_servico (
    id BIGSERIAL PRIMARY KEY,
    numero_ordem BIGINT NOT NULL UNIQUE,
    cliente_id BIGINT NOT NULL,
    defeito TEXT,
    servicos TEXT,
    assinatura TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'ABERTA',
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_os_cliente
        FOREIGN KEY(cliente_id)
        REFERENCES clientes(id)
        ON DELETE RESTRICT
);


CREATE INDEX idx_os_numero_ordem
ON ordens_servico(numero_ordem);


CREATE INDEX idx_os_cliente_id
ON ordens_servico(cliente_id);


CREATE INDEX idx_os_status
ON ordens_servico(status);



-- =========================
-- AUDITORIA
-- =========================

CREATE TABLE auditorias (
    id BIGSERIAL PRIMARY KEY,
    usuario_id BIGINT NOT NULL,
    acao VARCHAR(50) NOT NULL,
    tabela VARCHAR(100) NOT NULL,
    registro_id BIGINT NOT NULL,
    valores_antigos JSONB,
    valores_novos JSONB,
    horario TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_auditoria_usuario
        FOREIGN KEY(usuario_id)
        REFERENCES usuarios(id)
        ON DELETE RESTRICT
);

CREATE TABLE movimentacoes_estoque (
    id BIGSERIAL PRIMARY KEY,
    produto_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    tipo VARCHAR(20) NOT NULL, -- 'PERDA', 'AVARIA', 'BRINDE', 'ENTRADA_COMPRA'
    quantidade INT NOT NULL,
    motivo TEXT,
    custo_unitario NUMERIC(10, 2) NOT NULL, -- Custo do produto para calcular o prejuízo real
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_movimentacao_produto FOREIGN KEY (produto_id) REFERENCES produtos(id)
);



CREATE INDEX idx_auditoria_usuario
ON auditorias(usuario_id);


CREATE INDEX idx_auditoria_tabela
ON auditorias(tabela);


CREATE INDEX idx_auditoria_registro
ON auditorias(registro_id);


CREATE INDEX idx_auditoria_horario
ON auditorias(horario);

-- Índice útil para consultas de histórico
CREATE INDEX idx_auditoria_tabela_registro
ON auditorias(tabela, registro_id);