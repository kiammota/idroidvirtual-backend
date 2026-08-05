# IDroid Virtual — Gestão

Frontend mobile-first para gestão de assistências técnicas / lojas de celulares.

## Identidade visual
- **Branco** — fundo principal
- **Laranja `#f06522`** — destaque, header, botões
- **Preto** — textos e suporte

## Estrutura

```
/
├── index.html              # Redirect inteligente
├── assets/logo.svg
├── css/main.css            # Design system mobile-first
├── css/login.css
├── js/app.js               # Auth, toast, PIN modal, utils
├── client/                 # HTTP client + services (API)
│   ├── env.js
│   ├── http.js
│   ├── models/
│   └── services/
├── models/
│   ├── go/                 # Modelos Go (referência backend)
│   └── js/                 # Modelos JS
└── pages/
    ├── login/
    ├── me/
    ├── home/
    ├── ordens/             # OS + assinatura canvas
    ├── clientes/
    ├── trocas/             # Trade-in
    ├── estoque/            # Entrada + Brindes
    ├── entregas/
    ├── financeiro/         # Fechamento + contas (PIN)
    └── venda/              # Fluxo completo de venda
```

## Como rodar

ES modules exigem servidor HTTP (não use `file://`).

```bash
cd idroid-virtual
python -m http.server 8080
# ou: npx serve .
```

Acesse: `http://localhost:8080`

## Autenticação
- Token JWT em `sessionStorage` (`gv_jwt_token`)
- Páginas protegidas → redirect para login
- Login em **modo simulado** (e-mail + senha válidos)
- Descomente as `// CHAMADA DA API AQUI` quando a API estiver pronta

## Áreas com PIN
- Estoque & Aparelhos
- Financeiro & Faturamento

## Integração
Use os services em `client/services/`. Todos os pontos estão comentados com:

```javascript
// CHAMADA DA API AQUI
```
