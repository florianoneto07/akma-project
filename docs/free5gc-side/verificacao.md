# Verificação do MVP-2 e do MVP-3a: UDM, AUSF e AAnF juntos

**Data:** 17 de outubro de 2026 (seções 1 a 3, MVP-2) e 18 de outubro de 2026 (seção 4, MVP-3a) · VM do free5GC (192.168.16.130)

## 1. Banco de testes do free5GC (critérios 1, 3 e 7)

`[MVP]` O *patch* `free5gc/patches/test/0001-akma-tests.patch` muda o banco de testes em `~/free5gc/test`:

- o AAnF do projeto passa a rodar no mesmo processo das NFs, em `127.0.0.40:8000`, com política que autoriza `af1.akma.example.FF00000001` a partir de `127.0.0.1`;
- o AUSF de teste fica com `akma.enable: true`;
- o UE de teste guarda a própria K_AUSF (`RanUeContext.Kausf`), nos dois métodos;
- entram três testes novos em `akma_test.go`.

O `go.mod` do teste aponta, em `~/akma-free5gc`, para o AAnF, o `shared/akma` e a cópia do `openapi` v1.2.5 corrigida (`free5gc/build/openapi-test`).

| Teste | O que faz | Critério |
|---|---|---|
| `TestAkma5GAKA` | registra `imsi-208930000007490` (5G AKA, SUCI e `routingId` `"0000"`). Deriva, da K_AUSF do UE, o A-KID e a K_AF de `af1.akma.example.FF00000001`, pede a K_AF ao AAnF e compara | 1 (5G AKA) |
| `TestAkmaEAPAKAPrime` | o mesmo com EAP-AKA', `imsi-208930000007491` e RID `"012"` (zero à esquerda de ponta a ponta) | 1 (EAP-AKA') |
| `TestAkmaWithoutSubscription` | assinante sem `akmaAllowed`: o A-KID que o UE derivaria é desconhecido do AAnF (403) | 3 |

Os testes AKMA param no registro do UE e não sobem o SMF. Rodam sem `sudo`, com o free5GC parado, porque usam os mesmos endereços:

```bash
cd ~/free5gc/test && go test -v -count=1 -run TestAkma .
```

Resultado em 17/10/2026: os três passam. O log do teste não contém nenhuma cadeia de 64 caracteres hexadecimais, onde K_AKMA ou K_AF apareceriam (critério 7).

## 2. Regressão do free5GC (critério 5)

Com UDM, AUSF, NRF e banco de testes corrigidos, e com o AAnF iniciado pelo banco de testes, rodamos na máquina do free5GC (o `test.sh` pede `sudo`):

```bash
cd ~/free5gc && ./test.sh TestRegistration
```

O mesmo para `TestGUTIRegistration`, `TestReSynchronization` e `TestEAPAKAPrimeAuthentication`.

Resultado em 17/10/2026: os quatro passam. Nos quatro, o AAnF subiu no processo de teste, e o AUSF não fez nenhuma ação AKMA, porque esses assinantes não têm `akmaAllowed`.

## 3. UERANSIM (critérios 2, 4 e 6)

Assinante `imsi-208930000000001` provisionado com `routingId` `"0000"` (`assinatura-akma.md`). A sequência:

1. iniciar o free5GC com `./run.sh`;
2. iniciar o AAnF com `run-aanf.sh start`, com `logger.level: debug` no `aanfcfg.yaml` para que o A-KID apareça no log;
3. iniciar o gNB e o UE do UERANSIM.

| Critério | Verificação | Resultado em 17/10/2026 |
|---|---|---|
| 2 | o AAnF guarda A-KID `rid0000.atid<44 base64>@5gc.mnc093.mcc208.3gppnetwork.org` | `rid0000.atidfTiRYw4XQKRmE8qEUX9JDLIbElNilKMynonfqfxgnjI=@5gc.mnc093.mcc208.3gppnetwork.org`. O AUSF descobriu o AAnF no NRF com `routing-indicator=0000` (200) e registrou (200). O UDM não registrou divergência de RID |
| 4 | reiniciar o UE (autenticação nova): o AAnF substitui o contexto | `AKMA context registered (replaced previous: true)`, com A-KID novo (`rid0000.atid0xvH7kR99…`) no lugar do anterior (`rid0000.atidfTiRY…`) |
| 6 | AAnF parado: o UE se registra normalmente e o AUSF registra a falha no log (D9) | 5G AKA confirmado às 12:17:50.805. O AUSF registrou `Naanf_AKMA_AnchorKey_Register for imsi-208930000000001 failed: AAnF discovery: no uri for naanf-akma found`. `Registration Accept` às 12:17:51.052, sessão PDU estabelecida |
| 7 | nenhuma chave nos logs do free5GC e do AAnF | nenhuma cadeia de 64 hexadecimais nos logs do free5GC das duas execuções e no log do AAnF, mesmo em `debug` |

