# 010. Métricas da plataforma: números por estabelecimento, nunca por pessoa

- Status: Aceita
- Data: 2026-11-26

## Contexto
Quem é dono do software precisa saber se a plataforma está funcionando como negócio: quantos estabelecimentos entraram, quais estão vendendo, se a assinatura é de fato o caminho principal e quanto dinheiro passa por ela. Sem isso, as decisões sobre o produto ficam no escuro — e o produto ainda está validando se vale a pena (ver `market-analysis.md`).

Ao mesmo tempo, os dados dos clientes pertencem à relação entre o cliente e o estabelecimento. A plataforma ainda não decidiu o que pode ver deles nem quem é o controlador dos dados (questões em aberto na 001 e na 007). Uma visão da plataforma que mostrasse pessoas tomaria essa decisão sem querer.

## Decisão
- O dono do software vê números agregados por estabelecimento e por mês, nunca uma pessoa.
- Os números são calculados na hora, a partir dos dados do sistema, sem uma segunda cópia que possa divergir.
- O que é mostrado é dinheiro recebido, não receita: a receita segue a regra da 008.
- O banco de dados continua sem acesso pela internet.

Na prática:
- **Views no schema `analytics`** são o contrato dos painéis: `analytics.tenants` (estabelecimentos e desde quando existem), `analytics.tenant_monthly` (por estabelecimento e mês: compradores, novos compradores, pedidos e assinaturas criados e pagos, renovações, valores pagos) e `analytics.tenant_current` (assinaturas ativas, compradores, última atividade). São consultas com nome, executadas a cada leitura; não guardam dados. Os meses são contados no horário de Brasília.
- **O painel só lê essas views**, com o papel `analytics_reader`, que não lê as tabelas nem a view interna que contém o identificador de cada comprador. Nomes, e-mails e telefones não chegam ao Grafana.
- **"Comprador" é quem teve uma venda paga no estabelecimento.** O vínculo formal de cliente (`iam.tenant_customers`) ainda não é preenchido por nenhum fluxo (007).
- **Uma assinatura é um checkout**, como o cliente a enxerga: assinar três produtos de uma vez conta como uma assinatura. Renovações são contadas à parte; quando forem construídas, cada renovação de um checkout recebe uma chave de checkout nova, comum a todos os seus produtos, para continuar contando como uma.
- **O nome do estabelecimento aparece**, porque é público na vitrine. Para um MEI ele pode ser o nome de uma pessoa; aceitamos, já que é o nome com que o negócio se apresenta.
- **Valores pagos ≠ receita.** Os painéis somam pedidos e assinaturas pagos pela data do pagamento — a visão de caixa. A receita por competência (entrega) depende da visão unificada da 008, que ainda não existe; quando existir, os painéis de receita a usarão.
- **Acesso:** o Grafana Cloud consulta o Postgres por um agente (Private Data Source Connect) que roda junto da API, dentro da rede da AWS. O usuário de leitura tem limite de conexões e de tempo por consulta, para não pesar no banco da aplicação.
- **Logs também não levam pessoas ao Grafana:** os logs do serviço de autenticação (GoTrue) registram e-mail, identificador e IP de quem entra; esses valores são mascarados antes do envio.

## Consequências
Fica mais fácil:
- Acompanhar todos os estabelecimentos num só painel, sem código novo na API e com o histórico completo que o banco já guarda.
- Mudar uma tabela sem quebrar os painéis: basta manter as colunas das views, e um teste executa todas as consultas do painel.

Fica mais difícil, e aceitamos:
- Ainda não dá para medir cancelamentos, falhas de pagamento e receita entregue: esses dados não são gravados hoje. Um reembolso só muda o status do pagamento, então a venda some do mês em que foi paga em vez de aparecer como reembolso no mês em que aconteceu.
- Quem comprou em dois estabelecimentos conta como comprador em cada um; não há total de pessoas únicas na plataforma, por desenho.
- As consultas rodam no mesmo banco da aplicação. Com volume grande, precisarão de uma réplica de leitura.
- Ver dados de uma pessoa (por exemplo, para suporte) exige uma nova decisão sobre o acesso da plataforma a dados pessoais.
- O mascaramento dos logs do GoTrue depende de reconhecer o formato dos dados (e-mail, identificador, IP); um dado pessoal em outro formato passaria.
