import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { FinanceiroService } from '../../client/services/FinanceiroService.js';
import { DespesaService } from '../../client/services/DespesaService.js';
import { NotificationService } from '../../client/services/notificationService.js';

window.NotificationService = NotificationService;

/**
 * Auxiliar para formatar datas no formato DD/MM/YYYY
 */
function formatDate(date) {
  if (!date) return '--/--/----';
  const d = new Date(date);
  if (isNaN(d.getTime())) return '--/--/----';
  return d.toLocaleDateString('pt-BR');
}

/**
 * Página Financeiro / Fechamento & Contas a Pagar
 */
async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();

  let despesas = [];

  // Alternância das abas
  document.querySelectorAll('#fin-tabs .tab').forEach((tab) => {
    tab.onclick = () => {
      document.querySelectorAll('#fin-tabs .tab').forEach((t) => t.classList.remove('active'));
      tab.classList.add('active');
      document.getElementById('tab-fechamento').classList.toggle('hidden', tab.dataset.tab !== 'fechamento');
      document.getElementById('tab-contas').classList.toggle('hidden', tab.dataset.tab !== 'contas');
    };
  });

  // GET: Carrega o Fechamento Financeiro
  async function loadFechamento(periodo) {
    try {
      const response =
        periodo === 'dia' 
          ? await FinanceiroService.getHoje() 
          : await FinanceiroService.getMesAtual();

      if (response.success && response.data) {
        const data = response.data;

        document.getElementById('f-periodo-inicio').textContent = formatDate(data.periodo_inicio);
        document.getElementById('f-periodo-fim').textContent = formatDate(data.periodo_fim);

        document.getElementById('f-servicos').textContent = formatBRL(data.total_servicos);
        document.getElementById('f-qtd-os').textContent = data.quantidade_os;
        document.getElementById('f-ticket-servico').textContent = formatBRL(data.ticket_medio_servico);

        document.getElementById('f-vendas').textContent = formatBRL(data.total_vendas);
        document.getElementById('f-qtd-vendas').textContent = data.quantidade_vendas;
        document.getElementById('f-ticket-venda').textContent = formatBRL(data.ticket_medio_venda);

        if (Array.isArray(data.formas_pagamento)) {
          const pix = data.formas_pagamento.find(fp => fp.forma_pagamento === 'PIX')?.total || 0;
          const cartao = data.formas_pagamento.find(fp => fp.forma_pagamento === 'CARTAO_CREDITO')?.total || 0;
          const dinheiro = data.formas_pagamento.find(fp => fp.forma_pagamento === 'DINHEIRO')?.total || 0;

          document.getElementById('f-pgto-pix').textContent = formatBRL(pix);
          document.getElementById('f-pgto-cartao').textContent = formatBRL(cartao);
          document.getElementById('f-pgto-dinheiro').textContent = formatBRL(dinheiro);
        }

        document.getElementById('f-valor-perdas').textContent = formatBRL(data.valor_estoque_perdido);
        document.getElementById('f-qtd-perdas').textContent = data.quantidade_perdas;

        document.getElementById('f-bruto').textContent = formatBRL(data.faturamento_bruto);
        document.getElementById('f-despesas').textContent = formatBRL(data.total_despesas);
        document.getElementById('f-margem').textContent = `${(data.margem_lucro || 0).toFixed(2)}%`;
        document.getElementById('f-lucro').textContent = formatBRL(data.lucro_liquido);
      } else {
        showToast('Não foi possível carregar os dados do fechamento', 'warning');
      }
    } catch (err) {
      showToast('Erro ao carregar fechamento', 'error');
    }
  }

  // Alternância dos botões Dia / Mês
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

  // GET: Busca as despesas reais do backend
  async function carregarContas() {
    try {
      const response = await DespesaService.listar();

      if (response.success && Array.isArray(response.data)) {
        despesas = response.data;
      } else {
        showToast(response.message || 'Erro ao processar dados das despesas', 'warning');
      }
    } catch (err) {
      console.error('Erro detalhado no carregarContas:', err);
      showToast('Erro ao carregar despesas', 'error');
    }
    renderContas();
  }

  // Renderiza a lista de despesas na tela
  function renderContas() {
    const listContainer = document.getElementById('contas-list');
    if (!listContainer) return;
    
    if (despesas.length === 0) {
      listContainer.innerHTML = `<p class="text-muted" style="font-size:0.85rem; padding: 8px 0;">Nenhuma despesa cadastrada.</p>`;
      return;
    }

    listContainer.innerHTML = despesas
      .map((c) => {
        const isPago = c.status === 'PAGO';
        const diaVencimento = c.data_vencimento 
          ? new Date(c.data_vencimento).getUTCDate() 
          : '--';

        return `
          <div class="card list-item mb-2 ${isPago ? 'paid' : ''}" data-id="${c.id}" style="cursor: pointer; display: flex; justify-content: space-between; align-items: center; padding: 12px 16px;">
            <div class="info">
              <div class="title" style="font-weight: 600; ${isPago ? 'text-decoration: line-through; opacity: 0.6;' : ''}">
                ${c.descricao}
              </div>
              <div class="sub" style="font-size: 0.8rem; color: #888; margin-top: 2px;">
                Vence dia ${diaVencimento} ${isPago ? '• <strong style="color: #4caf50;">PAGO</strong>' : ''}
              </div>
            </div>
            <div class="fw-bold" style="${isPago ? 'text-decoration: line-through; opacity: 0.6;' : ''}">
              ${formatBRL(c.valor)}
            </div>
          </div>
        `;
      })
      .join('');

    // PUT: Paga a despesa ao clicar no card
    listContainer.querySelectorAll('.list-item').forEach((card) => {
      card.onclick = async () => {
        const id = card.dataset.id;
        const despesa = despesas.find((item) => String(item.id) === String(id));

        if (!despesa) return;

        if (despesa.status === 'PAGO') {
          showToast('Esta despesa já foi paga.', 'info');
          return;
        }

        const confirmou = confirm(`Pagar ${despesa.descricao}?`);
        if (!confirmou) return;

        const payload = {
          ...despesa,
          status: 'PAGO',
          valor_pago: despesa.valor,
          data_pagamento: new Date().toISOString()
        };

        try {
          const response = await DespesaService.atualizar(despesa.id, payload);
          if (response.success) {
            showToast(`"${despesa.descricao}" marcada como PAGA!`, 'success');

            if (NotificationService && typeof NotificationService.notificar === 'function') {
              NotificationService.notificar(
                'Despesa Paga',
                `A conta "${despesa.descricao}" de ${formatBRL(despesa.valor)} foi marcada como paga.`
              );
            }

            // Atualiza tela e recarrega fechamento
            await carregarContas();
            const btnMesAtivo = document.getElementById('btn-mes').classList.contains('btn-primary');
            await loadFechamento(btnMesAtivo ? 'mes' : 'dia');
          } else {
            showToast(response.message || 'Erro ao atualizar despesa', 'error');
          }
        } catch (err) {
          showToast('Erro ao comunicar com o servidor', 'error');
        }
      };
    });
  }

  // POST: Envia nova despesa para o backend ao salvar
  document.getElementById('form-conta').onsubmit = async (e) => {
    e.preventDefault();

    const descricao = document.getElementById('conta-nome').value.trim();
    const valor = parseFloat(document.getElementById('conta-valor').value) || 0;
    const diaVencimento = parseInt(document.getElementById('conta-dia').value) || 1;

    const hoje = new Date();
    const dataVencimento = new Date(Date.UTC(hoje.getFullYear(), hoje.getMonth(), diaVencimento)).toISOString();

    const payload = {
      descricao,
      categoria: 'OUTROS',
      valor,
      data_vencimento: dataVencimento,
      status: 'PENDENTE'
    };

    try {
      const response = await DespesaService.criar(payload);
      if (response.success) {
        showToast('Despesa cadastrada com sucesso', 'success');

        if (NotificationService && typeof NotificationService.notificar === 'function') {
          NotificationService.notificar(
            'Nova Despesa Cadastrada',
            `Lembrete: "${descricao}" no valor de ${formatBRL(valor)} foi cadastrado.`
          );
        }

        document.getElementById('form-conta').reset();
        await carregarContas();

        const btnMesAtivo = document.getElementById('btn-mes').classList.contains('btn-primary');
        await loadFechamento(btnMesAtivo ? 'mes' : 'dia');
      } else {
        showToast(response.message || 'Erro ao criar despesa', 'error');
      }
    } catch (err) {
      showToast('Erro ao comunicar com o servidor', 'error');
    }
  };

  // Carregamento inicial das duas abas ao entrar na página
  await loadFechamento('mes');
  await carregarContas();
}

init();