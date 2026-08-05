export class Usuario {
    constructor(data = {}) {
        this.id = Number(data.id) || 0;
        this.empresa_id = Number(data.empresa_id) || 0;
        this.email = String(data.email ?? "").trim().toLowerCase();
        this.nome = String(data.nome ?? "").trim();
        this.senha = data.senha ?? "";
        this.criado_em = data.criado_em ?? null;
    }

    /**
     * Valida os campos do usuário
     * @returns {string[]} Lista com os erros encontrados
     */
    validar() {
        const erros = [];

        if (!this.nome) {
            erros.push("O nome é obrigatório.");
        }

        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!this.email || !emailRegex.test(this.email)) {
            erros.push("Informe um e-mail válido.");
        }

        if (this.senha && this.senha.length < 6) {
            erros.push("A senha deve ter pelo menos 6 caracteres.");
        }

        return erros;
    }

    /**
     * Verifica se o objeto é válido
     * @returns {boolean}
     */
    isValid() {
        return this.validar().length === 0;
    }

    /**
     * Retorna o usuário em formato JSON sem dados sensíveis (ex: senha)
     * @returns {Object}
     */
    toJSON() {
        const { senha, ...usuarioSemSenha } = this;
        return usuarioSemSenha;
    }
}