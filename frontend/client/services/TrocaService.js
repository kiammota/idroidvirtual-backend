import ENDPOINTS from "../endpoints.js";
import { HttpClient } from "../http.js";
import { Troca } from "../models/TrocaModel.js";

class TrocaServiceClass extends HttpClient {
  /**
   * POST /api/v1/trocas
   */
  async create(params) {
    const body =
      params && typeof params.toCreatePayload === "function"
        ? params.toCreatePayload()
        : {
            venda_id: Number(params.venda_id),
            modelo_aparelho: String(params.modelo_aparelho ?? "").trim(),
            valor_avaliacao: Number(params.valor_avaliacao) || 0,
            ...(params.imei ? { imei: String(params.imei).replace(/\D/g, "") } : {}),
            ...(params.condicao ? { condicao: params.condicao } : {}),
            ...(params.observacoes != null && params.observacoes !== ""
              ? { observacoes: params.observacoes }
              : {}),
          };

    const response = await this.post(ENDPOINTS.TROCAS.BASE, body);
    if (response.success && response.data) {
      response.data = new Troca(response.data);
    }
    return response;
  }

  /**
   * GET /api/v1/trocas — listar todas (paginado no backend)
   */
  async listAll() {
    const response = await this.get(ENDPOINTS.TROCAS.BASE);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Troca(data));
    }
    return response;
  }

  /**
   * GET /api/v1/trocas/ultimo-mes
   */
  async listLastMonth() {
    const response = await this.get(ENDPOINTS.TROCAS.ULTIMO_MES);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Troca(data));
    }
    return response;
  }

  /**
   * GET /api/v1/trocas/tres-meses
   */
  async listLastThreeMonths() {
    const response = await this.get(ENDPOINTS.TROCAS.TRES_MESES);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Troca(data));
    }
    return response;
  }

  /**
   * GET /api/v1/trocas/seis-meses
   */
  async listLastSixMonths() {
    const response = await this.get(ENDPOINTS.TROCAS.SEIS_MESES);
    if (response.success && Array.isArray(response.data)) {
      response.data = response.data.map((data) => new Troca(data));
    }
    return response;
  }

  /**
   * GET /api/v1/trocas/venda/:vendaId
   */
  async getBySaleId(vendaId) {
    const response = await this.get(ENDPOINTS.TROCAS.BY_VENDA_ID(vendaId));
    if (response.success && response.data) {
      if (Array.isArray(response.data)) {
        response.data = response.data.map((d) => new Troca(d));
      } else {
        response.data = new Troca(response.data);
      }
    }
    return response;
  }
}

/** Singleton */
export const TrocaService = new TrocaServiceClass();
export default TrocaService;
