export class Usuario {
    constructor(data = {}) {
        this.id = data.id ?? 0;
        this.empresa_id = data.empresa_id ?? 0;
        this.email = data.email ?? "";
        this.nome = data.nome ?? "";
        this.senha = data.senha ?? "";
        this.criado_em = data.criado_em ?? null;
    }
}