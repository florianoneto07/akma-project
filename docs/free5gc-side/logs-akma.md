# Mensagens AKMA no log do free5GC

**MVP:** MVP-3a · **Data:** 18 de outubro de 2026

## 1. Objetivo

`[MVP]` Os eventos principais do AKMA aparecem na tela do `run.sh`, no formato e nos níveis de log que as NFs do free5GC já usam. Assim, dá para acompanhar o fluxo sem abrir arquivos de log nem mudar o nível para `debug`. A 3GPP não especifica mensagens de log; tudo o que está aqui é decisão do projeto (D14 a D16, no plano).

O lado do UE (UERANSIM) está em `docs/af-side/ue-ueransim.md`.

## 2. Formato

`[FREE5GC]` As NFs usam o `logrus` com dois campos fixos, a NF e a categoria:

```
2026-10-18T10:15:42.123456789-07:00 [INFO][AUSF][AKMA] AKMA key material of imsi-208930000000001 registered at the AAnF
```

`[FREE5GC]` No arquivo de log (`-l`, no `run.sh` o `free5gc.log` de cada execução), as mesmas entradas ficam no formato `chave=valor`:

```
time="2026-10-18T10:15:42.123456789-07:00" level="info" msg="AKMA key material of imsi-208930000000001 registered at the AAnF" CAT="AKMA" NF="AUSF"
```

`[MVP]` As mensagens AKMA usam a categoria `AKMA` nas três NFs (`[UDM][AKMA]`, `[AUSF][AKMA]`, `[AANF][AKMA]`). Para filtrar no arquivo:

```bash
grep 'CAT="AKMA"' ~/free5gc/log/<data>/free5gc.log
```

O nível segue o costume do free5GC: `info` para o evento normal, `warn` para algo fora do esperado que não impede o fluxo, `error` para falha, `debug` para detalhe que só interessa a quem investiga.

## 3. Catálogo

### 3.1 UDM

| Nível | Mensagem | Quando | Referência |
|---|---|---|---|
| `info` | `AKMA indication sent to the AUSF for <SUPI>, Routing Indicator <RID>` | assinante com `akmaAllowed` e `routingId` válido | TS 33.501 §6.1.3.1 e §6.1.3.2.0, passo 2 |
| `debug` | `No AKMA subscription for <SUPI>` | assinante sem `akmaAllowed` | — |
| `warn` | `Routing Indicator of the SUCI (<RID>) differs from the stored routingId (<RID>) for <SUPI>` | o RID do SUCI e o do UDR diferem | D10 |
| `error` | `AKMA allowed for <SUPI>, but the stored routingId "<valor>" is not a Routing Indicator: AKMA indication not sent` | `routingId` ausente ou fora do padrão | TS 29.505 Tab. 5.4.2.2-1 |
| `info` | `AuthTrigger for <SUPI> from <tipo de NF> at <endereço>: re-authentication notified to the AMF <id>`, com `, pending until the UE is reachable` quando o AMF responde 202 | `trigger-auth` aceito (MVP-12, `udm.md` §6) | TS 33.501 §6.1.5.2 passos 1 e 2 |
| `warn` | `AuthTrigger for <SUPI> from <tipo de NF> at <endereço> refused: <status> <causa> (<detalhe>)` | `trigger-auth` recusado | TS 29.503 tabela 6.2.3.13.4.2.2-2 |

### 3.2 AUSF

| Nível | Mensagem | Quando | Referência |
|---|---|---|---|
| `info` | `AKMA indication received from the UDM for <SUPI>, Routing Indicator <RID>` | resposta do UDM com `akmaInd` | TS 33.535 §6.1 passo 2 |
| `info` | `AKMA indication for <SUPI> ignored: akma.enable is false` | autenticação concluída, AKMA desligado no `ausfcfg.yaml` | — |
| `info` | `K_AKMA and A-KID derived for <SUPI>: A-KID <A-KID>` | autenticação primária concluída com sucesso | TS 33.535 §6.1 passo 3 |
| `info` | `AAnF for Routing Indicator <RID>: <URI> (discovered at the NRF)` ou `(configured)` | seleção do AAnF | TS 33.535 §6.7 |
| `info` | `AAnF for the default Routing Indicator: <URI> (discovered at the NRF, any AAnF)` | RID `"0"`, descoberta sem `routing-indicator` (MVP-15) | TS 33.535 §6.7; TS 23.003 §2.2B |
| `info` | `AAnF for Routing Indicator <RID>: <URI> (search result reused)` | resultado guardado ainda válido (MVP-15) | TS 29.510 §5.3.2.2.1 |
| `info` | `AKMA key material of <SUPI> registered at the AAnF` | resposta 200 do AAnF | TS 33.535 §6.1 passos 4 e 5 |
| `error` | `AKMA key material for <SUPI> not generated: <motivo>` | SUPI, RID ou domínio de origem inválidos | — |
| `error` | `Naanf_AKMA_AnchorKey_Register for <SUPI> failed: <motivo>` | descoberta ou registro falhou; a autenticação não muda | D9 |

