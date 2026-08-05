export class Assinatura {
    constructor(data = {}) {
        this.id = Number(data.id) || 0;
        this.empresa_id = Number(data.empresa_id) || 0;
        this.valor = Number(data.valor) || 0;
        this.status = String(data.status ?? "").trim().toUpperCase(); // Ex: ATIVA, CANCELADA, PENDENTE, ATRASADA
        this.inicio_em = data.inicio_em ? new Date(data.inicio_em) : null;
        this.proxima_cobranca_em = data.proxima_cobranca_em ? new Date(data.proxima_cobranca_em) : null;
        this.cancelada_em = data.cancelada_em ? new Date(data.cancelada_em) : null;
        this.criado_em = data.criado_em ? new Date(data.criado_em) : null;
    }

    /**
     * Valida as regras da assinatura
     * @returns {string[]} Lista com os erros de validação
     */
    validar() {
        const erros = [];

        if (!this.empresa_id) {
            erros.push("A empresa_id é obrigatória.");
        }

        if (this.valor <= 0) {
            erros.push("O valor da assinatura deve ser maior que zero.");
        }

        const statusPermitidos = ["ATIVA", "CANCELADA", "PENDENTE", "ATRASADA"];
        if (!this.status || !statusPermitidos.includes(this.status)) {
            erros.push(`Status inválido. Permitidos: ${statusPermitidos.join(", ")}.`);
        }

        if (this.status === "ATIVA" && !this.proxima_cobranca_em) {
            erros.push("Uma assinatura ativa precisa ter uma data de próxima cobrança.");
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

    /**
     * Verifica se a assinatura está ativa
     * @returns {boolean}
     */
    isAtiva() {
        return this.status === "ATIVA";
    }

    /**
     * Retorna o valor formatado em moeda (ex: R$ 99,90)
     * @returns {string}
     */
    getValorFormatado() {
        return this.valor.toLocaleString("pt-BR", {
            style: "currency",
            currency: "BRL"
        });
    }
}