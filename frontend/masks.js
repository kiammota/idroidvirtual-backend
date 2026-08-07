/**
 * Máscaras de input — telefone, IMEI, CNPJ, valor monetário
 */

export function onlyDigits(value) {
  return String(value || '').replace(/\D/g, '');
}

/** (11) 98765-4321 ou (11) 3456-7890 */
export function maskPhone(value) {
  let d = onlyDigits(value).slice(0, 11);
  if (d.length <= 2) return d.length ? `(${d}` : '';
  if (d.length <= 6) return `(${d.slice(0, 2)}) ${d.slice(2)}`;
  if (d.length <= 10) return `(${d.slice(0, 2)}) ${d.slice(2, 6)}-${d.slice(6)}`;
  return `(${d.slice(0, 2)}) ${d.slice(2, 7)}-${d.slice(7)}`;
}

/** IMEI: apenas dígitos, máx 15 */
export function maskImei(value) {
  return onlyDigits(value).slice(0, 15);
}

/** 00.000.000/0001-00 */
export function maskCnpj(value) {
  let d = onlyDigits(value).slice(0, 14);
  if (d.length <= 2) return d;
  if (d.length <= 5) return `${d.slice(0, 2)}.${d.slice(2)}`;
  if (d.length <= 8) return `${d.slice(0, 2)}.${d.slice(2, 5)}.${d.slice(5)}`;
  if (d.length <= 12) return `${d.slice(0, 2)}.${d.slice(2, 5)}.${d.slice(5, 8)}/${d.slice(8)}`;
  return `${d.slice(0, 2)}.${d.slice(2, 5)}.${d.slice(5, 8)}/${d.slice(8, 12)}-${d.slice(12)}`;
}

/**
 * Aplica máscara em um input ao digitar
 * @param {HTMLInputElement|string} el
 * @param {'phone'|'imei'|'cnpj'} type
 */
export function applyMask(el, type) {
  const input = typeof el === 'string' ? document.getElementById(el) : el;
  if (!input) return;

  const handlers = {
    phone: maskPhone,
    imei: maskImei,
    cnpj: maskCnpj
  };
  const fn = handlers[type];
  if (!fn) return;

  input.addEventListener('input', () => {
    const start = input.selectionStart;
    const oldLen = input.value.length;
    input.value = fn(input.value);
    const newLen = input.value.length;
    // tenta manter cursor
    const pos = Math.max(0, (start || 0) + (newLen - oldLen));
    try { input.setSelectionRange(pos, pos); } catch { /* ignore */ }
  });

  // aplica valor inicial se houver
  if (input.value) input.value = fn(input.value);
}

/** Atalho: aplica máscara em vários campos de uma vez */
export function bindMasks(map) {
  // map: { 'os-telefone': 'phone', 'trade-imei': 'imei' }
  Object.entries(map).forEach(([id, type]) => applyMask(id, type));
}

window.IDroidMasks = { onlyDigits, maskPhone, maskImei, maskCnpj, applyMask, bindMasks };
