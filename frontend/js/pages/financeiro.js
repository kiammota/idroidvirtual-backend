import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { ResumoFinanceiro } from '../../client/models/ResumoFinanceiroModel.js';
import { ClosingService } from '../../client/services/closingService.js';
import { SaleService } from '../../client/services/VendasService.js';
import { NotificationService } from '../../client/services/notificationService.js';

window.NotificationService = NotificationService;

/**
 * Página Financeiro / Fechamento
 */
async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();

  let contas = [
    { nome: 'Aluguel', valor: 2500, dia: 5 },
    { nome: 'Luz', valor: 380, dia: 15 },
    { nome: 'Internet', valor: 120, dia: 10 },
  ];

  document.querySelectorAll('#fin-tabs .tab').forEach((tab) => {
    tab.onclick = () => {
      document.querySelectorAll('#fin-tabs .tab').forEach((t) => t.classList.remove('active'));
      tab.classList.add('active');
      document.getElementById('tab-fechamento').classList.toggle('hidden', tab.dataset.tab !== 'fechamento');
      document.getElementById('tab-contas').classList.toggle('hidden', tab.dataset.tab !== 'contas');
    };
  });

  async function loadFechamento(periodo) {
    try {
      const response =
        periodo === 'dia' ? await ClosingService.getByDay() : await ClosingService.getByMonth();

      let data;
      if (response.success && response.data) {
        const d = response.data;
        data = new ResumoFinanceiro({
          total_servicos: d.total_servicos ?? d.servicos ?? 0,
          quantidade_os: d.quantidade_os ?? d.qtd_os ?? 0,
          total_vendas: d.total_vendas ?? d.vendas ?? 0,
          quantidade_vendas: d.quantidade_vendas ?? d.qtd_vendas ?? 0,
          faturamento_bruto: d.faturamento_bruto ?? d.bruto ?? 0,
          total_despesas: d.total_despesas ?? d.despesas ?? 0,
          lucro_liquido: d.lucro_liquido ?? d.lucro ?? 0,
        });
      } else {
        // Fallback: tenta agregar vendas
        const salesRes =
          periodo === 'dia' ? await SaleService.getByDay() : await SaleService.getByMonth();
        const vendas = salesRes.success && Array.isArray(salesRes.data) ? salesRes.data : [];
        const totalVendas = vendas.reduce((s, v) => s + (Number(v.total) || 0), 0);
        data = new ResumoFinanceiro({
          total_servicos: 0,
          quantidade_os: 0,
          total_vendas: totalVendas,
          quantidade_vendas: vendas.length,
          faturamento_bruto: totalVendas,
          total_despesas: 0,
          lucro_liquido: totalVendas,
        });
      }

      document.getElementById('f-servicos').textContent = formatBRL(data.total_servicos);
      document.getElementById('f-qtd-os').textContent = data.quantidade_os;
      document.getElementById('f-vendas').textContent = formatBRL(data.total_vendas);
      document.getElementById('f-qtd-vendas').textContent = data.quantidade_vendas;
      document.getElementById('f-bruto').textContent = formatBRL(data.faturamento_bruto);
      document.getElementById('f-despesas').textContent = formatBRL(data.total_despesas);
      document.getElementById('f-lucro').textContent = formatBRL(data.lucro_liquido);
    } catch (err) {
      showToast('Erro ao carregar fechamento', 'error');
    }
  }

  document.getElementById('btn-dia').onclick = function () {
    this.className = 'btn btn-primary btn-sm';
    document.getElementById('btn-mes').className = 'btn btn-outline btn-sm';
    loadFechamento('dia');
  };
  document.getElementById('btn-mes').onclick = function () {
    this.className = 'btn btn-primary btn-sm';
    document.getElementById('btn-dia').className = 'btn btn-outline btn-sm';
    loadFechamento('mes');
  };

  function renderContas() {
    document.getElementById('contas-list').innerHTML = contas
      .map(
        (c) => `
      <div class="card list-item">
        <div class="info">
          <div class="title">${c.nome}</div>
          <div class="sub">Vence dia ${c.dia}</div>
        </div>
        <div class="fw-bold">${formatBRL(c.valor)}</div>
      </div>
    `
      )
      .join('');
  }

  document.getElementById('form-conta').onsubmit = (e) => {
    e.preventDefault();
    contas.push({
      nome: document.getElementById('conta-nome').value.trim(),
      valor: parseFloat(document.getElementById('conta-valor').value) || 0,
      dia: parseInt(document.getElementById('conta-dia').value) || 1,
    });
    showToast('Lembrete adicionado', 'success');
    document.getElementById('form-conta').reset();
    renderContas();
  };

  await loadFechamento('mes');
  renderContas();
}

init();
