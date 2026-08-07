import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { ProductService } from '../../client/services/ProdutoService.js';
import { NotificationService } from '../../client/services/notificationService.js';

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
  let historico = [];
  let categoriaAtual = null;

  // Tabs
  document.querySelectorAll('#est-tabs .tab').forEach((tab) => {
    tab.onclick = () => {
      document.querySelectorAll('#est-tabs .tab').forEach((t) => t.classList.remove('active'));
      tab.classList.add('active');
      ['produtos', 'entrada', 'brindes'].forEach((t) => {
        document.getElementById('tab-' + t).classList.toggle('hidden', tab.dataset.tab !== t);
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
    document.getElementById('cat-grid').style.display = 'grid';
    document.getElementById('est-list-header').style.display = 'none';
    document.getElementById('est-list').innerHTML =
      `<div class="empty-state"><div class="icon">📦</div><p>Selecione uma categoria acima.</p></div>`;
  }

  function renderLista(lista, titulo) {
    const el = document.getElementById('est-list');
    document.getElementById('cat-title').textContent = titulo;
    document.getElementById('est-list-header').style.display = 'block';
    document.getElementById('cat-grid').style.display = 'none';

    if (!lista.length) {
      el.innerHTML = `<div class="empty-state"><div class="icon">📦</div><p>Nenhum produto nesta categoria.</p></div>`;
      return;
    }

    el.innerHTML = lista
      .map(
        (p) => `
      <div class="card list-item">
        <div class="info">
          <div class="title">${p.nome}</div>
          <div class="sub">Estoque: ${p.estoque}</div>
        </div>
        <div class="fw-bold text-primary">${formatBRL(p.preco)}</div>
      </div>
      `
      )
      .join('');
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

  // ===== BRINDES (local) =====
  document.getElementById('form-brinde').onsubmit = (e) => {
    e.preventDefault();
    const item = {
      produto: document.getElementById('bri-produto').value.trim(),
      qtd: parseInt(document.getElementById('bri-qtd').value) || 1,
      motivo: document.getElementById('bri-motivo').value.trim() || 'Brinde',
      data: new Date().toLocaleString('pt-BR'),
    };

    if (!item.produto) {
      showToast('Informe o produto', 'error');
      return;
    }

    historico.unshift(item);
    document.getElementById('bri-hist').innerHTML = historico
      .map(
        (h) => `
        <div class="list-item">
          <div class="info">
            <div class="title">${h.produto} ×${h.qtd}</div>
            <div class="sub">${h.motivo} · ${h.data}</div>
          </div>
        </div>
      `
      )
      .join('');

    showToast('Saída de brinde registrada', 'success');
    document.getElementById('form-brinde').reset();
    document.getElementById('bri-qtd').value = 1;
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