As três últimas mensagens `info` saem da *goroutine* do AKMA (D9). Por isso aparecem depois da resposta do AUSF ao AMF.

### 3.3 AAnF

| Nível | Mensagem | Quando | Referência |
|---|---|---|---|
| `info` | `Routing Indicators [<RID> …]; K_AF lifetime <duração>; AF_IDs authorized by local policy: [<AF_ID> …]` | partida | — |
| `info` | `AKMA context of <SUPI> registered, A-KID <A-KID>`, com ` (previous context replaced)` quando substitui | `register-anchorkey` | TS 33.535 §6.1 passo 5 |
| `info` | `K_AF of <SUPI> provided to AF_ID <AF_ID> for A-KID <A-KID>, expires <instante>`, com `, through the NEF at <endereço>` pelo NEF e, desde o MVP-15, `, without the SUPI` quando o SUPI não vai na resposta | `retrieve-applicationkey` atendido | TS 33.535 §6.2.1 passos 7 e 8; §6.3 passo 4 |
| `info` | `K_AKMA not present for A-KID <A-KID>, requested by AF_ID <AF_ID>` | A-KID desconhecido (403) | TS 33.535 §6.2.1 passo 2; TS 29.535 Tab. 5.1.7.3-1 |
| `warn` | `NF <NF Instance ID> at <endereço> not authorized for AF_ID <AF_ID>` | consumidor fora da política local (403); desde o MVP-11, a política é pelo certificado, e o endereço só fica no log | TS 33.535 §6.2.1 passo 2 |
| `warn` | `NEF <NF Instance ID> at <endereço> not authorized for AF_ID <AF_ID>` | AF_ID não aceito daquele NEF (403) | TS 33.535 §6.3 passo 3 |
| `info` | `AKMA context of <SUPI> removed` ou `AKMA context of <SUPI> to remove not found` | `remove-context` | TS 33.535 §6.6.1 |
| `error` | `K_AF derivation failed: <motivo>` | falha interna | — |
| `info` | `K_AKMA refresh requested from the UDM when a K_AF is asked for a K_AKMA older than <duração> (TS 33.535 cl. 6.4.4)` | partida, com `akma.kakmaRefresh` (MVP-12, `aanf.md` §11) | TS 33.535 §6.4.4 |
| `info` | `K_AKMA of <SUPI> registered <idade> ago: refresh requested from the UDM (TS 33.535 cl. 6.4.4)` | K_AF entregue para uma K_AKMA mais velha que a política | TS 33.535 §6.4.4 |
| `info` | `UDM for Nudm_UECM_AuthTrigger: <URI> (discovered at the NRF)` | primeira renovação | TS 29.510 |
| `info` | `Refresh of K_AKMA of <SUPI> accepted by the UDM: new primary authentication requested` | `trigger-auth` com 204 | TS 33.501 §6.1.5.2 passo 1 |
| `warn` | `Refresh of K_AKMA of <SUPI> not accepted by the UDM: <status> <causa> (<detalhe>)` | `trigger-auth` recusado ou sem resposta | TS 29.503 tabela 6.2.3.13.4.2.2-2 |
| `info` | `<método> <caminho>, User-Agent <valor>, 3gpp-Sbi-NF-Peer-Info <valor>` (categoria `SBI`; `none` sem o cabeçalho) | cada pedido recebido (MVP-15) | TS 29.500 tab. 5.2.2.2-1 e §5.2.3.2.21 |
| `info` | `Heart-beat sent to the NRF, next in <intervalo>` (categoria `Consumer`) | cada *heart-beat* aceito (MVP-15) | TS 29.510 §5.2.2.3.2 |
| `warn` | `Heart-beat: the NRF does not have the profile (404); registering again` ou `Heart-beat failed: <motivo>` (categoria `Consumer`) | *heart-beat* recusado (MVP-15) | TS 29.510 §5.2.2.3.2 |

### 3.4 AMF

Categoria `AKMA` nova no AMF, do *patch* `amf/0001` (MVP-12, `amf.md`).

| Nível | Mensagem | Quando | Referência |
|---|---|---|---|
| `info` | `Re-authentication of <SUPI> requested by the UDM: started (ngKSI <n>)` | notificação do UDM, UE conectado (204) | TS 33.501 §6.1.5.2 passos 3 e 4 |
| `info` | `Re-authentication of <SUPI> requested by the UDM: pending (UE not connected)` | UE em CM-IDLE (202) | TS 33.501 §6.1.5.2 passo 3 |
| `info` | `Re-authentication of <SUPI> requested by the UDM <tempo> ago: started (ngKSI <n>)` | autenticação pendente, depois de `Service Request` ou registro | TS 33.501 §6.1.5.2 passo 3 |
| `info` | `Re-authentication of <SUPI> requested by the UDM: completed, new KAMF in use (ngKSI <n>)` | `Security Mode Complete` | TS 33.501 §6.1.3.2 e §6.9.4.2 |
| `warn` | `Re-authentication of <SUPI> requested by the UDM: refused, <status> <causa>` | UE em outro procedimento ou desconhecido | TS 29.503 tabela 6.2.5.6-3 |
| `warn` | `Re-authentication of <SUPI> requested by the UDM: not started` | o AUSF não respondeu | — |

