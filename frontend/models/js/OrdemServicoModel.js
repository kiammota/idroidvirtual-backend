export class OrdemServico {
    constructor(data = {}) {
        this.id = data.id ?? 0;
        this.numero_ordem = data.numero_ordem ?? 0;
        this.produto_id = data.produto_id ?? 0;
        this.cliente_id = data.cliente_id ?? 0;
        this.defeito = data.defeito ?? null;
        this.servicos = data.servicos ?? null;
        this.assinatura = data.assinatura ?? null;
        this.status = data.status ?? "";
        this.valor = data.valor != null ? Number(data.valor) : 0;
        this.custo = data.custo != null ? Number(data.custo) : 0;
        this.criado_em = data.criado_em ?? null;
        this.atualizado_em = data.atualizado_em ?? null;
        // Campos auxiliares retornados pela API (join com cliente)
        this.cliente_nome = data.cliente_nome ?? data.nome_cliente ?? null;
        this.telefone = data.telefone ?? null;
    }

    /**
     * Valida os campos da ordem de serviço
     * @returns {string[]} Lista com mensagens de erro
     */
    validar() {
        const erros = [];


        if (!this.defeito || !String(this.defeito).trim()) {
            erros.push("O defeito relatado é obrigatório.");
        }

        if (!this.servicos || !String(this.servicos).trim()) {
            erros.push("O serviço a realizar é obrigatório.");
        }

        if (this.valor < 0) {
            erros.push("O valor não pode ser negativo.");
        }

        if (this.custo < 0) {
            erros.push("O custo não pode ser negativo.");
        }

        return erros;
    }

    /**
     * Verifica se os dados são válidos
     * @returns {boolean}
     */
    isValid() {
        return this.validar().length === 0;
    }
}
