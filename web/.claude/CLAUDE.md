# rabi-food-core/web — Padrões de código frontend

## Regras que nunca quebrar

### Render — renderização condicional
- **Sem ternário no render.** Nunca `condition ? <A /> : <B />` nem `condition ? x : y` dentro do JSX.
- **Bloco único opcional** → short-circuit: `{condition && <Component />}`
- **Dois blocos mutuamente exclusivos** → early return antes do JSX principal. Nunca `else` — o `if` retorna e o fluxo continua.
- **If aninhado no render** → mova para uma função auxiliar com early return interno.

### Código
- **Sem `else`** em nenhum lugar. Prefira early return.
- **If aninhado** → extraia para função com early return interno.

### Tipos
- **`readonly` em todos os tipos e arrays.** `readonly string[]`, `Readonly<Props>`, etc.
- **Sem tipagem inline** no JSX ou em funções — defina um tipo nomeado.
- **Sem números mágicos** — extraia como constante nomeada com intenção clara.
- **Listas e valores reutilizáveis** declarados com `as const` e tipados explicitamente.

### Estrutura de diretórios de componentes
- **Componente usado só dentro de uma página** → reside no diretório da própria página (ex: `app/routes/subscription/_components/product-sheet.tsx`).
- **Componente usado em mais de uma página** → sobe para `app/components/`.
- **Componente de UI genérico (sem lógica de domínio)** → reside em `app/components/ui/` e deve ser construído reutilizando os primitivos base existentes (`Button`, `Sheet`, etc.).
- **Nunca criar um componente de domínio direto em `app/components/ui/`** — UI é para primitivos reutilizáveis, não para componentes de feature.

### Internacionalização
- **Todo texto visível ao usuário via i18n** — sem string literal no JSX.
- **Placeholders de formato** (ex: `(11) 99999-9999`, `00000-000`) são constantes do componente, não chaves de tradução — o formato não muda entre idiomas.
- **Locale no servidor via middleware** (React Router 8 + remix-i18next 8): o `i18nextMiddleware` está registrado no `root.tsx`. Em loaders/actions use `getLocale(context)` e `getInstance(context)` de `~/i18n.server` — nunca crie uma instância i18next manualmente nem leia o locale direto do request.

### Autenticação
- **Tokens nunca chegam ao navegador.** Access e refresh token vivem só em cookies `HttpOnly`; o `authMiddleware` (`~/lib/session.server`) renova o access token no servidor antes de qualquer loader/action.
- Em loaders/actions use `getAccessToken(context)` — nunca leia o cookie manualmente, nunca receba token no body, nunca devolva token para o client.
- Para iniciar/encerrar sessão use `createSessionHeaders(tokens)` / `destroySessionHeaders()`. O client só sabe *se* está logado (`isLoggedIn: boolean` vindo do loader).
- Sem `localStorage`/`sessionStorage` para nada relacionado a sessão.

## Comandos úteis
```bash
npm run dev        # iniciar servidor de desenvolvimento
npm run build      # build de produção
npm run typecheck  # checar tipos
```

## Exemplos

```tsx
// ERRADO — ternário no render
{isNew ? <FieldA /> : <FieldB />}

// CERTO — early return para blocos mutuamente exclusivos
if (!isNew) return <FieldB />
return <FieldA />
```

```tsx
// CERTO — short-circuit para bloco único opcional
{isNew && <NameField />}
{!isNew && <EmailOnlyHint />}
```

```tsx
// ERRADO
const items = ["a", "b", "c"]

// CERTO
const ITEMS = ["a", "b", "c"] as const
type Item = typeof ITEMS[number]
```

```tsx
// ERRADO
if (code.length < 6) return

// CERTO — nomeie a intenção
const OTP_CODE_LENGTH = 6
if (code.length < OTP_CODE_LENGTH) return
```
