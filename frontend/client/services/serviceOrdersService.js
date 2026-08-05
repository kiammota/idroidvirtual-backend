import { HttpClient } from "../http.js";

export class ServiceOrderService extends HttpClient {
  /**
   * Cria uma ordem de serviço
   * 
   * | Parâmetro    | Tipo    | Descrição         |
   * |--------------|---------|-------------------|
   * | cliente_id   | integer | ID do cliente     |
   * | descricao    | string  | Descrição do serviço|
   * | status       | string  | Status            |
   * 
   * @param {Object} params
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async create(params) {
    const { cliente_id, descricao, status } = params;
    return this.post('/service-orders', { cliente_id, descricao, status });
  }

  /**
   * Lista todas as ordens de serviço
   * 
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async listAll() {
    return this.get('/service-orders');
  }

  /**
   * Busca ordem de serviço por ID
   * 
   * | Parâmetro | Tipo    | Descrição     |
   * |-----------|---------|---------------|
   * | id        | integer | ID da ordem   |
   * 
   * @param {number} id - ID da ordem de serviço
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async getById(id) {
    return this.get(`/service-orders/${id}`);
  }

  /**
   * Atualiza parcialmente uma ordem de serviço
   * 
   * | Parâmetro | Tipo    | Descrição     |
   * |-----------|---------|---------------|
   * | id        | integer | ID da ordem   |
   * | descricao | string  | Descrição     |
   * | status    | string  | Status        |
   * 
   * @param {Object} params
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async update(params) {
    const { id, descricao, status } = params;
    return this.patch(`/service-orders/${id}`, { descricao, status });
  }

  /**
   * Deleta uma ordem de serviço
   * 
   * | Parâmetro | Tipo    | Descrição     |
   * |-----------|---------|---------------|
   * | id        | integer | ID da ordem   |
   * 
   * @param {number} id - ID da ordem de serviço
   * @returns {Promise<{success: boolean, data: null, message: string}>}
   */
  async delete(id) {
    return this.delete(`/service-orders/${id}`);
  }

  /**
   * Lista ordens por cliente ID
   * 
   * | Parâmetro  | Tipo    | Descrição      |
   * |------------|---------|----------------|
   * | cliente_id | integer | ID do cliente  |
   * 
   * @param {number} cliente_id - ID do cliente
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async listByClientId(cliente_id) {
    return this.get(`/service-orders/client-id/${cliente_id}`);
  }

  /**
   * Lista ordens por nome do cliente
   * 
   * | Parâmetro   | Tipo   | Descrição       |
   * |-------------|--------|-----------------|
   * | client_name | string | Nome do cliente |
   * 
   * @param {string} client_name - Nome do cliente
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async listByClientName(client_name) {
    return this.get(`/service-orders/client-name/${encodeURIComponent(client_name)}`);
  }
}
