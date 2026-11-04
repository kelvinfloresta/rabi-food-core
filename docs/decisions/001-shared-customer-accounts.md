# 001. Conta do cliente compartilhada entre estabelecimentos

- Status: Aceita
- Data: 2026-10-23

## Contexto
O cliente compra de vários negócios locais: a padaria, o hortifrúti, a marmita da semana. Se cada um exigir um cadastro novo, o cliente desiste justamente no momento em que estava pronto para assinar — o perfil Decidido chega pelo QR code ou por um link e quer resolver na hora. Cada cadastro a mais é uma venda a menos para o negócio local.

## Decisão
- O cliente se cadastra uma vez e é bem-vindo em qualquer estabelecimento da plataforma.
- O que o cliente faz como cliente é diferente do que alguém faz em nome do estabelecimento: o segundo exige ser membro dele.
- Ser cliente de um estabelecimento é uma relação comercial, não uma permissão.

Na prática:
- Um mesmo usuário (GoTrue) é reutilizado em todos os estabelecimentos. O cliente pode navegar, pedir e assinar em qualquer um sem ser membro dele; as rotas do cliente recebem o estabelecimento pelo header `X-Tenant-ID`, que apenas escolhe onde o cliente está agindo.
- Agir **como** o estabelecimento — alterar os dados dele ou ler dados de vários clientes dele — exige um token com escopo do estabelecimento, emitido por `exchange-token` / `refresh-scoped` somente para membros (`iam.tenant_members`).
- `iam.tenant_customers` registra a relação comercial (quem é cliente de qual estabelecimento, via `enroll-customers`). Ela não controla acesso.

## Consequências
Fica mais fácil:
- Um cadastro, um perfil e uma forma de pagamento para todos os estabelecimentos.
- O cliente de um estabelecimento está a um passo de outro, então os estabelecimentos se beneficiam dos clientes uns dos outros.
- Uma única base de usuários: sem autenticação por estabelecimento nem contas duplicadas.

Fica mais difícil, e aceitamos:
- **A segurança depende de cada consulta.** O header não autoriza nada, então toda consulta das rotas do cliente precisa ser filtrada pelo usuário da sessão; esquecer isso expõe dados de outros clientes (já aconteceu uma vez, na listagem de assinaturas). Garantido pela regra em `.claude/rules/api/architecture.md` e por testes de integração.
- **Propriedade dos dados e LGPD.** A plataforma guarda a conta do cliente; quais dados cada estabelecimento pode ver, quem é o controlador dos dados e o que acontece com eles quando um estabelecimento sai precisam ser definidos antes do lançamento (ver 007).
- **Os dados do perfil são globais.** Mudanças de telefone, endereço e nome valem para todos os estabelecimentos; não existe uma versão do cliente por estabelecimento.
- **Ainda não há bloqueio por estabelecimento.** Um estabelecimento não pode banir um cliente sem um novo estado (provavelmente em `iam.tenant_customers`).
- **Marca da plataforma.** Os e-mails de login e OTP são da plataforma, não do estabelecimento.

Questão em aberto: o produto se posiciona contra os marketplaces, e contas compartilhadas com descoberta entre estabelecimentos se parecem com um. Onde fica o limite (ex.: sem comissão, sem estabelecimentos disputando o mesmo cliente, sem promoção cruzada) ainda não está decidido.
