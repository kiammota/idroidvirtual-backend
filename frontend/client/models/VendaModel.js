export class Venda {
    constructor(data = {}) {
        this.id = Number(data.id) || 0;
        this.produto_id = Number(data.produto_id) || 0;
        this.usuario_id = Number(data.usuario_id) || 0;
        this.quantidade = Number(data.quantidade) || 0;
        this.total = Number(data.total) || 0;
        this.tipo_pagamento = data.tipo_pagamento ?? "";
        this.status = data.status ?? "PENDENTE";
        this.criado_em = data.criado_em ?? null;

        this.validar();
    }

    validar() {
        if (this.quantidade <= 0) {
            throw new Error("A quantidade da venda deve ser maior que zero.");
        }

        if (this.total < 0) {
            throw new Error("O valor total da venda não pode ser negativo.");
        }

        const statusPermitidos = ["PENDENTE", "CONCLUIDA", "CANCELADA"];
        if (this.status && !statusPermitidos.includes(this.status)) {
            throw new Error(`Status inválido: ${this.status}. Permitidos: ${statusPermitidos.join(", ")}`);
        }
    }
}