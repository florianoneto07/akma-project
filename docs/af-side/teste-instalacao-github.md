# Instalação pelo repositório público numa VM limpa (MVP-7, critério 6)

**Data:** 20 de outubro de 2026 · repositório `github.com/florianoneto07/akma-free5gc`, commit `1527b46` · free5GC na VM 192.168.16.132, recém-revertida ao *snapshot* limpo; AF na VM 192.168.16.129

O teste segue só o `README.md` e o `docs/install.md` do repositório público, como faria uma pessoa que chega pelo fórum do free5GC. Nenhum arquivo do repositório do projeto foi usado nas duas máquinas.

## 1. Montagem

| Máquina | Estado inicial | O que roda |
|---|---|---|
| 192.168.16.132 | *snapshot* limpo: Ubuntu 18.04.6, kernel 5.4, free5GC v4.2.3 (NRF, UDM e AUSF v1.4.5, NEF v1.2.5) e UERANSIM `85a0fbfd` instalados pelo guia oficial, sem nenhum *patch*; Go 1.26.2 em `/usr/local/go/bin`; assinante `imsi-208930000000001` sem assinatura AKMA | free5GC com o AAnF; gNB e UE no *namespace* `ueransim-ue`; clone do GitHub em `~/akma-free5gc` |
| 192.168.16.129 | — | clone do GitHub em `~/akma-repro/akma-free5gc`, separado do repositório do projeto; AF interno, AF externo, `e2e.sh` |

O *snapshot* subiu por DHCP com o endereço 192.168.16.132, e não com o .131 de antes. O endereço não importa para o teste: a máquina do free5GC passou a ser a .132, e os endereços foram trocados só nos lugares que a tabela da seção 1.2 do guia lista:

| Onde | Mudança |
|---|---|
| `free5gc/config/aanfcfg.yaml` (máquina do AF) | `sbi.registerIPv4`, `sbi.bindingIPv4` e `authorizedNefs` → 192.168.16.132 |
| `free5gc/patches/config/0002-nef-akma.patch` (clone da .132, antes de aplicar) | `registerIPv4` e `bindingIPv4` do NEF → 192.168.16.132 |
| `af/config/afcfg.yaml` | `aanfUri` → `http://192.168.16.132:8010` |
| `af/config/afcfg-external.yaml` | `nefUri` → `http://192.168.16.132:8000` |
| `seed`, `deploy.sh`, `e2e.sh` | `--aanf http://192.168.16.132:8010`; argumento `woca@192.168.16.132` |

A máquina do AF continuou em 192.168.16.129, o endereço do exemplo, e `authorizedAfs` não mudou.

## 2. Passos do guia e resultado

| Parte | O que foi feito | Resultado |
|---|---|---|
| 2 | oráculo, `shared/akma`, `shared/crosscheck` | K_AKMA do exemplo trabalhado igual à do guia; testes passam |
| 3.1 | `build-aanf.sh`, testes do AAnF | compila; testes passam |
| 3.3 | `patch-nf.sh nrf apply`; `./run.sh` | os três *patches* aplicam; o NRF recompila; as NFs registram com `OAuth2 setting receive from NRF: true` |
| 3.4 | `deploy.sh woca@192.168.16.132`; `run-aanf.sh start` | AAnF em `192.168.16.132:8010`, OAuth2 ligado |
| 3.5 | perfil no MongoDB; `seed discover` | `nfType: 'AANF'`, `naanf-akma` `v1`, `aanfInfoList` com `0000`; RID `0000` encontra o AAnF, RID `0123` dá `null` |
| 3.6 | `seed register`, `retrieve` (da máquina do AF), `remove` | A-KID do vetor `testbed-ue`; K_AF igual à do vetor (comparada sem imprimir); depois do `remove`, `403 (K_AKMA_NOT_PRESENT)` |
| 4.1, 4.2 | as duas cópias do `openapi`; *patches* do UDM, do AUSF, da configuração e da bancada | aplicam; UDM e AUSF recompilam; `nefcfg.yaml` com 192.168.16.132 |
| 4.3 | `provision-akma-subscriber.sh imsi-208930000000001 0000` | `akmaAllowed: true`, `routingId: '0000'` |
| 4.4 | `go test -run TestAkma` na bancada | `TestAkma5GAKA`, `TestAkmaEAPAKAPrime` e `TestAkmaWithoutSubscription` passam |
| 5.1, 5.2 | *patch* do `run.sh`; `run-aanf.sh stop`; *patches* do UERANSIM; `linkIp: 10.100.0.1`; `ue-netns.sh up` | aplicam; `nr-ue` e `nr-akma-selftest` recompilam |
| 9.1 | `patch-nf.sh nef apply` | o NEF recompila; com `./run.sh -aanf`, NEF em `192.168.16.132:8000` e `AF_IDs accepted from a NEF: [af1.akma.example.FF00000001]` |
| 6.2, 6.3 | AF interno; `uaclient` com o vetor `testbed-ue` | três sessões com `rspauth`; K_AF errada `401`; A-KID desconhecido `403 AKMA_KEY_NOT_AVAILABLE`; outro `Host` `421` |
| 5.3 | gNB e UE | NG Setup; registro; sessão PDU; as seis linhas `[AKMA]` na ordem do guia |
| 7.1, 7.2 | `n6-nat.sh apply`; `ping` pelo `uesimtun0`; autoteste do UE | `ping` ao AF sem perda; `5 passed, 3 refused (FQDN not ASCII), 0 failed` |
| 7.3, 9.2 | Ua\* do UE com o AF interno (8080) e com o AF externo via NEF (8083), pelo túnel até o NRF | três sessões em cada um, a partir de 10.60.0.1 |
| 7.4 | `deregister normal` | nova autenticação; A-KID novo no UE, no AUSF e no AAnF, com `(previous context replaced)` |
| 8 | `e2e.sh woca@192.168.16.132` | **83 verificações, nenhuma falha** (execução `20261020-173639`: 36 na .132 e 47 na máquina do AF) |

