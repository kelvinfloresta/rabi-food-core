# Registros de decisão

Um arquivo por decisão com consequências duradouras, de produto ou de arquitetura: regras de produto e seus motivos, limites entre funcionalidades, nova infraestrutura, padrões transversais. Não serve para funcionalidades ou correções — essas são planejadas na conversa e registradas pelo código e pelos testes.

Liste as decisões de produto no índice ao final de `docs/product-context.md`, no grupo do tema a que pertencem (as substituídas vão para o grupo "Substituídas"). Os números são identificadores permanentes: nunca renumere um registro.

Nomeie os arquivos como `NNN-titulo-curto.md` (ex.: `001-payment-subscription-boundary.md`). Nunca reescreva um registro aceito: quando uma decisão muda, crie um novo e marque o antigo como `Substituída pela NNN`. Os documentos são escritos em português.

## Modelo

Escreva pelo porquê (Golden Circle). O contexto conta o propósito — que problema do cliente ou do estabelecimento existe e por que importa para o posicionamento — sem detalhes de implementação. A decisão começa pelos princípios que seguem desse porquê e só depois chega às regras concretas, cada uma com o seu motivo. As regras mudam com o tempo; o motivo dura mais. Se o contexto não convence sozinho, a decisão ainda não está madura.

```markdown
# NNN. Título

- Status: Proposta | Aceita | Substituída pela NNN
- Data: AAAA-MM-DD

## Contexto
Por que isso importa: o problema e o propósito, em poucas frases.

## Decisão
Os princípios e, na prática, as regras.

## Consequências
O que fica mais fácil, o que fica mais difícil e o que aceitamos.
```
