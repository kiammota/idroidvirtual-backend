import { requireAuth, getCurrentUser, clearToken, goTo, requestPin, showToast } from '../app.js';
import { NotificationService } from '../../client/services/notificationService.js';
import { SaleService } from '../../client/services/VendasService.js';
import { ServiceOrderService } from '../../client/services/OrdensServicoService.js';
import { ClientService } from '../../client/services/clientService.js';

window.NotificationService = NotificationService;

/**
 * Página Home / Dashboard
 */
async function init() {
  if (!requireAuth()) return;

  const user = getCurrentUser();
  document.getElementById('greeting').textContent = `Olá, ${user.nome || 'Usuário'}`;

  document.getElementById('btn-me').onclick = () => goTo('me');
  document.getElementById('btn-logout').onclick = () => {
    clearToken();
    sessionStorage.removeItem('gv_user');
    showToast('Logout realizado', 'success');
    setTimeout(() => goTo('login'), 300);
  };

  document.getElementById('btn-nova-venda').onclick = () => goTo('venda');

  async function unlockAndGo(path) {
    const ok = await requestPin('Área protegida');
    if (ok) goTo(path);
  }

  document.getElementById('card-estoque').addEventListener('click', (e) => {
    e.preventDefault();
    unlockAndGo('estoque');
  });
  document.getElementById('card-financeiro').addEventListener('click', (e) => {
    e.preventDefault();
    unlockAndGo('financeiro');
  });

  // Carregar resumo via API
  try {
    const [vendasRes, osRes, clientesRes] = await Promise.all([
      SaleService.getByDay(),
      ServiceOrderService.listAll(),
      ClientService.listAll(),
    ]);

    const qtdVendas =
      vendasRes.success && Array.isArray(vendasRes.data) ? vendasRes.data.length : 0;
    const qtdOs =
      osRes.success && Array.isArray(osRes.data) ? osRes.data.length : 0;
    const qtdClientes =
      clientesRes.success && Array.isArray(clientesRes.data) ? clientesRes.data.length : 0;

    document.getElementById('resumo-vendas').textContent = String(qtdVendas);
    document.getElementById('resumo-os').textContent = String(qtdOs);
    document.getElementById('resumo-clientes').textContent = String(qtdClientes);
  } catch {
    document.getElementById('resumo-vendas').textContent = '0';
    document.getElementById('resumo-os').textContent = '0';
    document.getElementById('resumo-clientes').textContent = '0';
  }
}

init();
