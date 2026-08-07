import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { bindMasks } from '../masks.js';
import { Troca } from '../../client/models/TrocaModel.js';
import { TrocaService } from '../../client/services/TrocaService.js';
import { NotificationService } from '../../client/services/notificationService.js';

window.NotificationService = NotificationService;

/**
 * Página de Trocas (Trade-in)
 */
async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();
  bindMasks({ 'troca-imei': 'imei' });

  const toggleDesc = document.getElementById('toggle-desconto');
  const toggleJuros = document.getElementById('toggle-juros');
  let useDesc = true;
  let useJuros = false;

  toggleDesc.onclick = () => {
    useDesc = !useDesc;
    toggleDesc.classList.toggle('on', useDesc);
    toggleDesc.setAttribute('aria-checked', useDesc);
    calc();
  };

  toggleJuros.onclick = () => {
    useJuros = !useJuros;
    toggleJuros.classList.toggle('on', useJuros);
    toggleJuros.setAttribute('aria-checked', useJuros);
    document.getElementById('juros-group').classList.toggle('hidden', !useJuros);
    calc();
  };

  function calc() {
    const aval = parseFloat(document.getElementById('troca-avaliacao').value) || 0;
    const novo = parseFloat(document.getElementById('troca-novo').value) || 0;
    const desc = useDesc ? aval * 0.1 : 0;
    const avalLiq = aval - desc;
    let base = Math.max(0, novo - avalLiq);
    const taxa = useJuros ? (parseFloat(document.getElementById('juros-pct').value) || 0) / 100 : 0;
    const juros = base * taxa;
    const final = base + juros;

    document.getElementById('r-avaliacao').textContent = formatBRL(aval);
    document.getElementById('r-desconto').textContent = formatBRL(desc);
    document.getElementById('r-novo').textContent = formatBRL(novo);
    document.getElementById('r-juros').textContent = formatBRL(juros);
    document.getElementById('r-final').textContent = formatBRL(final);
  }

  ['troca-avaliacao', 'troca-novo', 'juros-pct'].forEach((id) => {
    document.getElementById(id).addEventListener('input', calc);
  });

  document.getElementById('form-troca').onsubmit = async (e) => {
    e.preventDefault();

    const modelo = document.getElementById('troca-modelo').value.trim();
    const imei = document.getElementById('troca-imei').value.replace(/\D/g, '');
    const valorAvaliacao = parseFloat(document.getElementById('troca-avaliacao').value) || 0;

    if (imei.length < 14) {
      showToast('IMEI deve ter 15 dígitos', 'error');
      return;
    }

    if (!modelo) {
      showToast('Informe o modelo do aparelho', 'error');
      return;
    }

    // venda_id é obrigatório na API — aqui usamos 0 se não houver contexto de venda
    // Em fluxo real, a página de venda deve passar o ID ou o usuário seleciona a venda
    const vendaIdInput = document.getElementById('troca-venda-id');
    const vendaId = vendaIdInput ? Number(vendaIdInput.value) || 0 : 0;

    const troca = new Troca({
      venda_id: vendaId || 0,
      modelo_aparelho: modelo,
      imei,
      valor_avaliacao: valorAvaliacao,
      condicao: 'USADO',
    });

    // Se não houver venda_id, ainda tentamos registrar (backend valida)
    try {
      const response = await TrocaService.create(troca.toCreatePayload());

      if (response.success) {
        showToast('Troca registrada · aparelho preparado no estoque', 'success');
        document.getElementById('form-troca').reset();
        calc();
      } else {
        showToast(response.message || 'Erro ao registrar troca', 'error');
      }
    } catch (err) {
      showToast(err.message || 'Erro de conexão ao registrar troca', 'error');
    }
  };

  calc();
}

init();
