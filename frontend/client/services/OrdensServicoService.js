import { HttpClient } from "../http.js";
import { OrdemServico } from "../models/OrdemServicoModel.js";

class ServiceOrderServiceClass extends HttpClient {
  async create(params) {
    const body = {
      cliente_id: params.cliente_id,
      defeito: params.defeito ?? null,
      servicos: params.servicos ?? null,
      assinatura: params.assinatura ?? null,
      status: params.status || "ABERTA",
    };
    if (params.telefone != null) body.telefone = params.telefone;
    if (params.nome_servico != null) body.nome_servico = params.nome_servico;
    if (params.valor != null) body.valor = Number(params.valor);

    const response = await this.post("/service-orders", body);
    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }
    return response;
  }

  async listAll() {
    const response = await this.get("/service-orders");
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async getById(id) {
    const response = await this.get(`/service-orders/${id}`);
    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }
    return response;
  }

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
    if (params.telefone !== undefined) body.telefone = params.telefone;

    const response = await this.patch(`/service-orders/${id}`, body);
    if (response.success && response.data) {
      response.data = new OrdemServico(response.data);
    }
    return response;
  }

  async remove(id) {
    return this.delete(`/service-orders/${id}`);
  }

  async deleteOrder(id) {
    return this.remove(id);
  }

  async listByClientId(cliente_id) {
    const response = await this.get(`/service-orders/client-id/${cliente_id}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listByClientName(client_name) {
    const response = await this.get(`/service-orders/client-name/${encodeURIComponent(client_name)}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listByClientEmail(client_email) {
    const response = await this.get(`/service-orders/client-email/${encodeURIComponent(client_email)}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listByClientNumber(client_number) {
    const response = await this.get(`/service-orders/client-number/${encodeURIComponent(client_number)}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listLastThreeMonths() {
    const response = await this.get("/service-orders/three-months");
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listLastSixMonths() {
    const response = await this.get("/service-orders/six-months");
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }

  async listByStatus(status) {
    const response = await this.get(`/service-orders/status/${encodeURIComponent(status)}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new OrdemServico(data));
    }
    return response;
  }
}

/** Singleton */
export const ServiceOrderService = new ServiceOrderServiceClass();
export default ServiceOrderService;
