import { HttpClient } from "../http.js";
import { Usuario } from "../models/UsuarioModel.js"; // Ajuste o caminho de importação conforme a estrutura do seu projeto

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
   * @returns {Promise<{success: boolean, data: Usuario, message: string}>}
   */
  async create(params) {
    const { nome, email, senha } = params;
    const response = await this.post('/users', { nome, email, senha });

    if (response.success && response.data) {
      response.data = new Usuario(response.data);
    }

    return response;
  }

  /**
   * Busca usuário por ID
   * 
   * | Parâmetro | Tipo    | Descrição   |
   * |-----------|---------|-------------|
   * | id        | integer | ID do usuário|
   * 
   * @param {number} id - ID do usuário
   * @returns {Promise<{success: boolean, data: Usuario, message: string}>}
   */
  async getById(id) {
    const response = await this.get(`/users/${id}`);

    if (response.success && response.data) {
      response.data = new Usuario(response.data);
    }

    return response;
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