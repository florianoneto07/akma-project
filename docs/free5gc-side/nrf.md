# NRF: correções para o AAnF (MVP-1b)

**MVP:** MVP-1b · **Decisão:** D13 do plano · **Data:** 17 de outubro de 2026 · **Aplicado e verificado na VM do free5GC** (192.168.16.130) em 17/10/2026

## 1. Por que mexer no NRF

O MVP-1 encontrou três pontos em que o NRF do free5GC v4.2.3 (`free5gc/nrf` v1.4.5) não segue TS 29.510 (`registro-nrf.md` §5). Nos três, quem diverge da norma é o NRF; o AAnF segue a norma. O responsável pelo projeto decidiu corrigi-los antes do MVP-2. É no MVP-2 que o AUSF passa a escolher o AAnF, e TS 33.535 §6.7 manda que ele *"shall utilize the NRF to discover the AAnF instance(s) unless AAnF information is available by other means"*.

Cada correção é um *patch* separado, com seu teste, para poder ir ao free5GC isoladamente.

## 2. Arquivos tocados

Árvore: `~/free5gc/NFs/nrf` (checkout do `free5gc/nrf` v1.4.5, *commit* `8e567d4`). Binário: `~/free5gc/bin/nrf`, com o original guardado em `~/free5gc/bin/nrf.orig`.

| *Patch* (repositório) | Arquivo do NRF | Mudança |
|---|---|---|
| `free5gc/patches/nrf/0001-store-aanf-info-list.patch` | `internal/context/management_data.go` | `nnrfNFManagementOption` copia `aanfInfoList` para o perfil guardado |
| | `internal/context/management_data_aanf_test.go` (novo) | teste: o mapa é guardado, como cópia, com os RIDs como cadeias |
| `free5gc/patches/nrf/0002-discover-aanf-by-routing-indicator.patch` | `internal/sbi/processor/nf_discovery.go` | `AANF` entra na lista de tipos aceitos em `target-nf-type` e `requester-nf-type`. O `routing-indicator` com `target-nf-type=AANF` filtra por `aanfInfoList` |
| | `internal/sbi/processor/nf_discovery_aanf_test.go` (novo) | testes da validação e do filtro |
| `free5gc/patches/nrf/0003-access-token-audience.patch` | `internal/sbi/processor/access_token.go` | `aud` do token = `targetNfInstanceId` do pedido, ou o `targetNfType` quando o consumidor não pede instância |
| | `internal/sbi/processor/access_token_aud_test.go` (novo) | teste dos dois casos |
| `free5gc/patches/nrf/0004-sbi-tls.patch` | `pkg/factory/`, `pkg/service/`, `internal/sbi/` | MVP-11, fase (a): `sbi.tls` com o certificado da SBI, separado do `sbi.cert` que assina os *tokens*; o `sbitls` no servidor e no cliente (`tls-sbi.md` §6.5) |
| `free5gc/patches/nrf/0005-token-identity.patch` | `internal/context/`, `internal/sbi/`, `internal/util/peer.go` (novo), `pkg/factory/` | MVP-11, fase (b): o certificado da sessão TLS confere o pedido de *token*, o registro, a atualização, o descadastro e a descoberta; *tokens* em `ES256`, com `kid`, `x5t#S256` e `iss`; `allowedOperationsPerNfType`; saem `SignNFCert` e o `rootcert` (`tls-sbi.md` §7.1 a §7.3) |

## 3. O que cada correção faz e por quê

### N1 — guardar `aanfInfoList`

`[3GPP]` TS 29.510 Tab. 6.1.6.2.2-1 traz `aanfInfoList` (mapa de `AanfInfo`) no `NFProfile`. §5.2.2.2.2 manda o NRF guardar no perfil até os atributos que não conhece (*"those attributes shall be stored as part of the NF's profile data in NRF"*). `[FREE5GC]` `nnrfNFManagementOption` copia um a um os `*Info` que conhece (`udrInfo`, `udmInfo`, `ausfInfo`…) e descartava os demais. A correção acrescenta `aanfInfoList` a essa lista.

### N2 — descobrir `AANF` pelo Routing Indicator

`[3GPP]` TS 29.510 Tab. 6.2.3.2.3.1-1: `target-nf-type` é um `NFType`, e `routing-indicator` *"may be included if the target NF type is "AUSF", "AANF" or "UDM""*. Tab. 6.1.6.2.73-1: *"If not provided, the AAnf can serve any Routing Indicator"*.

`[FREE5GC]` `validateQueryParameters` aceita só uma lista fixa de tipos, sem `AANF`, e respondia 400. O filtro de `routing-indicator` tratava só AUSF e UDM.

A correção acrescenta `AANF` à lista. `[INTERPRETAÇÃO]` A troca da lista pela validação dos tipos aceitos no registro (`validNfType`) foi considerada e descartada: essa outra lista não tem `N3IWF`, e a troca faria a descoberta recusar um tipo que hoje aceita. O filtro casa um perfil quando:

1. não há `aanfInfoList`; ou
2. alguma entrada do mapa lista o RID; ou
3. alguma entrada não tem `routingIndicators`.

Como as chaves do mapa são livres, o filtro usa `$expr` com `$objectToArray` (MongoDB 3.6 ou mais recente; a VM tem 5.0.34). O RID é comparado como cadeia, então `"0000"` não casa com `"000"`.

