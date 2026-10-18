# AUSF: K_AKMA, A-KID e registro no AAnF

**MVP:** MVP-2 · **Data:** 17 de outubro de 2026 · **Aplicado e verificado na VM do free5GC** (192.168.16.130) em 17/10/2026

## 1. O que o AUSF precisa fazer

`[3GPP]` TS 33.535 §6.1:

- **passo 3:** *"If the AUSF receives the AKMA indication from the UDM, the AUSF shall store the KAUSF and generate the AKMA Anchor Key (KAKMA) and the A-KID from KAUSF after the primary authentication procedure is successfully completed"*;
- **passo 4:** o AUSF seleciona o AAnF (§6.7) e envia A-KID, K_AKMA e SUPI com Naanf_AKMA_AnchorKey_Register;
- *"The AUSF shall use the RID received from the UDM"*;
- **NOTA 1:** o AUSF não precisa guardar o material AKMA depois de entregá-lo.

Derivações: K_AKMA por TS 33.535 A.2, A-TID por A.3; SUPI na KDF sem `imsi-` (TS 33.501 A.7.0). O A-KID segue o formato NAI de TS 29.522 (ADR-0002). Tudo vem do `shared/akma` (MVP-0).

`[FREE5GC]` O AUSF v1.4.5 recebia a resposta do UDM, mas ignorava `akmaInd` e `routingId`.

## 2. Quando o AKMA é acionado

O gatilho é o fim bem-sucedido da autenticação primária, do ponto de vista da rede de origem:

| Método | Ponto no AUSF | Cláusula |
|---|---|---|
| 5G AKA | `Auth5gAkaComfirmRequestProcedure`: RES\* = XRES\* e resultado aceito pelo UDM, antes da resposta 200 ao AMF | TS 33.501 §6.1.3.2.0 passo 11 |
| EAP-AKA' | `EapAuthComfirmRequestProcedure`: RES = XRES, MAC correto, resultado aceito pelo UDM; junto com o EAP-Success | TS 33.501 §6.1.3.1 passo 10 |

`[INTERPRETAÇÃO]` Decisão D9: o registro no AAnF roda numa *goroutine*, sem atrasar nem alterar a resposta ao AMF. Falha no AAnF vai para o log (categoria `AKMA`) e não muda o resultado da autenticação. Cada registro tem prazo de 10 s, descoberta incluída.

## 3. Seleção do AAnF

`[3GPP]` TS 33.535 §6.7: o AUSF *"shall utilize the NRF to discover the AAnF instance(s) unless AAnF information is available by other means, e.g. locally configured"*. A seleção considera o RID (NOTA 1: o AUSF o recebe do UDM).

- **Sem `akma.aanfUri`:** descoberta no NRF com `target-nf-type=AANF`, `requester-nf-type=AUSF`, `service-names=naanf-akma` e `routing-indicator=<RID do UDM>` (TS 29.510 Tab. 6.2.3.2.3.1-1). Exige o NRF corrigido (`nrf.md`). Usa a primeira instância devolvida.
- **Com `akma.aanfUri`:** usa o endereço configurado.

O token OAuth2 tem escopo `naanf-akma` e destino `AANF`. O NRF corrigido preenche `aud` com `AANF`, que o AAnF exige.

## 4. Domínio de origem do A-KID

`[3GPP]` TS 33.535 §6.1: o *realm* do A-KID *"shall include Home Network Identifier"*; o formato é `5gc.mnc<MNC>.mcc<MCC>.3gppnetwork.org` (TS 23.003 §28.2). TS 23.003 §2.2: o comprimento do MNC (2 ou 3 dígitos) depende do MCC e não pode ser lido do IMSI.

`[INTERPRETAÇÃO]` O AUSF toma MCC e MNC do SUCI do pedido (`suci-0-<MCC>-<MNC>-…`). Quando o pedido trouxe SUPI, usa a única PLMN de `plmnSupportList` cujos MCC e MNC prefixam o IMSI. Se nenhuma ou mais de uma casar, o AKMA não é feito para aquela autenticação, e o motivo vai para o log.

## 5. Arquivos tocados

Árvore: `~/free5gc/NFs/ausf` (`free5gc/ausf` v1.4.5). Binário: `~/free5gc/bin/ausf`, com o original em `bin/ausf.orig`. Configuração: `~/free5gc/config/ausfcfg.yaml`. *Patches*: `free5gc/patches/ausf/0001-akma-anchor-key-register.patch` e `free5gc/patches/config/0001-ausf-akma.patch`.

| Arquivo | Mudança |
|---|---|
| `go.mod`, `go.sum` | `require` de `shared/akma`. `replace` de `shared/akma` e de `free5gc/openapi` para `../../../akma-free5gc/…`, isto é, `~/akma-free5gc`, com o `openapi` corrigido (`openapi.md`) |
| `pkg/factory/config.go` | seção `akma`: `enable` (padrão `false`) e `aanfUri` (opcional) |
| `internal/context/context.go`, `ausf_context_init.go` | `AkmaEnable`, `AkmaAanfUri`; no contexto do UE: `AkmaInd`, `RoutingId`, `SupiOrSuci` |
| `internal/logger/logger.go` | categoria `AKMA` |
| `internal/sbi/processor/ue_authentication.go` | guarda a indicação do UDM; chama `startAkma` nos dois pontos da seção 2 |
| `internal/sbi/processor/akma.go` (novo) | `startAkma`, `buildAkmaKeyInfo`, `homeNetwork` |
| `internal/sbi/processor/akma_test.go` (novo) | vetor `testbed-ue` do oráculo; descoberta e registro com NRF e AAnF simulados; condições de disparo; D9 |
| `internal/sbi/consumer/aanf_service.go` (novo), `consumer.go` | descoberta do AAnF, token, `register-anchorkey` |
| `config/ausfcfg.yaml` | `akma: { enable: true }` |

Logs (categoria `AKMA`): desde o MVP-3a, o AUSF registra em `info` cinco eventos: a indicação recebida do UDM, a derivação com o A-KID, o AAnF selecionado, o registro aceito e, se for o caso, a indicação ignorada por `akma.enable: false`. Falhas vão em `error`. K_AUSF e K_AKMA nunca aparecem. Catálogo em `logs-akma.md`.

## 6. Como reproduzir numa instalação limpa

Pré-requisitos: MVP-1 e MVP-1b, e o repositório em `~/akma-free5gc`. Os `replace` do `go.mod` são relativos a `~/free5gc/NFs/ausf`. Prepare a cópia do `openapi` com `free5gc/scripts/prepare-openapi.sh`. Com o free5GC parado:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh ausf apply
```

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh config apply
```

O primeiro aplica o *patch*, guarda `bin/ausf.orig` e recompila. O segundo acrescenta `akma.enable: true` ao `ausfcfg.yaml`. Arquivos não rastreados, como o `aanfcfg.yaml`, não bloqueiam o *patch*. Para desfazer, use `revert` nos dois alvos.

## 7. Como verificar

```bash
cd ~/free5gc/NFs/ausf && go test ./internal/sbi/processor/
```

Depois de registrar um UE com assinatura AKMA, o log do free5GC deve ter:

- `[AUSF][AKMA]`, nesta ordem: `AKMA indication received from the UDM`, `K_AKMA and A-KID derived`, `AAnF for Routing Indicator … (discovered at the NRF)` e `AKMA key material of imsi-… registered at the AAnF`;
- no NRF, a descoberta `target-nf-type=AANF` com `routing-indicator`, respondida com 200.

O fluxo completo está em `verificacao.md`.
