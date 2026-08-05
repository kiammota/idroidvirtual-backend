
import { HttpClient } from "../http.js"

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
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async create(params) {
    const { nome, telefone, email, endereco } = params;
    return this.post('/clients/', { nome, telefone, email, endereco });
  }

  /**
   * Lista todos os clientes
   * 
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async listAll() {
    return this.get('/clients/');
  }

  /**
   * Busca cliente por ID
   * 
   * | Parâmetro | Tipo    | Descrição    |
   * |-----------|---------|--------------|
   * | id        | integer | ID do cliente|
   * 
   * @param {number} id - ID do cliente
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async getById(id) {
    return this.get(`/clients/${id}`);
  }

  /**
   * Busca clientes por telefone
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | telefone  | string | Telefone  |
   * 
   * @param {string} telefone - Número de telefone
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async searchByPhone(telefone) {
    return this.get(`/clients/number?telefone=${encodeURIComponent(telefone)}`);
  }

  /**
   * Busca clientes por endereço
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | endereco  | string | Endereço  |
   * 
   * @param {string} endereco - Endereço
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async searchByAddress(endereco) {
    return this.get(`/clients/address?endereco=${encodeURIComponent(endereco)}`);
  }

  /**
   * Busca clientes por nome
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | nome      | string | Nome      |
   * 
   * @param {string} nome - Nome do cliente
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async searchByName(nome) {
    return this.get(`/clients/name?nome=${encodeURIComponent(nome)}`);
  }

  /**
   * Busca clientes por email
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | email     | string | Email     |
   * 
   * @param {string} email - Email do cliente
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async searchByEmail(email) {
    return this.get(`/clients/email?email=${encodeURIComponent(email)}`);
  }
}
