import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { bindMasks } from '../masks.js';
import { Produto } from '../../client/models/ProdutoModel.js';

if (!requireAuth()) throw new Error('redirect');
document.getElementById('btn-back').onclick = () => goBack();

bindMasks({ 'ent-imei': 'imei' });

let produtos = [];
let historico = [];

document.querySelectorAll('#est-tabs .tab').forEach(tab => {
  tab.onclick = () => {
    document.querySelectorAll('#est-tabs .tab').forEach(t => t.classList.remove('active'));
    tab.classList.add('active');
    ['produtos', 'entrada', 'brindes'].forEach(t => {
      document.getElementById('tab-' + t).classList.toggle('hidden', tab.dataset.tab !== t);
    });
  };
});

function renderProdutos() {
  const el = document.getElementById('est-list');
  if (!produtos.length) {
    el.innerHTML = `<div class="empty-state"><div class="icon">📦</div><p>Estoque vazio.</p></div>`;
    return;
  }
  el.innerHTML = produtos.map(p => `
    <div class="card list-item">
      <div class="info">
        <div class="title">${p.nome}</div>
        <div class="sub">${p.categoria} · Estoque: ${p.estoque} ${p.imei ? '· IMEI ' + p.imei : ''}</div>
      </div>
      <div class="fw-bold text-primary">${formatBRL(p.preco)}</div>
    </div>
  `).join('');
}

document.getElementById('form-entrada').onsubmit = async (e) => {
  e.preventDefault();
  const cat = document.getElementById('ent-categoria').value;
  const imei = document.getElementById('ent-imei').value.replace(/\D/g, '');
  if (cat === 'CELULAR' && imei.length < 14) {
    showToast('IMEI obrigatório para celular', 'error');
    return;
  }
  const model = new Produto({
    nome: document.getElementById('ent-nome').value,
    categoria: cat,
    preco: parseFloat(document.getElementById('ent-preco').value) || 0,
    estoque: parseInt(document.getElementById('ent-qtd').value) || 1
  });
  if (!model.isValid()) {
    showToast(model.validar()[0], 'error');
    return;
  }

  // CHAMADA DA API AQUI
  produtos.unshift({ ...model, id: Date.now(), imei: imei || null });
  showToast('Entrada registrada', 'success');
  document.getElementById('form-entrada').reset();
  document.querySelector('[data-tab="produtos"]').click();
  renderProdutos();
};

document.getElementById('form-brinde').onsubmit = async (e) => {
  e.preventDefault();
  const item = {
    produto: document.getElementById('bri-produto').value.trim(),
    qtd: parseInt(document.getElementById('bri-qtd').value) || 1,
    motivo: document.getElementById('bri-motivo').value.trim() || 'Brinde',
    data: new Date().toLocaleString('pt-BR')
  };
  // CHAMADA DA API AQUI
  historico.unshift(item);
  document.getElementById('bri-hist').innerHTML = historico.map(h =>
    `<div class="list-item"><div class="info"><div class="title">${h.produto} ×${h.qtd}</div><div class="sub">${h.motivo} · ${h.data}</div></div></div>`
  ).join('');
  showToast('Saída de brinde registrada', 'success');
  document.getElementById('form-brinde').reset();
};

// CHAMADA DA API AQUI
renderProdutos();
