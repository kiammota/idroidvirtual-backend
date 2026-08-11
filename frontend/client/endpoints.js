export const ENDPOINTS = {
  // Rotas Púbicas e Gerais
  BASE: "/",
  LOGIN: "/login",
  SWAGGER: (filepath) => `/swagger/${filepath}`,

  // Autenticação e Usuários
  AUTH: {
    LOGIN: "/auth/login",
  },
  USERS: {
    BASE: "/users",
    BY_ID: (id) => `/users/${id}`,
    ELEVATE: "/users/elevate",
  },

  // Ordens de Serviço
  SERVICE_ORDERS: {
    BASE: "/service-orders",
    THREE_MONTHS: "/service-orders/three-months",
    SIX_MONTHS: "/service-orders/six-months",
    BY_CLIENT_ID: (clienteId) => `/service-orders/client-id/${clienteId}`,
    BY_CLIENT_NAME: (clientName) => `/service-orders/client-name/${encodeURIComponent(clientName)}`,
    BY_CLIENT_EMAIL: (clientEmail) => `/service-orders/client-email/${encodeURIComponent(clientEmail)}`,
    BY_CLIENT_NUMBER: (clientNumber) => `/service-orders/client-number/${clientNumber}`,
    BY_ID: (id) => `/service-orders/${id}`,
    BY_STATUS: (status) => `/service-orders/status/${encodeURIComponent(status)}`,
  },

  // Clientes
  CLIENTS: {
    BASE: "/clients/",
    NUMBER: "/clients/number",
    ADDRESS: "/clients/address",
    NAME: "/clients/name",
    EMAIL: "/clients/email",
    BY_ID: (id) => `/clients/${id}`,
  },

  // Produtos
  PRODUCTS: {
    BASE: "/products",
    BY_CATEGORY: (category) => `/products/category/${encodeURIComponent(category)}`,
    BY_ID: (id) => `/products/${id}`,
    SET_STOCK: (id, num) => `/products/${id}/set/${num}`,
  },

  // Vendas
  SALES: {
    BASE: "/sales",
    DIA: "/sales/dia",
    MES: "/sales/mes",
    BY_ID: (id) => `/sales/${id}`,
  },


  // Trocas
  TROCAS: {
    BASE: "/trocas",
    ULTIMO_MES: "/trocas/ultimo-mes",
    TRES_MESES: "/trocas/tres-meses",
    SEIS_MESES: "/trocas/seis-meses",
    BY_ID: (id) => `/trocas/venda/${id}`,
  },

  // Financeiro
  FINANCIAL: {
    TODAY: "/financial/today",
    MONTH: "/financial/month",
    PERIOD: "/financial/period",
  },

  // Despesas
  EXPENSES: {
    BASE: "/expenses",
    BY_ID: (id) => `/expenses/${id}`,
  },
  AUDIT: {
    BASE: "/audit"
  }
};

export default ENDPOINTS;