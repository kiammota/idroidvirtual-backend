import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { ProductService } from '../../client/services/ProdutoService.js';
import { NotificationService } from '../../client/services/notificationService.js';
import { openInspect } from '../components/InspectModal.js';

window.NotificationService = NotificationService;

const CATEGORIAS = [
  { key: 'CELULAR', label: 'Celulares', icon: '📱' },
  { key: 'ACESSORIO', label: 'Acessórios', icon: '🎧' },
  { key: 'PECA', label: 'Peças', icon: '🔧' },
  { key: 'OUTRO', label: 'Outros', icon: '📦' },
];

async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();

  let produtos = [];
  let categoriaAtual = null;
  let listaAtual = [];

  // Tabs (apenas Produtos e Entrada)
  document.querySelectorAll('#est-tabs .tab').forEach((tab) => {
    tab.onclick = () => {
      document.querySelectorAll('#est-tabs .tab').forEach((t) => t.classList.remove('active'));
      tab.classList.add('active');
      ['produtos', 'entrada'].forEach((t) => {
        const el = document.getElementById('tab-' + t);
        if (el) el.classList.toggle('hidden', tab.dataset.tab !== t);
      });
    };
  });

  // ===== CATEGORIAS =====
  function atualizarContadores() {
    const counts = { CELULAR: 0, ACESSORIO: 0, PECA: 0, OUTRO: 0 };
    produtos.forEach((p) => {
      if (counts[p.categoria] !== undefined) counts[p.categoria]++;
    });
    Object.keys(counts).forEach((cat) => {
      const el = document.getElementById(`count-${cat}`);
      if (el) el.textContent = `${counts[cat]} item${counts[cat] !== 1 ? 's' : ''}`;
    });
  }

  function mostrarCategorias() {
    categoriaAtual = null;
    listaAtual = [];
    document.getElementById('cat-grid').style.display = 'grid';
    document.getElementById('est-list-header').style.display = 'none';
    document.getElementById('est-list').innerHTML =
      `<div class="empty-state"><div class="icon">📦</div><p>Selecione uma categoria acima.</p></div>`;
  }

  function abrirInspectProduto(p) {
    const criado = p.criado_em
      ? new Date(p.criado_em).toLocaleString('pt-BR')
      : '—';
    const atualizado = p.atualizado_em
      ? new Date(p.atualizado_em).toLocaleString('pt-BR')
      : '—';

    openInspect({
      title: 'Produto',
      subtitle: 'IDroid Virtual',
      fields: [
        { label: 'ID', value: `#${p.id}` },
        { label: 'Nome', value: p.nome || '—' },
        { label: 'Categoria', value: p.categoria || '—' },
        { label: 'Preço', value: formatBRL(p.preco), highlight: true },
        { label: 'Estoque', value: String(p.estoque ?? 0) },
        { label: 'Cadastrado em', value: criado },
        { label: 'Atualizado em', value: atualizado },
      ],
      footer: 'IDroid Virtual · Gestão',
    });
  }

  async function aplicarEstoque(produto, novoValor, inputEl) {
    const atual = Number(produto.estoque) || 0;
    const num = Number(novoValor);

    if (Number.isNaN(num) || num < 0) {
      showToast('Valor de estoque inválido', 'error');
      if (inputEl) inputEl.value = atual;
      return;
    }

    if (num === atual) {
      showToast('Estoque já está nesse valor', 'info');
      if (inputEl) inputEl.value = atual;
      return;
    }

    const acao = num > atual ? 'aumentar' : 'reduzir';
    const nome = produto.nome || 'produto';
    const confirmar = confirm(
      `Você tem certeza que quer ${acao} o ${nome}?\n\nEstoque atual: ${atual}\nNovo estoque: ${num}`
    );
    if (!confirmar) {
      if (inputEl) inputEl.value = atual;
      return;
    }

    try {
      const response = await ProductService.setStock(produto.id, num);
      if (response.success) {
        const novoEstoque =
          response.data && response.data.estoque != null
            ? Number(response.data.estoque)
            : num;

        // Atualiza na lista local e na lista global
        produto.estoque = novoEstoque;
        const idxGlobal = produtos.findIndex((x) => Number(x.id) === Number(produto.id));
        if (idxGlobal >= 0) produtos[idxGlobal].estoque = novoEstoque;

        // Atualiza o texto na UI sem re-render completo
        const card = document.querySelector(`.list-item[data-id="${produto.id}"]`);
        if (card) {
          const sub = card.querySelector('.sub');
          if (sub) sub.textContent = `Estoque: ${novoEstoque}`;
          const input = card.querySelector('.stock-input');
          if (input) input.value = novoEstoque;
        }

        showToast(
          `Estoque ${acao === 'aumentar' ? 'aumentado' : 'reduzido'} com sucesso`,
          'success'
        );
      } else {
        showToast(response.message || 'Erro ao atualizar estoque', 'error');
        if (inputEl) inputEl.value = atual;
      }
    } catch (err) {
      showToast(err.message || 'Erro de conexão', 'error');
      if (inputEl) inputEl.value = atual;
    }
  }

  function renderLista(lista, titulo) {
    const el = document.getElementById('est-list');
    document.getElementById('cat-title').textContent = titulo;
    document.getElementById('est-list-header').style.display = 'flex';
    document.getElementById('cat-grid').style.display = 'none';
    listaAtual = lista || [];

    if (!lista.length) {
      el.innerHTML = `<div class="empty-state"><div class="icon">📦</div><p>Nenhum produto nesta categoria.</p></div>`;
      return;
    }

    el.innerHTML = lista
      .map(
        (p) => `
      <div class="card list-item" data-id="${p.id}" style="cursor:pointer">
        <div class="info" style="flex:1; min-width:0">
          <div class="title" title="${(p.nome || '').replace(/"/g, '&quot;')}">${p.nome || '—'}</div>
          <div class="sub">Estoque: ${p.estoque ?? 0}</div>
        </div>
        <div class="fw-bold text-primary" style="margin-right:8px">${formatBRL(p.preco)}</div>
        <div class="stock-controls" style="display:flex;align-items:center;gap:4px;flex-shrink:0" onclick="event.stopPropagation()">
          <button type="button" class="btn btn-outline btn-sm stock-btn-minus" data-id="${p.id}" style="padding:4px 10px;min-width:32px" aria-label="Diminuir">−</button>
          <input type="number" class="stock-input" data-id="${p.id}" value="${p.estoque ?? 0}" min="0" inputmode="numeric" style="width:56px;padding:6px 4px;text-align:center;font-size:0.9rem" />
          <button type="button" class="btn btn-outline btn-sm stock-btn-plus" data-id="${p.id}" style="padding:4px 10px;min-width:32px" aria-label="Aumentar">+</button>
        </div>
      </div>
      `
      )
      .join('');

    // Clique no card (exceto controles) → Inspect
    el.querySelectorAll('.list-item').forEach((card) => {
      card.onclick = (e) => {
        if (e.target.closest('.stock-controls')) return;
        const id = Number(card.dataset.id);
        const produto = listaAtual.find((x) => Number(x.id) === id) ||
          produtos.find((x) => Number(x.id) === id);
        if (produto) abrirInspectProduto(produto);
      };
    });

    // Botões + / −
    el.querySelectorAll('.stock-btn-plus').forEach((btn) => {
      btn.onclick = (e) => {
        e.stopPropagation();
        const id = Number(btn.dataset.id);
        const produto = listaAtual.find((x) => Number(x.id) === id) ||
          produtos.find((x) => Number(x.id) === id);
        if (!produto) return;
        const input = el.querySelector(`.stock-input[data-id="${id}"]`);
        const atual = Number(input?.value ?? produto.estoque) || 0;
        const novo = atual + 1;
        if (input) input.value = novo;
        aplicarEstoque(produto, novo, input);
      };
    });

    el.querySelectorAll('.stock-btn-minus').forEach((btn) => {
      btn.onclick = (e) => {
        e.stopPropagation();
        const id = Number(btn.dataset.id);
        const produto = listaAtual.find((x) => Number(x.id) === id) ||
          produtos.find((x) => Number(x.id) === id);
        if (!produto) return;
        const input = el.querySelector(`.stock-input[data-id="${id}"]`);
        const atual = Number(input?.value ?? produto.estoque) || 0;
        const novo = Math.max(0, atual - 1);
        if (input) input.value = novo;
        aplicarEstoque(produto, novo, input);
      };
    });

    // Enter no input → aplica o valor digitado
    el.querySelectorAll('.stock-input').forEach((input) => {
      input.onclick = (e) => e.stopPropagation();
      input.onkeydown = (e) => {
        if (e.key === 'Enter') {
          e.preventDefault();
          e.stopPropagation();
          const id = Number(input.dataset.id);
          const produto = listaAtual.find((x) => Number(x.id) === id) ||
            produtos.find((x) => Number(x.id) === id);
          if (produto) aplicarEstoque(produto, input.value, input);
        }
      };
      input.onchange = (e) => {
        e.stopPropagation();
        const id = Number(input.dataset.id);
        const produto = listaAtual.find((x) => Number(x.id) === id) ||
          produtos.find((x) => Number(x.id) === id);
        if (produto) aplicarEstoque(produto, input.value, input);
      };
    });
  }

  // Clique nas caixinhas de categoria
  document.querySelectorAll('#cat-grid .grid-card').forEach((card) => {
    card.onclick = async () => {
      const cat = card.dataset.categoria;
      const meta = CATEGORIAS.find((c) => c.key === cat);
      categoriaAtual = cat;

      try {
        const response = await ProductService.listByCategory(cat);
        if (response.success && Array.isArray(response.data)) {
          renderLista(response.data, meta ? meta.label : cat);
        } else {
          renderLista([], meta ? meta.label : cat);
        }
      } catch (err) {
        showToast('Erro ao carregar categoria', 'error');
        renderLista([], meta ? meta.label : cat);
      }
    };
  });

  // Botão voltar para categorias
  document.getElementById('btn-voltar-cats').onclick = () => {
    mostrarCategorias();
  };

  // ===== ENTRADA =====
  document.getElementById('form-entrada').onsubmit = async (e) => {
    e.preventDefault();

    const nome = document.getElementById('ent-nome').value.trim();
    const categoria = document.getElementById('ent-categoria').value;
    const preco = parseFloat(document.getElementById('ent-preco').value) || 0;
    const estoque = parseInt(document.getElementById('ent-qtd').value) || 1;

    if (!nome) {
      showToast('Nome é obrigatório', 'error');
      return;
    }

    try {
      const response = await ProductService.create({ nome, categoria, preco, estoque });

      if (response.success) {
        showToast('Entrada registrada', 'success');
        document.getElementById('form-entrada').reset();
        document.getElementById('ent-qtd').value = 1;
        document.querySelector('[data-tab="produtos"]').click();
        await carregarProdutos();
        mostrarCategorias();
      } else {
        showToast(response.message || 'Erro ao registrar entrada', 'error');
      }
    } catch (err) {
      showToast(err.message || 'Erro de conexão', 'error');
    }
  };

  // ===== CARREGAR =====
  async function carregarProdutos() {
    try {
      const response = await ProductService.listAll();
      if (response.success && Array.isArray(response.data)) {
        produtos = response.data;
      } else {
        produtos = [];
      }
    } catch (err) {
      showToast('Erro ao carregar estoque', 'error');
      produtos = [];
    } finally {
      atualizarContadores();
    }
  }

  await carregarProdutos();
  mostrarCategorias();
}

init();
