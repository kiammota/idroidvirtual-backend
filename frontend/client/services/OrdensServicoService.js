import { HttpClient } from "../http.js";
import { OrdemServico } from "../models/OrdemServicoModel.js"; // Ajuste o caminho conforme a estrutura do projeto

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
   * @returns {Promise<{success: boolean, data: OrdemServico, message: string}>}
   */
  async create(params) {
    const { cliente_id, descricao, status } = params;
    const response = await this.post('/service-orders', { cliente_id, descricao, status });

    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }

    return response;
  }

  /**
   * Lista todas as ordens de serviço
   * 
   * @returns {Promise<{success: boolean, data: OrdemServico[], message: string}>}
   */
  async listAll() {
    const response = await this.get('/service-orders');

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new OrdemServico(data));
    }

    return response;
  }

  /**
   * Busca ordem de serviço por ID
   * 
   * | Parâmetro | Tipo    | Descrição     |
   * |-----------|---------|---------------|
   * | id        | integer | ID da ordem   |
   * 
   * @param {number} id - ID da ordem de serviço
   * @returns {Promise<{success: boolean, data: OrdemServico, message: string}>}
   */
  async getById(id) {
    const response = await this.get(`/service-orders/${id}`);

    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }

    return response;
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
   * @returns {Promise<{success: boolean, data: OrdemServico, message: string}>}
   */
  async update(params) {
    const { id, descricao, status } = params;
    const response = await this.patch(`/service-orders/${id}`, { descricao, status });

    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }

    return response;
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
   * @returns {Promise<{success: boolean, data: OrdemServico[], message: string}>}
   */
  async listByClientId(cliente_id) {
    const response = await this.get(`/service-orders/client-id/${cliente_id}`);

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new OrdemServico(data));
    }

    return response;
  }

  /**
   * Lista ordens por nome do cliente
   * 
   * | Parâmetro   | Tipo   | Descrição       |
   * |-------------|--------|-----------------|
   * | client_name | string | Nome do cliente |
   * 
   * @param {string} client_name - Nome do cliente
   * @returns {Promise<{success: boolean, data: OrdemServico[], message: string}>}
   */
  async listByClientName(client_name) {
    const response = await this.get(`/service-orders/client-name/${encodeURIComponent(client_name)}`);

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new OrdemServico(data));
    }

    return response;
  }
}