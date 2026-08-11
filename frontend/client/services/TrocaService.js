import ENDPOINTS from "../endpoints.js";
import { HttpClient } from "../http.js";
import { Troca } from "../models/TrocaModel.js";

function mapList(response) {
  // Backend devolve array puro: [{...}, ...]
  if (Array.isArray(response)) {
    return { success: true, data: response.map((d) => new Troca(d)) };
  }
  if (response && Array.isArray(response.data)) {
    response.data = response.data.map((data) => new Troca(data));
    response.success = response.success !== false;
    return response;
  }
  if (response && response.error) {
    return { success: false, message: response.error, data: [] };
  }
  return response || { success: false, data: [] };
}

function mapOne(response) {
  if (!response) return { success: false };
  // objeto Troca direto
  if (response.id != null && !response.data) {
    return { success: true, data: new Troca(response) };
  }
  if (response.data) {
    response.data = new Troca(response.data);
    response.success = response.success !== false;
    return response;
  }
  if (response.error || response.message) {
    return { success: false, message: response.error || response.message };
  }
  return response;
}

class TrocaServiceClass extends HttpClient {
  /**
   * POST /trocas
   * Body: { modelo_aparelho, valor_avaliacao, imei?, condicao?, observacoes?, venda_id? }
   */
  async create(params) {
    const body =
      params && typeof params.toCreatePayload === "function"
        ? params.toCreatePayload()
        : {
            modelo_aparelho: String(params.modelo_aparelho ?? "").trim(),
            valor_avaliacao: Number(params.valor_avaliacao) || 0,
            ...(params.venda_id > 0 ? { venda_id: Number(params.venda_id) } : {}),
            ...(params.imei ? { imei: String(params.imei).replace(/\D/g, "") } : {}),
            ...(params.condicao ? { condicao: params.condicao } : {}),
            ...(params.observacoes != null && params.observacoes !== ""
              ? { observacoes: params.observacoes }
              : {}),
          };

    const response = await this.post(ENDPOINTS.TROCAS.BASE, body);
    return mapOne(response);
  }

  /**
   * GET /trocas?limit=&offset=
   */
  async listAll(limit = 50, offset = 0) {
    const response = await this.get(
      `${ENDPOINTS.TROCAS.BASE}?limit=${limit}&offset=${offset}`
    );
    return mapList(response);
  }

  /**
   * GET /trocas/ultimo-mes?limit=&offset=
   */
  async listLastMonth(limit = 50, offset = 0) {
    const response = await this.get(
      `${ENDPOINTS.TROCAS.ULTIMO_MES}?limit=${limit}&offset=${offset}`
    );
    return mapList(response);
  }

  /**
   * GET /trocas/tres-meses?limit=&offset=
   */
  async listLastThreeMonths(limit = 50, offset = 0) {
    const response = await this.get(
      `${ENDPOINTS.TROCAS.TRES_MESES}?limit=${limit}&offset=${offset}`
    );
    return mapList(response);
  }

  /**
   * GET /trocas/seis-meses?limit=&offset=
   */
  async listLastSixMonths(limit = 50, offset = 0) {
    const response = await this.get(
      `${ENDPOINTS.TROCAS.SEIS_MESES}?limit=${limit}&offset=${offset}`
    );
    return mapList(response);
  }

  /**
   * GET /trocas/venda/:id  → ObterPorID
   */
  async getById(id) {
    const response = await this.get(ENDPOINTS.TROCAS.BY_ID(id));
    return mapOne(response);
  }
}

/** Singleton */
export const TrocaService = new TrocaServiceClass();
export default TrocaService;
