import { requireAuth, goBack, showToast, formatBRL, getCurrentUser } from '../app.js';
import { SaleService } from '../../client/services/VendasService.js';
import { ClientService } from '../../client/services/clientService.js';
import { ProductService } from '../../client/services/ProdutoService.js';
import { NotificationService } from '../../client/services/notificationService.js';
import { openInspect } from '../components/InspectModal.js';

window.NotificationService = NotificationService;

/**
 * Página de Vendas
 * - Aba Nova Venda
 * - Aba Vendas do Dia (GET /sales/dia)
 */
async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();

  let produtosCache = [];
  let clientesCache = [];
  let vendasHoje = [];

  // ===== TABS =====
  document.querySelectorAll('#venda-tabs .tab').forEach((tab) => {
    tab.onclick = () => {
      document.querySelectorAll('#venda-tabs .tab').forEach((t) => t.classList.remove('active'));
      tab.classList.add('active');

      const target = tab.dataset.tab;
      document.getElementById('tab-nova').classList.toggle('hidden', target !== 'nova');
      document.getElementById('tab-hoje').classList.toggle('hidden', target !== 'hoje');

      if (target === 'hoje') {
        carregarVendasDoDia();
      }
    };
  });

  // ===== ELEMENTOS =====
  const selectProduto = document.getElementById('venda-produto');
  const inputQtd = document.getElementById('venda-qtd');
  const inputPreco = document.getElementById('venda-preco');
  const inputCliente = document.getElementById('venda-cliente');
  const hiddenClienteId = document.getElementById('venda-cliente-id');
  const selectPagamento = document.getElementById('forma-pag');

  // ===== CARREGAR PRODUTOS =====
  async function carregarProdutos() {
    try {
      const res = await ProductService.listAll();
      if (res.success && Array.isArray(res.data)) {
        produtosCache = res.data;
        selectProduto.innerHTML =
          `<option value="">Selecione um produto</option>` +
          produtosCache
            .map(
              (p) =>
                `<option value="${p.id}" data-preco="${p.preco}">${p.nome} — ${formatBRL(p.preco)} (est: ${p.estoque})</option>`
            )
            .join('');
      } else {
        selectProduto.innerHTML = `<option value="">Nenhum produto encontrado</option>`;
      }
    } catch {
      selectProduto.innerHTML = `<option value="">Erro ao carregar produtos</option>`;
      showToast('Erro ao carregar produtos', 'error');
    }
  }

  // ===== CARREGAR CLIENTES =====
  async function carregarClientes() {
    try {
      const res = await ClientService.listAll();
      if (res.success && Array.isArray(res.data)) {
        clientesCache = res.data;
        const dl = document.getElementById('lista-clientes');
        if (dl) {
          dl.innerHTML = clientesCache
            .map((c) => `<option value="${c.nome}" data-id="${c.id}">`)
            .join('');
        }
      }
    } catch {
      // silencioso
    }
  }

  // Cliente opcional
  inputCliente.addEventListener('input', () => {
    const nome = inputCliente.value.trim();
    const encontrado = clientesCache.find(
      (c) => c.nome.toLowerCase() === nome.toLowerCase()
    );
    hiddenClienteId.value = encontrado ? encontrado.id : 0;
  });

  // Preço automático
  selectProduto.addEventListener('change', () => {
    const option = selectProduto.selectedOptions[0];
    const preco = option ? parseFloat(option.dataset.preco) || 0 : 0;
    inputPreco.value = preco.toFixed(2);
    calc();
  });

  inputQtd.addEventListener('input', calc);

  function calc() {
    const qtd = parseInt(inputQtd.value) || 1;
    const preco = parseFloat(inputPreco.value) || 0;
    const total = qtd * preco;

    document.getElementById('r-qtd').textContent = qtd;
    document.getElementById('r-preco').textContent = formatBRL(preco);
    document.getElementById('r-total').textContent = formatBRL(total);

    return total;
  }

  // ===== SUBMIT NOVA VENDA =====
  document.getElementById('form-venda').onsubmit = async (e) => {
    e.preventDefault();

    const produtoId = Number(selectProduto.value);
    const quantidade = parseInt(inputQtd.value) || 1;
    const tipoPagamento = selectPagamento.value || 'PIX';
    const clienteId = Number(hiddenClienteId.value) || 0;
    const total = calc();

    const currentUser = getCurrentUser();
    const usuarioId = Number(currentUser?.id) || 0;

    if (!produtoId) {
      showToast('Selecione um produto', 'error');
      return;
    }
    if (quantidade < 1) {
      showToast('Quantidade inválida', 'error');
      return;
    }
    if (!usuarioId) {
      showToast('Usuário não identificado. Faça login novamente.', 'error');
      return;
    }

    try {
      const response = await SaleService.create({
        produto_id: produtoId,
        quantidade,
        tipo_pagamento: tipoPagamento,
        usuario_id: usuarioId,
        cliente_id: clienteId,
        total,
      });

      if (response.success) {
        showToast('Venda registrada com sucesso!', 'success');
        // Limpa o formulário
        document.getElementById('form-venda').reset();
        inputQtd.value = 1;
        inputPreco.value = '';
        hiddenClienteId.value = 0;
        calc();
        // Atualiza a lista do dia (caso o usuário mude de aba)
        vendasHoje = [];
      } else {
        showToast(response.message || 'Erro ao registrar venda', 'error');
      }
    } catch (err) {
      showToast(err.message || 'Erro de conexão', 'error');
    }
  };


  function abrirComprovante(venda) {
    if (!venda) return;

    const data = venda.criado_em
      ? new Date(venda.criado_em).toLocaleString('pt-BR')
      : '—';

    const produto = produtosCache.find((p) => p.id === venda.produto_id);
    const nomeProduto = produto ? produto.nome : `Produto #${venda.produto_id}`;

    let nomeCliente = 'Consumidor Final';
    if (venda.cliente_id > 0) {
      const cli = clientesCache.find((c) => c.id === venda.cliente_id);
      nomeCliente = cli ? cli.nome : `Cliente #${venda.cliente_id}`;
    }

    openInspect({
      title: 'Comprovante de Venda',
      subtitle: 'IDroid Virtual',
      status: {
        label: venda.status === 'CANCELADA' ? 'CANCELADA' : 'CONCLUÍDA',
        color: venda.status === 'CANCELADA' ? '#c0392b' : '#27ae60',
      },
      fields: [
        { label: 'Venda', value: `#${venda.id}` },
        { label: 'Data', value: data },
        { label: 'Produto', value: nomeProduto },
        { label: 'Quantidade', value: venda.quantidade },
        { label: 'Cliente', value: nomeCliente },
        { label: 'Pagamento', value: venda.tipo_pagamento || '—' },
        { label: 'Total', value: formatBRL(venda.total), highlight: true },
      ],
      footer: 'IDroid Virtual · Gestão',
    });
  }

  // ===== VENDAS DO DIA =====
  async function carregarVendasDoDia() {
    const lista = document.getElementById('lista-vendas');
    lista.innerHTML = `<div class="empty-state"><div class="icon">⏳</div><p>Carregando...</p></div>`;

    try {
      const res = await SaleService.getByDay();
      if (res.success && Array.isArray(res.data)) {
        vendasHoje = res.data;
      } else {
        vendasHoje = [];
      }
    } catch {
      showToast('Erro ao carregar vendas do dia', 'error');
      vendasHoje = [];
    }

    renderVendasDoDia();
  }

    function renderVendasDoDia() {
    const lista = document.getElementById('lista-vendas');
    const totalEl = document.getElementById('total-dia');
    const qtdEl = document.getElementById('qtd-dia');

    const total = vendasHoje.reduce((s, v) => s + (Number(v.total) || 0), 0);
    totalEl.textContent = formatBRL(total);
    qtdEl.textContent = vendasHoje.length;

    if (!vendasHoje.length) {
      lista.innerHTML = `
        <div class="empty-state">
          <div class="icon">🛒</div>
          <p>Nenhuma venda hoje.</p>
        </div>`;
      return;
    }

    lista.innerHTML = vendasHoje
      .map((v) => {
        const statusClass = v.status === 'CANCELADA' ? 'text-muted' : '';
        const statusLabel = v.status === 'CANCELADA' ? ' (Cancelada)' : '';
        const data = v.criado_em
          ? new Date(v.criado_em).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
          : '';

        return `
          <div class="card list-item ${statusClass}" data-id="${v.id}" style="cursor:pointer">
            <div class="info">
              <div class="title">Venda #${v.id}${statusLabel}</div>
              <div class="sub">
                Qtd: ${v.quantidade} · ${v.tipo_pagamento || '—'}
                ${data ? ' · ' + data : ''}
              </div>
            </div>
            <div style="text-align:right">
              <div class="fw-bold text-primary">${formatBRL(v.total)}</div>
              ${
                v.status !== 'CANCELADA'
                  ? `<button class="btn btn-ghost btn-sm btn-cancelar" style="color:var(--danger);margin-top:4px">Cancelar</button>`
                  : ''
              }
            </div>
          </div>
        `;
      })
      .join('');

    // Clique no card → abre comprovante
    lista.querySelectorAll('.list-item').forEach((card) => {
      card.onclick = (e) => {
        // Não abre se clicou no botão cancelar
        if (e.target.closest('.btn-cancelar')) return;

        const id = Number(card.dataset.id);
        const venda = vendasHoje.find((v) => v.id === id);
        if (venda) abrirComprovante(venda);
      };
    });

    // Botões de cancelar
    lista.querySelectorAll('.btn-cancelar').forEach((btn) => {
      btn.onclick = async (e) => {
        e.stopPropagation();
        const card = btn.closest('[data-id]');
        const id = Number(card.dataset.id);
        if (!id) return;

        if (!confirm('Deseja realmente cancelar esta venda?')) return;

        const currentUser = getCurrentUser();
        const usuarioId = Number(currentUser?.id) || 0;
        if (!usuarioId) {
          showToast('Usuário não identificado', 'error');
          return;
        }

        try {
          const res = await SaleService.remove(id, usuarioId, 'Cancelamento pelo app');
          if (res.success) {
            showToast('Venda cancelada', 'success');
            await carregarVendasDoDia();
          } else {
            showToast(res.message || 'Erro ao cancelar', 'error');
          }
        } catch (err) {
          showToast(err.message || 'Erro ao cancelar venda', 'error');
        }
      };
    });
  }

  // Inicialização
  await Promise.all([carregarProdutos(), carregarClientes()]);
  calc();
}

init();