import { HttpClient } from "../http.js";

export class SaleService extends HttpClient {
  /**
   * Cria uma venda
   * 
   * | Parâmetro      | Tipo    | Descrição          |
   * |----------------|---------|-------------------|
   * | cliente_id     | integer | ID do cliente      |
   * | produto_id     | integer | ID do produto      |
   * | quantidade     | integer | Quantidade         |
   * | preco_unitario | number  | Preço unitário     |
   * 
   * @param {Object} params
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async create(params) {
    const { cliente_id, produto_id, quantidade, preco_unitario } = params;
    return this.post('/sales', { cliente_id, produto_id, quantidade, preco_unitario });
  }

  /**
   * Deleta uma venda
   * 
   * | Parâmetro | Tipo    | Descrição   |
   * |-----------|---------|-------------|
   * | id        | integer | ID da venda |
   * 
   * @param {number} id - ID da venda
   * @returns {Promise<{success: boolean, data: null, message: string}>}
   */
  async delete(id) {
    return this.delete(`/sales/${id}`);
  }

  /**
   * Lista vendas do dia
   * 
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async getByDay() {
    return this.get('/sales/dia');
  }

  /**
   * Lista vendas do mês
   * 
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async getByMonth() {
    return this.get('/sales/mes');
  }
}
