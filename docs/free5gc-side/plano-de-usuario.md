# Plano de usuário: do UE ao AF pelo N6

**MVP:** MVP-4 · **Data:** 18 de outubro de 2026 · **Aplicado e verificado** na VM do free5GC em 18/10/2026

## 1. O problema

Desde o MVP-1, o UE do UERANSIM se registrava e abria a sessão PDU, mas nada passava pelo `uesimtun0`: nem 8.8.8.8 nem a VM do AF respondiam. O MVP-4 precisa desse caminho (D5): o cliente Ua\* do UE vai ao AF pela sessão PDU (`uesimtun0` → gNB → UPF → N6 → 192.168.16.129).

Havia duas causas:

1. **UE no *namespace* do host.** `[FREE5GC]` Com o UERANSIM e o free5GC na mesma máquina, o endereço da sessão (10.60.0.1) ficava local ao host. O UPF entregava os pacotes do UE pela `upfgtp` com esse endereço de origem, e o kernel os descartava. A correção é do lado UE: o UE roda num *namespace* próprio (`docs/af-side/ue-ueransim.md` §5).
2. **Sem NAT no N6.** O `status` antes da correção mostrou o `ip_forward` ligado e a política de FORWARD em `ACCEPT`, o `ufw` inativo e nenhuma regra de MASQUERADE para 10.60.0.0/16. Sem NAT, a resposta da VM do AF ou da Internet vai para 10.60.0.x, que só a VM do free5GC sabe alcançar.

## 2. O que muda no host do free5GC

`[FREE5GC]` O guia de instalação do free5GC pede, para o N6, o encaminhamento IP, o MASQUERADE da sub-rede dos UEs na interface de saída e o ajuste de MSS. O *script* `free5gc/scripts/n6-nat.sh` aplica isso:

| Regra | Para quê |
|---|---|
| `net.ipv4.ip_forward=1` | o host encaminha entre a `upfgtp` e a interface do N6 |
| `-t nat POSTROUTING -s 10.60.0.0/16 -o ens33 -j MASQUERADE` | a resposta volta para o endereço do host, que a desfaz para o UE |
| `-t filter FORWARD -s 10.60.0.0/16 -j ACCEPT` e `-d 10.60.0.0/16 -j ACCEPT` | o encaminhamento passa mesmo com política `DROP` |
| `-t mangle FORWARD … SYN … -j TCPMSS --set-mss 1400` | segmentos TCP cabem no túnel GTP-U |

- **Interface:** a do N6 é a da rota padrão (`N6_IFACE` muda); a sub-rede vem de `upfcfg.yaml` (`UE_SUBNET` muda).
- **Idempotente:** cada regra só é acrescentada uma vez.
- **Não persistente:** nada sobrevive a uma reinicialização.

```bash
sudo ~/akma-free5gc/free5gc/scripts/n6-nat.sh status
```

```bash
sudo ~/akma-free5gc/free5gc/scripts/n6-nat.sh apply
```

Para desfazer: `remove`, que não mexe no `ip_forward`.

`[MVP]` Com o NAT, o AF vê a requisição do UE vindo de 192.168.16.130, e não de 10.60.0.x. O Ua\* não depende do endereço de origem, ao contrário da política do AAnF, que vale entre o AF e o AAnF. A alternativa sem NAT é uma rota para 10.60.0.0/16 via 192.168.16.130 na VM do AF; ela não foi usada.

## 3. Arquivos

| Onde | Arquivo | Mudança |
|---|---|---|
| repositório | `free5gc/scripts/n6-nat.sh` (novo) | `apply`, `remove`, `status` |
| host do free5GC | `iptables` (nat, filter, mangle), `sysctl` | regras da seção 2, não persistentes |

Nenhum arquivo do free5GC é alterado.

## 4. Como verificar

Com o free5GC, o gNB e o UE no *namespace* em execução:

```bash
sudo ip netns exec ueransim-ue ping -I uesimtun0 -c 3 8.8.8.8
```

```bash
sudo ip netns exec ueransim-ue ping -I uesimtun0 -c 3 192.168.16.129
```

Os contadores da `upfgtp` (`ip -s link show upfgtp`) devem subir.

Resultado em 18/10/2026:

- 8.8.8.8: 3/3 respostas, cerca de 22 ms;
- VM do AF: 3/3 respostas, cerca de 2 ms;
- `upfgtp`: de 0 para 18 pacotes recebidos e 13 enviados, depois dos `ping` e do Ua\*.
