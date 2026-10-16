# akma-free5gc — documentação normativa

Análise normativa e decisões de projeto para a implementação do **3GPP AKMA (Authentication and Key Management for Applications)** integrado ao **free5GC**.

Este repositório contém a documentação que precede a implementação. O código vive em [`florianoneto07/akma-free5gc`](https://github.com/florianoneto07/akma-free5gc).

**Release alvo:** 3GPP Release 19.

---

## Princípio do trabalho

A prioridade é a especificação 3GPP. Implementações existentes — free5GC, AKMA+, trabalhos acadêmicos — são referências, não autoridade normativa. Onde uma fonte secundária diverge do padrão, segue-se o padrão e registra-se a divergência.

Toda afirmação nos documentos é classificada:

| Marcação | Significado |
|---|---|
| `[3GPP]` | Requisito extraído de especificação, com cláusula indicada |
| `[FREE5GC]` | Fato observado no código ou na documentação do free5GC |
| `[MVP]` | Simplificação proposta, com justificativa e impacto declarados |
| `[INTERPRETAÇÃO]` | Leitura adotada quando o texto admite mais de uma |
| `[LACUNA]` | Informação necessária não determinável pelas fontes |
| `[DIVERGÊNCIA]` | Inconsistência entre fontes ou dentro de uma fonte |

Nenhuma decisão do projeto é apresentada como exigência do 3GPP.

---

## Estrutura

```
docs/
├── 00-baseline-normativo.md          análise normativa completa
├── 01-adendo-A-stage3-udm.md         interfaces AUSF↔UDM e AAnF↔UDM
├── 02-adendo-B-sba-nrf.md            camada SBA, registro no NRF, tipos comuns
├── 03-adendo-C-ts29122.md            tipos de TS 29.122
├── 04-adendo-D-codificacao-uastar.md codificação de chaves e perfis Ua* Stage 3
├── baseline/especificacoes.md        versões fixadas
├── decisions/                        ADRs
├── defects/                          defeitos identificados em especificações
└── traceability/                     matriz cláusula → código → teste
reference/kdf/                        implementação de referência da KDF
```

Ordem de leitura sugerida: `00` → `01` → `02` → `03` → `04`. Os adendos são incrementais e não substituem o baseline.

---

## Decisões arquiteturais

| ADR | Assunto |
|---|---|
| [ADR-0001](docs/decisions/ADR-0001-topologia-de-implantacao.md) | Topologia de implantação — AAnF em VM separada do free5GC |
| [ADR-0002](docs/decisions/ADR-0002-codificacao-de-identificadores-e-chaves.md) | Codificação de identificadores e material de chave |
| [ADR-0003](docs/decisions/ADR-0003-identificador-ua-star.md) | Identificador de protocolo de segurança Ua* do MVP |

---

## Achados principais

### Duas advertências capazes de produzir falha silenciosa

**O prefixo `imsi-`.** TS 33.501 Anexo A.7.0 fixa o parâmetro P0 da KDF como o IMSI conforme TS 23.003 §2.2 — apenas os dígitos. O free5GC carrega `imsi-208930000000003` internamente. Usar a cadeia interna produz K_AKMA e A-TID incorretos e um K_AF divergente entre UE e rede, sem erro visível. Confirmado por três fontes independentes: TS 33.501, TS 29.503 e TS 29.571.

**O `AF_ID` tem duas formas.** A de transporte (TS 29.522) é `"<FQDN>.<10 hex>"`, com ponto. A de KDF (TS 33.535 A.4) é `FQDN ‖ 5 octetos brutos`, sem ponto. Nenhuma das duas especificações declara explicitamente que a forma de transporte não é o P0.

### Uma restrição de tipo pouco evidente

`routingId` é `string` com padrão `^[0-9]{1,4}$`, e TS 23.003 §2.2B estabelece que cada dígito é significativo — `"012"` não equivale a `"12"`. Convertê-lo para inteiro em qualquer ponto da cadeia UDM → AUSF → A-KID → AAnF descarta zeros à esquerda. A falha atinge apenas assinantes cujo RID começa com zero.

### Cinco defeitos em especificações

Consolidados em [`docs/defects/registro-de-defeitos-3gpp.md`](docs/defects/registro-de-defeitos-3gpp.md), com decisão adotada para cada um. O conjunto tem valor próprio como contribuição, independentemente da implementação.

### O Stage 3 do Ua* existe

TS 24.109 V19.0.0, Anexo K, especifica três perfis — HTTP Digest, TLS 1.2 PSK e TLS 1.3 PSK. Dois requisitos aparecem apenas ali e são omitidos por TS 33.535: o separador `";"` na `psk_identity` e a codificação base64 de K_AF no Digest.

---

## Questões em aberto

| Nº | Assunto | Estado |
|---|---|---|
| 1 | Obrigatoriedade da subscrição de `ROAMING_STATUS` — TS 33.535 §6.2.1 usa `shall`, TS 29.535 §5.1.8 torna `RoamingRestriction` opcional | Aberta |
| 5 | Presença de `AkmaInd`, `RoutingId`, `AANF` e `naanf-akma` em `free5gc/openapi` | Aberta — requer inspeção do repositório |
| A–E | Defeitos em especificações | Ver registro de defeitos |

---

## Fases

| Fase | Conteúdo | Estado |
|---|---|---|
| F0 | Aquisição documental e fixação do baseline | **Concluída** |
| F1 | Decisões de codificação | **Concluída** — ADR-0002, ADR-0003 |
| F2 | KDF, A-KID, AF_ID, com validação cruzada Go × Python | Próxima |
| F3 | AAnF autônomo — NRF, três *endpoints*, contexto, erros | |
| F4 | Integração AUSF e UDM | |
| F5 | AF e Ua* simplificado | |
| F6 | UE | |
| F7 | Ponta a ponta e testes negativos | |
| F8+ | Acesso anônimo, remoção de contexto, NEF, *roaming*, desempenho, verificação formal, *upstream* | |

---

## Referência da KDF

`reference/kdf/akma_kdf.py` implementa as três derivações do Anexo A de TS 33.535 diretamente a partir de TS 33.220 B.2, e serve de oráculo para validação cruzada com a implementação em Go.

TS 33.535 V19.0.0 **não fornece vetores de teste**. Todo vetor deste projeto é gerado pelo projeto e assim rotulado. A validação de correção obtém-se por implementações independentes, não por comparação com valores publicados pelo 3GPP.

---

## Recursos do projeto

| Categoria | Recurso |
|---|---|
| Código | `florianoneto07/akma-free5gc` |
| Core base | `free5gc/free5gc` |
| Ambiente Docker | `free5gc/free5gc-compose` |
| Solicitação *upstream* | `free5gc/free5gc` issue #1146 — 3GPP AKMA support |
| Discussão técnica | `open5gs/open5gs` discussion #4710 |
| AKMA+ (referência, não baseline) | Zenodo, DOI 10.5281/zenodo.14649847 |

**Sobre o AKMA+.** É uma modificação do protocolo, não uma implementação de TS 33.535. Seu uso neste projeto restringe-se a referência de organização de *testbed*, verificação cruzada de derivações comprovadamente idênticas às do padrão, e *related work*. Nenhum elemento que altere o protocolo entra no baseline.

---

## Licença

Documentação sob CC BY 4.0. As especificações 3GPP não são redistribuídas neste repositório; obtenha-as em https://www.3gpp.org/ftp/Specs/archive/.
