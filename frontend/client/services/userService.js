import { HttpClient } from "../http.js";

export class UserService extends HttpClient {
  /**
   * Cria um novo usuário
   * 
   * | Parâmetro | Tipo   | Descrição    |
   * |-----------|--------|--------------|
   * | nome      | string | Nome completo|
   * | email     | string | Email único  |
   * | senha     | string | Senha        |
   * 
   * @param {Object} params
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async create(params) {
    const { nome, email, senha } = params;
    return this.post('/users', { nome, email, senha });
  }

  /**
   * Busca usuário por ID
   * 
   * | Parâmetro | Tipo    | Descrição   |
   * |-----------|---------|-------------|
   * | id        | integer | ID do usuário|
   * 
   * @param {number} id - ID do usuário
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async getById(id) {
    return this.get(`/users/${id}`);
  }

  /**
   * Valida senha administrativa do usuário
   * 
   * | Parâmetro | Tipo   | Descrição      |
   * |-----------|--------|----------------|
   * | senha     | string | Senha admin    |
   * 
   * @param {string} senha - Senha administrativa
   * @returns {Promise<{success: boolean, data: null, message: string}>}
   */
  async elevate(senha) {
    return this.post('/users/elevate', { senha });
  }
}
