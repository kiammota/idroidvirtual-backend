import { HttpClient } from "../http.js";
import { Venda } from "../models/VendaModel.js"; // Ajuste o caminho conforme a estrutura do seu projeto

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
   * @returns {Promise<{success: boolean, data: Venda, message: string}>}
   */
  async create(params) {
    const { cliente_id, produto_id, quantidade, preco_unitario } = params;
    const response = await this.post('/sales', { cliente_id, produto_id, quantidade, preco_unitario });

    if (response.success && response.data) {
      response.data = new Venda(response.data);
    }

    return response;
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
  async remove(id) {
    return this.delete(`/sales/${id}`);
  }

  async delete(id) {
    return this.remove(id);
  }

  /**
   * Lista vendas do dia
   * 
   * @returns {Promise<{success: boolean, data: Venda[], message: string}>}
   */
  async getByDay() {
    const response = await this.get('/sales/dia');

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new Venda(data));
    }

    return response;
  }

  /**
   * Lista vendas do mês
   * 
   * @returns {Promise<{success: boolean, data: Venda[], message: string}>}
   */
  async getByMonth() {
    const response = await this.get('/sales/mes');

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new Venda(data));
    }

    return response;
  }
}