import { setToken, setCurrentUser, showToast, isAuthenticated, goTo } from '../app.js';
import { LoginRequest } from '../../client/models/LoginModel.js';
import { NotificationService } from '../../client/services/notificationService.js';
import { AuthService } from '../../client/services/authService.js'; // Ajuste o caminho se necessário

window.NotificationService = NotificationService;

if (isAuthenticated()) {
  goTo('home');
}

const form = document.getElementById('login-form');
const btn = document.getElementById('btn-login');
const authService = new AuthService();

form.addEventListener('submit', async (e) => {
  e.preventDefault();
  const email = document.getElementById('email').value.trim();
  const senha = document.getElementById('senha').value;

  const model = new LoginRequest({ email, senha });
  if (!model.isValid()) {
    showToast(model.validar()[0] || 'Dados inválidos', 'error');
    return;
  }

  btn.disabled = true;
  btn.textContent = 'Entrando...';

  try {
    const response = await authService.login(email, senha);

    if (response.success) {
      if (response.data) {
        setToken(response.data);
      }
      
      const nomeUsuario = email.split('@')[0].replace(/[._]/g, ' ').replace(/\b\w/g, c => c.toUpperCase()) || 'Usuário';
      
      setCurrentUser({
        id: 1,
        nome: nomeUsuario,
        email,
        empresa_id: 1
      });

      showToast('Login realizado com sucesso', 'success');
      setTimeout(() => goTo('home'), 400);
    } else {
      throw new Error(response.message || 'Erro ao realizar login');
    }
  } catch (err) {
    showToast(err.message || 'Erro ao entrar', 'error');
    btn.disabled = false;
    btn.textContent = 'Entrar';
  }
});