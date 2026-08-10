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

  /**
   * Decodifica o payload do JWT e retorna o empresa_id
   * @returns {number|null}
   */
  getEmpresaId() {
    const token = this.getToken();
    if (!token) return null;

    try {
      const payloadBase64 = token.split('.')[1];
      if (!payloadBase64) return null;

      const base64 = payloadBase64.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      );

      const payload = JSON.parse(jsonPayload);
      return payload.empresa_id || payload.empresaId || null;
    } catch (e) {
      console.error('Erro ao decodificar token JWT para obter empresa_id:', e);
      return null;
    }
  }

  /**
   * Decodifica o payload do JWT e retorna o usuario_id (ou sub / id)
   * @returns {number|string|null}
   */
  getUsuarioId() {
    const token = this.getToken();
    if (!token) return null;

    try {
      const payloadBase64 = token.split('.')[1];
      if (!payloadBase64) return null;

      const base64 = payloadBase64.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      );

      const payload = JSON.parse(jsonPayload);
      
      // Procura pelas chaves mais comuns de ID de usuário em tokens JWT
      return payload.usuario_id || payload.usuarioId || payload.sub || payload.id || null;
    } catch (e) {
      console.error('Erro ao decodificar token JWT para obter usuario_id:', e);
      return null;
    }
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
        // Notifica só em 401 (sessão); demais erros ficam a cargo da página
        if (typeof NotificationService !== 'undefined') {
          NotificationService.error(
            'Sua sessão expirou. Por favor, faça login novamente.',
            'Autenticação'
          );
        }
        return { success: false, data: null, message: 'Não autenticado' };
      }

      const responseData = await response.json();
      // Não dispara toast aqui — evita notificação duplicada com as páginas
      return responseData;
    } catch (error) {
      // Erro de rede: retorna para a página tratar (evita toast duplo)
      return { success: false, data: null, message: error.message || 'Erro de conexão com o servidor' };
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

  async put(endpoint, body) {
    return this.request('PUT', endpoint, body);
  }
}