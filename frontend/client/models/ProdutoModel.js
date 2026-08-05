export class Produto {
    constructor(data = {}) {
        this.id = Number(data.id) || 0;
        this.nome = String(data.nome ?? "").trim().toUpperCase(); // Força o nome a ficar inteiramente em MAIÚSCULAS
        this.preco = Number(data.preco) || 0;
        this.estoque = Number(data.estoque) || 0;
        this.categoria = String(data.categoria ?? "").trim().toUpperCase();
        this.criado_em = data.criado_em ?? null;
        this.atualizado_em = data.atualizado_em ?? null;
    }

    /**
     * Valida as regras de negócio do produto
     * @returns {string[]} Lista de mensagens de erro
     */
    validar() {
        const erros = [];

        if (!this.nome) {
            erros.push("O nome do produto é obrigatório.");
        }

        if (this.preco <= 0) {
            erros.push("O preço do produto deve ser maior que zero.");
        }

        if (this.estoque < 0) {
            erros.push("O estoque não pode ser um número negativo.");
        }

        if (!this.categoria) {
            erros.push("A categoria do produto é obrigatória.");
        }

        return erros;
    }

    /**
     * Verifica se o produto está com os dados válidos
     * @returns {boolean}
     */
    isValid() {
        return this.validar().length === 0;
    }

    /**
     * Verifica se há estoque disponível do produto
     * @param {number} quantidade
     * @returns {boolean}
     */
    temEstoque(quantidade = 1) {
        return this.estoque >= quantidade;
    }

    /**
     * Retorna o preço formatado em moeda (ex: R$ 12,50)
     * @returns {string}
     */
    getPrecoFormatado() {
        return this.preco.toLocaleString("pt-BR", {
            style: "currency",
            currency: "BRL"
        });
    }
}