# rabi-food-core/api — Contexto para Claude Code

## Visão Geral
Monolito Go para gestão de pedidos e delivery multi-tenant.
Arquitetura: Feature-based Modular Monolith.
Princípio central: dependências entre features devem ser **unidirecionais e intencionais**.

## Stack
- **Runtime**: Go 1.26
- **HTTP**: Fiber v2
- **DB**: pgx/v5 (Postgres)
- **DI**: samber/do
- **Logs**: zerolog + Wide Event pattern
- **Validação**: go-playground/validator
- **Testes**: testify suite + httpexpect + fixtures próprias

## Estrutura de Pacotes
```
api/
├── config/
├── app_context/        — session.go (UserSession)
├── domain/             — tipos compartilhados (auth/role, payment_status)
├── features/           — auth, category, order, product, subscription, tenant
│   └── [feature]/
│       ├── controller/ — fiber_*.go (handlers HTTP)
│       ├── gateway/    — interface.go + pgx_*.go
│       ├── model/
│       ├── routes/
│       └── usecases/   — lógica de negócio + integration_test.go
├── fixtures/           — helpers para testes de integração
└── libs/
    ├── database/       — PgxAdapter, PaginateInput, DefaultPage/DefaultPageSize
    │   └── filter/     — filter.Bool (optional boolean filter type)
    ├── di/             — new.go (wiring DI)
    ├── errs/           — erros de domínio por feature
    ├── http/           — FiberAdapter + middlewares
    ├── logger/         — WideEvent pattern
    └── validator/
```

## Migrations
Migrations rodam externamente — **não são executadas pelo binário da aplicação**. Para criar uma nova migration, adicione um arquivo `.sql` em `libs/database/migrations/` com anotações goose e rode `task migrate`.

- **Pré-produção (situação atual):** o projeto ainda não está em produção, então mudanças de schema vão direto na migration `create` da tabela — sem migrations incrementais (`add_x_to_y`). Depois de alterar, recrie o banco local (`task db-down && task db && task migrate && task seed`).
- **Depois do go-live:** nunca edite uma migration já aplicada; toda mudança vira uma migration nova. Editar uma antiga faz bancos novos e existentes divergirem.
- Os testes de integração aplicam todas as migrations do zero em um banco novo — se `task test` falhar no startup com erro de SQL, a cadeia de migrations está quebrada para qualquer ambiente novo.

## Testes de integração
- Rodam com **testcontainers-go**: `fixtures.Api.Start` sobe Postgres + GoTrue por pacote, aplica as migrations e só então constrói o DI (os adapters de GoTrue, Stripe e HTTP copiam valores de `config` na criação; o de banco guarda o ponteiro e só lê no `Connect`).
- O ambiente é hermético: segredos (`AUTH_SECRET`, service key do GoTrue) são gerados em `fixtures/environment.go`, nunca lidos de `.env`.
- Só precisa de Docker. Imagens de Postgres/GoTrue devem ficar em sincronia com `infra/api/docker-compose.db.yaml`.

## Datas e fuso horário
- Colunas de data são sempre `TIMESTAMPTZ`; gateways gravam `time.Now().UTC()`.
- O `PgxAdapter` registra no `AfterConnect` um `TimestamptzCodec{ScanLocation: time.UTC}`: tudo que é lido do banco volta em UTC, independente do fuso do servidor.
- **Nunca** altere `time.Local` — é estado global do processo. Regras de negócio que dependem de fuso (entregas, cutoff, cron) usam `time.LoadLocation` + `.In(loc)` explicitamente.
- Em testes, compare instantes com `.Equal()`, não com `==`.

## Health check
- `GET /health` (sondado pelo ALB) é uma rota em `libs/http/health_check.go`, registrada em `fiber_adapter.go` **antes** dos `Use`: os middlewares nunca rodam para ela, então não gera log de request, métrica nem CORS.
- `http.HealthCheck` (construtor `http.NewHealthCheck`) é registrado no DI e injetado em `http.New`, como os controllers. Verifica `database` (`PgxAdapter.Ping`) e `auth` (`AuthGateway.Health`), com timeout de 2s cada. Sucesso responde só `200` (sem JSON); falha responde `503` com `{"unavailable": [...]}` — o erro vai para o log, nunca para a resposta.
- Só entram dependências sem as quais a API não opera: nova dependência crítica → novo campo em `HealthCheck`, resolvido no provider de `*http.HealthCheck` em `libs/di/new.go`, e uma chamada a `probe` em `Handle`. Serviços de terceiros (ex.: Stripe) não entram: uma instabilidade deles derrubaria todas as tasks.

## Issues Conhecidos (não replicar)
1. `gateway.Create` em Order não recebe `context.Context` — inconsistente, novos gateways devem receber
2. Módulo de pagamento externo ainda não implementado
3. Scheduling de entregas de Subscription ainda não implementado

## Comandos úteis
```bash
task test          # rodar testes de integração (testcontainers, só precisa de Docker)
task db            # subir Postgres + GoTrue + Mailpit locais (desenvolvimento)
task migrate       # rodar migrations pendentes
task build         # build
```
