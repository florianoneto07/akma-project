# Registro de alterações

## [Não lançado]

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
