export class Empresa {
    constructor(data = {}) {
        this.id = Number(data.id) || 0;
        this.nome = String(data.nome ?? "").trim().toUpperCase();
        this.cnpj = String(data.cnpj ?? "").replace(/\D/g, ""); // Salva apenas os números do CNPJ
        this.criado_em = data.criado_em ?? null;
    }

    /**
     * Valida os campos da empresa
     * @returns {string[]} Lista de erros
     */
    validar() {
        const erros = [];

        if (!this.nome) {
            erros.push("O nome da empresa é obrigatório.");
        }

        if (!this.cnpj) {
            erros.push("O CNPJ é obrigatório.");
        } else if (this.cnpj.length !== 14) {
            erros.push("O CNPJ deve conter exatamente 14 dígitos.");
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
     * Retorna o CNPJ formatado no padrão 00.000.000/0001-00
     * @returns {string}
     */
    getCnpjFormatado() {
        if (this.cnpj.length !== 14) return this.cnpj;
        return this.cnpj.replace(/^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$/, "$1.$2.$3/$4-$5");
    }
}