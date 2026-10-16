# Registro de alterações

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
