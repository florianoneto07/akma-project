# Registro de alterações

## [Não lançado]

### MVP-4 — lado UE: K_AF no ME e Ua\* pelo plano de usuário (18/10/2026)
- `af/ue/ueransim/patches/0003-akma-kaf-nr-cli-selftest.patch`:
  - o `nr-ue` deriva a K_AF sob demanda (TS 33.535 A.4, §6.2.1 passo 1) e a entrega, com o A-KID, pelo comando `akma-kaf` do `nr-cli`; a K_AKMA nunca sai do ME (§4.4.3);
  - FQDN não ASCII é recusado (DV-4);
  - `nr-akma-selftest` confere K_AKMA, A-KID e K_AF contra os vetores do MVP-0.
- `0004-free5gc-ue-gnb-on-veth.patch`: o UE procura o gNB no `veth` do seu *namespace*.
- `af/scripts/ue-netns.sh`: *namespace* do UE (D5), com `up`, `ue`, `cli`, `ua` e `down`.
- `af/scripts/patch-ueransim.sh`: confere a série numa cópia temporária, compila um alvo por vez com o `cmake` que configurou a árvore e também compila o autoteste.
- Cliente Ua\*: obtém A-KID e K_AF do ME pelo `nr-cli` (`-nr-cli`, `-ue`) e liga as conexões ao endereço da sessão PDU (`-bind`).
- Plano de usuário corrigido. O UE passou a rodar no próprio *namespace*, porque no do host o endereço da sessão é local e o kernel descartava os pacotes que saíam da `upfgtp`. `free5gc/scripts/n6-nat.sh` aplica o encaminhamento e o NAT do N6.
- Verificação na VM:
  - autoteste com 5 vetores aprovados e 3 de NFKC recusados (DV-4);
  - UE alcança 8.8.8.8 e o AF pelo `uesimtun0`;
  - A-KID igual no UE, no AUSF, no AAnF e no AF; Ua\* com a K_AF do UE, pelo plano de usuário;
  - nova autenticação troca o A-KID, e o antigo recebe 403 do AAnF.
- DV-4 implementado. Documentação:
  - `docs/af-side/{ue-ueransim,ua-client,teste-ue}.md` e `docs/free5gc-side/plano-de-usuario.md`;
  - tutorial, partes 5 (*namespace*) e 7; matriz, revisão 3.7.

### MVP-3 — AF e Ua\* simplificado (18/10/2026)
- `docs/af-side/ua-star.md`: especificação do Ua\* simplificado, aprovada antes do código. É o HTTP Digest de TS 24.109 K.2 sem TLS, com:
  - `username` = A-KID e senha = base64(K_AF);
  - `realm = "3GPP-bootstrapping-akma@<FQDN>"` (DV-6), SHA-256 (DV-5) e `rspauth` obrigatório;
  - FQDN conferido no `Host` (TS 33.220 §4.5.2) e iniciação do AKMA pelo desafio (§6.5);
  - política `refetch` ou `reject` para a K_AF vencida (§6.4.2).
- `af/`: módulo Go com:
  - `uastar`: Digest de RFC 7616, compartilhado;
  - o AF (`af/af`): servidor Ua\*, `Naanf_AKMA_ApplicationKey_Get` pelo cliente do `openapi` corrigido e K_AF guardada até `expiry`;
  - o cliente Ua\* (`af/ue/uaclient`), que recusa AF sem `rspauth` válido;
  - `af/config/afcfg.yaml` e `af/scripts/build-af.sh`.
- `af/uastar/testdata/gen_digest_vectors.py`: implementação Python independente do Digest; reproduz o exemplo de RFC 7616 §3.9.1 e gera os vetores do projeto.
- Verificação com o AAnF da VM do free5GC e um contexto semeado:
  - a K_AF do oráculo autentica;
  - K_AF errada → 401; A-KID desconhecido → 403; outro `Host` → 421; AF_ID fora da política do AAnF → 503;
  - nenhuma chave nos logs.
- DV-5 e DV-6 implementados. Documentação: `docs/af-side/{ua-star,af,teste-af}.md`; tutorial, parte 6; matriz, revisão 3.6.

