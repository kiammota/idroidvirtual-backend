export class LoginRequest {
    constructor(data = {}) {
        this.email = String(data.email ?? "").trim().toLowerCase();
        this.senha = String(data.senha ?? "");
    }

    /**
     * Valida os campos do request de login
     * @returns {string[]} Lista de erros de validação
     */
    validar() {
        const erros = [];

        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!this.email || !emailRegex.test(this.email)) {
            erros.push("Informe um e-mail válido.");
        }

        if (!this.senha) {
            erros.push("A senha é obrigatória.");
        }

        return erros;
    }

    /**
     * Verifica se os dados de login são válidos antes do envio
     * @returns {boolean}
     */
    isValid() {
        return this.validar().length === 0;
    }
}