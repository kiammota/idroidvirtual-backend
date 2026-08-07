import { HttpClient } from "../http.js";
import { Produto } from "../models/ProdutoModel.js";

class ProductServiceClass extends HttpClient {
  async listAll() {
    const response = await this.get("/products");
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Produto(data));
    }
    return response;
  }

  async create(params) {
    const { nome, categoria, preco, estoque } = params;
    const body = { nome, categoria, preco, estoque };
    const response = await this.post("/products", body);
    if (response.success && response.data) {
      response.data = new Produto(response.data);
    }
    return response;
  }

  async listByCategory(categoria) {
    const response = await this.get(`/products/category/${encodeURIComponent(categoria)}`);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Produto(data));
    }
    return response;
  }

  async update(params) {
    const { id, nome, categoria, preco, estoque } = params;
    const body = { nome, categoria, preco, estoque };
    const response = await this.post(`/products/${id}`, body);
    if (response.success && response.data) {
      response.data = new Produto(response.data);
    }
    return response;
  }

  async remove(id) {
    return this.delete(`/products/${id}`);
  }

  async deleteProduct(id) {
    return this.remove(id);
  }
}

/** Singleton */
export const ProductService = new ProductServiceClass();
export default ProductService;