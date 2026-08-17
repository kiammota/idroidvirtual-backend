CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================
-- EMPRESAS
-- =========================
CREATE TABLE empresas (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(150) NOT NULL,
    cnpj VARCHAR(20) NOT NULL UNIQUE,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW()
);

-- =========================
-- USUÁRIOS
-- =========================
CREATE TABLE usuarios (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    nome VARCHAR(150) NOT NULL,
    senha TEXT NOT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_usuario_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_usuarios_empresa ON usuarios(empresa_id);

-- =========================
-- PRODUTOS
-- =========================
CREATE TABLE produtos (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    nome VARCHAR(150) NOT NULL,
    preco NUMERIC(12,2) NOT NULL DEFAULT 0,
    categoria TEXT NOT NULL,
    estoque INTEGER NOT NULL DEFAULT 0,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_produto_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_produtos_nome ON produtos(nome);
CREATE INDEX idx_produtos_estoque ON produtos(estoque);
CREATE INDEX idx_produtos_empresa ON produtos(empresa_id);

-- =========================
-- TROCAS (aparelhos usados recebidos)
-- =========================
CREATE TABLE trocas (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    modelo_aparelho VARCHAR(150) NOT NULL,
    imei VARCHAR(50) NOT NULL,
    valor_avaliacao NUMERIC(12,2) NOT NULL,
    condicao VARCHAR(50) DEFAULT 'USADO',
    observacoes TEXT,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_troca_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_trocas_imei ON trocas(imei);
CREATE INDEX idx_trocas_empresa ON trocas(empresa_id);

CREATE TABLE ordens_servico (
    id BIGSERIAL PRIMARY KEY,
    numero_ordem SERIAL UNIQUE,
    usuario_id BIGINT NOT NULL,
    cliente_id BIGINT NOT NULL,
    nome_servico VARCHAR(255) NOT NULL,
    defeito TEXT,
    servicos TEXT,
    valor NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    custo NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    assinatura TEXT,
    status VARCHAR(50) NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios(id)
        ON DELETE RESTRICT
);

-- =========================
-- VENDAS
-- =========================
CREATE TABLE vendas (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    produto_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    troca_id BIGINT,                          -- NULL se for venda direta
    quantidade INTEGER NOT NULL DEFAULT 1,
    valor_venda NUMERIC(12,2) NOT NULL,
    valor_troca NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    total NUMERIC(12,2) NOT NULL,
    tipo VARCHAR(30) NOT NULL DEFAULT 'VENDA_DIRETA',
    tipo_pagamento VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'CONCLUIDA',
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_venda_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_venda_produto
        FOREIGN KEY (produto_id)
        REFERENCES produtos(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_venda_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_venda_troca
        FOREIGN KEY (troca_id)
        REFERENCES trocas(id)
        ON DELETE SET NULL,

    CONSTRAINT ck_venda_quantidade CHECK (quantidade > 0),
    CONSTRAINT ck_venda_valor_venda CHECK (valor_venda >= 0),
    CONSTRAINT ck_venda_valor_troca CHECK (valor_troca >= 0),
    CONSTRAINT ck_venda_total CHECK (total >= 0),
    CONSTRAINT ck_venda_tipo CHECK (tipo IN ('VENDA_DIRETA', 'TROCA')),
    CONSTRAINT ck_venda_status CHECK (status IN ('CONCLUIDA', 'CANCELADA'))
);

CREATE INDEX idx_vendas_produto ON vendas(produto_id);
CREATE INDEX idx_vendas_usuario ON vendas(usuario_id);
CREATE INDEX idx_vendas_empresa ON vendas(empresa_id);
CREATE INDEX idx_vendas_criado_em ON vendas(criado_em);
CREATE INDEX idx_vendas_tipo_pagamento ON vendas(tipo_pagamento);

-- =========================
-- DESPESAS
-- =========================
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
    data_pagamento DATE,
    comprovante_url TEXT,
    numero_documento VARCHAR(100),
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_despesa_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas(id)
        ON DELETE CASCADE,

    CONSTRAINT ck_despesa_valor CHECK (valor >= 0),
    CONSTRAINT ck_despesa_status CHECK (status IN ('PENDENTE', 'PAGO', 'ATRASADO', 'CANCELADO'))
);

CREATE INDEX idx_despesas_empresa ON despesas(empresa_id);
CREATE INDEX idx_despesas_status ON despesas(status);
CREATE INDEX idx_despesas_vencimento ON despesas(data_vencimento);

-- =========================
-- AUDITORIA
-- =========================
CREATE TABLE auditorias (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    acao TEXT NOT NULL,
    tabela VARCHAR(100),
    registro_id BIGINT,
    detalhes JSONB,
    horario TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_auditoria_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_auditoria_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_auditorias_empresa ON auditorias(empresa_id);
CREATE INDEX idx_auditorias_usuario ON auditorias(usuario_id);
CREATE INDEX idx_auditorias_horario ON auditorias(horario);
CREATE INDEX idx_auditorias_tabela ON auditorias(tabela);
CREATE INDEX idx_auditorias_registro ON auditorias(registro_id);