### MVP-3a — mensagens AKMA na tela e derivação no ME (18/10/2026)
- Pedido do responsável: acompanhar o AKMA pelas telas do `run.sh` e do UE. Decisões D14 (escopo; derivação no ME antecipada do MVP-4), D15 (`run.sh -aanf`) e D16 (SUPI, A-KID e AF_ID em `info`).
- `free5gc/patches/udm/0001-*.patch` e `free5gc/patches/ausf/0001-*.patch` regenerados:
  - categoria `AKMA` no UDM;
  - no UDM, a indicação enviada em `info`;
  - no AUSF, em `info`: indicação recebida, K_AKMA e A-KID derivados (com o A-KID), AAnF selecionado e registro.
- AAnF: identificadores em `info` e uma linha de partida com RIDs, tempo de vida da K_AF e AF_ID autorizados.
- `free5gc/patches/run/0001-aanf-option.patch`: opção `-aanf` no `run.sh` e `aanf` no `force_kill.sh`. `patch-nf.sh` ganha o alvo `run` e passa a checar alterações locais só nos arquivos que os *patches* tocam. `run-aanf.sh` exporta `GIN_MODE=release`.
- `af/ue/ueransim/patches/`: o `nr-ue` deriva K_AKMA, A-TID e A-KID (TS 33.535 §6.1, A.2, A.3) ao ativar o contexto de segurança da autenticação primária, com a opção `akma`, e mostra o A-KID em `[nas] [info]`. `af/scripts/patch-ueransim.sh` aplica, reverte e recompila só o `nr-ue`.
- Verificação na VM: o A-KID do UE é igual ao do AUSF e do AAnF em duas autenticações. Sem `akma`, o UE não deriva nada. Não há chave nos logs nem nas telas. Os testes `TestAkma*` continuam passando.
- Documentação: `docs/free5gc-side/logs-akma.md`, `docs/af-side/ue-ueransim.md`, `verificacao.md` §4; tutorial, parte 5 (as seguintes renumeradas); matriz, revisão 3.5.

### MVP-2 — UDM e AUSF (17/10/2026)
- `free5gc/patches/udm/0001-akma-indication-and-routing-id.patch`: o UDM preenche `akmaInd` e `routingId` (TS 33.501 §6.1.3.1 e §6.1.3.2.0; TS 29.503 Tab. 6.3.6.2.3-1) a partir de `akmaAllowed` e `routingId` do UDR (TS 29.505 Tab. 5.4.2.2-1). Compara com o RID do SUCI e registra divergências (D10).
- `free5gc/patches/ausf/0001-akma-anchor-key-register.patch`: depois do 5G AKA ou do EAP-AKA' bem-sucedido, o AUSF deriva K_AKMA e o A-KID e os registra no AAnF (TS 33.535 §6.1), sem bloquear a autenticação (D9). O AAnF é descoberto no NRF pelo Routing Indicator ou configurado (`akma.aanfUri`).
- `free5gc/patches/config/0001-ausf-akma.patch`: `akma.enable` no `ausfcfg.yaml`.
- `free5gc/patches/test/0001-akma-tests.patch`: AAnF no banco de testes do free5GC; testes `TestAkma5GAKA`, `TestAkmaEAPAKAPrime` e `TestAkmaWithoutSubscription`, que rodam sem `sudo`.
- `free5gc/scripts/provision-akma-subscriber.sh`; `prepare-openapi.sh` com argumento de versão; `patch-nf.sh` com os alvos `test` e `config`.
- Verificação: K_AF do AAnF igual à derivada pelo UE de teste nos dois métodos; A-KID `rid0000…@5gc.mnc093.mcc208.3gppnetwork.org` com UERANSIM; reautenticação substitui o contexto; AAnF parado não impede o registro; `TestRegistration`, `TestGUTIRegistration`, `TestReSynchronization` e `TestEAPAKAPrimeAuthentication` continuam passando.
- Documentação: `docs/free5gc-side/{udm,ausf,assinatura-akma,verificacao}.md`; tutorial, parte 4; matriz, revisão 3.4.

### MVP-1b — correção do NRF do free5GC (17/10/2026)
- `free5gc/patches/nrf/`: três *patches* com testes sobre o `free5gc/nrf` v1.4.5:
  - 0001: guarda `aanfInfoList` (TS 29.510 §5.2.2.2.2, Tab. 6.1.6.2.2-1);
  - 0002: descoberta de `AANF` com filtro por `routing-indicator` sobre `aanfInfoList` (Tab. 6.2.3.2.3.1-1, Tab. 6.1.6.2.73-1);
  - 0003: `aud` do token com a instância ou o tipo do produtor (Tab. 6.3.5.2.4-1).
