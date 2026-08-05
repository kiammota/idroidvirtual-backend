import { requireAuth, getCurrentUser, getToken, clearToken, goTo, goBack, showToast } from '../app.js';

if (!requireAuth()) throw new Error('redirect');

const user = getCurrentUser();
const token = getToken() || '';

document.getElementById('btn-back').onclick = () => goBack();
document.getElementById('user-nome').textContent = user.nome || '—';
document.getElementById('user-email').textContent = user.email || '—';
document.getElementById('user-id').textContent = user.id ?? '—';
document.getElementById('user-empresa').textContent = user.empresa_id ?? '—';
document.getElementById('user-token').textContent = token ? token.slice(0, 40) + '…' : '—';
document.getElementById('avatar').textContent = (user.nome || '?').charAt(0).toUpperCase();

// CHAMADA DA API AQUI — buscar dados atualizados do usuário

document.getElementById('btn-logout').onclick = () => {
  clearToken();
  sessionStorage.removeItem('gv_user');
  showToast('Logout realizado', 'success');
  setTimeout(() => goTo('login'), 300);
};
