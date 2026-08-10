import ENDPOINTS from "../endpoints.js";
import { HttpClient } from "../http.js";
import { Venda } from "../models/VendaModel.js";

class SaleServiceClass extends HttpClient {
  async create(params) {
    const body = {
      empresa_id: Number(params.empresa_id), // <-- ADICIONADO AQUI
      produto_id: Number(params.produto_id),
      quantidade: Number(params.quantidade) || 1,
      tipo_pagamento: params.tipo_pagamento || "PIX",
      usuario_id: Number(params.usuario_id),
      cliente_id: Number(params.cliente_id) || 0,
    };
    if (params.total != null) body.total = Number(params.total);

    const response = await this.post(ENDPOINTS.SALES.BASE, body);
    if (response.success && response.data) {
      response.data = new Venda(response.data);
    }
    return response;
  }

  /**
   * Cancela uma venda
   * O backend exige empresa_id, usuario_id (e opcionalmente motivo) no body
   */
  async remove(id, usuarioId, empresaId, motivo = "Cancelamento manual") {
    const body = {
      empresa_id: Number(empresaId), // <-- ADICIONADO AQUI
      usuario_id: Number(usuarioId),
      motivo: motivo || "Cancelamento manual",
    };
    return this.delete(ENDPOINTS.SALES.BY_ID(id), body);
  }

  async deleteSale(id, usuarioId, empresaId, motivo) {
    return this.remove(id, usuarioId, empresaId, motivo);
  }

  async getByDay() {
    const response = await this.get(ENDPOINTS.SALES.DIA);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Venda(data));
    }
    return response;
  }

  async getByMonth() {
    const response = await this.get(ENDPOINTS.SALES.MES);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Venda(data));
    }
    return response;
  }
}

export const SaleService = new SaleServiceClass();
export default SaleService;