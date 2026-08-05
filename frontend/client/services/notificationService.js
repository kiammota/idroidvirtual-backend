export class NotificationService {
  static error(message, title = 'Erro') {
    console.error(`[${title}]`, message);
    if (window.IDroid && window.IDroid.showToast) {
      window.IDroid.showToast(message, 'error');
    } else {
      alert(`${title}: ${message}`);
    }
  }

  static success(message, title = 'Sucesso') {
    console.log(`[${title}]`, message);
    if (window.IDroid && window.IDroid.showToast) {
      window.IDroid.showToast(message, 'success');
    }
  }

  static info(message, title = 'Info') {
    console.info(`[${title}]`, message);
    if (window.IDroid && window.IDroid.showToast) {
      window.IDroid.showToast(message, 'info');
    }
  }
}

if (typeof window !== 'undefined') {
  window.NotificationService = NotificationService;
}
