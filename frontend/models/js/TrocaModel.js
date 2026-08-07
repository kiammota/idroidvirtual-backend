/**
 * Modelo de Troca (Trade-in)
 * Espelha o struct Go: models.Troca
 */
export class Troca {
  constructor(data = {}) {
    this.id = Number(data.id) || 0;
    this.venda_id = Number(data.venda_id) || 0;
    this.modelo_aparelho = String(data.modelo_aparelho ?? "").trim();
    this.imei = data.imei ? String(data.imei).replace(/\D/g, "") : "";
    this.valor_avaliacao = Number(data.valor_avaliacao) || 0;
    this.condicao = data.condicao ? String(data.condicao).trim() : "";
    this.observacoes = data.observacoes != null ? String(data.observacoes) : null;
    this.criado_em = data.criado_em ?? null;
  }

  /**
   * Valida os campos obrigatórios da troca
   * @returns {string[]}
   */
  validar() {
    const erros = [];

    if (!this.venda_id || this.venda_id <= 0) {
      erros.push("O ID da venda é obrigatório.");
    }

    if (!this.modelo_aparelho) {
      erros.push("O modelo do aparelho é obrigatório.");
    }

    if (!this.imei || this.imei.length < 14 || this.imei.length > 15) {
      erros.push("O IMEI deve ter 14 ou 15 dígitos.");
    }

    if (this.valor_avaliacao < 0) {
      erros.push("O valor de avaliação não pode ser negativo.");
    }

    return erros;
  }

  /**
   * @returns {boolean}
   */
  isValid() {
    return this.validar().length === 0;
  }

  /**
   * Payload para POST /trocas
   * @returns {Object}
   */
  toCreatePayload() {
    const payload = {
      venda_id: this.venda_id,
      modelo_aparelho: this.modelo_aparelho,
      imei: this.imei,
      valor_avaliacao: this.valor_avaliacao,
    };
    if (this.condicao) payload.condicao = this.condicao;
    if (this.observacoes != null) payload.observacoes = this.observacoes;
    return payload;
  }
}
