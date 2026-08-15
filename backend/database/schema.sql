-- =====================================================================
-- SCHEMA COMPLETO - FINANCEIRO ROBUSTO
-- =====================================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================
-- EMPRESAS E USUÁRIOS
-- =========================

CREATE TABLE IF NOT EXISTS empresas (
    id BIGSERIAL PRIMARY KEY,
    cnpj VARCHAR(150) NOT NULL UNIQUE,
    razao_social VARCHAR(150) NOT NULL,
    nome_fantasia VARCHAR(150) NOT NULL,
    inscricao_estadual VARCHAR(150) NOT NULL,
    inscricao_municipal VARCHAR(150) NOT NULL,
    regime_tributario VARCHAR(150) NOT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS usuarios (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    nome VARCHAR(150) NOT NULL,
    senha TEXT NOT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_usuario_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas (id)
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_usuarios_email ON usuarios (email);
CREATE INDEX IF NOT EXISTS idx_usuarios_empresa ON usuarios (empresa_id);

CREATE TABLE IF NOT EXISTS auditorias (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    acao TEXT NOT NULL,
    horario TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_auditoria_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_auditoria_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas (id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_auditoria_usuario ON auditorias (usuario_id);
CREATE INDEX IF NOT EXISTS idx_auditoria_empresa ON auditorias (empresa_id);
CREATE INDEX IF NOT EXISTS idx_auditoria_horario ON auditorias (horario);

-- =========================
-- PRODUTOS E CLIENTES
-- =========================

CREATE TABLE IF NOT EXISTS produtos (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    nome VARCHAR(150) NOT NULL,
    descricao TEXT,
    preco NUMERIC(12, 2) NOT NULL DEFAULT 0,
    categoria TEXT NOT NULL,
    estoque INTEGER NOT NULL DEFAULT 0,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_produto_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas (id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_produtos_empresa ON produtos (empresa_id);
CREATE INDEX IF NOT EXISTS idx_produtos_nome ON produtos (nome);
CREATE INDEX IF NOT EXISTS idx_produtos_estoque ON produtos (estoque);

-- Cliente NÃO tem empresa_id (cadastro global / compartilhado)
CREATE TABLE IF NOT EXISTS clientes (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(150) NOT NULL,
    cpf_cnpj VARCHAR(30),
    telefone VARCHAR(30),
    email VARCHAR(255),
    endereco TEXT,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_clientes_nome ON clientes (nome);
CREATE INDEX IF NOT EXISTS idx_clientes_email ON clientes (email);
CREATE INDEX IF NOT EXISTS idx_clientes_telefone ON clientes (telefone);

-- =========================
-- VENDAS E TROCAS
-- =========================

CREATE TABLE IF NOT EXISTS vendas (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    produto_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    cliente_id BIGINT,
    quantidade INTEGER NOT NULL DEFAULT 1,
    valor_unitario NUMERIC(12, 2) NOT NULL,
    valor_venda NUMERIC(12, 2) NOT NULL,
    valor_troca_cedida NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    desconto NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    total NUMERIC(12, 2) NOT NULL,
    tipo_pagamento VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'CONCLUIDA',
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_venda_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_venda_produto
        FOREIGN KEY (produto_id)
        REFERENCES produtos (id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_venda_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_venda_cliente
        FOREIGN KEY (cliente_id)
        REFERENCES clientes (id)
        ON DELETE SET NULL,
    CONSTRAINT ck_venda_quantidade CHECK (quantidade > 0),
    CONSTRAINT ck_venda_valor_venda CHECK (valor_venda >= 0),
    CONSTRAINT ck_venda_troca_cedida CHECK (valor_troca_cedida >= 0),
    CONSTRAINT ck_venda_desconto CHECK (desconto >= 0),
    CONSTRAINT ck_venda_total CHECK (total >= 0),
    CONSTRAINT ck_venda_status CHECK (status IN ('CONCLUIDA', 'CANCELADA'))
);

CREATE INDEX IF NOT EXISTS idx_vendas_empresa ON vendas (empresa_id);
CREATE INDEX IF NOT EXISTS idx_vendas_produto ON vendas (produto_id);
CREATE INDEX IF NOT EXISTS idx_vendas_usuario ON vendas (usuario_id);
CREATE INDEX IF NOT EXISTS idx_vendas_criado_em ON vendas (criado_em);
CREATE INDEX IF NOT EXISTS idx_vendas_tipo_pagamento ON vendas (tipo_pagamento);

CREATE TABLE IF NOT EXISTS trocas (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    modelo_aparelho VARCHAR(150) NOT NULL,
    imei VARCHAR(50),
    valor_avaliacao NUMERIC(12, 2) NOT NULL,
    condicao VARCHAR(50) DEFAULT 'USADO',
    observacoes TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'RECEBIDA',
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_troca_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_troca_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id)
        ON DELETE RESTRICT,
    CONSTRAINT ck_troca_valor CHECK (valor_avaliacao >= 0),
    CONSTRAINT ck_troca_status CHECK (status IN ('RECEBIDA', 'AVALIADA', 'EM_ESTOQUE', 'VENDIDA', 'DESCARTADA'))
);

CREATE INDEX IF NOT EXISTS idx_trocas_empresa ON trocas (empresa_id);
CREATE INDEX IF NOT EXISTS idx_trocas_usuario ON trocas (usuario_id);
CREATE INDEX IF NOT EXISTS idx_trocas_criado_em ON trocas (criado_em);
CREATE INDEX IF NOT EXISTS idx_trocas_status ON trocas (status);

-- =========================
-- ORDENS DE SERVIÇO
-- =========================

CREATE SEQUENCE IF NOT EXISTS ordens_servico_numero_ordem_seq;

CREATE TABLE IF NOT EXISTS ordens_servico (
    id BIGSERIAL PRIMARY KEY,
    numero_ordem BIGINT NOT NULL DEFAULT nextval('ordens_servico_numero_ordem_seq'),
    empresa_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    cliente_id BIGINT NOT NULL,
    nome_servico VARCHAR(255) NOT NULL,
    defeito TEXT,
    servicos TEXT,
    valor NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    custo NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    assinatura TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'ABERTA',
    data_conclusao TIMESTAMP,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_os_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_os_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_os_cliente
        FOREIGN KEY (cliente_id)
        REFERENCES clientes (id)
        ON DELETE RESTRICT,
    CONSTRAINT ck_os_valor CHECK (valor >= 0),
    CONSTRAINT ck_os_custo CHECK (custo >= 0),
    CONSTRAINT ck_os_status CHECK (status IN ('ABERTA', 'EM_ANDAMENTO', 'FINALIZADA', 'CANCELADA'))
);

ALTER SEQUENCE ordens_servico_numero_ordem_seq OWNED BY ordens_servico.numero_ordem;

SELECT setval(
    'ordens_servico_numero_ordem_seq',
    COALESCE((SELECT MAX(numero_ordem) FROM ordens_servico), 0) + 1,
    false
);

CREATE INDEX IF NOT EXISTS idx_os_empresa ON ordens_servico (empresa_id);
CREATE INDEX IF NOT EXISTS idx_os_numero_ordem ON ordens_servico (numero_ordem);
CREATE INDEX IF NOT EXISTS idx_os_usuario ON ordens_servico (usuario_id);
CREATE INDEX IF NOT EXISTS idx_os_cliente ON ordens_servico (cliente_id);
CREATE INDEX IF NOT EXISTS idx_os_status ON ordens_servico (status);
CREATE INDEX IF NOT EXISTS idx_os_criado_em ON ordens_servico (criado_em);
CREATE INDEX IF NOT EXISTS idx_os_data_conclusao ON ordens_servico (data_conclusao);

-- =========================
-- DESPESAS
-- =========================

CREATE TABLE IF NOT EXISTS despesas (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    usuario_id BIGINT,
    descricao VARCHAR(255) NOT NULL,
    categoria VARCHAR(50) NOT NULL,
    valor NUMERIC(12, 2) NOT NULL,
    valor_pago NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    forma_pagamento VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE',
    data_vencimento DATE NOT NULL,
    data_pagamento DATE,
    comprovante_url TEXT,
    numero_documento VARCHAR(100),
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_despesa_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_despesa_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id)
        ON DELETE SET NULL,
    CONSTRAINT ck_despesa_valor CHECK (valor >= 0),
    CONSTRAINT ck_despesa_valor_pago CHECK (valor_pago >= 0),
    CONSTRAINT ck_despesa_status CHECK (status IN ('PENDENTE', 'PAGO', 'ATRASADO', 'CANCELADO'))
);

CREATE INDEX IF NOT EXISTS idx_despesas_empresa ON despesas (empresa_id);
CREATE INDEX IF NOT EXISTS idx_despesas_usuario ON despesas (usuario_id);
CREATE INDEX IF NOT EXISTS idx_despesas_status ON despesas (status);
CREATE INDEX IF NOT EXISTS idx_despesas_data_vencimento ON despesas (data_vencimento);
CREATE INDEX IF NOT EXISTS idx_despesas_data_pagamento ON despesas (data_pagamento);

-- =========================
-- ASSINATURAS
-- =========================

CREATE TABLE IF NOT EXISTS assinaturas (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL UNIQUE,
    valor NUMERIC(12, 2) NOT NULL DEFAULT 140.00,
    status VARCHAR(30) NOT NULL DEFAULT 'ATIVA',
    inicio_em TIMESTAMP NOT NULL DEFAULT NOW(),
    proxima_cobranca_em TIMESTAMP NOT NULL,
    cancelada_em TIMESTAMP,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_assinatura_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas (id)
        ON DELETE RESTRICT,
    CONSTRAINT ck_assinatura_status CHECK (status IN ('ATIVA', 'CANCELADA', 'SUSPENSA'))
);

CREATE INDEX IF NOT EXISTS idx_assinaturas_empresa ON assinaturas (empresa_id);
CREATE INDEX IF NOT EXISTS idx_assinaturas_status ON assinaturas (status);

-- =========================
-- MOVIMENTAÇÃO DE ESTOQUE
-- =========================

CREATE TABLE IF NOT EXISTS movimentacoes_estoque (
    id BIGSERIAL PRIMARY KEY,
    empresa_id BIGINT NOT NULL,
    produto_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    quantidade INT NOT NULL,
    motivo TEXT,
    custo_unitario NUMERIC(12, 2) NOT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_movimentacao_empresa
        FOREIGN KEY (empresa_id)
        REFERENCES empresas (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_movimentacao_produto
        FOREIGN KEY (produto_id)
        REFERENCES produtos (id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_movimentacao_usuario
        FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id)
        ON DELETE RESTRICT,
    CONSTRAINT ck_movimento_tipo CHECK (tipo IN ('ENTRADA_COMPRA', 'PERDA', 'AVARIA', 'BRINDE', 'VENDA'))
);

CREATE INDEX IF NOT EXISTS idx_movimentacao_empresa ON movimentacoes_estoque (empresa_id);
CREATE INDEX IF NOT EXISTS idx_movimentacao_produto ON movimentacoes_estoque (produto_id);
CREATE INDEX IF NOT EXISTS idx_movimentacao_tipo ON movimentacoes_estoque (tipo);
CREATE INDEX IF NOT EXISTS idx_movimentacao_criado_em ON movimentacoes_estoque (criado_em);
