import { HttpClient } from "../http.js";
import { Venda } from "../models/VendaModel.js";

class SaleServiceClass extends HttpClient {
  /**
   * Cria uma venda
   * Obrigatórios: produto_id, quantidade, tipo_pagamento, usuario_id
   * Opcional: cliente_id (0 = Consumidor Final)
   */
  async create(params) {
    const body = {
      produto_id: Number(params.produto_id),
      quantidade: Number(params.quantidade) || 1,
      tipo_pagamento: params.tipo_pagamento || "PIX",
      usuario_id: Number(params.usuario_id), // obrigatório
      cliente_id: Number(params.cliente_id) || 0, // 0 = Consumidor Final
    };

    if (params.total != null) {
      body.total = Number(params.total);
    }

    const response = await this.post("/sales", body);
    if (response.success && response.data) {
      response.data = new Venda(response.data);
    }
    return response;
  }

  async remove(id) {
    return this.delete(`/sales/${id}`);
  }

  async deleteSale(id) {
    return this.remove(id);
  }

  async getByDay() {
    const response = await this.get("/sales/dia");
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Venda(data));
    }
    return response;
  }

  async getByMonth() {
    const response = await this.get("/sales/mes");
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Venda(data));
    }
    return response;
  }
}

/** Singleton */
export const SaleService = new SaleServiceClass();
export default SaleService;