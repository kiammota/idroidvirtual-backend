export class Empresa {
    constructor(data = {}) {
        this.id = data.id ?? 0;
        this.nome = data.nome ?? "";
        this.cnpj = data.cnpj ?? "";
        this.criado_em = data.criado_em ?? null;
    }
}