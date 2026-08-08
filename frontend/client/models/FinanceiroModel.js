export class DetalhePagamento {
  /**
   * @param {Object} data
   * @param {string} [data.formaPagamento] - PIX, CARTAO_CREDITO, DINHEIRO
   * @param {number} [data.total]
   */
  constructor({ formaPagamento = '', total = 0.0 } = {}) {
    this.forma_pagamento = formaPagamento;
    this.total = Number(total);
  }
}

export class ResumoFinanceiro {
  /**
   * @param {Object} data
   * @param {Date|string} [data.periodoInicio]
   * @param {Date|string} [data.periodoFim]
   * @param {number} [data.totalVendas]
   * @param {number} [data.quantidadeVendas]
   * @param {number} [data.ticketMedioVenda]
   * @param {number} [data.totalServicos]
   * @param {number} [data.quantidadeOS]
   * @param {number} [data.ticketMedioServico]
   * @param {number} [data.faturamentoBruto]
   * @param {Array<DetalhePagamento>} [data.formasPagamento]
   * @param {number} [data.totalDespesas]
   * @param {number} [data.valorEstoquePerdido]
   * @param {number} [data.quantidadePerdas]
   * @param {number} [data.lucroLiquido]
   * @param {number} [data.margemLucro]
   */
  constructor(data = {}) {
    this.periodo_inicio = data.periodoInicio ? new Date(data.periodoInicio) : null;
    this.periodo_fim = data.periodoFim ? new Date(data.periodoFim) : null;

    // Vendas de Produtos
    this.total_vendas = Number(data.totalVendas || 0);
    this.quantidade_vendas = Number(data.quantidadeVendas || 0);
    this.ticket_medio_venda = Number(data.ticketMedioVenda || 0);

    // Ordens de Serviço
    this.total_servicos = Number(data.totalServicos || 0);
    this.quantidade_os = Number(data.quantidadeOS || 0);
    this.ticket_medio_servico = Number(data.ticketMedioServico || 0);

    // Faturamento e Formas de Pagamento
    this.faturamento_bruto = Number(data.faturamentoBruto || 0);
    this.formas_pagamento = Array.isArray(data.formasPagamento)
      ? data.formasPagamento.map(fp => new DetalhePagamento(fp))
      : [];

    // Custos e Deduções
    this.total_despesas = Number(data.totalDespesas || 0);
    this.valor_estoque_perdido = Number(data.valorEstoquePerdido || 0);
    this.quantidade_perdas = Number(data.quantidadePerdas || 0);

    // Indicadores
    this.lucro_liquido = Number(data.lucroLiquido || 0);
    this.margem_lucro = Number(data.margemLucro || 0); // Em porcentagem (%)
  }
}