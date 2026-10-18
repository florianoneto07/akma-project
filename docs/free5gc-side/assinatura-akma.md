# Assinatura AKMA de um assinante

**MVP:** MVP-2 · **Data:** 17 de outubro de 2026

## 1. O dado

`[3GPP]` TS 33.535 §3.1 define *"AKMA subscription data"* como o dado que indica se o assinante pode usar AKMA; §4.2.5 o coloca no UDM. Stage 3 (TS 29.505 Tab. 5.4.2.2-1): dois atributos do `AuthenticationSubscription`, guardado no UDR.

| Atributo | Tipo | Regra |
|---|---|---|
| `akmaAllowed` | booleano | `true`: o assinante pode usar AKMA; ausente ou `false`: não pode |
| `routingId` | cadeia, `^[0-9]{1,4}$` | obrigatório com `akmaAllowed: true`; o Routing Indicator do SUCI do assinante. Cada dígito conta: `"0000"` não é `"0"` (TS 23.003 §2.2B) |

No free5GC, o documento fica no MongoDB, base `free5gc`, coleção `subscriptionData.authenticationData.authenticationSubscription`, chave `ueId`.

## 2. Provisionar

`[FREE5GC]` A *webconsole* do free5GC v4.2.3 não tem campos para esses atributos. Ao editar o assinante, ela regrava o documento e os remove. Por isso, provisione depois de criar o assinante na *webconsole*, e de novo após cada edição.

Na máquina do free5GC:

```bash
~/akma-free5gc/free5gc/scripts/provision-akma-subscriber.sh imsi-208930000000001 0000
```

O *script* faz um `$set` só desses dois atributos e mostra o documento resultante. O RID tem de ser o mesmo configurado no UE: em UERANSIM, `routingIndicator` de `free5gc-ue.yaml`, que é `'0000'` no testbed. Para retirar a assinatura:

```bash
~/akma-free5gc/free5gc/scripts/provision-akma-subscriber.sh imsi-208930000000001 --revoke
```

## 3. Verificar

```bash
mongosh --quiet free5gc --eval 'printjson(db.getCollection("subscriptionData.authenticationData.authenticationSubscription").findOne({ueId: "imsi-208930000000001"}, {_id: 0, ueId: 1, akmaAllowed: 1, routingId: 1}))'
```

O UDR devolve o documento inteiro ao UDM (`QueryAuthSubsData`), sem filtrar campos. Nenhuma mudança no UDR é necessária.

Resultado em 17/10/2026: `imsi-208930000000001` com `akmaAllowed: true` e `routingId: '0000'`.
