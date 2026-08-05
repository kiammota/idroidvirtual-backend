import { HttpClient } from "../http.js";

export class ClosingService extends HttpClient {
  /**
   * Fechamento do dia
   * 
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async getByDay() {
    return this.get('/closing/dia');
  }

  /**
   * Fechamento do mês
   * 
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async getByMonth() {
    return this.get('/closing/mes');
  }
}
