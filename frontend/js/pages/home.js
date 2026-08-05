import { requireAuth, getCurrentUser, clearToken, goTo, requestPin, showToast } from '../app.js';
import { NotificationService } from '../../client/services/notificationService.js';

window.NotificationService = NotificationService;

if (!requireAuth()) throw new Error('redirect');

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

// CHAMADA DA API AQUI — carregar resumo
document.getElementById('resumo-vendas').textContent = '0';
document.getElementById('resumo-os').textContent = '0';
document.getElementById('resumo-clientes').textContent = '0';
