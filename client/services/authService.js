import { HttpClient } from "../http.js"

export class AuthService extends HttpClient {
  /**
   * Realiza login
   * 
   * | Parâmetro | Tipo   | Descrição      |
   * |-----------|--------|----------------|
   * | login     | string | Email/usuário  |
   * | senha     | string | Senha          |
   * 
   * @param {string} login - Email ou usuário
   * @param {string} senha - Senha
   * @returns {Promise<{success: boolean, data: string, message: string}>}
   */
  async login(login, senha) {
    const response = await this.post('/auth/login', { login, senha });

    if (response.success && response.data) {
      this.setToken(response.data);
      NotificationService.success('Login realizado com sucesso');
    }

    return response;
  }

  /**
   * Faz logout
   */
  logout() {
    this.clearToken();
    NotificationService.success('Logout realizado');
  }

  /**
   * Valida senha administrativa
   * 
   * | Parâmetro | Tipo   | Descrição      |
   * |-----------|--------|----------------|
   * | senha     | string | Senha admin    |
   * 
   * @param {string} senha - Senha administrativa
   * @returns {Promise<{success: boolean, data: null, message: string}>}
   */
  async elevate(senha) {
    return this.post('/auth/elevate', { senha });
  }
}
