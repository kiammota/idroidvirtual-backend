import { HttpClient } from "../http.js";
import { Produto } from "../models/ProdutoModel.js"; // Ajuste o caminho de acordo com o seu projeto

export class ProductService extends HttpClient {
  /**
   * Lista todos os produtos
   * 
   * @returns {Promise<{success: boolean, data: Produto[], message: string}>}
   */
  async listAll() {
    const response = await this.get('/products');

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new Produto(data));
    }

    return response;
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
   * @returns {Promise<{success: boolean, data: Produto, message: string}>}
   */
  async create(params) {
    const { nome, categoria, descricao, preco, estoque } = params;
    const response = await this.post('/products', { nome, categoria, descricao, preco, estoque });

    if (response.success && response.data) {
      response.data = new Produto(response.data);
    }

    return response;
  }

  /**
   * Lista produtos por categoria
   * 
   * | Parâmetro | Tipo   | Descrição |
   * |-----------|--------|-----------|
   * | categoria | string | Categoria |
   * 
   * @param {string} categoria - Categoria (ex: CELULAR)
   * @returns {Promise<{success: boolean, data: Produto[], message: string}>}
   */
  async listByCategory(categoria) {
    const response = await this.get(`/products/category/${encodeURIComponent(categoria)}`);

    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map(data => new Produto(data));
    }

    return response;
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
   * @returns {Promise<{success: boolean, data: Produto, message: string}>}
   */
  async update(params) {
    const { id, nome, categoria, descricao, preco, estoque } = params;
    const response = await this.post(`/products/${id}`, { nome, categoria, descricao, preco, estoque });

    if (response.success && response.data) {
      response.data = new Produto(response.data);
    }

    return response;
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