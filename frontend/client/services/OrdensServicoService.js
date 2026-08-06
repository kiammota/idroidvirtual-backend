import { HttpClient } from "../http.js";
import { OrdemServico } from "../models/OrdemServicoModel.js";

export class ServiceOrderService extends HttpClient {
  /**
   * Cria uma ordem de serviço
   * @param {Object|OrdemServico} params
   * @returns {Promise<{success: boolean, data: OrdemServico, message: string}>}
   */
  async create(params) {
    const body = {
      cliente_id: params.cliente_id,
      nome_servico: params.nome_servico || "",
      defeito: params.defeito,
      servicos: params.servicos,
      assinatura: params.assinatura || null,
      status: params.status || "ABERTA",
      valor: params.valor != null ? Number(params.valor) : 0
    };
    const response = await this.post('/service-orders', body);

    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }

    return response;
  }

  /**
   * Lista todas as ordens de serviço
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
   * @param {number} id
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
   * Atualiza uma ordem de serviço
   * @param {Object} params
   * @returns {Promise<{success: boolean, data: OrdemServico, message: string}>}
   */
  async update(params) {
    const id = params.id;
    const body = {};
    if (params.defeito !== undefined) body.defeito = params.defeito;
    if (params.servicos !== undefined) body.servicos = params.servicos;
    if (params.assinatura !== undefined) body.assinatura = params.assinatura;
    if (params.status !== undefined) body.status = params.status;
    if (params.valor !== undefined) body.valor = Number(params.valor);
    if (params.nome_servico !== undefined) body.nome_servico = params.nome_servico;
    if (params.cliente_id !== undefined) body.cliente_id = params.cliente_id;

    const response = await this.patch(`/service-orders/${id}`, body);

    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }

    return response;
  }

  /**
   * Deleta uma ordem de serviço
   * @param {number} id
   * @returns {Promise<{success: boolean, data: null, message: string}>}
   */
  async remove(id) {
    return this.delete(`/service-orders/${id}`);
  }

  /**
   * Alias para compatibilidade
   */
  async delete(id) {
    return this.remove(id);
  }

  /**
   * Lista ordens por cliente ID
   * @param {number} cliente_id
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
   * @param {string} client_name
   * @returns {Promise<{success: boolean, data: OrdemServico[], message: string}>}
   */
  async listByClientName(client_name) {
    const response = await this.get(`/service-orders/client-name/${encodeURIComponent(client_name)}`);

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new OrdemServico(data));
    }

    return response;
  }

  /**
   * Lista ordens por status (ex: EM_ANDAMENTO, FINALIZADA, CANCELADA)
   * @param {string} status
   * @returns {Promise<{success: boolean, data: OrdemServico[], message: string}>}
   */
  async listByStatus(status) {
    const response = await this.get(`/service-orders/status/${encodeURIComponent(status)}`);

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new OrdemServico(data));
    }

    return response;
  }
}