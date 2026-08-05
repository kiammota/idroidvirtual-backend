export class Cliente {
    constructor(data = {}) {
        this.id = data.id ?? 0;
        this.nome = data.nome ?? "";
        this.telefone = data.telefone ?? null;
        this.endereco = data.endereco ?? null;
        this.email = data.email ?? null;
        this.criado_em = data.criado_em ?? null;
        this.atualizado_em = data.atualizado_em ?? null;
    }
}