Depois da verificação, o `aanfcfg.yaml` voltou a `logger.level: info`, igual ao do repositório.

## 4. Mensagens AKMA na tela (MVP-3a)

Pré-requisitos:

- *patches* do UDM e do AUSF regenerados e o alvo `run` aplicado (`logs-akma.md` §7);
- AAnF recompilado;
- UERANSIM com `af/scripts/patch-ueransim.sh apply` (`docs/af-side/ue-ueransim.md`);
- assinante `imsi-208930000000001` com `routingId` `"0000"`.

A sequência:

1. iniciar o free5GC com `./run.sh -aanf`, sem o AAnF do `run-aanf.sh`;
2. iniciar o gNB e o UE do UERANSIM, como na seção 3;
3. pedir a K_AF da VM do AF com o `seed retrieve`, para o AF autorizado e para um AF_ID fora da política;
4. parar o UE e rodar o `nr-ue` duas vezes: com uma cópia da configuração sem `akma` e com a configuração normal;
5. parar tudo com Ctrl-C no `run.sh`.

Os testes `TestAkma*` do banco de testes (seção 1) rodaram antes, com o free5GC parado. Passaram, e a saída deles já mostra as mesmas mensagens, nos dois métodos de autenticação.

Na tela do UE, às 08:33:01.451:

```
[nas] [debug] Security Mode Command received
[nas] [debug] Selected integrity[2] ciphering[0]
[nas] [info] AKMA key material derived, A-KID[rid0000.atidXlyySARFJ4fb/dIvhrDFIw9s87pSbkajjuuB41NTr0Q=@5gc.mnc093.mcc208.3gppnetwork.org]
[nas] [debug] Registration accept received
```

Na tela do `run.sh`, filtrando a categoria `AKMA` (horário, NF, mensagem):

```
08:32:20.956 [INFO][AANF][AKMA] Routing Indicators [0000]; K_AF lifetime 1h0m0s; AF_IDs authorized by local policy: [af1.akma.example.FF00000001]
08:33:01.412 [INFO][UDM][AKMA]  AKMA indication sent to the AUSF for imsi-208930000000001, Routing Indicator 0000
08:33:01.414 [INFO][AUSF][AKMA] AKMA indication received from the UDM for imsi-208930000000001, Routing Indicator 0000
08:33:01.448 [INFO][AUSF][AKMA] K_AKMA and A-KID derived for imsi-208930000000001: A-KID rid0000.atidXlyySARFJ4fb/dIvhrDFIw9s87pSbkajjuuB41NTr0Q=@5gc.mnc093.mcc208.3gppnetwork.org
08:33:01.480 [INFO][AUSF][AKMA] AAnF for Routing Indicator 0000: http://192.168.16.130:8000 (discovered at the NRF)
08:33:01.494 [INFO][AANF][AKMA] AKMA context of imsi-208930000000001 registered, A-KID rid0000.atidXlyySARFJ4fb/dIvhrDFIw9s87pSbkajjuuB41NTr0Q=@5gc.mnc093.mcc208.3gppnetwork.org
08:33:01.494 [INFO][AUSF][AKMA] AKMA key material of imsi-208930000000001 registered at the AAnF
08:35:31.055 [INFO][AANF][AKMA] K_AF of imsi-208930000000001 provided to AF_ID af1.akma.example.FF00000001 for A-KID rid0000.atidXlyy…=@…, expires 2026-10-18T16:35:31Z
08:35:31.073 [WARN][AANF][AKMA] AF at 192.168.16.129 not authorized for AF_ID af2.akma.example.FF00000001
```

