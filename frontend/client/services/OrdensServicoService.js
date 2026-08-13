import ENDPOINTS from "../endpoints.js";
import { HttpClient } from "../http.js";
import { OrdemServico } from "../models/OrdemServicoModel.js";

class ServiceOrderServiceClass extends HttpClient {
  
  async create(params) {
    const usuarioId = params.usuario_id || this.getUsuarioId();

    const body = {
      cliente_id: params.cliente_id,
      usuario_id: usuarioId ? Number(usuarioId) : null, // <-- Preenchido via JWT caso não passe no params
      defeito: params.defeito ?? null,
      servicos: params.servicos ?? null,
      assinatura: params.assinatura ?? null,
      status: params.status || "ABERTA",
    };

    if (params.empresa_id != null) body.empresa_id = Number(params.empresa_id);
    if (params.telefone != null) body.cliente_telefone = params.telefone;
    if (params.nome_servico != null) body.nome_servico = params.nome_servico;
    if (params.valor != null) body.valor = Number(params.valor);
    if (params.custo != null) body.custo = Number(params.custo);

    const response = await this.post(ENDPOINTS.SERVICE_ORDERS.BASE, body);
    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }
    return response;
  }

  async listAll() {
    const response = await this.get(ENDPOINTS.SERVICE_ORDERS.BASE);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async getById(id) {
    const response = await this.get(ENDPOINTS.SERVICE_ORDERS.BY_ID(id));
    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }
    return response;
  }

  async update(params) {
    const id = params.id;
    const body = {};

    if (params.usuario_id !== undefined) body.usuario_id = Number(params.usuario_id);
    if (params.empresa_id !== undefined) body.empresa_id = Number(params.empresa_id);
    if (params.defeito !== undefined) body.defeito = params.defeito;
    if (params.servicos !== undefined) body.servicos = params.servicos;
    if (params.assinatura !== undefined) body.assinatura = params.assinatura;
    if (params.status !== undefined) body.status = params.status;
    if (params.valor !== undefined) body.valor = Number(params.valor);
    if (params.custo !== undefined) body.custo = Number(params.custo);
    if (params.nome_servico !== undefined) body.nome_servico = params.nome_servico;
    if (params.cliente_id !== undefined) body.cliente_id = params.cliente_id;
    if (params.telefone !== undefined) body.telefone = params.telefone;

    const response = await this.patch(ENDPOINTS.SERVICE_ORDERS.BY_ID(id), body);
    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }
    return response;
  }

  async remove(id) {
    return this.delete(ENDPOINTS.SERVICE_ORDERS.BY_ID(id));
  }

  async deleteOrder(id) {
    return this.remove(id);
  }

  async listByClientId(cliente_id) {
    const response = await this.get(ENDPOINTS.SERVICE_ORDERS.BY_CLIENT_ID(cliente_id));
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listByClientName(client_name) {
    const response = await this.get(ENDPOINTS.SERVICE_ORDERS.BY_CLIENT_NAME(client_name));
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listByClientEmail(client_email) {
    const response = await this.get(ENDPOINTS.SERVICE_ORDERS.BY_CLIENT_EMAIL(client_email));
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listByClientNumber(client_number) {
    const response = await this.get(ENDPOINTS.SERVICE_ORDERS.BY_CLIENT_NUMBER(client_number));
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listLastThreeMonths() {
    const response = await this.get(ENDPOINTS.SERVICE_ORDERS.THREE_MONTHS);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listLastSixMonths() {
    const response = await this.get(ENDPOINTS.SERVICE_ORDERS.SIX_MONTHS);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listByStatus(status) {
    const response = await this.get(ENDPOINTS.SERVICE_ORDERS.BY_STATUS(status));
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }
}

/** Singleton */
export const ServiceOrderService = new ServiceOrderServiceClass();
export default ServiceOrderService;