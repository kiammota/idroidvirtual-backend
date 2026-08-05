import { requireAuth, goBack, showToast, formatBRL } from '../app.js';
import { bindMasks } from '../masks.js';

if (!requireAuth()) throw new Error('redirect');
document.getElementById('btn-back').onclick = () => goBack();

bindMasks({ 'troca-imei': 'imei' });

const toggleDesc = document.getElementById('toggle-desconto');
const toggleJuros = document.getElementById('toggle-juros');
let useDesc = true, useJuros = false;

toggleDesc.onclick = () => {
  useDesc = !useDesc;
  toggleDesc.classList.toggle('on', useDesc);
  toggleDesc.setAttribute('aria-checked', useDesc);
  calc();
};
toggleJuros.onclick = () => {
  useJuros = !useJuros;
  toggleJuros.classList.toggle('on', useJuros);
  toggleJuros.setAttribute('aria-checked', useJuros);
  document.getElementById('juros-group').classList.toggle('hidden', !useJuros);
  calc();
};

function calc() {
  const aval = parseFloat(document.getElementById('troca-avaliacao').value) || 0;
  const novo = parseFloat(document.getElementById('troca-novo').value) || 0;
  const desc = useDesc ? aval * 0.10 : 0;
  const avalLiq = aval - desc;
  let base = Math.max(0, novo - avalLiq);
  const taxa = useJuros ? (parseFloat(document.getElementById('juros-pct').value) || 0) / 100 : 0;
  const juros = base * taxa;
  const final = base + juros;

  document.getElementById('r-avaliacao').textContent = formatBRL(aval);
  document.getElementById('r-desconto').textContent = formatBRL(desc);
  document.getElementById('r-novo').textContent = formatBRL(novo);
  document.getElementById('r-juros').textContent = formatBRL(juros);
  document.getElementById('r-final').textContent = formatBRL(final);
}

['troca-avaliacao', 'troca-novo', 'juros-pct'].forEach(id => {
  document.getElementById(id).addEventListener('input', calc);
});

document.getElementById('form-troca').onsubmit = async (e) => {
  e.preventDefault();
  const imei = document.getElementById('troca-imei').value.replace(/\D/g, '');
  if (imei.length < 14) {
    showToast('IMEI deve ter 15 dígitos', 'error');
    return;
  }

  // CHAMADA DA API AQUI
  showToast('Troca registrada · aparelho preparado no estoque', 'success');
  document.getElementById('form-troca').reset();
  calc();
};
