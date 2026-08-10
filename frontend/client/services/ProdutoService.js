import ENDPOINTS from "../endpoints.js";
import { HttpClient } from "../http.js";
import { Produto } from "../models/ProdutoModel.js";

class ProductServiceClass extends HttpClient {
  async listAll() {
    const response = await this.get(ENDPOINTS.PRODUCTS.BASE);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Produto(data));
    }
    return response;
  }

  async create(params) {
    const { nome, categoria, preco, estoque } = params;
    const body = { nome, categoria, preco, estoque };
    const response = await this.post(ENDPOINTS.PRODUCTS.BASE, body);
    if (response.success && response.data) {
      response.data = new Produto(response.data);
    }
    return response;
  }

  async listByCategory(categoria) {
    const response = await this.get(ENDPOINTS.PRODUCTS.BY_CATEGORY(categoria));
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Produto(data));
    }
    return response;
  }

  async update(params) {
    const { id, nome, categoria, preco, estoque } = params;
    const body = { nome, categoria, preco, estoque };
    const response = await this.post(ENDPOINTS.PRODUCTS.BY_ID(id), body);
    if (response.success && response.data) {
      response.data = new Produto(response.data);
    }
    return response;
  }

  async remove(id) {
    return this.delete(ENDPOINTS.PRODUCTS.BY_ID(id));
  }

  async deleteProduct(id) {
    return this.remove(id);
  }

  /**
   * Define o estoque absoluto do produto para o valor informado.
   * Endpoint: /products/{id}/set/{num}
   */
  async setStock(id, num) {
    const response = await this.post(ENDPOINTS.PRODUCTS.SET_STOCK(id, num));
    if (response.success && response.data) {
      response.data = new Produto(response.data);
    }
    return response;
  }
}

/** Singleton */
export const ProductService = new ProductServiceClass();
export default ProductService;