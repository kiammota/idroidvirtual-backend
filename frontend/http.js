import { API_BASE_URL, TOKEN_STORAGE_KEY } from './env.js';

export class HttpClient {
  /**
   * @param {string} baseUrl - URL base da API
   */
  constructor(baseUrl = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  getToken() {
    return sessionStorage.getItem(TOKEN_STORAGE_KEY);
  }

  setToken(token) {
    sessionStorage.setItem(TOKEN_STORAGE_KEY, token);
  }

  clearToken() {
    sessionStorage.removeItem(TOKEN_STORAGE_KEY);
  }

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

      if (response.status === 401) {
        this.clearToken();
        if (typeof NotificationService !== 'undefined') {
          NotificationService.error(
            'Sua sessão expirou. Por favor, faça login novamente.',
            'Autenticação'
          );
        }
        return { success: false, data: null, message: 'Não autenticado' };
      }

      const responseData = await response.json();

      if (!response.ok && !responseData.success) {
        if (typeof NotificationService !== 'undefined') {
          NotificationService.error(
            responseData.message || 'Erro na requisição',
            `Erro ${response.status}`
          );
        }
      }

      return responseData;
    } catch (error) {
      if (typeof NotificationService !== 'undefined') {
        NotificationService.error(
          error.message || 'Erro de conexão com o servidor',
          'Erro de Rede'
        );
      }
      return { success: false, data: null, message: error.message };
    }
  }

  async get(endpoint) {
    return this.request('GET', endpoint);
  }

  async post(endpoint, body) {
    return this.request('POST', endpoint, body);
  }

  async delete(endpoint) {
    return this.request('DELETE', endpoint);
  }

  async patch(endpoint, body) {
    return this.request('PATCH', endpoint, body);
  }
}
