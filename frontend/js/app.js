/**
 * IDroid Virtual - Gestão
 * Shared utilities & auth
 */

import { TOKEN_STORAGE_KEY, API_BASE_URL } from '../client/env.js';
import { UserService } from '../client/services/UsuarioService.js';

/** Base relativa: sempre sobe para /pages/ a partir de qualquer subpasta */
export function pagesBase() {
  // Ex: /pages/home/index.html -> ../
  // Ex: /pages/login/ -> ../
  return '../';
}

export function getToken() {
  return sessionStorage.getItem(TOKEN_STORAGE_KEY);
}

export function setToken(token) {
  sessionStorage.setItem(TOKEN_STORAGE_KEY, token);
}

export function clearToken() {
  sessionStorage.removeItem(TOKEN_STORAGE_KEY);
}

export function isAuthenticated() {
  return !!getToken();
}

/**
 * Protege páginas: redireciona para login se não autenticado
 */
export function requireAuth() {
  if (!isAuthenticated()) {
    window.location.href = pagesBase() + 'login/';
    return false;
  }
  return true;
}

export function getCurrentUser() {
  const stored = sessionStorage.getItem('gv_user');
  if (stored) {
    try { return JSON.parse(stored); } catch { /* ignore */ }
  }
  return {
    id: 1,
    nome: 'Administrador',
    email: 'admin@idroid.local',
    empresa_id: 1
  };
}

export function setCurrentUser(user) {
  sessionStorage.setItem('gv_user', JSON.stringify(user));
}

const _toastRecent = new Map();

export function showToast(message, type = 'info') {
  const key = `${type}::${String(message || '').trim()}`;
  const now = Date.now();
  const last = _toastRecent.get(key) || 0;
  if (now - last < 2500) return; // evita notificações duplicadas
  _toastRecent.set(key, now);
  if (_toastRecent.size > 30) {
    const cutoff = now - 5000;
    for (const [k, t] of _toastRecent) {
      if (t < cutoff) _toastRecent.delete(k);
    }
  }

  let container = document.querySelector('.toast-container');
  if (!container) {
    container = document.createElement('div');
    container.className = 'toast-container';
    document.body.appendChild(container);
  }
  const toast = document.createElement('div');
  toast.className = `toast ${type}`;
  toast.textContent = message;
  container.appendChild(toast);
  setTimeout(() => {
    toast.style.opacity = '0';
    toast.style.transition = 'opacity .3s';
    setTimeout(() => toast.remove(), 300);
  }, 3200);
}

export function formatBRL(value) {
  return Number(value || 0).toLocaleString('pt-BR', {
    style: 'currency',
    currency: 'BRL'
  });
}

export function formatPhone(tel) {
  if (!tel) return '';
  const d = String(tel).replace(/\D/g, '');
  if (d.length === 11) return d.replace(/^(\d{2})(\d{5})(\d{4})$/, '($1) $2-$3');
  if (d.length === 10) return d.replace(/^(\d{2})(\d{4})(\d{4})$/, '($1) $2-$3');
  return tel;
}

export function openModal(id) {
  const el = document.getElementById(id);
  if (el) el.classList.add('open');
}

export function closeModal(id) {
  const el = document.getElementById(id);
  if (el) el.classList.remove('open');
}

export function requestPin(title = 'Área protegida') {
  return new Promise((resolve) => {
    let overlay = document.getElementById('pin-modal');
    if (!overlay) {
      overlay = document.createElement('div');
      overlay.id = 'pin-modal';
      overlay.className = 'modal-overlay center';
      overlay.innerHTML = `
        <div class="modal">
          <div class="modal-title" id="pin-title">${title}</div>
          <p class="text-muted text-center mb-2" style="font-size:0.9rem">Digite o PIN ou senha administrativa</p>
          <div class="form-group">
            <input type="password" id="pin-input" placeholder="••••" maxlength="12" inputmode="numeric" autocomplete="off">
          </div>
          <div class="modal-actions">
            <button class="btn btn-ghost" id="pin-cancel">Cancelar</button>
            <button class="btn btn-primary" id="pin-confirm">Confirmar</button>
          </div>
        </div>
      `;
      document.body.appendChild(overlay);
    }
    document.getElementById('pin-title').textContent = title;
    const input = document.getElementById('pin-input');
    input.value = '';
    openModal('pin-modal');
    setTimeout(() => input.focus(), 100);

    const cleanup = (result) => {
      closeModal('pin-modal');
      document.getElementById('pin-confirm').onclick = null;
      document.getElementById('pin-cancel').onclick = null;
      resolve(result);
    };

    document.getElementById('pin-confirm').onclick = async () => {
      const pin = input.value.trim();
      if (pin.length < 4) {
        showToast('PIN inválido', 'error');
        return;
      }
      try {
        const res = await UserService.elevate(pin);
        if (res.success) {
          showToast('Acesso liberado', 'success');
          cleanup(true);
        } else {
          showToast(res.message || 'Senha administrativa inválida', 'error');
        }
      } catch (err) {
        showToast(err.message || 'Erro ao validar PIN', 'error');
      }
    };
    document.getElementById('pin-cancel').onclick = () => cleanup(false);
  });
}

/** Navegação relativa entre páginas irmãs em /pages/ */
export function goTo(page) {
  // page: 'login', 'home', 'me', 'venda', etc. (sem barra inicial)
  const clean = String(page).replace(/^\/?(pages\/)?/, '').replace(/\/$/, '');
  window.location.href = pagesBase() + clean + '/';
}

export function goBack() {
  if (window.history.length > 1) {
    window.history.back();
  } else {
    goTo('home');
  }
}

window.IDroid = {
  getToken, setToken, clearToken, isAuthenticated, requireAuth,
  getCurrentUser, setCurrentUser, showToast, formatBRL, formatPhone,
  openModal, closeModal, requestPin, goTo, goBack, pagesBase, API_BASE_URL
};
