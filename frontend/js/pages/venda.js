import { requireAuth, goBack, showToast, formatBRL, getCurrentUser } from '../app.js';
import { SaleService } from '../../client/services/VendasService.js';
import { ClientService } from '../../client/services/clientService.js';
import { ProductService } from '../../client/services/ProdutoService.js';
import { NotificationService } from '../../client/services/notificationService.js';

window.NotificationService = NotificationService;

async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();

  let produtosCache = [];
  let clientesCache = [];

  const selectProduto = document.getElementById('venda-produto');
  const inputQtd = document.getElementById('venda-qtd');
  const inputPreco = document.getElementById('venda-preco');
  const inputCliente = document.getElementById('venda-cliente');
  const hiddenClienteId = document.getElementById('venda-cliente-id');
  const selectPagamento = document.getElementById('forma-pag');

  // ===== Carregar produtos =====
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

  // ===== Carregar clientes =====
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

  inputCliente.addEventListener('input', () => {
    const nome = inputCliente.value.trim();
    const encontrado = clientesCache.find(
      (c) => c.nome.toLowerCase() === nome.toLowerCase()
    );
    hiddenClienteId.value = encontrado ? encontrado.id : 0;
  });

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

  // ===== Submit =====
  document.getElementById('form-venda').onsubmit = async (e) => {
    e.preventDefault();

    const produtoId = Number(selectProduto.value);
    const quantidade = parseInt(inputQtd.value) || 1;
    const tipoPagamento = selectPagamento.value || 'PIX';
    const clienteId = Number(hiddenClienteId.value) || 0;
    const total = calc();

    // Pega o usuário logado (já salvo no sessionStorage)
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
        setTimeout(() => goBack(), 800);
      } else {
        showToast(response.message || 'Erro ao registrar venda', 'error');
      }
    } catch (err) {
      showToast(err.message || 'Erro de conexão', 'error');
    }
  };

  await Promise.all([carregarProdutos(), carregarClientes()]);
  calc();
}

init();