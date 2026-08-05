
export class HttpClient {
  /**
   * @param {string} baseUrl - URL base da API
   */
  constructor(baseUrl = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  /**
   * Obtém o token JWT do sessionStorage
   * @returns {string|null}
   */
  getToken() {
    return sessionStorage.getItem(TOKEN_STORAGE_KEY);
  }

  /**
   * Armazena o token JWT no sessionStorage
   * @param {string} token - Token JWT
   */
  setToken(token) {
    sessionStorage.setItem(TOKEN_STORAGE_KEY, token);
  }

  /**
   * Remove o token JWT do sessionStorage
   */
  clearToken() {
    sessionStorage.removeItem(TOKEN_STORAGE_KEY);
  }

  /**
   * Cria headers padrão com Authorization
   * @returns {Object}
   */
  getHeaders() {
    const headers = {
      'Content-Type': 'application/json',
    };

    const token = this.getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    return headers;
  }

  /**
   * Realiza uma requisição HTTP
   * @param {string} method - GET, POST, DELETE, PATCH, etc.
   * @param {string} endpoint - Caminho do endpoint (sem baseUrl)
   * @param {Object} body - Corpo da requisição (opcional)
   * @returns {Promise<{success: boolean, data: *, message: string}>}
   */
  async request(method, endpoint, body = null) {
    try {
      const url = `${this.baseUrl}${endpoint}`;
      const options = {
        method,
        headers: this.getHeaders(),
      };

      if (body) {
        options.body = JSON.stringify(body);
      }

      const response = await fetch(url, options);

      // Trata erro 401 (Unauthorized)
      if (response.status === 401) {
        this.clearToken();
        NotificationService.error(
          'Sua sessão expirou. Por favor, faça login novamente.',
          'Autenticação'
        );
        // Você pode redirecionar para login aqui
        // window.location.href = '/login';
        return { success: false, data: null, message: 'Não autenticado' };
      }

      const responseData = await response.json();

      // Se a resposta não for bem-sucedida e não for um envelope válido
      if (!response.ok && !responseData.success) {
        NotificationService.error(
          responseData.message || 'Erro na requisição',
          `Erro ${response.status}`
        );
      }

      return responseData;
    } catch (error) {
      NotificationService.error(
        error.message || 'Erro de conexão com o servidor',
        'Erro de Rede'
      );
      return { success: false, data: null, message: error.message };
    }
  }

  /**
   * GET
   */
  async get(endpoint) {
    return this.request('GET', endpoint);
  }

  /**
   * POST
   */
  async post(endpoint, body) {
    return this.request('POST', endpoint, body);
  }

  /**
   * DELETE
   */
  async delete(endpoint) {
    return this.request('DELETE', endpoint);
  }

  /**
   * PATCH
   */
  async patch(endpoint, body) {
    return this.request('PATCH', endpoint, body);
  }
}
