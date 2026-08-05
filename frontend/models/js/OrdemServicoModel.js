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
        this.criado_em = data.criado_em ?? null;
        this.atualizado_em = data.atualizado_em ?? null;
    }
}