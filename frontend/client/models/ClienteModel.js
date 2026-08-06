export class Cliente {
    constructor(data = {}) {
        this.id = Number(data.id) || 0;
        this.nome = String(data.nome ?? "").trim().toUpperCase();
        this.telefone = data.telefone ? String(data.telefone).replace(/\D/g, "") : null; // Remove caracteres não numéricos
        this.endereco = data.endereco ? String(data.endereco).trim().toUpperCase() : null;
        this.email = data.email ? String(data.email).trim().toLowerCase() : null;
        this.criado_em = data.criado_em ?? null;
        this.atualizado_em = data.atualizado_em ?? null;
    }

    /**
     * Valida os campos do cliente
     * @returns {string[]} Lista com mensagens de erro
     */
    validar() {
        const erros = [];


        if (this.email) {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(this.email)) {
                erros.push("O e-mail informado não é válido.");
            }
        }

        if (this.telefone && (this.telefone.length < 10 || this.telefone.length > 11)) {
            erros.push("O telefone deve ter entre 10 e 11 dígitos (DDD + número).");
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
     * Retorna o telefone formatado como (XX) XXXXX-XXXX ou (XX) XXXX-XXXX
     * @returns {string|null}
     */
    getTelefoneFormatado() {
        if (!this.telefone) return null;

        if (this.telefone.length === 11) {
            return this.telefone.replace(/^(\d{2})(\d{5})(\d{4})$/, "($1) $2-$3");
        }

        if (this.telefone.length === 10) {
            return this.telefone.replace(/^(\d{2})(\d{4})(\d{4})$/, "($1) $2-$3");
        }

        return this.telefone;
    }
}