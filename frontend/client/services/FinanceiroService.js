import ENDPOINTS from "../endpoints.js";
import { HttpClient } from "../http.js";
import { ResumoFinanceiro } from "../models/FinanceiroModel.js";

class FinanceiroServiceClass extends HttpClient {
  /**
   * Obtém o resumo financeiro do dia atual
   * GET /financial/today
   */
  async getHoje() {
    const response = await this.get(ENDPOINTS.FINANCIAL.TODAY);
    if (response.success && response.data) {
      response.data = new ResumoFinanceiro(response.data);
    }
    return response;
  }

  /**
   * Obtém o resumo financeiro do mês corrente
   * GET /financial/month
   */
  async getMesAtual() {
    const response = await this.get(ENDPOINTS.FINANCIAL.MONTH);
    if (response.success && response.data) {
      response.data = new ResumoFinanceiro(response.data);
    }
    return response;
  }

  /**
   * Obtém o resumo financeiro por intervalo de datas
   * GET /financial/period?inicio=YYYY-MM-DD&fim=YYYY-MM-DD
   * @param {string} inicio - Data no formato YYYY-MM-DD
   * @param {string} fim - Data no formato YYYY-MM-DD
   */
  async getPorIntervalo(inicio, fim) {
    const params = new URLSearchParams({ inicio, fim }).toString();
    const response = await this.get(`${ENDPOINTS.FINANCIAL.PERIOD}?${params}`);

    if (response.success && response.data) {
      response.data = new ResumoFinanceiro(response.data);
    }
    return response;
  }
}

/** Singleton */
export const FinanceiroService = new FinanceiroServiceClass();
export default FinanceiroService;