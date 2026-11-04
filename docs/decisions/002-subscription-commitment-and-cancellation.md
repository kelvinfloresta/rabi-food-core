# 002. A assinatura é um compromisso; cancelamento e reembolso

- Status: Aceita
- Data: 2026-05-01

## Contexto
O negócio local só consegue cobrar um preço justo e servir comida fresca se souber, com antecedência, quanto vai produzir. É o compromisso do cliente que dá essa previsibilidade: sem ele, a assinatura vira só mais um pedido, e o estabelecimento volta a produzir no escuro, desperdiçando ou faltando. Ao mesmo tempo, o cliente tem medo de ficar preso — e com razão, depois de tantos serviços que dificultam o cancelamento. A assinatura precisa ser um compromisso de verdade e, ainda assim, ter uma saída clara e honesta.

## Decisão
- O compromisso é real: o cliente paga pelo que escolheu e recebe o que escolheu.
- A saída é sempre visível e sem pegadinhas: parar de renovar é fácil a qualquer momento.
- Trocar de produto e pedir o dinheiro de volta são coisas diferentes e são tratadas de forma diferente.
- O direito de arrependimento do CDC é respeitado, contado a partir do momento em que o cliente de fato conhece o produto.

Na prática:
- **Sempre um compromisso, nunca um teste.** O cliente se inscreve e paga o plano escolhido; as entregas acontecem como agendadas.
- **A renovação automática não é uma opção liga/desliga.** Ela é um fato da assinatura. A interface diz isso com clareza: *"sua assinatura renova automaticamente ao fim do ciclo. Você pode cancelar a renovação a qualquer momento."*
- **Cancelamento sem reembolso:** a qualquer momento. As entregas já pagas continuam até acabar; sem multa, sem reembolso — a assinatura apenas deixa de renovar.
- **Cancelamento com reembolso (CDC):** somente até 7 dias após a **primeira entrega** (não a data da inscrição). O reembolso é automático e devolve o valor dos pedidos ainda não entregues.
- **Substituição não é reembolso.** "Prefiro outro produto" é um pedido de troca, sem custo nem atrito. "Quero meu dinheiro de volta" é reembolso, válido só na janela de 7 dias. A interface comunica os dois separadamente.

## Consequências
- A assinatura continua sendo um compromisso real, e o cliente sempre tem uma saída clara.
- As regras de reembolso dependem da data da primeira entrega, então ela precisa ser conhecida em cada assinatura.
- Todo fluxo que fala de cancelamento precisa manter substituição e reembolso separados.
