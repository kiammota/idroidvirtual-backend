export class Produto {
    constructor(data = {}) {
        this.id = data.id ?? 0;
        this.nome = data.nome ?? "";
        this.preco = data.preco ?? 0;
        this.estoque = data.estoque ?? 0;
        this.categoria = data.categoria ?? "";
        this.criado_em = data.criado_em ?? null;
        this.atualizado_em = data.atualizado_em ?? null;
    }
}