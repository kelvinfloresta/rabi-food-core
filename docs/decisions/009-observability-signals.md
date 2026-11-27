# 009. Sinais de observabilidade: wide event, spans e métricas

- Status: Aceita
- Data: 2026-11-26

## Contexto
Quando algo falha em produção, precisamos descobrir rápido o que aconteceu com aquele pedido, aquela assinatura ou aquele cliente, e de onde veio o problema: do nosso código, do banco ou de um serviço externo. Ao mesmo tempo, logs são cobrados por volume e por consulta. Uma linha de log a cada passo, ou dezenas de campos indexados, custa caro e ainda deixa a resposta espalhada.

## Decisão
- Cada pergunta tem um sinal responsável, e um sinal não faz o trabalho do outro.
- Uma requisição gera uma única linha de log, com tudo o que se sabe sobre ela.
- O código do erro diz onde a falha aconteceu, sem precisar ler texto nem ter traces.
- Os nomes dos campos seguem o padrão do OpenTelemetry, para que a adoção futura de traces não mude os logs.

Na prática:
- **Wide event (log):** uma linha por requisição ou execução de cron, com o contexto do domínio (IDs da entidade, loja) e, quando falha, o código do erro (`error.type`) e a causa. Responde "o que aconteceu com esta requisição". Poucos campos indexados, com poucos valores possíveis (`event`, `status`, `error.type`); a mensagem do erro (`exception.message`) é guardada, mas não indexada.
- **O código do erro é atribuído uma vez**, onde a falha é detectada (no gateway, ao falhar o banco ou um serviço externo; no caso de uso, ao violar uma regra). Quem chama só acrescenta texto, para não esconder o ponto da falha nem mudar a resposta. A exceção é quando o código de quem chama muda o que o cliente deve fazer (por exemplo, `SIGNUP_SIGNIN_FAILED`: a conta foi criada, basta entrar).
- **A causa é lida no topo:** quando o erro vem do banco, a linha registra o banco (`db.system.name`) e o código do Postgres (`db.response.status_code`, ex.: `23505` duplicado, `57014` timeout), sem que cada gateway precise informar. O detalhe do erro do banco, que traz os valores da linha que falhou, nunca vai para o log.
- **Spans (traces):** o caminho e o tempo dentro de uma requisição — qual consulta ou chamada externa rodou, em que ordem, quanto demorou. Ainda não adotados; a adoção é uma decisão própria, por exigir um backend de traces e ter custo de armazenamento. Até lá, tempos por chamada e listas de chamadas não entram no log.
- **Métricas (Prometheus):** contagens e taxas para painéis e alertas. Alertas olham métricas, não consultas de log.

## Consequências
Fica mais fácil:
- Encontrar e contar falhas por um único campo (`error.type`) e saber quando a causa é o banco, sem busca textual.
- Adotar OpenTelemetry depois: os campos de erro e de banco já têm os nomes dos atributos dos spans.

Fica mais difícil, e aceitamos:
- Consultas e alertas que usavam `error_code` e `error` precisam passar a usar `error.type` e `exception.message`.
- Os demais campos (`method`, `path`, `status`, `auth_user_id`, `version`) mantêm os nomes atuais até a adoção do OpenTelemetry, então a linha mistura os dois padrões por enquanto.
- Sem traces, a lentidão de uma chamada específica dentro de uma requisição não é visível; só a duração total.
- Detectar a causa no topo depende de todo erro ser encadeado com `%w`. Erros de conversão de dados e timeouts genéricos não são marcados como do banco.
