export class ResumoFinanceiro {
    constructor(data = {}) {
        this.periodo_inicio = data.periodo_inicio ?? null;
        this.periodo_fim = data.periodo_fim ?? null;

        // Vendas de Produtos
        this.total_vendas = data.total_vendas ?? 0;
        this.quantidade_vendas = data.quantidade_vendas ?? 0;

        // Ordens de Serviço
        this.total_servicos = data.total_servicos ?? 0;
        this.quantidade_os = data.quantidade_os ?? 0;

        // Faturamento Bruto
        this.faturamento_bruto = data.faturamento_bruto ?? 0;

        // Custos e Perdas
        this.total_despesas = data.total_despesas ?? 0;
        this.valor_estoque_perdido = data.valor_estoque_perdido ?? 0;
        this.quantidade_perdas = data.quantidade_perdas ?? 0;

        // Resultado Financeiro Real
        this.lucro_liquido = data.lucro_liquido ?? 0;
    }
}