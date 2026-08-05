import { requireAuth, goBack, showToast } from '../app.js';

if (!requireAuth()) throw new Error('redirect');
document.getElementById('btn-back').onclick = () => goBack();

let entregas = [];

document.querySelectorAll('#ent-tabs .tab').forEach(tab => {
  tab.onclick = () => {
    document.querySelectorAll('#ent-tabs .tab').forEach(t => t.classList.remove('active'));
    tab.classList.add('active');
    document.getElementById('tab-lista').classList.toggle('hidden', tab.dataset.tab !== 'lista');
    document.getElementById('tab-nova').classList.toggle('hidden', tab.dataset.tab !== 'nova');
  };
});

function badgeClass(s) {
  if (s === 'ENTREGUE') return 'badge-entregue';
  if (s === 'EM_TRANSPORTE') return 'badge-em-transporte';
  return 'badge-pendente';
}

function render() {
  const el = document.getElementById('ent-list');
  if (!entregas.length) {
    el.innerHTML = `<div class="empty-state"><div class="icon">🛵</div><p>Nenhuma entrega em andamento.</p></div>`;
    return;
  }
  el.innerHTML = entregas.map((e, i) => `
    <div class="card">
      <div class="flex justify-between items-center mb-1">
        <span class="fw-bold">${e.cliente}</span>
        <span class="badge ${badgeClass(e.status)}">${e.status.replace('_', ' ')}</span>
      </div>
      <div class="text-muted" style="font-size:0.85rem">${e.produto}</div>
      <div class="text-muted" style="font-size:0.85rem">Motoboy: ${e.motoboy}</div>
      ${e.status !== 'ENTREGUE' ? `
        <button class="btn btn-primary btn-sm mt-1" data-idx="${i}" data-action="concluir">Confirmar entrega</button>
      ` : ''}
    </div>
  `).join('');

  el.querySelectorAll('[data-action="concluir"]').forEach(btn => {
    btn.onclick = () => {
      const idx = +btn.dataset.idx;
      entregas[idx].status = 'ENTREGUE';
      // CHAMADA DA API AQUI
      showToast('Entrega confirmada!', 'success');
      render();
    };
  });
}

document.getElementById('form-entrega').onsubmit = (e) => {
  e.preventDefault();
  const item = {
    id: Date.now(),
    cliente: document.getElementById('ent-cliente').value.trim(),
    produto: document.getElementById('ent-produto').value.trim(),
    motoboy: document.getElementById('ent-motoboy').value.trim(),
    status: document.getElementById('ent-status').value
  };
  // CHAMADA DA API AQUI
  entregas.unshift(item);
  showToast('Entrega registrada', 'success');
  document.getElementById('form-entrega').reset();
  document.querySelector('[data-tab="lista"]').click();
  render();
};

render();
