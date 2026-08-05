import { HttpClient } from "../http.js";

export class NotificationService {
  /**
   * Exibe uma notificação de erro
   * @param {string} message - Mensagem de erro
   * @param {string} title - Título (opcional)
   */
  static error(message, title = 'Erro') {
    console.error(`[${title}]`, message);
    // Substitua por seu sistema de notificação real
    // Ex: toast.error(message), showNotification(message), etc.
    alert(`${title}: ${message}`);
  }

  /**
   * Exibe notificação de sucesso
   * @param {string} message - Mensagem
   * @param {string} title - Título (opcional)
   */
  static success(message, title = 'Sucesso') {
    console.log(`[${title}]`, message);
  }

  /**
   * Exibe notificação de informação
   * @param {string} message - Mensagem
   * @param {string} title - Título (opcional)
   */
  static info(message, title = 'Info') {
    console.info(`[${title}]`, message);
  }
}
