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
    categoria TEXT NOT NULL,
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
    telefone VARCHAR(30) UNIQUE,
    endereco TEXT,
    email VARCHAR(255) UNIQUE,
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
    usuario_id BIGINT NOT NULL,
    cliente_id BIGINT DEFAULT NULL,
    quantidade INTEGER NOT NULL DEFAULT 1,
    valor_venda NUMERIC(12,2) NOT NULL,
    aparelho_troca_valor NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    total NUMERIC(12,2) NOT NULL,
    tipo VARCHAR(30) NOT NULL DEFAULT 'VENDA_DIRETA',
    aparelho_troca_modelo VARCHAR(100),
    aparelho_troca_imei VARCHAR(50),
    tipo_pagamento VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'CONCLUIDA',
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_venda_produto
        FOREIGN KEY (produto_id)
        REFERENCES produtos(id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_venda_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios(id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_venda_cliente
        FOREIGN KEY (cliente_id)
        REFERENCES clientes(id)
        ON DELETE SET NULL,
    CONSTRAINT ck_venda_quantidade
        CHECK (quantidade > 0),
    CONSTRAINT ck_venda_valor_venda
        CHECK (valor_venda >= 0),
    CONSTRAINT ck_venda_aparelho_troca_valor
        CHECK (aparelho_troca_valor >= 0),
    CONSTRAINT ck_venda_total
        CHECK (total >= 0),
    CONSTRAINT ck_venda_tipo
        CHECK (tipo IN ('VENDA_DIRETA', 'TROCA')),
    CONSTRAINT ck_venda_status
        CHECK (status IN ('CONCLUIDA', 'CANCELADA'))
);

CREATE TABLE trocas (
    id BIGSERIAL PRIMARY KEY,
    modelo_aparelho VARCHAR(150) NOT NULL,        -- Ex: 'iPhone 12 Pro Max 128GB'
    imei VARCHAR(50) NOT NULL,                    -- IMEI do aparelho usado recebido
    valor_avaliacao NUMERIC(12,2) NOT NULL,       -- Valor abatido na venda (custo de entrada no estoque)
    condicao VARCHAR(50) DEFAULT 'USADO',         -- Ex: 'SEMINOVO', 'DEFEITO', 'EXCELENTE'
    observacoes TEXT,                             -- Avarias, marcas de uso, etc.
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
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
    numero_ordem BIGSERIAL NOT NULL UNIQUE,
    usuario_id BIGINT NOT NULL,
    cliente_id BIGINT NOT NULL,
    nome_servico VARCHAR(255) NOT NULL,
    defeito TEXT,
    servicos TEXT,
    valor NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    assinatura TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'ABERTA',
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_os_cliente
        FOREIGN KEY(cliente_id)
        REFERENCES clientes(id)
        ON DELETE restrict,
    CONSTRAINT fk_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios(id)
        ON DELETE restrict
);


-- Índices recomendados
CREATE INDEX idx_os_numero_ordem ON ordens_servico(numero_ordem);
CREATE INDEX idx_os_cliente_id ON ordens_servico(cliente_id);
CREATE INDEX idx_os_status ON ordens_servico(status);

-- 1. Cria a sequência vinculada à tabela
CREATE SEQUENCE IF NOT EXISTS ordens_servico_numero_ordem_seq;

-- 2. Define o valor padrão da coluna para puxar o próximo número da sequência
ALTER TABLE ordens_servico 
ALTER COLUMN numero_ordem SET DEFAULT nextval('ordens_servico_numero_ordem_seq');

-- 3. Sincroniza o contador da sequência com o maior numero_ordem já existente no banco
SELECT setval(
    'ordens_servico_numero_ordem_seq', 
    COALESCE((SELECT MAX(numero_ordem) FROM ordens_servico), 0) + 1, 
    false
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
    empresa_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    acao TEXT NOT NULL,
    horario TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_auditoria_usuario
        FOREIGN KEY(usuario_id)
        REFERENCES usuarios(id)
        ON DELETE restrict,
    CONSTRAINT fk_empresa
        FOREIGN KEY(empresa_id)
        REFERENCES empresas(id)
        ON DELETE CASCADE
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

CREATE TABLE despesas (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    descricao VARCHAR(255) NOT NULL,
    categoria VARCHAR(50) NOT NULL,
    valor NUMERIC(12,2) NOT NULL,
    valor_pago NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    forma_pagamento VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE',
    data_vencimento DATE NOT NULL,
    data_pagamento DATE DEFAULT NULL,
    comprovante_url TEXT DEFAULT NULL,
    numero_documento VARCHAR(100) DEFAULT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_despesa_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas(id)
        ON DELETE CASCADE,
    CONSTRAINT ck_despesa_valor
        CHECK (valor >= 0),
    CONSTRAINT ck_despesa_status
        CHECK (status IN ('PENDENTE', 'PAGO', 'ATRASADO', 'CANCELADO'))
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