### N3 — `aud` no token

`[3GPP]` TS 29.510 Tab. 6.3.5.2.4-1: `aud` é obrigatório e *"shall contain the NF service producer's NF instance ID(s) (if the exact NF instance(s) of the NF service producer is known) or the NF type of NF service producers"*. TS 33.501 §13.4.1.1.2 passo 2: o produtor confere o `aud`.

`[FREE5GC]` O NRF copiava `targetNfInstanceId`, que os clientes do free5GC não enviam, e o `aud` saía vazio.

A correção usa o tipo de NF (em maiúsculas) quando não há instância. Nenhuma NF do free5GC lê `aud` (`oauth.VerifyOAuth` não o confere), então nada muda para elas. O AAnF passa a exigir `aud` igual ao seu `nfInstanceId` ou a `AANF`, o que encerra o DV-8.

## 4. Como reproduzir numa instalação limpa

Na máquina do free5GC, com este repositório disponível (ou só `free5gc/patches/nrf/` e `free5gc/scripts/patch-nf.sh` copiados para ela):

```bash
free5gc/scripts/patch-nf.sh nrf apply
```

O *script*:

1. confere que `~/free5gc/NFs/nrf` não tem alterações locais;
2. confere e aplica os três *patches* com `git apply`;
3. guarda o binário original em `bin/nrf.orig` (só na primeira vez);
4. recompila com `make nrf`.

Se os *patches* já estiverem aplicados, ele só recompila. Para outra raiz do free5GC, use `FREE5GC_DIR`. Para outro diretório de *patches*, passe-o como terceiro argumento. Reinicie o NRF (ou o free5GC inteiro) para usar o binário novo.

Para desfazer:

```bash
free5gc/scripts/patch-nf.sh nrf revert
```

Reverte os *patches* em ordem inversa e restaura `bin/nrf` a partir de `bin/nrf.orig`. `patch-nf.sh nrf status` informa o estado.

## 5. Como verificar

1. Testes do NRF, na árvore corrigida:

```bash
cd ~/free5gc/NFs/nrf && go test ./internal/context/ ./internal/sbi/processor/
```

2. Com NRF e AAnF em execução: o perfil guardado tem `aanfInfoList`. Use a consulta de `registro-nrf.md` §7.
3. Descoberta com a ferramenta do projeto, na máquina do free5GC:

```bash
~/akma-tools/seed discover --nrf https://nrf.5gc.mnc093.mcc208.3gppnetwork.org:8000 --routing-indicator 0000
```

Deve listar o AAnF. Com `--routing-indicator 0001` ou `000`, a lista vem vazia (`null`).

4. Token: `seed register --nrf …` (`aanf.md` §9) é aceito pelo AAnF, que agora exige `aud`.

## 6. Resultado em 17/10/2026

| Verificação | Resultado |
|---|---|
| `go test` do NRF, na VM do AF e na VM do free5GC | passa |
| `aanfInfoList` no perfil guardado (MongoDB) | `{'1': {routingIndicators: ['0000']}}` |
| descoberta `routing-indicator=0000` / `0001` / `000` | AAnF / `null` / `null` |
| descoberta sem `routing-indicator`, `requester-nf-type=AF` | AAnF |
| semântica do filtro no MongoDB 5.0, banco temporário `akma_mvp1b_test` (apagado), quatro perfis: A sem lista, B com `["0000"]`, C com entrada sem RIDs, D com `["0001"]` e `["0000","12"]` | `0000` → A, B, C, D; `12` → A, C, D; `0001` → A, C, D; `9` → A, C; `012` → A, C |
| `register-anchorkey` com token do NRF corrigido e AAnF estrito | 200 |
| K_AF pedida da VM do AF | igual ao vetor `testbed-ue` |
| regressão parcial: `amf smf udr pcf udm nssf ausf chf nef` com o NRF corrigido | as 9 NFs se registram (201); 12 tokens emitidos (200); nenhum aviso nem erro do NRF |

**Regressão completa pendente.** O registro de um UE (UPF, gNB e UE do UERANSIM) exige `sudo` na VM e fica para o responsável. Sequência: `./run.sh` em `~/free5gc`; depois o gNB e o UE do UERANSIM, como no teste de instalação; confirmar o registro e o PDU *session* do UE.

## 7. Fora do escopo, registrado

- **`complexQuery`.** A descoberta por `complexQuery` (segundo trecho `[Query-27]` de `nf_discovery.go`) continua sem tratar `routing-indicator` para `AANF`, como já acontece para os tipos que não são AUSF nem UDM.
- **NEF.** `[FREE5GC]` O NEF v4.2.3 pede o token de `nnrf-nfm` com `targetNfType=NEF` em vez de `NRF`. O NRF recusa (`invalid_scope`), e o NEF não se descadastra ao encerrar. A falha aparece igual no log de 15/10/2026, com o NRF original, e não tem relação com estes *patches*.
- **Versão do binário.** O `make nrf` grava no binário a versão `v4.2.3` e o *commit* do submódulo, sem indicar os *patches*. O estado da árvore (`git -C ~/free5gc/NFs/nrf status`) e a presença de `bin/nrf.orig` mostram que o NRF está corrigido.
