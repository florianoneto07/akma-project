# Teste do AF via NEF, do acesso anônimo e do UERANSIM em outra máquina (MVP-6)

**Data:** 19 de outubro de 2026 · free5GC na VM 192.168.16.130; AF na VM 192.168.16.129; gNB e UE na VM 192.168.16.131

## 1. Montagem

| Máquina | O que roda |
|---|---|
| 192.168.16.130 | free5GC v4.2.3 com `run.sh -aanf`: NRF, UDM e AUSF corrigidos; NEF com o *patch* `nef/0001-akma-api.patch` em 192.168.16.130:8000; AAnF em 192.168.16.130:8010; N2 e N3 em 192.168.16.130 |
| 192.168.16.129 | túnel SSH até o NRF; três AF com o mesmo AF_ID, `af1.akma.example.FF00000001` |
| 192.168.16.131 | gNB (`ngapIp` e `gtpIp` 192.168.16.131, AMF 192.168.16.130) e UE no *namespace* `ueransim-ue` |

Os três AF:

| Porta | Modo | Configuração |
|---|---|---|
| 8080 | dentro da rede do operador: K_AF do AAnF | `af/config/afcfg.yaml` |
| 8083 | fora da rede: K_AF do NEF | `af/config/afcfg-external.yaml` |
| 8084 | fora da rede, anônimo | a mesma, com `anonymous: true` e `listen` na 8084 |

## 2. Resultado

**gNB e UE na máquina 192.168.16.131:**

```
[sctp] [info] SCTP connection established (192.168.16.130:38412)
[ngap] [info] NG Setup procedure is successful
[nas] [info] AKMA key material derived, A-KID[rid0000.atido5Mg9o72+sT+ikGfCjuti0v92WQSqhNUwn/6nblbW90=@5gc.mnc093.mcc208.3gppnetwork.org]
[nas] [info] Initial Registration is successful
[nas] [info] PDU Session establishment is successful PSI[1]
```

**Ua\* do UE com os três AF**, dois pedidos cada, pelo plano de usuário: a origem é 10.60.0.1, a sessão PDU do UE.

```
== Ua* to the AF inside (AAnF) at 192.168.16.129:8080
session established (nc 00000001): AF authenticated by rspauth; K_AF valid until 2026-10-20T01:17:28Z
session established (nc 00000002): AF authenticated by rspauth; K_AF valid until 2026-10-20T01:17:28Z
== Ua* to the AF outside (NEF) at 192.168.16.129:8083
session established (nc 00000001): AF authenticated by rspauth; K_AF valid until 2026-10-20T01:17:28Z
session established (nc 00000002): …
== Ua* to the AF outside (NEF), anonymous at 192.168.16.129:8084
session established (nc 00000001): AF authenticated by rspauth; K_AF valid until 2026-10-20T01:17:28Z
session established (nc 00000002): …
```

**Linhas AKMA do free5GC** (`free5gc.log`, hora local da VM):

```
17:17:26 UDM  AKMA indication sent to the AUSF for imsi-208930000000001, Routing Indicator 0000
17:17:26 AUSF K_AKMA and A-KID derived for imsi-208930000000001: A-KID rid0000.atido5Mg9o72+sT+ikGfCjuti0v92WQSqhNUwn/6nblbW90=@…
17:17:26 AUSF AAnF for Routing Indicator 0000: http://192.168.16.130:8010 (discovered at the NRF)
17:17:26 AANF AKMA context of imsi-208930000000001 registered, A-KID rid0000.atido5Mg9o72+sT+ikGfCjuti0v92WQSqhNUwn/6nblbW90=@…
17:17:28 AANF K_AF of imsi-208930000000001 provided to AF_ID af1.akma.example.FF00000001 for A-KID …, expires 2026-10-20T01:17:28Z
17:17:28 NEF  K_AF request from AF_ID af1.akma.example.FF00000001 for A-KID rid0000.atido5Mg9o72+…
17:17:28 AANF K_AF of imsi-208930000000001 provided to AF_ID af1.akma.example.FF00000001 for A-KID …, expires …, through the NEF at 192.168.16.130
17:17:28 NEF  K_AF for A-KID rid0000.atido5Mg9o72+… relayed to AF_ID af1.akma.example.FF00000001, expires 2026-10-20T01:17:28Z
17:17:28 NEF  K_AF request from AF_ID af1.akma.example.FF00000001 for A-KID rid0000.atido5Mg9o72+…, anonymous
17:17:28 AANF K_AF of imsi-208930000000001 provided to AF_ID af1.akma.example.FF00000001 for A-KID …, through the NEF at 192.168.16.130
17:17:28 NEF  K_AF for A-KID rid0000.atido5Mg9o72+… relayed to AF_ID af1.akma.example.FF00000001, expires …
```

O que as linhas mostram:

- **A-KID igual** no UE da máquina 192.168.16.131, no AUSF e no AAnF.
- **NEF sem SUPI:** as linhas do NEF não trazem o SUPI (§6.3 passo 5).
- **Origem marcada no AAnF:** as entregas que passaram pelo NEF aparecem como `through the NEF at 192.168.16.130`.
- **Mesma K_AF:** as três sessões têm a mesma validade. O AAnF entrega ao AF_ID a mesma K_AF, direto ou via NEF, e o UE, que a deriva sozinho, autenticou os três AF.

**AF externos.** Os dois se registraram no NRF pelo túnel e se descadastraram ao encerrar (`Deregistered from the NRF`).

**Nenhuma chave nos logs.** Não há cadeia de 64 hexadecimais no `free5gc.log`, nos logs do gNB e do UE na máquina 192.168.16.131, nem nos logs dos três AF.

**NEF sem *token*.** O pedido sem `Authorization` recebe `401` com `{"error":"verify OAuth Authorization header invalid"}` em `application/json` e sem `WWW-Authenticate`. É o *middleware* do free5GC, mantido no MVP-6 (F12 de `docs/free5gc-side/correcoes-free5gc.md`).

## 3. Critérios do MVP-6 cobertos

| Critério | Situação |
|---|---|
| 1. Mesma K_AF via NEF e direto, sem `supi` na resposta do NEF | atendido aqui pelo Ua\* com os dois AF; a comparação direta das duas respostas está no `e2e.sh` (E01, E02) |
| 2. Ua\* do UE com o AF externo, pelo plano de usuário | atendido |
| 4. Acesso anônimo via NEF | atendido pelo Ua\*; ausência de `supi` e `gpsi` nas respostas verificada no `e2e.sh` (E03, N01c) |
| 6. UERANSIM em outra máquina | atendido |

Os critérios 3 e 5 (negativos e roteiro completo) estão em `docs/af-side/e2e.md`.
