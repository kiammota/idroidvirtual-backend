import { HttpClient } from "../http.js";
import { Cliente } from "../models/ClienteModel.js";

class ClientServiceClass extends HttpClient {
  async create(params) {
    const { nome, telefone, email, endereco } = params;
    const response = await this.post("/clients/", { nome, telefone, email, endereco });
    if (response.success && response.data) {
      response.data = new Cliente(response.data);
    }
    return response;
  }

  async listAll() {
    const response = await this.get("/clients/");
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Cliente(data));
    }
    return response;
  }

  async getById(id) {
    const response = await this.get(`/clients/${id}`);
    if (response.success && response.data) {
      response.data = new Cliente(response.data);
    }
    return response;
  }

  async searchByPhone(telefone) {
    const response = await this.get(`/clients/number?telefone=${encodeURIComponent(telefone)}`);
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
