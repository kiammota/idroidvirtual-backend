import { HttpClient } from "../http.js";
import { LoginRequest } from "../models/LoginModel.js";

class AuthServiceClass extends HttpClient {
  /**
   * Realiza login
   * @param {string} email
   * @param {string} senha
   * @returns {Promise<{success: boolean, data: string, message: string}>}
   */
  async login(email, senha) {
    const model = new LoginRequest({ email, senha });
    const response = await this.post("/auth/login", model);

    if (response.success && response.data) {
      this.setToken(response.data);
      if (typeof NotificationService !== "undefined") {
        NotificationService.success("Login realizado com sucesso");
      }
    }

    return response;
  }

  /**
   * Faz logout
   */
  logout() {
    this.clearToken();
    if (typeof NotificationService !== "undefined") {
      NotificationService.success("Logout realizado");
    }
  }

  /**
   * Valida senha administrativa (OpenAPI: POST /users/elevate)
   * @param {string} senha
   * @returns {Promise<{success: boolean, data: null, message: string}>}
   */
  async elevate(senha) {
    return this.post("/users/elevate", { senha });
  }
}

/** Singleton */
export const AuthService = new AuthServiceClass();
export default AuthService;
