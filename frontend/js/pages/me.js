import { requireAuth, getCurrentUser, getToken, clearToken, goTo, goBack, showToast } from '../app.js';
import { UserService } from '../../client/services/UsuarioService.js';
import { NotificationService } from '../../client/services/notificationService.js';

window.NotificationService = NotificationService;

/**
 * Página Meu Perfil
 */
async function init() {
  if (!requireAuth()) return;

  const user = getCurrentUser();
  const token = getToken() || '';

  document.getElementById('btn-back').onclick = () => goBack();
  document.getElementById('user-nome').textContent = user.nome || '—';
  document.getElementById('user-email').textContent = user.email || '—';
  document.getElementById('user-id').textContent = user.id ?? '—';
  document.getElementById('user-empresa').textContent = user.empresa_id ?? '—';
  document.getElementById('user-token').textContent = token ? token.slice(0, 40) + '…' : '—';
  document.getElementById('avatar').textContent = (user.nome || '?').charAt(0).toUpperCase();

  // Buscar dados atualizados do usuário
  if (user.id) {
    try {
      const response = await UserService.getById(user.id);
      if (response.success && response.data) {
        const u = response.data;
        document.getElementById('user-nome').textContent = u.nome || user.nome || '—';
        document.getElementById('user-email').textContent = u.email || user.email || '—';
        document.getElementById('user-id').textContent = u.id ?? user.id ?? '—';
        document.getElementById('user-empresa').textContent = u.empresa_id ?? user.empresa_id ?? '—';
        document.getElementById('avatar').textContent = (u.nome || user.nome || '?').charAt(0).toUpperCase();
      }
    } catch {
      // mantém dados locais
    }
  }

  document.getElementById('btn-logout').onclick = () => {
    clearToken();
    sessionStorage.removeItem('gv_user');
    showToast('Logout realizado', 'success');
    setTimeout(() => goTo('login'), 300);
  };
}

init();
