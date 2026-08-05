import { HttpClient } from "../http.js";
import { Cliente } from "../models/ClienteModel.js"; // Importe a classe no caminho correto do seu projeto

export class ClientService extends HttpClient {
  /**
   * Cria um novo cliente
   * 
   * | Parâmetro | Tipo   | Descrição    |
   * |-----------|--------|--------------|
   * | nome      | string | Nome         |
   * | telefone  | string | Telefone     |
   * | email     | string | Email        |
   * | endereco  | string | Endereço     |
   * 
   * @param {Object} params
   * @returns {Promise<{success: boolean, data: Cliente, message: string}>}
   */
  async create(params) {
    const { nome, telefone, email, endereco } = params;
    const response = await this.post('/clients/', { nome, telefone, email, endereco });
    
    if (response.success && response.data) {
      response.data = new Cliente(response.data);
    }
    
    return response;
  }

  /**
   * Lista todos os clientes
   * 
   * @returns {Promise<{success: boolean, data: Cliente[], message: string}>}
   */
  async listAll() {
    const response = await this.get('/clients/');
    
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new Cliente(data));
    }
    
    return response;
  }

  /**
   * Busca cliente por ID
   * 
   * | Parâmetro | Tipo    | Descrição    |
   * |-----------|---------|--------------|
   * | id        | integer | ID do cliente|
   * 
   * @param {number} id - ID do cliente
   * @returns {Promise<{success: boolean, data: Cliente, message: string}>}
   */
  async getById(id) {
    const response = await this.get(`/clients/${id}`);
    
    if (response.success && response.data) {
      response.data = new Cliente(response.data);
    }
    
    return response;
  }

  /**
   * Busca clientes por telefone
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | telefone  | string | Telefone  |
   * 
   * @param {string} telefone - Número de telefone
   * @returns {Promise<{success: boolean, data: Cliente[], message: string}>}
   */
  async searchByPhone(telefone) {
    const response = await this.get(`/clients/number?telefone=${encodeURIComponent(telefone)}`);
    
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new Cliente(data));
    }
    
    return response;
  }

  /**
   * Busca clientes por endereço
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | endereco  | string | Endereço  |
   * 
   * @param {string} endereco - Endereço
   * @returns {Promise<{success: boolean, data: Cliente[], message: string}>}
   */
  async searchByAddress(endereco) {
    const response = await this.get(`/clients/address?endereco=${encodeURIComponent(endereco)}`);
    
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new Cliente(data));
    }
    
    return response;
  }

  /**
   * Busca clientes por nome
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | nome      | string | Nome      |
   * 
   * @param {string} nome - Nome do cliente
   * @returns {Promise<{success: boolean, data: Cliente[], message: string}>}
   */
  async searchByName(nome) {
    const response = await this.get(`/clients/name?nome=${encodeURIComponent(nome)}`);
    
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new Cliente(data));
    }
    
    return response;
  }

  /**
   * Busca clientes por email
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | email     | string | Email     |
   * 
   * @param {string} email - Email do cliente
   * @returns {Promise<{success: boolean, data: Cliente[], message: string}>}
   */
  async searchByEmail(email) {
    const response = await this.get(`/clients/email?email=${encodeURIComponent(email)}`);
    
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new Cliente(data));
    }
    
    return response;
  }
}