Além dessas, a categoria `GIN` registra cada requisição HTTP com método, caminho e código, como nas outras NFs (`| 200 | 127.0.0.1 | POST | /naanf-akma/v1/register-anchorkey |`). No NRF, a descoberta do AAnF aparece do mesmo jeito, com `target-nf-type=AANF&routing-indicator=<RID>` no caminho.

## 4. O que aparece e o que nunca aparece

- **Chaves:** K_AUSF, K_AKMA e K_AF nunca aparecem, em nível algum.
- **Identificadores:** `[MVP]` Decisão D16. SUPI, A-KID e AF_ID aparecem em `info`.
  - As NFs do free5GC já registram o SUPI em `info`.
  - O A-KID não é chave: o UE o envia em claro ao AF no Ua\*.
  - Com os três identificadores, dá para seguir o mesmo UE nas telas do UE, do UDM, do AUSF e do AAnF.
  - Até o MVP-2, o AUSF e o AAnF só mostravam os identificadores em `debug`.

## 5. AAnF na mesma tela

`[FREE5GC]` O `run.sh` do free5GC inicia uma lista fixa de NFs, com opções para as opcionais (`-bsf`, `-n3iwf`, `-tngf`). Todas escrevem no mesmo terminal e no mesmo `free5gc.log`.

`[MVP]` Decisão D15. O *patch* `free5gc/patches/run/0001-aanf-option.patch` faz duas mudanças:

- acrescenta ao `run.sh` a opção `-aanf`, no molde do `-bsf`, que põe `aanf` no fim da lista, depois do NRF;
- acrescenta `aanf` à lista do `force_kill.sh`.

Com a opção, o AAnF passa a usar o mesmo terminal e o mesmo arquivo de log das outras NFs, e para junto com elas no Ctrl-C:

```bash
cd ~/free5gc && ./run.sh -aanf
```

O `run-aanf.sh` continua servindo para rodar o AAnF à parte. Os dois não podem rodar juntos, porque o AAnF escuta em `192.168.16.130:8010`. O `run-aanf.sh` também passa a exportar `GIN_MODE=release`, como o `run.sh`, para não imprimir as linhas `[GIN-debug]` do Gin.

## 6. Arquivos tocados

| Onde | Arquivo | Mudança |
|---|---|---|
| UDM (`NFs/udm`) | `internal/logger/logger.go` | categoria `AKMA` (`AkmaLog`) |
| UDM | `internal/sbi/processor/akma.go` | mensagens da seção 3.1 |
| AUSF (`NFs/ausf`) | `internal/sbi/processor/ue_authentication.go` | indicação recebida do UDM |
| AUSF | `internal/sbi/processor/akma.go` | derivação, registro, AKMA desligado |
| AUSF | `internal/sbi/consumer/aanf_service.go` | AAnF selecionado (antes em `debug`, categoria `Consumer`) |
| AAnF (`free5gc/aanf`) | `internal/context/context.go` | linha de partida; `AfPolicy.AfIds` |
| AAnF | `internal/sbi/processor/{register_anchorkey,retrieve_applicationkey,remove_context}.go` | identificadores em `info` |
| free5GC (raiz) | `run.sh`, `force_kill.sh` | opção `-aanf` |
| repositório | `free5gc/scripts/run-aanf.sh` | `GIN_MODE=release` |
| repositório | `free5gc/scripts/patch-nf.sh` | alvo `run`; checagem de alterações locais só nos arquivos que os *patches* tocam |

As mudanças do UDM e do AUSF estão nos *patches* do MVP-2 (`free5gc/patches/{udm,ausf}/0001-*.patch`), regenerados.

## 7. Como reproduzir numa instalação limpa

Numa instalação limpa, os *patches* do MVP-2 já trazem as mensagens; basta aplicar o novo alvo:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh run apply
```

Numa instalação que já tem os *patches* do MVP-2, é preciso trocar os do UDM e do AUSF. O `revert` precisa dos *patches* antigos, então a ordem é esta, com o free5GC parado:

1. reverter com a cópia antiga do repositório:

   ```bash
   cd ~/akma-free5gc && free5gc/scripts/patch-nf.sh udm revert && free5gc/scripts/patch-nf.sh ausf revert
   ```

2. atualizar `~/akma-free5gc`;
3. aplicar de novo (recompila o UDM e o AUSF):

   ```bash
   cd ~/akma-free5gc && free5gc/scripts/patch-nf.sh udm apply && free5gc/scripts/patch-nf.sh ausf apply && free5gc/scripts/patch-nf.sh run apply
   ```

4. recompilar o AAnF e implantar (`build-aanf.sh` e `deploy.sh`, na máquina do repositório).

## 8. Como verificar

Com o AAnF parado no `run-aanf.sh`, inicie o free5GC com `./run.sh -aanf`, depois o gNB e o UE do UERANSIM. A sequência esperada na tela, filtrada por `[AKMA]`, está em `verificacao.md`, seção 4.
