import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { bindMasks, applyMask } from '../masks.js';

if (!requireAuth()) throw new Error('redirect');
document.getElementById('btn-back').onclick = () => goBack();

bindMasks({ 'rapido-tel': 'phone', 'trade-imei': 'imei' });

let itens = [];
let useDesc = true, useJuros = false, useTrade = false, useEntrega = false, useBrinde = false;

function setupToggle(id, flagSetter, extra) {
  const el = document.getElementById(id);
  el.onclick = () => {
    const on = !el.classList.contains('on');
    el.classList.toggle('on', on);
    el.setAttribute('aria-checked', on);
    flagSetter(on);
    if (extra) extra(on);
    calc();
  };
}
setupToggle('toggle-desc', v => useDesc = v);
setupToggle('toggle-juros', v => useJuros = v, on => document.getElementById('juros-group').classList.toggle('hidden', !on));
setupToggle('toggle-trade', v => useTrade = v, on => document.getElementById('trade-fields').classList.toggle('hidden', !on));
setupToggle('toggle-entrega', v => useEntrega = v, on => document.getElementById('entrega-fields').classList.toggle('hidden', !on));
setupToggle('toggle-brinde', v => useBrinde = v);

document.getElementById('btn-cadastro-rapido').onclick = () => {
  document.getElementById('cadastro-rapido').classList.toggle('hidden');
};

function addItem(data = {}) {
  const id = Date.now() + Math.random();
  itens.push({ id, tipo: data.tipo || 'PRODUTO', nome: data.nome || '', qtd: data.qtd || 1, valor: data.valor || 0, imei: data.imei || '' });
  renderItens();
  calc();
}

function renderItens() {
  const el = document.getElementById('itens-venda');
  if (!itens.length) {
    el.innerHTML = '<p class="text-muted" style="font-size:0.85rem">Nenhum item adicionado.</p>';
    return;
  }
  el.innerHTML = itens.map((it) => `
    <div class="card" style="padding:12px;margin-bottom:8px;background:#fafafa" data-id="${it.id}">
      <div class="form-group" style="margin-bottom:8px">
        <label>Tipo</label>
        <select class="item-tipo">
          <option value="PRODUTO" ${it.tipo === 'PRODUTO' ? 'selected' : ''}>Produto</option>
          <option value="SERVICO" ${it.tipo === 'SERVICO' ? 'selected' : ''}>Serviço</option>
        </select>
      </div>
      <div class="form-group" style="margin-bottom:8px">
        <label>Descrição</label>
        <input type="text" class="item-nome" value="${it.nome}" placeholder="Nome do produto/serviço">
      </div>
      <div class="flex gap-1">
        <div class="form-group" style="flex:1;margin-bottom:8px">
          <label>Qtd</label>
          <input type="number" class="item-qtd" value="${it.qtd}" min="1" inputmode="numeric">
        </div>
        <div class="form-group" style="flex:2;margin-bottom:8px">
          <label>Valor unit. (R$)</label>
          <input type="number" class="item-valor" value="${it.valor}" step="0.01" min="0" inputmode="decimal">
        </div>
      </div>
      <div class="form-group item-imei-wrap" style="margin-bottom:8px;${it.tipo === 'PRODUTO' ? '' : 'display:none'}">
        <label>IMEI (celular)</label>
        <input type="text" class="item-imei" value="${it.imei}" maxlength="15" inputmode="numeric" placeholder="Obrigatório se celular">
      </div>
      <button type="button" class="btn btn-ghost btn-sm item-remove" style="color:var(--danger)">Remover</button>
    </div>
  `).join('');

  el.querySelectorAll('[data-id]').forEach(card => {
    const id = +card.dataset.id;
    const it = itens.find(x => x.id === id);
    card.querySelector('.item-tipo').onchange = (e) => {
      it.tipo = e.target.value;
      card.querySelector('.item-imei-wrap').style.display = it.tipo === 'PRODUTO' ? '' : 'none';
    };
    card.querySelector('.item-nome').oninput = (e) => { it.nome = e.target.value; };
    card.querySelector('.item-qtd').oninput = (e) => { it.qtd = parseInt(e.target.value) || 1; calc(); };
    card.querySelector('.item-valor').oninput = (e) => { it.valor = parseFloat(e.target.value) || 0; calc(); };
    const imeiInput = card.querySelector('.item-imei');
    applyMask(imeiInput, 'imei');
    imeiInput.oninput = (e) => { it.imei = e.target.value; };
    card.querySelector('.item-remove').onclick = () => {
      itens = itens.filter(x => x.id !== id);
      renderItens();
      calc();
    };
  });
}

document.getElementById('btn-add-item').onclick = () => addItem();

function calc() {
  const subtotal = itens.reduce((s, it) => s + (it.qtd * it.valor), 0);
  const desc = useDesc ? subtotal * 0.10 : 0;
  const trade = useTrade ? (parseFloat(document.getElementById('trade-valor').value) || 0) : 0;
  let base = Math.max(0, subtotal - desc - trade);
  const taxa = useJuros ? (parseFloat(document.getElementById('juros-pct').value) || 0) / 100 : 0;
  const juros = base * taxa;
  const total = base + juros;

  document.getElementById('r-subtotal').textContent = formatBRL(subtotal);
  document.getElementById('r-desc').textContent = formatBRL(desc);
  document.getElementById('r-trade').textContent = formatBRL(trade);
  document.getElementById('r-juros').textContent = formatBRL(juros);
  document.getElementById('r-total').textContent = formatBRL(total);
  return { subtotal, desc, trade, juros, total };
}
document.getElementById('trade-valor').addEventListener('input', calc);
document.getElementById('juros-pct').addEventListener('input', calc);

document.getElementById('form-venda').onsubmit = async (e) => {
  e.preventDefault();
  if (!itens.length) {
    showToast('Adicione ao menos um item', 'error');
    return;
  }
  for (const it of itens) {
    if (!it.nome) { showToast('Preencha a descrição dos itens', 'error'); return; }
    const dig = (it.imei || '').replace(/\D/g, '');
    if (it.tipo === 'PRODUTO' && dig.length > 0 && dig.length < 14) {
      showToast('IMEI inválido em um dos itens', 'error'); return;
    }
  }
  const totals = calc();

  // CHAMADA DA API AQUI
  showToast('Venda registrada com sucesso!', 'success');
  setTimeout(() => goBack(), 800);
};

// CHAMADA DA API AQUI — popular datalist de clientes
addItem();
