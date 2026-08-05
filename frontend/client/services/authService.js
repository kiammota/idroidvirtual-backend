import { HttpClient } from "../http.js";
import { LoginRequest } from "../models/LoginModel.js";

export class AuthService extends HttpClient {
  /**
   * Realiza login
   * 
   * @param {string} email - Email do usuário
   * @param {string} senha - Senha do usuário
   * @returns {Promise<{success: boolean, data: string, message: string}>}
   */
  async login(email, senha) {
    const model = new LoginRequest({
      email,
      senha
    });

    const response = await this.post('/auth/login', model);

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
   * @param {string} senha - Senha administrativa
   * @returns {Promise<{success: boolean, data: null, message: string}>}
   */
  async elevate(senha) {
    return this.post('/auth/elevate', { senha });
  }
}