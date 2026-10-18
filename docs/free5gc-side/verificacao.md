# Verificação do MVP-2: UDM, AUSF e AAnF juntos

**Data:** 17 de outubro de 2026 · VM do free5GC (192.168.16.130)

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

## 4. Observações de operação

- **UE em `LIMITED-SERVICE`.** O UE do UERANSIM às vezes fica em `MM-DEREGISTERED/LIMITED-SERVICE` na primeira partida, por uma corrida na seleção de célula ("no SI"). Basta reiniciá-lo; não tem relação com o AKMA.
- **Gatilho de nova autenticação.** O critério 4 foi verificado reiniciando o UE. O `nr-cli` (`UERANSIM/build/nr-cli`) consulta o UE sem `sudo`. O `deregister normal` dele, como gatilho alternativo, não foi verificado: na tentativa, o UE estava sendo reiniciado (`No node found`).
- **Encerramento abrupto do `run.sh`.** Se o `run.sh` cair sem encerrar as NFs (sessão SSH perdida, por exemplo), ficam perfis órfãos no `NfProfile` e a interface `upfgtp`. O `./force_kill.sh` do free5GC limpa os dois antes de um novo `./run.sh`.
