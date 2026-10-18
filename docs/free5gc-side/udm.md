# UDM: indicação AKMA e Routing Indicator

**MVP:** MVP-2 · **Data:** 17 de outubro de 2026 · **Aplicado e verificado na VM do free5GC** (192.168.16.130) em 17/10/2026

## 1. O que o UDM precisa fazer

`[3GPP]` TS 33.501 §6.1.3.1 passo 2 (EAP-AKA') e §6.1.3.2.0 passo 2 (5G AKA): *"If a subscriber has an AKMA subscription, the UDM shall include the AKMA indication and Routing indicator in the Nudm_UEAuthentication_Get Response."* TS 29.503 Tab. 6.3.6.2.3-1, `AuthenticationInfoResult`:

- `akmaInd`: *"true: AKMA keys need to be generated for the UE"*;
- `routingId`: *"shall be present if the akmaInd is set to "true""* e contém *"the Routing Indicator contained in the SUCI"*.

`[FREE5GC]` O UDM v1.4.5 não preenchia nenhum dos dois. O `free5gc/openapi` v1.2.4 já tem os campos, tanto no resultado (`UdmUeauAuthenticationInfoResult.AkmaInd`, `RoutingId`) quanto na assinatura lida do UDR (`AuthenticationSubscription.AkmaAllowed`, `RoutingId`).

## 2. De onde vem o RID (decisão D10)

`[3GPP]` TS 29.505 Tab. 5.4.2.2-1, `AuthenticationSubscription` no UDR:

- `akmaAllowed` indica se o assinante pode usar AKMA;
- `routingId` *"shall be present if the akmaAllowed is set to "true""* e *"shall include the Routing Indicator contained in the SUCI"*. O YAML define o padrão `^[0-9]{1,4}$`, uma cadeia de dígitos.

`[MVP]` Decisão D10 (17/10/2026): o RID enviado ao AUSF é o `routingId` guardado no UDR. Ele existe também quando o UDM recebe SUPI em vez de SUCI, por exemplo numa reautenticação depois de registro por GUTI. Monitoramento pedido pelo responsável: quando o pedido traz SUCI, o UDM compara o RID do SUCI com o do UDR e registra no log (`warning`) qualquer diferença. O valor enviado continua sendo o do UDR.

Com `akmaAllowed: true` e `routingId` ausente ou fora do padrão, o dado viola TS 29.505, e o UDM não envia a indicação AKMA. `[MVP]` Isso fica registrado como `error` no log. A autenticação segue normalmente.

## 3. Arquivos tocados

Árvore: `~/free5gc/NFs/udm` (`free5gc/udm` v1.4.5). Binário: `~/free5gc/bin/udm`, com o original em `bin/udm.orig`. *Patch*: `free5gc/patches/udm/0001-akma-indication-and-routing-id.patch`.

| Arquivo | Mudança |
|---|---|
| `pkg/suci/suci.go` | nova função exportada `RoutingIndicator(suci)`: o RID do SUCI como cadeia, com todos os dígitos (TS 23.003 §2.2B) |
| `pkg/suci/suci_test.go` | RIDs `"0"`, `"012"`, `"0000"` e `"12"` (MNC de 3 dígitos) |
| `internal/sbi/processor/akma.go` (novo) | `setAkmaIndication`: preenche `akmaInd` e `routingId` e faz a comparação de D10 |
| `internal/sbi/processor/generate_auth_data.go` | chama `setAkmaIndication` antes de responder, nos dois métodos de autenticação |
| `internal/sbi/processor/akma_test.go` (novo) | casos da função; procedimento completo com UDR simulado, com e sem `akmaAllowed` |
| `internal/logger/logger.go` | categoria `AKMA` (MVP-3a): as mensagens AKMA do UDM saem como `[UDM][AKMA]`; até o MVP-2 saíam como `[UDM][UEAU]` |

## 4. Como reproduzir numa instalação limpa

Com o repositório em `~/akma-free5gc`, ao lado de `~/free5gc`, e o free5GC parado:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh udm apply
```

Aplica o *patch*, guarda `bin/udm.orig` e recompila. Para desfazer, use `patch-nf.sh udm revert`.

## 5. Como verificar

```bash
cd ~/free5gc/NFs/udm && go test ./internal/sbi/processor/ ./pkg/suci/
```

Com o banco de testes: `TestAkma5GAKA`, `TestAkmaEAPAKAPrime` e `TestAkmaWithoutSubscription` (`verificacao.md`). No log do free5GC, `[UDM][AKMA] AKMA indication sent to the AUSF for imsi-…, Routing Indicator …` confirma o envio (MVP-3a, `logs-akma.md`). A ausência de `Routing Indicator of the SUCI ... differs` confirma que o RID do SUCI e o do UDR coincidem.

Resultado em 17/10/2026: os testes do UDM passam. Com UERANSIM (SUCI `suci-0-208-93-0000-…`), o UDR tinha `routingId: "0000"` e não houve divergência registrada.
