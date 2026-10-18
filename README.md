# akma-free5gc

Análise normativa, decisões de projeto e implementação do **3GPP AKMA (Authentication and Key Management for Applications)** integrado ao **free5GC**.

Este repositório reúne a documentação e, a partir do MVP-0, o código (decisão D2 do [plano de desenvolvimento](docs/plano/plano-de-desenvolvimento.md), aprovado em 16/10/2026). O objetivo final é um tutorial: o que fazer depois de instalar o free5GC para usar AKMA.

**Release alvo:** 3GPP Release 19.

---

## Princípio do trabalho

A prioridade é a especificação 3GPP. Implementações existentes — free5GC, AKMA+ de Yang et al., trabalhos acadêmicos — são referências, não autoridade normativa. Onde uma fonte secundária diverge do padrão, segue-se o padrão e registra-se a divergência. Todo ponto em que o projeto não segue literalmente a norma é decidido e registrado em [`docs/conformidade/desvios.md`](docs/conformidade/desvios.md).

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
├── 05-adendo-E-ts23502-v19-9-0.md    versão de TS 23.502 no acervo
├── baseline/
│   ├── especificacoes.md             versões fixadas
│   ├── ambiente-verificado.md        free5GC, VMs e free5gc/openapi verificados
│   └── conferencia-normativa.md      conferência da documentação contra as normas
├── conformidade/desvios.md           desvios e escolhas entre textos normativos
├── decisions/                        ADRs
├── defects/                          defeitos identificados em especificações
├── literatura/                       análise dos artigos sobre AKMA
├── plano/                            plano de desenvolvimento por MVPs
├── traceability/                     matriz cláusula → código → teste
├── shared/                           KDF, identificadores, codificação, vetores (comum aos dois lados)
├── free5gc-side/                     o que se faz dentro do free5GC
├── af-side/                          AF, Ua* e lado UE
└── tutorial/                         o que fazer depois de instalar o free5GC
shared/akma/                          biblioteca Go comum (KDF, SUPI, A-KID, AF_ID, codecs)
shared/crosscheck/                    comparação com a KDF do free5GC (módulo próprio)
reference/kdf/                        oráculo Python e gerador de vetores
specs/  papers/                       normas e artigos — locais, fora do controle de versão
```

Ordem de leitura sugerida: `00` → `01` → `02` → `03` → `04` → `05`. Os adendos são incrementais e não substituem o baseline. As correções de 16/10/2026 estão marcadas no texto como *errata* e detalhadas em `docs/baseline/conferencia-normativa.md`.

---

## Decisões arquiteturais

| ADR | Assunto |
|---|---|
| [ADR-0001](docs/decisions/ADR-0001-topologia-de-implantacao.md) | Topologia de implantação — AAnF em VM separada do free5GC (revisto pelo ADR-0001-R2) |
| [ADR-0001-R2](docs/decisions/ADR-0001-R2-aanf-na-vm-do-free5gc.md) | AAnF na VM do free5GC; AF, Ua\* e repositório na VM do AF |
| [ADR-0002](docs/decisions/ADR-0002-codificacao-de-identificadores-e-chaves.md) | Codificação de identificadores e material de chave |
| [ADR-0003](docs/decisions/ADR-0003-identificador-ua-star.md) | Identificador de protocolo de segurança Ua* do MVP |

---

## Achados principais

### Duas advertências capazes de produzir falha silenciosa

**O prefixo `imsi-`.** TS 33.501 Anexo A.7.0 fixa o parâmetro P0 da KDF como o IMSI conforme TS 23.003 §2.2 — apenas os dígitos. O free5GC carrega `imsi-208930000000003` internamente. Usar a cadeia interna produz K_AKMA e A-TID incorretos e um K_AF divergente entre UE e rede, sem erro visível. Confirmado por três fontes independentes: TS 33.501, TS 29.503 e TS 29.571.

**O `AF_ID` tem duas formas.** A de transporte (TS 29.522) é `"<FQDN>.<identificador Ua*>"`, com ponto; a norma descreve o identificador como *"a string of 5 octet"*, e a forma em 10 caracteres hexadecimais vem do exemplo. A de KDF (TS 33.535 A.4) é `FQDN ‖ 5 octetos brutos`, sem ponto. Nenhuma das duas especificações declara explicitamente que a forma de transporte não é o P0.

### Uma restrição de tipo pouco evidente

`routingId` é `string` com padrão `^[0-9]{1,4}$`, e TS 23.003 §2.2B estabelece que cada dígito é significativo — `"012"` não equivale a `"12"`. Convertê-lo para inteiro em qualquer ponto da cadeia UDM → AUSF → A-KID → AAnF descarta zeros à esquerda. A falha atinge apenas assinantes cujo RID começa com zero.

### Nove defeitos em especificações

Consolidados em [`docs/defects/registro-de-defeitos-3gpp.md`](docs/defects/registro-de-defeitos-3gpp.md) (D-A a D-I), com decisão adotada para cada um. O conjunto tem valor próprio como contribuição, independentemente da implementação.

### O Stage 3 do Ua* existe

TS 24.109 V19.0.0, Anexo K, especifica três perfis — HTTP Digest, TLS 1.2 PSK e TLS 1.3 PSK. Dois requisitos aparecem apenas ali e são omitidos por TS 33.535: o separador `";"` na `psk_identity` e a codificação base64 de K_AF no Digest. O Anexo K também diverge de TS 33.535 no `realm` do Digest; o projeto segue TS 33.535 (`3GPP-bootstrapping-akma@<FQDN>`, desvio DV-6).

---

## Questões em aberto

| Nº | Assunto | Estado |
|---|---|---|
| 1 | Obrigatoriedade da subscrição de `ROAMING_STATUS` — TS 33.535 §6.2.1 usa `shall`, TS 29.535 §5.1.8 torna `RoamingRestriction` opcional | Aberta — no MVP, etapa configurável e desabilitada (desvio DV-2) |
| 5 | Presença de `AkmaInd`, `RoutingId`, `AANF` e `naanf-akma` em `free5gc/openapi` | **Encerrada** — todos presentes no `openapi` v1.2.4; faltam só os modelos do Naanf_AKMA (`docs/baseline/ambiente-verificado.md`) |
| A–E | Questões ligadas a defeitos em especificações | Ver o registro de defeitos, que traz a correspondência entre as letras A–E e os itens D-A a D-I |
| D10 | Origem do RID quando o UDM recebe SUPI | Depende de TS 29.505, ainda fora do acervo (plano, D8) |

---

## Fases

As fases passaram a MVPs no [plano de desenvolvimento](docs/plano/plano-de-desenvolvimento.md), aprovado em 16/10/2026.

| Fase | MVP | Conteúdo | Estado |
|---|---|---|---|
| F0 | — | Aquisição documental e fixação do baseline | **Concluída** |
| F1 | — | Decisões de codificação | **Concluída** — ADR-0002, ADR-0003 |
| F2 | MVP-0 | KDF, A-KID, AF_ID, com validação cruzada Go × Python × KDF do free5GC | **Concluída** em 16/10/2026 (`go test -race` atendido em 17/10/2026) |
| F3 | MVP-1, MVP-1b | AAnF autônomo — NRF, três *endpoints*, contexto, erros; correção do NRF do free5GC (D13) | **Concluída** em 17/10/2026 |
| F4 | MVP-2, MVP-3a | Integração AUSF e UDM; mensagens AKMA na tela do free5GC e do UERANSIM, com a derivação no ME (D14) | **Concluída** em 17/10/2026; MVP-3a em 18/10/2026 |
| F5 | MVP-3 | AF e Ua* simplificado | **Concluída** em 18/10/2026 |
| F6 | MVP-4 | UE (K_AF no ME do UERANSIM e cliente Ua*) | **Concluída** em 18/10/2026 (derivação de K_AKMA e A-KID antecipada no MVP-3a) |
| F7 | MVP-5 | Ponta a ponta e testes negativos — primeiro MVP funcional | |
| F8+ | MVP-6+ | Acesso anônimo, remoção de contexto, NEF, *roaming*, TLS, desempenho, verificação formal, *upstream* | |

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
| AKMA+ de Yang et al. (referência, não baseline) | Zenodo, DOI 10.5281/zenodo.14649847 |

**Sobre o AKMA+.** Há dois trabalhos com esse nome; o citado aqui é o de Yang et al. O "AKMA+" de Boureanu et al. (ASIA CCS '25) é outra proposta (`docs/literatura/analise-dos-artigos.md`). O de Yang et al. é uma modificação do protocolo, não uma implementação de TS 33.535. Seu uso neste projeto restringe-se a referência de organização de *testbed*, verificação cruzada de derivações comprovadamente idênticas às do padrão, e *related work*. Nenhum elemento que altere o protocolo entra no baseline.

---

## Licença

Documentação sob CC BY 4.0. As especificações 3GPP não são redistribuídas neste repositório; obtenha-as em https://www.3gpp.org/ftp/Specs/archive/ e coloque-as em `specs/` (ignorado pelo git). Os artigos citados também não são redistribuídos; ficam em `papers/`, igualmente ignorado.
