# Registro de alterações

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
