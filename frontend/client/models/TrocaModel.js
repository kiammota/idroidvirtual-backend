/**
 * Modelo de Troca (Trade-in)
 * Espelha o struct Go: models.Troca
 *
 * Troca ≠ Venda — venda_id é opcional (pode vir 0 / ausente).
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
   * Valida campos obrigatórios.
   * - modelo_aparelho: obrigatório
   * - valor_avaliacao: >= 0
   * - imei: opcional (se preenchido, 14–15 dígitos)
   * - venda_id: opcional (troca independente de venda)
   * @returns {string[]}
   */
  validar() {
    const erros = [];

    if (!this.modelo_aparelho) {
      erros.push("O modelo do aparelho é obrigatório.");
    }

    if (this.imei && (this.imei.length < 14 || this.imei.length > 15)) {
      erros.push("O IMEI deve ter 14 ou 15 dígitos.");
    }

    if (this.valor_avaliacao < 0) {
      erros.push("O valor de avaliação não pode ser negativo.");
    }

    return erros;
  }

  isValid() {
    return this.validar().length === 0;
  }

  /**
   * Payload para POST /api/v1/trocas
   * Não envia venda_id se for 0 (troca sem venda vinculada).
   */
  toCreatePayload() {
    const payload = {
      modelo_aparelho: this.modelo_aparelho,
      valor_avaliacao: this.valor_avaliacao,
    };
    if (this.venda_id > 0) payload.venda_id = this.venda_id;
    if (this.imei) payload.imei = this.imei;
    if (this.condicao) payload.condicao = this.condicao;
    if (this.observacoes != null && this.observacoes !== "") {
      payload.observacoes = this.observacoes;
    }
    return payload;
  }
}
