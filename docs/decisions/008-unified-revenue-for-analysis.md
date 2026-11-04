# 008. Receita unificada para análise

- Status: Aceita
- Data: 2026-10-25

## Contexto
Muitos negócios locais quebram não por falta de clientes, mas por não saberem se estão tendo lucro: misturam o dinheiro que entrou com o que ganharam, esquecem custos, descobrem o prejuízo tarde demais — e acabam recorrendo a empréstimos e juros abusivos. Ajudar o negócio local a não depender disso faz parte do posicionamento do produto. Para isso, o estabelecimento precisa de uma resposta confiável para duas perguntas simples: "estou ganhando dinheiro?" e "tenho dinheiro para pagar as contas?".

A assinatura complica essas respostas: o cliente paga adiantado por várias entregas, então o dinheiro entra num mês e a comida sai em outros. Um relatório que olhe só para os pagamentos mostra um mês ótimo seguido de um mês ruim, sem que nada tenha mudado de verdade.

## Decisão
- Uma única fonte de verdade para todas as análises: toda venda, avulsa ou de assinatura, conta da mesma forma.
- Lucro e caixa são perguntas diferentes e têm respostas separadas.
- A receita é contada quando o estabelecimento cumpre sua parte — quando entrega a comida.
- O passado não muda: o que foi vendido a um custo continua com aquele custo.

Na prática:
- **Uma única visão de receita alimenta todas as análises** (DRE, totais do CRM, painéis). Cada pedido avulso (`commerce.orders`) e cada entrega de assinatura (`subscription.subscription_deliveries`) é uma linha de receita.
- **A receita conta quando a comida é entregue (regime de competência).** Uma entrega de assinatura conta no mês em que é entregue, pelo valor pago dividido pelo número de entregas que ele cobre, então os descontos da assinatura já ficam descontados. Um pedido avulso conta quando é retirado ou entregue; uma venda de balcão, no momento da venda. Entregas que nunca acontecem nunca viram receita, o que bate com a 002, em que o reembolso devolve o valor dos pedidos não entregues.
- **Ao lado, uma visão de caixa:** o dinheiro recebido em cada mês, incluindo o que ainda está pendente, como pedidos com pagamento presencial ainda não confirmados.
- **Os custos são gravados no momento da venda, como os preços.** Cada produto tem um custo, copiado para os itens do pedido e para as assinaturas quando vendido.
- **Outros custos:** as despesas operacionais são lançamentos manuais por categoria e mês (aluguel, folha, embalagem, entrega, outros); as taxas do gateway vêm do provedor de pagamento; os impostos usam uma alíquota configurada por estabelecimento.

Estrutura do DRE: receita bruta − impostos e reembolsos = receita líquida − custo dos produtos vendidos (CMV) = lucro bruto − despesas operacionais e taxas do gateway = resultado operacional.

## Consequências
- Os produtos ganham um custo; pedidos e assinaturas o guardam.
- A conta de receita e custo é uma regra de negócio: os testes de integração são a especificação dela (data de reconhecimento, valor por entrega, reembolsos, taxas).
- Mudar o regime de reconhecimento depois exigiria recalcular os relatórios históricos, e por isso ele é decidido agora.
- Uma alíquota única por estabelecimento é uma simplificação: as alíquotas do Simples Nacional mudam conforme a faixa de faturamento anual. Serve para um relatório gerencial, não para declarar impostos.