- `free5gc/scripts/patch-nf.sh`: aplica, reverte e informa *patches* de uma NF do free5GC e a recompila, guardando o binário original.
- AAnF: `aud` passa a ser exigido; DV-8 encerrado. `seed discover` aceita `--routing-indicator`.
- Aplicado na VM do free5GC. Verificados `aanfInfoList` guardado, descoberta por RID (com o filtro conferido no MongoDB 5.0) e token com `aud`. Regressão parcial com nove NFs; a regressão com UE fica com o responsável (exige `sudo`).
- Decisões: D13 revista (NRF corrigido agora); D8 (TS 29.505 V19.6.0 incorporada ao baseline); D10 (RID vem do `routingId` do UDR, com comparação ao RID do SUCI em log).
- Documentação: `docs/free5gc-side/nrf.md`; tutorial, seção 3.3; matriz, revisão 3.3.

### MVP-1 — AAnF autônomo (17/10/2026)
- `free5gc/patches/openapi/`: `0001-naanf-akma-models.patch` (`AkmaKeyInfo`, `CtxRemove`, `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif`) e `0002-naanf-akma-client.patch` (cliente `aanf/AKMA`) sobre `free5gc/openapi` v1.2.4; `free5gc/scripts/prepare-openapi.sh`.
- `free5gc/aanf/`: NF AAnF com a estrutura das NFs do free5GC. Três operações do Naanf_AKMA (TS 29.535); contexto em memória com substituição do par (A-KID, K_AKMA) e K_AF reaproveitada até expirar (TS 33.535 §6.1, §6.2.1, §5.2); política local que associa AF_ID a endereço de origem; OAuth2 com verificação de `exp` e `aud`; erros de TS 29.500; registro e cancelamento no NRF.
- `free5gc/aanf/tools/seed/`: ferramenta de teste nos papéis do AUSF (`register`) e do AF (`retrieve`), mais `remove` e `discover`.
- `free5gc/config/aanfcfg.yaml`; `free5gc/scripts/{build-aanf,run-aanf,deploy}.sh`.
- Verificação na VM do free5GC e a partir da VM do AF: K_AF igual ao vetor `testbed-ue`; 403 para AF fora da política e para AF_ID pedido de outra máquina.
- Decisões: DV-8 (`aud` vazio tolerado no token, por limitação do NRF do free5GC); D13 (limitações do NRF v1.4.5 — `aanfInfoList` descartado, descoberta de `AANF` recusada, `aud` vazio — a corrigir com *patch* futuro). DV-1, DV-2 e DV-7 implementados.
- Documentação: `docs/free5gc-side/{aanf,openapi,registro-nrf}.md`; tutorial, parte 3; matriz de rastreabilidade, revisão 3.2.
- MVP-0: `go test -race` passa em `shared/akma` (critério 1 atendido por inteiro).

### MVP-0 — biblioteca AKMA comum (16/10/2026)
- `shared/akma`: KDF de TS 33.220 B.2 com NFKC; SUPI como parâmetro (TS 33.501 A.7.0, só IMSI — DV-3); Routing Indicator, Home Network Domain e A-KID (TS 23.003, TS 29.522); AF_ID nas formas da KDF e de transporte; K_AKMA, A-TID e K_AF (TS 33.535 Anexo A); codecs por camada (ADR-0002).
- `shared/crosscheck`: comparação com `ueauth.GetKDFValue` do free5GC (`free5gc/util` v1.3.2).
- `reference/kdf/gen_vectors.py` e `shared/akma/testdata/vectors.json`: 8 vetores do projeto, igualdades esperadas e contraexemplos; não são vetores 3GPP.
- Decisão P9: identificador Ua\* no `AfId` em maiúsculas na escrita (errata no `00` §6.4 e no ADR-0003).
- Documentação: `docs/shared/` (KDF, identificadores, codificação, vetores), `docs/tutorial/` (pré-requisitos e conferência das derivações), `docs/free5gc-side/` e `docs/af-side/` criados.

## [2.1] — 2026-10-16

Verificação do ambiente, conferência contra o texto das normas, revisão da topologia e plano de desenvolvimento aprovado.

