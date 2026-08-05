export class Venda {
    constructor(data = {}) {
        this.id = data.id ?? 0;
        this.produto_id = data.produto_id ?? 0;
        this.usuario_id = data.usuario_id ?? 0;
        this.cliente_id = data.cliente_id ?? 0;
        this.quantidade = data.quantidade ?? 0;
        this.total = data.total ?? 0;
        this.tipo_pagamento = data.tipo_pagamento ?? "";
        this.status = data.status ?? ""; // CONCLUIDA ou CANCELADA
        this.criado_em = data.criado_em ?? null;
    }
}