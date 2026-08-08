export class Despesa {
  /**
   * @param {Object} data
   */
  constructor(data = {}) {
    this.id = Number(data.id || 0);
    this.empresa_id = Number(data.empresa_id || data.empresaId || 0);
    this.descricao = data.descricao || '';
    this.categoria = data.categoria || 'OUTROS';
    this.valor = Number(data.valor || 0);
    this.valor_pago = Number(data.valor_pago || data.valorPago || 0);
    this.forma_pagamento = data.forma_pagamento || data.formaPagamento || null;
    this.status = data.status || 'PENDENTE';
    
    this.data_vencimento = data.data_vencimento || data.dataVencimento ? new Date(data.data_vencimento || data.dataVencimento) : null;
    this.data_pagamento = data.data_pagamento || data.dataPagamento ? new Date(data.data_pagamento || data.dataPagamento) : null;
    
    this.comprovante_url = data.comprovante_url || data.comprovanteUrl || null;
    this.numero_documento = data.numero_documento || data.numeroDocumento || null;
    this.criado_em = data.criado_em || data.criadoEm ? new Date(data.criado_em || data.criadoEm) : new Date();
  }
}