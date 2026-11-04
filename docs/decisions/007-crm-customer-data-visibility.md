# 007. CRM: o que o estabelecimento vê dos seus clientes

- Status: Proposta
- Data: 2026-10-25

## Contexto
O negócio local sempre conheceu seus clientes pelo nome — é uma vantagem que nenhum marketplace tem. Mas, conforme cresce, ele perde essa memória: não sabe quem parou de comprar, quem compra toda semana sem assinar, quem merece um agradecimento. O CRM devolve essa memória ao estabelecimento.

Só que a conta do cliente pertence à plataforma (001), e o cliente compra em vários estabelecimentos. Se um estabelecimento pudesse ver o que o cliente faz nos outros, a plataforma trairia a confiança do cliente e viraria exatamente o marketplace que critica, com negócios disputando os mesmos clientes. O CRM precisa aproximar o estabelecimento do cliente sem expor o cliente.

## Decisão
- Cada estabelecimento conhece os seus clientes, e só o que aconteceu com ele.
- O cliente nunca é exposto à concorrência.
- O que o estabelecimento anota sobre um cliente é dele, e fica só com ele.
- O cliente sabe quando passa a ser cliente de um estabelecimento.

Na prática:
- **O estabelecimento vê apenas os seus clientes**: os cadastrados em `iam.tenant_customers` para ele. O cadastro acontece automaticamente no primeiro pedido ou na primeira assinatura, ou quando o atendente vincula o cliente no balcão (005). A 001 diz que essa tabela não controla acesso, e isso continua valendo para o cliente: ele compra em qualquer lugar sem se cadastrar. Aqui ela define o que o estabelecimento pode ver.
- **Apenas a atividade naquele estabelecimento:** pedidos avulsos, assinaturas e entregas. Nunca a atividade em outros estabelecimentos.
- **Dados do perfil: nome, telefone e endereço.** São globais (001): o estabelecimento vê os valores atuais do cliente, não uma cópia por estabelecimento.
- **Etiquetas e anotações internas pertencem ao estabelecimento.** São invisíveis para o cliente e para os outros estabelecimentos.
- **Vendas anônimas no balcão não são clientes** e não aparecem no CRM.

## Consequências
- Toda consulta do CRM é filtrada por estabelecimento e por cadastro; os testes de integração são a especificação, já que esquecer esse filtro já expôs dados uma vez (001).
- O cadastro automático precisa ser adicionado à criação de pedidos e de assinaturas. Hoje o cadastro é um endpoint separado (`enroll-customers`) que nenhum fluxo chama; sem isso, o CRM ficaria vazio.
- **Em aberto, antes do lançamento:**
  - Papel de controlador dos dados (plataforma ou estabelecimento) e o texto de consentimento que o cliente vê ao ser cadastrado.
  - Retenção: o que acontece com a visão que o estabelecimento tem de um cliente quando o estabelecimento sai ou o cliente pede a exclusão.
  - Vincular no balcão um telefone sem conta na plataforma: criar a conta, convidar o cliente ou manter a venda anônima.
  - Quais membros do estabelecimento podem ver o CRM (ainda não existem papéis).
