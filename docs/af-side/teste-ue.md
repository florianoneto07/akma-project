# Teste do lado UE (MVP-4)

**Data:** 18 de outubro de 2026 · UERANSIM e free5GC na VM do free5GC (192.168.16.130); AF na VM do AF (192.168.16.129)

## 1. Autoteste do ME

```bash
~/UERANSIM/build/nr-akma-selftest ~/akma-free5gc/shared/akma/testdata/vectors.json
```

Resultado em 18/10/2026:

```
PASS baseline-2.5
PASS testbed-ue
PASS rid-012
PASS mnc-3-digits
PASS ua-tls13
REFUSED nfkc-composed: K_AF not derived, FQDN not ASCII (DV-4)
REFUSED nfkc-decomposed: K_AF not derived, FQDN not ASCII (DV-4)
REFUSED nfkc-fullwidth: K_AF not derived, FQDN not ASCII (DV-4)
5 passed, 3 refused (DV-4), 0 failed
```

## 2. Ponta a ponta

Preparação, uma vez por inicialização da VM:

1. `patch-ueransim.sh apply` e `linkIp: 10.100.0.1` no gNB (`ue-ueransim.md` §9);
2. `sudo n6-nat.sh apply` (`docs/free5gc-side/plano-de-usuario.md`);
3. `sudo ue-netns.sh up`.

A sequência, uma aba por processo:

| # | Onde | Comando |
|---|---|---|
| 1 | VM do AF | `af/build/bin/af -c af/config/afcfg.yaml` |
| 2 | VM do free5GC | `cd ~/free5gc && ./run.sh -aanf` |
| 3 | VM do free5GC | `cd ~/UERANSIM && build/nr-gnb -c config/free5gc-gnb.yaml` |
| 4 | VM do free5GC | `sudo ~/akma-free5gc/af/scripts/ue-netns.sh ue` |
| 5 | VM do free5GC | `sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -n 3` |

O mesmo A-KID nas quatro telas:

| Tela | Linha |
|---|---|
| UE | `[nas] [info] AKMA key material derived, A-KID[rid0000.atidP4PEZ0q3jE3t7QxMHmxp5npr52GFKkmmXSQFIuKebwk=@5gc.mnc093.mcc208.3gppnetwork.org]` |
| free5GC, AUSF | `[INFO][AUSF][AKMA] K_AKMA and A-KID derived for imsi-208930000000001: A-KID rid0000.atidP4PEZ0q3…` |
| free5GC, AAnF | `[INFO][AANF][AKMA] AKMA context of imsi-208930000000001 registered, A-KID rid0000.atidP4PEZ0q3…` |
| UE, pedido do cliente | `[app] [info] AKMA K_AF derived for AF_ID[af1.akma.example.FF00000001], A-KID[rid0000.atidP4PEZ0q3…]` |
| free5GC, AAnF | `[INFO][AANF][AKMA] K_AF of imsi-208930000000001 provided to AF_ID af1.akma.example.FF00000001 for A-KID rid0000.atidP4PEZ0q3…` |
| AF | `[INFO][AF][UA] AKMA initiation: challenge sent to 192.168.16.130:59449`, depois `UE with A-KID rid0000.atidP4PEZ0q3… authenticated` (`nc` 1 a 3) |

Saída do cliente:

```
ME: K_AF derived for AF_ID af1.akma.example.FF00000001
AF af1.akma.example at 192.168.16.129:8080 from 10.60.0.1, A-KID rid0000.atidP4PEZ0q3jE3t7QxMHmxp5npr52GFKkmmXSQFIuKebwk=@5gc.mnc093.mcc208.3gppnetwork.org
session established (nc 00000001): AF authenticated by rspauth; K_AF valid until 2026-10-18T17:53:59Z
session established (nc 00000002): …
session established (nc 00000003): …
```

## 3. Nova autenticação

O comando foi `sudo ue-netns.sh cli imsi-208930000000001 -e "deregister normal"`:

- **UE:** depois do T3511, o UE fez um registro inicial com autenticação nova, ganhou o A-KID `rid0000.atid13yptlVzwk5lcD4yayAz8YuMfHJgKvV/+QUR/KRr4pQ=@…` e uma sessão PDU nova (10.60.0.2).
- **AAnF:** registrou `(previous context replaced)`.
- **Cliente Ua\*:** autenticou com o A-KID novo, a partir de 10.60.0.2.

Da VM do AF, no papel do AF, pedi a K_AF dos dois A-KID com o `seed retrieve`:

- **A-KID antigo:** `403 K_AKMA_NOT_PRESENT`;
- **A-KID novo:** 200.

## 4. Interface do `nr-cli`

| Comando | Resultado |
|---|---|
| `info` | SUPI, HPLMN, IMEI e parâmetros de acesso, sem chave |
| `akma-kaf af_1.example` | `ERROR: AF FQDN must be an ASCII DNS name (no NFKC in the UE)` |
| `akma-kaf af1.akma.example` | `a-kid`, `af-id: af1.akma.example.FF00000001` e `k-af`; nenhuma K_AKMA |

## 5. Critérios de aceitação do MVP-4

| Critério | Resultado |
|---|---|
| 1. A-KID do UE igual, caractere a caractere, ao do AAnF | atendido, nas duas autenticações (seções 2 e 3) |
| 2. Ua\* bem-sucedido com a K_AF derivada pelo UE | atendido: três sessões com `rspauth` conferido, pelo plano de usuário (seção 2) |
| 3. Derivações em C++ reproduzem os vetores do MVP-0 | atendido: 5 vetores conferidos; nos 3 de NFKC, K_AKMA e A-KID conferidos e K_AF recusada (DV-4) |
| 4. Nova autenticação troca o A-KID, e o antigo deixa de valer no AAnF | atendido (seção 3) |
| 5. K_AKMA não exposta pelo `nr-cli` | atendido (seção 4; nenhum comando a mostra) |

Nenhuma cadeia de 64 hexadecimais no `free5gc.log` da execução. A tela do UE não mostra chaves, e na do cliente a K_AF não aparece.
