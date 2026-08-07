import { HttpClient } from "../http.js";
import { Venda } from "../models/VendaModel.js";

class SaleServiceClass extends HttpClient {
  async create(params) {
    const body = {
      produto_id: Number(params.produto_id),
      quantidade: Number(params.quantidade) || 1,
      tipo_pagamento: params.tipo_pagamento || "PIX",
      usuario_id: Number(params.usuario_id),
      cliente_id: Number(params.cliente_id) || 0,
    };
    if (params.total != null) body.total = Number(params.total);

    const response = await this.post("/sales", body);
    if (response.success && response.data) {
      response.data = new Venda(response.data);
    }
    return response;
  }

  /**
   * Cancela uma venda
   * O backend exige usuario_id (e opcionalmente motivo) no body
   */
  async remove(id, usuarioId, motivo = "Cancelamento manual") {
    const body = {
      usuario_id: Number(usuarioId),
      motivo: motivo || "Cancelamento manual",
    };
    return this.delete(`/sales/${id}`, body); // precisa que o HttpClient aceite body no DELETE
  }

  async deleteSale(id, usuarioId, motivo) {
    return this.remove(id, usuarioId, motivo);
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

export const SaleService = new SaleServiceClass();
export default SaleService;