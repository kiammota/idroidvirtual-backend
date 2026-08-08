import ENDPOINTS from "../endpoints.js";
import { HttpClient } from "../http.js";
import { Auditoria } from "../models/AuditoriaModel.js";

class AuditoriaServiceClass extends HttpClient {
  /**
   * Lista os registros de auditoria por período e paginação
   * @param {Object} [options]
   * @param {string} [options.periodo='dia'] - 'dia', 'mes', 'ano', etc.
   * @param {number} [options.pagina=1]
   * @param {number} [options.itens_por_pagina=20]
   * @returns {Promise<Object>}
   */
  async listar({ periodo = "dia", pagina = 1, itens_por_pagina = 20 } = {}) {
    const params = new URLSearchParams({
      periodo,
      pagina: String(pagina),
      itens_por_pagina: String(itens_por_pagina)
    }).toString();

    const response = await this.get(`${ENDPOINTS.AUDIT.BASE}?${params}`);

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(item => new Auditoria(item));
    }
    return response;
  }
}

/** Singleton */
export const AuditoriaService = new AuditoriaServiceClass();
export default AuditoriaService;