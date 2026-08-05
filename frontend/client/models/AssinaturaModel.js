export class Assinatura {
    constructor(data = {}) {
        this.id = data.id ?? 0;
        this.empresa_id = data.empresa_id ?? 0;
        this.valor = data.valor ?? 0;
        this.status = data.status ?? "";
        this.inicio_em = data.inicio_em ?? null;
        this.proxima_cobranca_em = data.proxima_cobranca_em ?? null;
        this.cancelada_em = data.cancelada_em ?? null;
        this.criado_em = data.criado_em ?? null;
    }
}