**A-KID igual nas três pontas.** Na primeira autenticação, o mesmo A-KID aparece no UE (`AKMA key material derived`), no AUSF (`K_AKMA and A-KID derived`) e no AAnF (`AKMA context … registered`): `rid0000.atidSCGzJb07IA2hd8oL594+NU2qEymcBTqfwcF/Rzd/jcA=@5gc.mnc093.mcc208.3gppnetwork.org`. Depois do `deregister normal`, os três passam juntos para `rid0000.atidooRhLgsFtg66Dt5vuA5pz6V3Lc2kRMk2ugCXPIzSSG8=@…`.

**Linhas AKMA do Ua\* via NEF** (`free5gc.log` da .132, hora local da VM):

```
10:31:14 NEF  K_AF request from AF_ID af1.akma.example.FF00000001 for A-KID rid0000.atidSCGzJb07…
10:31:14 AANF K_AF of imsi-208930000000001 provided to AF_ID af1.akma.example.FF00000001 for A-KID …, expires …, through the NEF at 192.168.16.132
10:31:14 NEF  K_AF for A-KID rid0000.atidSCGzJb07… relayed to AF_ID af1.akma.example.FF00000001, expires 2026-10-20T18:31:14Z
```

As linhas do NEF não trazem o SUPI.

**`e2e.sh`.** Todas as verificações passaram, inclusive:
- R2.akid-ausf e R2.akid-aanf: o A-KID do UE é igual ao do AUSF e ao do AAnF;
- R2.ua e R2n: Ua\* direto, via NEF e anônimo;
- E00 a E12: a API do NEF;
- R7 e L1: nenhuma cadeia de 64 hexadecimais nos logs do free5GC, do UE, do gNB e dos AF.

Os cinco A-KID da execução:

| Autenticação | A-KID |
|---|---|
| 1 | `rid0000.atidC1G0c3FG7cCDfkuE9bT3FK7YsEdxZY9IphtOlvVPAwM=@…` |
| 2 (nova autenticação) | `rid0000.atid/S08Ehd4JkR6cbrxoitcQBv0KZkSgecb6BvbVGUmZYE=@…` |
| 3 (AAnF parado) | `rid0000.atidAwhknUKsTEMhsP82q7W8XitdEato5CBgjWki6i/JebA=@…` |
| 4 (sem assinatura) | `rid0000.atidAOi44prkbdudiH5G16s6FVgNVTA0MPpI0ufgWtt15KE=@…` |
| 5 (restaurada) | `rid0000.atiddlSUfLZuZJmigl2K0tuIGXJTQUCBiSEoehySi3FfVSM=@…` |

## 3. Onde o guia não bastou

| Ponto | O que aconteceu | Correção |
|---|---|---|
| `aanfUri` do `ausfcfg.yaml` | o *patch* de configuração deixa `aanfUri` comentado com o endereço do exemplo, 192.168.16.130. Nem a tabela da seção 1.2 nem a seção 4.2 diziam para trocá-lo. Não afetou o teste, porque o AUSF descobre o AAnF no NRF | `public/docs/install.md`: nova linha na tabela da seção 1.2; a seção 4.2 diz para usar o endereço da seção 3.2, porta 8010 |

Não são falhas do guia:

- **Endereço por DHCP.** O *snapshot* voltou com outro endereço; a tabela da seção 1.2 bastou para a troca.
- **Go num `ssh` não interativo.** Na .132 o Go está no `PATH` de um terminal interativo (o guia oficial do free5GC o põe no `~/.bashrc`), mas não no de um `ssh máquina comando`. Os comandos do guia são digitados num terminal.
- **Último Ua\* depois da nova autenticação.** A sessão SSH que encadeava os comandos da seção 7.4 fechou antes do último `ua`. O Ua\* com o A-KID novo está coberto pela verificação R3 do `e2e.sh`, que passou.

Não foram feitos à mão, por estarem cobertos pelo `e2e.sh`:

| Item do guia | Verificações |
|---|---|
| 6.3, quarto caso (AF_ID fora da política do AAnF) | N03: a recusa no AAnF, `403` sem causa |
| 9.3, acesso anônimo | E03, N01c, R2n.ua-nef-anon, R2n.ua-anon |
| 9.4, remoção do contexto | R6, C1 |

## 4. Critério 6 do MVP-7

| Exigência | Situação |
|---|---|
| o UE registra | atendida |
| o A-KID é igual no UE, no AUSF e no AAnF | atendida, nas duas autenticações feitas à mão e nas do `e2e.sh` |
| o UE autentica no AF pelo Ua\*, direto e via NEF | atendida |
| o `e2e.sh` termina com código 0 | atendida: 83 verificações, nenhuma falha |

## 5. Estado das máquinas depois do teste

- **.132.** Fica com tudo o que o guia instala:
  - os *patches* aplicados;
  - o AAnF, o `seed` e o assinante com assinatura AKMA;
  - o *namespace* do UE e as regras do NAT do N6, que o `e2e.sh` mantém no fim.

  A parte 11 do guia desfaz tudo; reverter o *snapshot* também.
- **.129.** O clone fica em `~/akma-repro/akma-free5gc`, com as três configurações apontadas para a .132.
