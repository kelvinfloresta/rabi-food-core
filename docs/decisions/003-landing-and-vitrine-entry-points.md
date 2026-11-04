# 003. Duas portas de entrada: landing page e vitrine

- Status: Substituída pela 006
- Data: 2026-05-01

## Contexto
Assinar comida é um hábito novo para a maioria das pessoas: o mercado está acostumado a pedir na hora, não a se comprometer com antecedência. Por isso o produto recebe dois tipos de pessoa com necessidades opostas. Quem ainda não conhece o modelo (o Descobridor) precisa ser convencido de que vale a pena — se for jogado direto numa lista de produtos, não entende o que está comprando e vai embora. Quem já conhece (o Decidido) está parado na frente da loja com o celular na mão — se tiver que ler uma explicação, desiste. Uma página só, tentando servir aos dois, falha com os dois.

## Decisão
- Cada pessoa chega pela porta feita para ela: quem precisa entender encontra explicação; quem já decidiu encontra os produtos.
- Cada página tem uma única função e não repete a da outra.
- Nenhuma das duas fecha o caminho: o Decidido sempre pode ir buscar a explicação, se quiser.

Na prática:
- **Landing page** para o Descobridor (busca, redes sociais, indicação). Ela comunica o valor antes do fluxo de assinatura: conveniência e organização, preço competitivo, comida fresca e apoio aos negócios locais.
- **Vitrine** para o Decidido (QR code no estabelecimento, link compartilhado por outro cliente). Ela *é* a listagem de produtos onde o cliente escolhe o que assinar: sem botão "assinar" que leve a outra listagem.
- A vitrine não explica o produto; ela leva para a landing page.

## Consequências
- Cada página tem uma função, e nenhuma duplica a da outra.
- Links compartilhados ou impressos (QR codes) apontam para a vitrine, não para a landing page.
