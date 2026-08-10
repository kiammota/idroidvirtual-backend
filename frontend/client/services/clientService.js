import ENDPOINTS from "../endpoints.js";
import { HttpClient } from "../http.js";
import { Cliente } from "../models/ClienteModel.js";

class ClientServiceClass extends HttpClient {
  async create(params) {
    const { nome, telefone, email, endereco } = params;
    const response = await this.post(ENDPOINTS.CLIENTS.BASE, { nome, telefone, email, endereco });
    if (response.success && response.data) {
      response.data = new Cliente(response.data);
    }
    return response;
  }

  async listAll() {
    const response = await this.get(ENDPOINTS.CLIENTS.BASE);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Cliente(data));
    }
    return response;
  }

  async getById(id) {
    const response = await this.get(ENDPOINTS.CLIENTS.BY_ID(id));
    if (response.success && response.data) {
      response.data = new Cliente(response.data);
    }
    return response;
  }

  async searchByPhone(telefone) {
    const params = new URLSearchParams({ telefone }).toString();
    const response = await this.get(`${ENDPOINTS.CLIENTS.NUMBER}?${params}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Cliente(data));
    }
    return response;
  }

  async searchByAddress(endereco) {
    const response = await this.get(`/clients/address?endereco=${encodeURIComponent(endereco)}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Cliente(data));
    }
    return response;
  }

  async searchByName(nome) {
    const response = await this.get(`/clients/name?nome=${encodeURIComponent(nome)}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Cliente(data));
    }
    return response;
  }

  async searchByEmail(email) {
    const response = await this.get(`/clients/email?email=${encodeURIComponent(email)}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Cliente(data));
    }
    return response;
  }
}

/** Singleton */
export const ClientService = new ClientServiceClass();
export default ClientService;
