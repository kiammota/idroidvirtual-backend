import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { bindMasks, applyMask } from '../masks.js';
import { SaleService } from '../../client/services/VendasService.js';
import { ClientService } from '../../client/services/clientService.js';
import { ProductService } from '../../client/services/ProdutoService.js';
import { TrocaService } from '../../client/services/TrocaService.js';
import { NotificationService } from '../../client/services/notificationService.js';

window.NotificationService = NotificationService;

/**
 * Página de Nova Venda
 */
async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();
  bindMasks({ 'rapido-tel': 'phone', 'trade-imei': 'imei' });

  let itens = [];
  let useDesc = true,
    useJuros = false,
    useTrade = false,
    useEntrega = false,
    useBrinde = false;
  let produtosCache = [];
  let clientesCache = [];

  function setupToggle(id, flagSetter, extra) {
    const el = document.getElementById(id);
    if (!el) return;
    el.onclick = () => {
      const on = !el.classList.contains('on');
      el.classList.toggle('on', on);
      el.setAttribute('aria-checked', on);
      flagSetter(on);
      if (extra) extra(on);
      calc();
    };
  }
  setupToggle('toggle-desc', (v) => (useDesc = v));
  setupToggle('toggle-juros', (v) => (useJuros = v), (on) =>
    document.getElementById('juros-group')?.classList.toggle('hidden', !on)
  );
  setupToggle('toggle-trade', (v) => (useTrade = v), (on) =>
    document.getElementById('trade-fields')?.classList.toggle('hidden', !on)
  );
  setupToggle('toggle-entrega', (v) => (useEntrega = v), (on) =>
    document.getElementById('entrega-fields')?.classList.toggle('hidden', !on)
  );
  setupToggle('toggle-brinde', (v) => (useBrinde = v));

  const btnCadastro = document.getElementById('btn-cadastro-rapido');
  if (btnCadastro) {
    btnCadastro.onclick = () => {
      document.getElementById('cadastro-rapido')?.classList.toggle('hidden');
    };
  }

  function addItem(data = {}) {
    const id = Date.now() + Math.random();
    itens.push({
      id,
      tipo: data.tipo || 'PRODUTO',
      nome: data.nome || '',
      qtd: data.qtd || 1,
      valor: data.valor || 0,
      imei: data.imei || '',
      produto_id: data.produto_id || 0,
    });
    renderItens();
    calc();
  }

  function renderItens() {
    const el = document.getElementById('itens-venda');
    if (!itens.length) {
      el.innerHTML = '<p class="text-muted" style="font-size:0.85rem">Nenhum item adicionado.</p>';
      return;
    }
    el.innerHTML = itens
      .map(
        (it) => `
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
        <input type="text" class="item-nome" value="${it.nome}" placeholder="Nome do produto/serviço" list="produtos-list">
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
  `
      )
      .join('');

    el.querySelectorAll('[data-id]').forEach((card) => {
      const id = +card.dataset.id;
      const it = itens.find((x) => x.id === id);
      card.querySelector('.item-tipo').onchange = (e) => {
        it.tipo = e.target.value;
        card.querySelector('.item-imei-wrap').style.display = it.tipo === 'PRODUTO' ? '' : 'none';
      };
      card.querySelector('.item-nome').oninput = (e) => {
        it.nome = e.target.value;
        const match = produtosCache.find((p) => p.nome === e.target.value);
        if (match) {
          it.produto_id = match.id;
          it.valor = match.preco;
          card.querySelector('.item-valor').value = match.preco;
          calc();
        }
      };
      card.querySelector('.item-qtd').oninput = (e) => {
        it.qtd = parseInt(e.target.value) || 1;
        calc();
      };
      card.querySelector('.item-valor').oninput = (e) => {
        it.valor = parseFloat(e.target.value) || 0;
        calc();
      };
      const imeiInput = card.querySelector('.item-imei');
      applyMask(imeiInput, 'imei');
      imeiInput.oninput = (e) => {
        it.imei = e.target.value;
      };
      card.querySelector('.item-remove').onclick = () => {
        itens = itens.filter((x) => x.id !== id);
        renderItens();
        calc();
      };
    });
  }

  document.getElementById('btn-add-item').onclick = () => addItem();

  function calc() {
    const subtotal = itens.reduce((s, it) => s + it.qtd * it.valor, 0);
    const desc = useDesc ? subtotal * 0.1 : 0;
    const trade = useTrade ? parseFloat(document.getElementById('trade-valor')?.value) || 0 : 0;
    let base = Math.max(0, subtotal - desc - trade);
    const taxa = useJuros ? (parseFloat(document.getElementById('juros-pct')?.value) || 0) / 100 : 0;
    const juros = base * taxa;
    const total = base + juros;

    const set = (id, val) => {
      const el = document.getElementById(id);
      if (el) el.textContent = formatBRL(val);
    };
    set('r-subtotal', subtotal);
    set('r-desc', desc);
    set('r-trade', trade);
    set('r-juros', juros);
    set('r-total', total);
    return { subtotal, desc, trade, juros, total };
  }

  document.getElementById('trade-valor')?.addEventListener('input', calc);
  document.getElementById('juros-pct')?.addEventListener('input', calc);

  document.getElementById('form-venda').onsubmit = async (e) => {
    e.preventDefault();
    if (!itens.length) {
      showToast('Adicione ao menos um item', 'error');
      return;
    }
    for (const it of itens) {
      if (!it.nome) {
        showToast('Preencha a descrição dos itens', 'error');
        return;
      }
      const dig = (it.imei || '').replace(/\D/g, '');
      if (it.tipo === 'PRODUTO' && dig.length > 0 && dig.length < 14) {
        showToast('IMEI inválido em um dos itens', 'error');
        return;
      }
    }

    const tipoPagamento =
      document.getElementById('tipo-pagamento')?.value ||
      document.querySelector('[name="tipo_pagamento"]')?.value ||
      'PIX';

    const clienteId =
      Number(document.getElementById('cliente-id')?.value) ||
      Number(document.getElementById('venda-cliente-id')?.value) ||
      0;

    try {
      // API cria uma venda por produto — envia o primeiro item (ou cada um)
      let lastSale = null;
      for (const it of itens) {
        if (it.tipo !== 'PRODUTO') continue;
        const produtoId = it.produto_id || 0;
        if (!produtoId) {
          // tenta resolver pelo nome
          const match = produtosCache.find((p) => p.nome === it.nome);
          if (match) it.produto_id = match.id;
        }
        if (!it.produto_id) {
          showToast(`Produto "${it.nome}" não encontrado no catálogo`, 'error');
          return;
        }

        const response = await SaleService.create({
          produto_id: it.produto_id,
          quantidade: it.qtd,
          tipo_pagamento: tipoPagamento,
          cliente_id: clienteId,
        });

        if (!response.success) {
          showToast(response.message || 'Erro ao registrar venda', 'error');
          return;
        }
        lastSale = response.data;
      }

      // Se trade-in ativo e venda criada, registra a troca
      if (useTrade && lastSale?.id) {
        const modelo = document.getElementById('trade-modelo')?.value || 'Aparelho troca';
        const imei = (document.getElementById('trade-imei')?.value || '').replace(/\D/g, '');
        const valor = parseFloat(document.getElementById('trade-valor')?.value) || 0;
        if (imei.length >= 14) {
          await TrocaService.create({
            venda_id: lastSale.id,
            modelo_aparelho: modelo,
            imei,
            valor_avaliacao: valor,
            condicao: 'USADO',
          });
        }
      }

      showToast('Venda registrada com sucesso!', 'success');
      setTimeout(() => goBack(), 800);
    } catch (err) {
      showToast(err.message || 'Erro ao registrar venda', 'error');
    }
  };

  // Popular datalist de produtos e clientes
  try {
    const [prodRes, cliRes] = await Promise.all([
      ProductService.listAll(),
      ClientService.listAll(),
    ]);
    if (prodRes.success && Array.isArray(prodRes.data)) {
      produtosCache = prodRes.data;
      let dl = document.getElementById('produtos-list');
      if (!dl) {
        dl = document.createElement('datalist');
        dl.id = 'produtos-list';
        document.body.appendChild(dl);
      }
      dl.innerHTML = produtosCache.map((p) => `<option value="${p.nome}">`).join('');
    }
    if (cliRes.success && Array.isArray(cliRes.data)) {
      clientesCache = cliRes.data;
      let dl = document.getElementById('clientes-list');
      if (!dl) {
        dl = document.createElement('datalist');
        dl.id = 'clientes-list';
        document.body.appendChild(dl);
      }
      dl.innerHTML = clientesCache.map((c) => `<option value="${c.nome}">`).join('');
    }
  } catch {
    // silencioso
  }

  addItem();
}

init();
