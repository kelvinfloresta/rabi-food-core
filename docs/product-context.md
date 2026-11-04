# Contexto do produto

Para quem é o produto e como ele se posiciona. As decisões específicas ficam em `docs/decisions/`; leia as que tocam a tarefa.

## Perfis de cliente

**Descobridor** — ainda não conhece o produto. Precisa entender antes de se comprometer. Chega por busca, redes sociais ou indicação.

**Decidido** — já sabe o que quer. Precisa do mínimo de atrito. Chega pelo QR code no estabelecimento ou por um link compartilhado por outro cliente.

## Posicionamento

Um modelo de assinatura, fundamentalmente diferente do que o mercado está acostumado, construído sob extrema incerteza (lean startup). O valor dele:

- Conveniência e organização
- Preço mais competitivo
- Comida fresca
- Ajuda os negócios locais a gerar valor para a sua região: reduz a dependência de empréstimos, evita juros abusivos, preço justo sem comissão por venda — ao contrário dos marketplaces que exploram os negócios locais. Este é um diferencial central e precisa ser mais elaborado.

A venda avulsa (no balcão e online) complementa a assinatura como porta de entrada, nunca como caminho principal. A plataforma também dá ao estabelecimento as ferramentas para gerir o negócio — CRM e análise financeira (DRE) sobre todas as suas vendas.

## Mercado

Concorrentes, preços, reclamações e a pontuação da oportunidade: [market-analysis.md](market-analysis.md).

## Decisões de produto

Agrupadas pelo caminho do cliente: como chega, como compra, de quem são os dados e como o estabelecimento gere o negócio. Os números são identificadores permanentes, não uma ordem.

**Chegada**
- [006 A vitrine oferece comprar e assinar](decisions/006-vitrine-offers-buy-and-subscribe.md)

**Compra**
- [002 A assinatura é um compromisso; cancelamento e reembolso](decisions/002-subscription-commitment-and-cancellation.md)
- [004 A tela de resumo é onde a confiança é construída](decisions/004-summary-screen-builds-trust.md)
- [005 A venda avulsa complementa a assinatura](decisions/005-one-off-sales-complement-subscriptions.md)

**Contas e dados do cliente**
- [001 Conta do cliente compartilhada entre estabelecimentos](decisions/001-shared-customer-accounts.md)
- [007 CRM: o que o estabelecimento vê dos seus clientes](decisions/007-crm-customer-data-visibility.md) — proposta

**Gestão do negócio**
- [008 Receita unificada para análise](decisions/008-unified-revenue-for-analysis.md)

**Substituídas**
- [003 Duas portas de entrada: landing page e vitrine](decisions/003-landing-and-vitrine-entry-points.md) — pela 006
