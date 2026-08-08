import ENDPOINTS from "../endpoints.js";
import { HttpClient } from "../http.js";
import { Despesa } from "../models/DespesaModel.js";

class DespesaServiceClass extends HttpClient {

  /**
   * Lista todas as despesas
   */
  async listar() {
    const response = await this.get(ENDPOINTS.EXPENSES.BASE);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(item => new Despesa(item));
    }
    return response;
  }

  /**
   * Cria um novo registro de despesa incluindo o empresa_id do token
   * @param {Object} payload
   */
  async criar(payload) {
    const empresaId = this.getEmpresaId();

    const body = {
      ...payload,
      empresa_id: empresaId || payload.empresa_id || payload.empresaId
    };

    const response = await this.post(ENDPOINTS.EXPENSES.BASE, body);
    if (response.success && response.data) {
      response.data = new Despesa(response.data);
    }
    return response;
  }

  /**
   * Atualiza uma despesa existente
   * @param {number|string} id
   * @param {Object} payload
   */
  async atualizar(id, payload) {
    const response = await this.put(`${ENDPOINTS.EXPENSES.BASE}/${id}`, payload);
    if (response.success && response.data) {
      response.data = new Despesa(response.data);
    }
    return response;
  }

  /**
   * Remove um lembrete/despesa
   * @param {number|string} id
   */
  async deletar(id) {
    return await this.delete(`${ENDPOINTS.EXPENSES.BASE}/${id}`);
  }
}

// Instancia a classe alterada e exporta
export const DespesaService = new DespesaServiceClass();
export default DespesaService;