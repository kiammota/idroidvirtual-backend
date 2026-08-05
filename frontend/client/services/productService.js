import { HttpClient } from "../http.js";

export class ProductService extends HttpClient {
  /**
   * Lista todos os produtos
   * 
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async listAll() {
    return this.get('/products');
  }

  /**
   * Cria um novo produto
   * 
   * | Parâmetro  | Tipo     | Descrição         |
   * |------------|----------|-------------------|
   * | nome       | string   | Nome do produto   |
   * | categoria  | string   | Categoria         |
   * | descricao  | string   | Descrição         |
   * | preco      | number   | Preço             |
   * | estoque    | integer  | Qtd em estoque    |
   * 
   * @param {Object} params
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async create(params) {
    const { nome, categoria, descricao, preco, estoque } = params;
    return this.post('/products', { nome, categoria, descricao, preco, estoque });
  }

  /**
   * Lista produtos por categoria
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | categoria | string | Categoria |
   * 
   * @param {string} categoria - Categoria (ex: CELULAR)
   * @returns {Promise<{success: boolean, data: Array, message: string}>}
   */
  async listByCategory(categoria) {
    return this.get(`/products/category/${encodeURIComponent(categoria)}`);
  }

  /**
   * Atualiza um produto completo
   * 
   * | Parâmetro  | Tipo     | Descrição         |
   * |------------|----------|-------------------|
   * | id         | integer  | ID do produto     |
   * | nome       | string   | Nome do produto   |
   * | categoria  | string   | Categoria         |
   * | descricao  | string   | Descrição         |
   * | preco      | number   | Preço             |
   * | estoque    | integer  | Qtd em estoque    |
   * 
   * @param {Object} params
   * @returns {Promise<{success: boolean, data: Object, message: string}>}
   */
  async update(params) {
    const { id, nome, categoria, descricao, preco, estoque } = params;
    return this.post(`/products/${id}`, { nome, categoria, descricao, preco, estoque });
  }

  /**
   * Deleta um produto
   * 
   * | Parâmetro | Tipo    | Descrição    |
   * |-----------|---------|--------------|
   * | id        | integer | ID do produto|
   * 
   * @param {number} id - ID do produto
   * @returns {Promise<{success: boolean, data: null, message: string}>}
   */
  async delete(id) {
    return this.delete(`/products/${id}`);
  }
}
