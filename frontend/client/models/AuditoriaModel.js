export class Auditoria {
  /**
   * @param {Object} data
   * @param {number} [data.id]
   * @param {number} [data.usuario_id]
   * @param {string} [data.acao]
   * @param {string} [data.tabela]
   * @param {number} [data.registro_id]
   * @param {Object} [data.valores_antigos]
   * @param {Object} [data.valores_novos]
   * @param {string|Date} [data.horario]
   */
  constructor(data = {}) {
    this.id = data.id || null;
    this.usuario_id = data.usuario_id || null;
    this.acao = data.acao || '';
    this.tabela = data.tabela || '';
    this.registro_id = data.registro_id || null;
    this.valores_antigos = data.valores_antigos || {};
    this.valores_novos = data.valores_novos || {};
    this.horario = data.horario ? new Date(data.horario) : null;
  }

  /**
   * Formata a data e hora do registro para o padrão PT-BR
   * @returns {string} Ex: "08/08/2026 14:30:00"
   */
  getHorarioFormatado() {
    if (!this.horario || isNaN(this.horario.getTime())) {
      return '--/--/---- --:--';
    }
    return this.horario.toLocaleString('pt-BR');
  }
}