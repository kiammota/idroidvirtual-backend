import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { bindMasks } from '../masks.js';
import { Troca } from '../../client/models/TrocaModel.js';
import { TrocaService } from '../../client/services/TrocaService.js';
import { NotificationService } from '../../client/services/notificationService.js';
import { openInspect } from '../components/InspectModal.js';

window.NotificationService = NotificationService;

/**
 * Página de Trocas (Trade-in)
 * Troca ≠ Venda — usuário não informa venda_id.
 */
async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();
  bindMasks({ 'troca-imei': 'imei' });

  let todasTrocas = [];

  // ─── Tabs ──────────────────────────────────────────────
  document.querySelectorAll('#troca-tabs .tab').forEach((tab) => {
    tab.onclick = () => {
      document.querySelectorAll('#troca-tabs .tab').forEach((t) => t.classList.remove('active'));
      tab.classList.add('active');
      const name = tab.dataset.tab;
      document.getElementById('tab-nova').classList.toggle('hidden', name !== 'nova');
      document.getElementById('tab-lista').classList.toggle('hidden', name !== 'lista');
      if (name === 'lista') carregarTrocas();
    };
  });

  // ─── Form submit ───────────────────────────────────────
  document.getElementById('form-troca').onsubmit = async (e) => {
    e.preventDefault();

    const modelo = document.getElementById('troca-modelo').value.trim();
    const imeiRaw = document.getElementById('troca-imei').value.replace(/\D/g, '');
    const valorAvaliacao = parseFloat(document.getElementById('troca-avaliacao').value) || 0;
    const condicao = document.getElementById('troca-condicao').value.trim();
    const observacoes = document.getElementById('troca-observacoes').value.trim() || null;

    if (!modelo) {
      showToast('Informe o modelo do aparelho', 'error');
      return;
    }
    if (imeiRaw && (imeiRaw.length < 14 || imeiRaw.length > 15)) {
      showToast('IMEI deve ter 14 ou 15 dígitos', 'error');
      return;
    }
    if (valorAvaliacao < 0) {
      showToast('Valor de avaliação inválido', 'error');
      return;
    }

    const troca = new Troca({
      modelo_aparelho: modelo,
      imei: imeiRaw || '',
      valor_avaliacao: valorAvaliacao,
      condicao: condicao || '',
      observacoes,
    });

    const erros = troca.validar();
    if (erros.length) {
      showToast(erros[0], 'error');
      return;
    }

    const btn = document.querySelector('#form-troca button[type="submit"]');
    if (btn) {
      btn.disabled = true;
      btn.textContent = 'Salvando…';
    }

    try {
      const response = await TrocaService.create(troca.toCreatePayload());
      if (response.success) {
        showToast('Troca registrada com sucesso', 'success');
        document.getElementById('form-troca').reset();
        carregarTrocas();
      } else {
        showToast(response.message || response.error || 'Erro ao registrar troca', 'error');
      }
    } catch (err) {
      showToast(err.message || 'Erro de conexão ao registrar troca', 'error');
    } finally {
      if (btn) {
        btn.disabled = false;
        btn.textContent = 'Registrar Troca';
      }
    }
  };

  // ─── Lista / Histórico ─────────────────────────────────
  function abrirInspectTroca(t) {
    const criado = t.criado_em
      ? new Date(t.criado_em).toLocaleString('pt-BR')
      : '—';

    const fields = [
      { label: 'ID', value: `#${t.id}` },
      { label: 'Modelo', value: t.modelo_aparelho || '—' },
      { label: 'IMEI', value: t.imei || '—' },
      { label: 'Avaliação', value: formatBRL(t.valor_avaliacao), highlight: true },
      { label: 'Condição', value: t.condicao || '—' },
      { label: 'Observações', value: t.observacoes || '—' },
      { label: 'Registrada em', value: criado },
    ];
    if (t.venda_id > 0) {
      fields.splice(1, 0, { label: 'Venda', value: `#${t.venda_id}` });
    }

    openInspect({
      title: `Troca #${t.id}`,
      subtitle: 'Trade-in · IDroid Virtual',
      fields,
      footer: 'IDroid Virtual · Gestão',
    });
  }

  function renderLista(items) {
    const el = document.getElementById('troca-list');
    if (!items || !items.length) {
      el.innerHTML = `<div class="empty-state"><div class="icon">🔄</div><p>Nenhuma troca registrada ainda.</p></div>`;
      return;
    }

    el.innerHTML = items
      .map((t) => {
        const imeiStr = t.imei ? ` · IMEI ${t.imei}` : '';
        const condStr = t.condicao ? ` · ${t.condicao}` : '';
        return `
      <div class="card list-item" data-id="${t.id}" style="cursor:pointer">
        <div class="info">
          <div class="title">${t.modelo_aparelho || 'Sem modelo'}</div>
          <div class="sub">${formatBRL(t.valor_avaliacao)}${imeiStr}${condStr}</div>
        </div>
        <div class="fw-bold text-primary">${formatBRL(t.valor_avaliacao)}</div>
      </div>`;
      })
      .join('');

    el.querySelectorAll('.list-item').forEach((card) => {
      card.onclick = () => {
        const id = Number(card.dataset.id);
        const troca = todasTrocas.find((x) => Number(x.id) === id);
        if (troca) abrirInspectTroca(troca);
      };
    });
  }

  async function carregarTrocas() {
    const el = document.getElementById('troca-list');
    el.innerHTML = `<div class="empty-state"><div class="icon">⏳</div><p>Carregando…</p></div>`;

    try {
      const response = await TrocaService.listAll();
      if (response.success && Array.isArray(response.data)) {
        todasTrocas = response.data;
      } else if (Array.isArray(response)) {
        // backend às vezes devolve array puro
        todasTrocas = response.map((d) => new Troca(d));
      } else {
        todasTrocas = [];
        if (response.message || response.error) {
          showToast(response.message || response.error, 'error');
        }
      }
    } catch (err) {
      todasTrocas = [];
      showToast(err.message || 'Erro ao carregar trocas', 'error');
    }
    renderLista(todasTrocas);
  }
}

init();
