# 005. A venda avulsa complementa a assinatura

- Status: Aceita
- Data: 2026-10-25

## Contexto
Nem todo mundo está pronto para assinar no primeiro contato — e ninguém assina algo que nunca provou. Hoje, quem entra na loja e compra um pão fica fora da plataforma: o estabelecimento não sabe quem é, não tem como convidá-lo a assinar, e sua receita fica espalhada entre o caderno, a maquininha e o sistema. Isso enfraquece justamente o que queremos dar ao negócio local: clareza sobre o próprio negócio (ver 008) e relação com seus clientes (ver 007). A compra avulsa é o primeiro passo de muitos futuros assinantes; ignorá-la é perdê-los.

O risco é o oposto: se o avulso virar o caminho principal, o produto vira mais um app de pedidos e perde o que o torna diferente.

## Decisão
- O avulso é a porta de entrada e o complemento da assinatura, nunca o caminho principal.
- Toda venda do estabelecimento acontece num só lugar, seja no balcão ou online.
- Ninguém fica sem comprar por dificuldade com tecnologia: se o cliente não consegue, o atendente faz por ele.
- O avulso respeita o planejamento do estabelecimento: nada de produzir ou entregar no improviso.
- Assinar precisa continuar valendo mais a pena do que comprar avulso.

Na prática:
- **Dois canais, um único fluxo de pedido.** `online`: o cliente pede pelo próprio celular, dentro ou fora da loja. `counter` (balcão): um atendente faz o pedido no tablet ou celular para o cliente que tem dificuldade. Todo pedido registra quem o criou.
- **No balcão, identificar o cliente é opcional.** O atendente pode vincular o cliente pelo telefone ou deixar a venda anônima, para não atrasar a fila. Vendas anônimas entram no DRE, mas não no CRM.
- **O primeiro pedido cadastra o cliente no estabelecimento** (`iam.tenant_customers`), online ou vinculado no balcão, igual a uma primeira assinatura (ver 007).
- **Retirada ou entrega agendada, nunca sob demanda.** Retirada na loja, ou entrega somente nos dias e horários em que o estabelecimento entrega as assinaturas, com a mesma antecedência mínima (`order_lead_minutes`).
- **Pagamento online ou presencial.** Online pelo gateway, ou presencial (dinheiro, Pix, maquininha) no balcão ou na entrega. O pedido com pagamento presencial fica pendente até um atendente confirmar o recebimento, para que os números batam com o caixa.
- **Preço normal.** A venda avulsa cobra o preço normal do produto. Os descontos da assinatura (`discount_rules`) valem só para assinaturas: são a recompensa pelo compromisso da 002.
- **O checkout mostra o que o cliente vai receber, como na 004:** data concreta de retirada ou entrega, total e forma de pagamento. Não tem conteúdo de renovação nem de substituição.
- **As regras de reembolso da assinatura (002) não se aplicam.** Os fluxos avulsos nunca reutilizam os textos de cancelamento da assinatura.

## Consequências
- Todas as vendas do estabelecimento ficam num só lugar, do qual o CRM (007) e o DRE (008) dependem.
- Quem compra avulso com frequência fica visível como candidato a assinante.
- Os pedidos ganham canal, criador, forma de atendimento (retirada ou entrega), horário de entrega e forma de pagamento, e podem existir sem cliente (vendas anônimas no balcão).
- O atendente agindo em nome do cliente exige um token com escopo do estabelecimento (001); toda rota desse tipo precisa ser isolada por estabelecimento e coberta por testes de integração.
- **O estabelecimento precisa de uma agenda de entregas.** Hoje o horário de entrega é escolhido em cada assinatura (`StartHour` na criação da assinatura); não existe uma lista de dias e horários de entrega do estabelecimento para o pedido avulso escolher. Se as assinaturas também passam a seguir essa agenda é uma decisão separada.
- **Em aberto, antes do lançamento:** se o direito de arrependimento do CDC (art. 49, compras fora do estabelecimento) vale para comida perecível comprada online, e para um pedido que o cliente faz no próprio celular estando dentro da loja.