O UE registra a derivação 3 ms depois do AUSF. O AUSF deriva depois do 5G AKA confirmado (`5g-aka-confirmation`, 200), e o UE deriva ao aceitar o *Security Mode Command* que vem em seguida.

| Critério | Verificação | Resultado em 18/10/2026 |
|---|---|---|
| 1 | AAnF sobe com `./run.sh -aanf`, escreve na mesma tela e no mesmo `free5gc.log`, para no Ctrl-C | sobe depois do CHF e registra no NRF (201). As linhas `[AANF]` vão para a tela e para o `free5gc.log` da execução. No Ctrl-C (08:50:21), para o servidor SBI e se descadastra: o NRF registra `NfDeregister Success: AANF [b6f417f0-…]` (204), e o AAnF termina com `AANF App is terminated`. Nenhuma NF fica em execução, e o `cert/aanf_<id>.pem` é apagado pelo NRF |
| 2 | sequência UDM → AUSF → AAnF em `info`, categoria `AKMA` | a sequência acima, sem nenhuma linha AKMA fora de `info`, salvo o `warning` do AF não autorizado, que é o esperado |
| 3 | A-KID do UE igual ao do AUSF e do AAnF | igual, caractere a caractere, nas duas autenticações com `akma: true`: `rid0000.atidXlyySARFJ4fb…` às 08:33 e `rid0000.atidmwKSp6kqYpxe72jqwHTlTnnU1n1dnq9wuGWXpvIOmgU=@…` às 08:36 |
| 4 | sem `akma`, o `nr-ue` age como o original | às 08:36:23, nenhuma linha AKMA no UE. A rede seguiu a assinatura: AUSF e AAnF registraram o A-KID `rid0000.atidOpP7…` `(previous context replaced)`. Às 08:36:32, com `akma: true`, o UE voltou a derivar, com A-KID igual ao da rede |
| 5 | nenhuma chave nas telas e nos logs | zero cadeias de 64 hexadecimais no `free5gc.log` e na tela do UE; a K_AF entregue, procurada em base64, também não aparece no log |
| 6 | testes continuam passando | UDM (`processor`, `suci`), AUSF (`processor`) e AAnF (todos os pacotes) localmente; `TestAkma5GAKA`, `TestAkmaEAPAKAPrime` e `TestAkmaWithoutSubscription` na VM |

As duas execuções do critério 4 rodaram sem `sudo`, com `timeout`: o registro e a derivação acontecem antes de o UE precisar do TUN. O `TUN interface could not be setup` que aparece em seguida é esperado.

## 5. Ponta a ponta (MVP-5)

O roteiro `af/scripts/e2e.sh` (`docs/af-side/e2e.md`) refaz, de uma vez, as verificações das seções 3 e 4 e as do lado UE, e acrescenta os testes negativos. Do lado free5GC:

- free5GC com `run.sh -aanf`;
- A-KID igual no UE, no AUSF e no AAnF;
- substituição do contexto numa nova autenticação;
- AAnF parado durante a autenticação (D9);
- assinante sem assinatura AKMA;
- `remove-context`;
- os códigos de erro da API do AAnF a partir da máquina do AF;
- nenhuma chave nos logs.

Resultado em 18/10/2026: 63 verificações e nenhuma falha.

## 6. Observações de operação

- **UE em `LIMITED-SERVICE`.** O UE do UERANSIM às vezes fica em `MM-DEREGISTERED/LIMITED-SERVICE` na primeira partida, por uma corrida na seleção de célula ("no SI"). Basta reiniciá-lo; não tem relação com o AKMA.
- **Gatilho de nova autenticação.** O critério 4 foi verificado reiniciando o UE. O `nr-cli` (`UERANSIM/build/nr-cli`) consulta o UE sem `sudo`. O `deregister normal` dele, como gatilho alternativo, não foi verificado: na tentativa, o UE estava sendo reiniciado (`No node found`).
- **Encerramento abrupto do `run.sh`.** Se o `run.sh` cair sem encerrar as NFs (sessão SSH perdida, por exemplo), ficam perfis órfãos no `NfProfile` e a interface `upfgtp`. O `./force_kill.sh` do free5GC limpa os dois antes de um novo `./run.sh`.
