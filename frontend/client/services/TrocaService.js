import { HttpClient } from "../http.js";
import { Troca } from "../models/TrocaModel.js";

class TrocaServiceClass extends HttpClient {
  async create(params) {
    const body =
      params && typeof params.toCreatePayload === "function"
        ? params.toCreatePayload()
        : {
            venda_id: Number(params.venda_id),
            modelo_aparelho: String(params.modelo_aparelho ?? "").trim(),
            imei: String(params.imei ?? "").replace(/\D/g, ""),
            valor_avaliacao: Number(params.valor_avaliacao) || 0,
            ...(params.condicao ? { condicao: params.condicao } : {}),
            ...(params.observacoes != null ? { observacoes: params.observacoes } : {}),
          };

    const response = await this.post("/trocas", body);
    if (response.success && response.data) {
      response.data = new Troca(response.data);
    }
    return response;
  }

  async getBySaleId(vendaId) {
    const response = await this.get(`/trocas/venda/${vendaId}`);
    if (response.success && response.data) {
      response.data = new Troca(response.data);
    }
    return response;
  }
}

/** Singleton */
export const TrocaService = new TrocaServiceClass();
export default TrocaService;
