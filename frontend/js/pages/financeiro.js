import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { ResumoFinanceiro } from '../../client/models/ResumoFinanceiroModel.js';

if (!requireAuth()) throw new Error('redirect');
document.getElementById('btn-back').onclick = () => goBack();

let contas = [
  { nome: 'Aluguel', valor: 2500, dia: 5 },
  { nome: 'Luz', valor: 380, dia: 15 },
  { nome: 'Internet', valor: 120, dia: 10 }
];

document.querySelectorAll('#fin-tabs .tab').forEach(tab => {
  tab.onclick = () => {
    document.querySelectorAll('#fin-tabs .tab').forEach(t => t.classList.remove('active'));
    tab.classList.add('active');
    document.getElementById('tab-fechamento').classList.toggle('hidden', tab.dataset.tab !== 'fechamento');
    document.getElementById('tab-contas').classList.toggle('hidden', tab.dataset.tab !== 'contas');
  };
});

async function loadFechamento(periodo) {
  // CHAMADA DA API AQUI
  const data = new ResumoFinanceiro({
    total_servicos: periodo === 'dia' ? 450 : 8200,
    quantidade_os: periodo === 'dia' ? 3 : 42,
    total_vendas: periodo === 'dia' ? 1890 : 28500,
    quantidade_vendas: periodo === 'dia' ? 5 : 87,
    faturamento_bruto: periodo === 'dia' ? 2340 : 36700,
    total_despesas: periodo === 'dia' ? 200 : 4500,
    lucro_liquido: periodo === 'dia' ? 2140 : 32200
  });

  document.getElementById('f-servicos').textContent = formatBRL(data.total_servicos);
  document.getElementById('f-qtd-os').textContent = data.quantidade_os;
  document.getElementById('f-vendas').textContent = formatBRL(data.total_vendas);
  document.getElementById('f-qtd-vendas').textContent = data.quantidade_vendas;
  document.getElementById('f-bruto').textContent = formatBRL(data.faturamento_bruto);
  document.getElementById('f-despesas').textContent = formatBRL(data.total_despesas);
  document.getElementById('f-lucro').textContent = formatBRL(data.lucro_liquido);
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
  document.getElementById('contas-list').innerHTML = contas.map(c => `
    <div class="card list-item">
      <div class="info">
        <div class="title">${c.nome}</div>
        <div class="sub">Vence dia ${c.dia}</div>
      </div>
      <div class="fw-bold">${formatBRL(c.valor)}</div>
    </div>
  `).join('');
}

document.getElementById('form-conta').onsubmit = (e) => {
  e.preventDefault();
  contas.push({
    nome: document.getElementById('conta-nome').value.trim(),
    valor: parseFloat(document.getElementById('conta-valor').value) || 0,
    dia: parseInt(document.getElementById('conta-dia').value) || 1
  });
  showToast('Lembrete adicionado', 'success');
  document.getElementById('form-conta').reset();
  renderContas();
};

loadFechamento('mes');
renderContas();