### Novo
- `docs/baseline/ambiente-verificado.md`: free5GC v4.2.3 e submódulos fixados; `free5gc/openapi` v1.2.4 já contém `AANF`, `naanf-akma`, `AanfInfo`, `AkmaInd` e `RoutingId` (questão 5 encerrada).
- `docs/baseline/conferencia-normativa.md`: conferência da documentação contra TS 33.535, TS 29.535 e as demais normas do baseline.
- ADR-0001-R2: AAnF na VM do free5GC.
- Adendo E: TS 23.502 do acervo na V19.9.0; baseline mantido em V19.8.0.
- `docs/literatura/analise-dos-artigos.md`: 14 trabalhos sobre AKMA.
- `docs/conformidade/desvios.md`: desvios DV-1 a DV-7.
- `docs/plano/plano-de-desenvolvimento.md`: MVP-0 a MVP-5, aprovado.

### Revisto
- `realm` do Digest do Ua*: passa a seguir TS 33.535 B.1.2.2, `3GPP-bootstrapping-akma@<FQDN>` (DV-6; defeito D-C).
- Registro de defeitos: correspondência entre as letras A–E e D-A a D-I; novos D-F (RFC 3310 em TS 24.109 K.2), D-G (204 inalcançável em TS 29.535), D-H e D-I (editoriais).
- ADR-0001 marcado como revisto; ADR-0002 e ADR-0003 aceitos; ADR-0003 registra por que não usa `(0x01,0x00,0x00,0x00,0x02)`.
- Matriz de rastreabilidade, revisão 3.1.
- Normas e artigos ficam fora do controle de versão (`/specs/`, `/papers/`).

### Corrigido (errata marcada no texto)
- `00` §2.5: espaço espúrio na cadeia `S` do A-TID.
- `00` §5.1: guarda de K_AUSF citada em TS 33.501 §6.1.1.1 e §6.2.2.1 (não §6.2.2).
- `00` §2.4: forma hexadecimal do identificador Ua\* no `AfId` reclassificada como interpretação; tabela de H.3 completada.
- Adendo A §A.5: `immediateFlag` em TS 29.503 §5.5.2.2.2, com relatório imediato não garantido.
- Adendo B: marcação `[HIPÓTESE]` substituída; lista de cabeçalhos obrigatórios de TS 29.500 completada.
- Adendo D: precedente ancorado de `Kausf` (TS 29.503); conflito RFC 3310 × RFC 7616 (DV-5); decisão do `realm` revista.

### Removido
- `scripts/bootstrap-repo.sh`: fazia commits e criava repositório público, o que contraria as regras do projeto; o repositório já está inicializado.

## [2.0] — 2026-09-02

Consolidação da fase F0 (aquisição documental) e da fase F1 (decisões de codificação).

### Baseline
- Release 19 fixado. Dezesseis especificações incorporadas.
- Par corpo/OpenAPI da TS 29.535 fixado em V19.5.0 + `TS29535_Naanf_AKMA.yaml` v1.2.0.

### Encerrado
- Definição da KDF — TS 33.220 Anexo B.2.0.
- Codificação do SUPI como parâmetro de entrada — TS 33.501 Anexo A.7.0.
- Identificador de protocolo de segurança Ua* — TS 33.220 Anexo H.
- Tipos `AKId`, `AfId`, `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif` — TS 29.522.
- Existência e Stage 3 de `Nudm_UECM_AuthTrigger` — TS 33.501 §14.2.6, TS 29.503 §5.3.2.16.
- Reconhecimento do AAnF como consumidor do UDM — TS 29.503 §5.2.2.1.
- `NFType` `AANF`, `AanfInfo` e descoberta por `routing-indicator` — TS 29.510.
- Codificação de `kAkma` e `kaf` — por precedente de TS 29.509.
- Equivalência de `DateTime` entre TS 29.122 e TS 29.571.

### Revisto
- Codificação do valor do A-TID: recomendação alterada de hexadecimal para base64, por analogia com o B-TID do GBA (TS 33.220 §4.5.2). Registrado em ADR-0002.

### Novo
- Stage 3 do ponto de referência Ua* — TS 24.109 Anexo K, três perfis.
- Registro de cinco defeitos em especificações.
- ADR-0001, ADR-0002 e ADR-0003.
- Matriz de rastreabilidade, revisão 3.
