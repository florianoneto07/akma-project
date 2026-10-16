# ADR-0001-R2 — AAnF na VM do free5GC (revisão do ADR-0001)

**Estado:** Aceito
**Data:** 15 de outubro de 2026
**Revê:** [ADR-0001](ADR-0001-topologia-de-implantacao.md), seções 3, 5, 6 e 8
**Mantém:** ADR-0001, seções 2 e 4.3

**Verificação normativa:** todas as citações deste ADR — TS 33.535, TS 29.535, TS 33.501, TS 29.522, TS 29.122 e TS 29.571 — foram conferidas no texto das normas do acervo (`specs/`), em 15 e 16/10/2026.

---

## 1. Contexto

O ADR-0001 recomendou executar o AAnF, o AF de teste e o oráculo de derivação em uma VM separada da que executa o free5GC (topologia B), deixando na VM do free5GC apenas os *patches* de AUSF, UDM e `free5gc/openapi`.

O objetivo final do projeto foi entretanto precisado: produzir um **tutorial** que descreva o que um usuário precisa fazer **depois de instalar o free5GC** para usar AKMA. Sob esse objetivo, a pergunta deixa de ser "onde o AAnF pode rodar" — o ADR-0001 mostrou que pode rodar em qualquer lugar — e passa a ser "do lado de quem ele está".

---

## 2. Decisão

**O AAnF é executado na VM do free5GC** (`192.168.16.130`), como mais uma NF da instalação do free5GC.

**O AF de teste, o protocolo Ua\*, o código do lado UE, o oráculo de derivação, o repositório e a documentação ficam na VM do AF** (`192.168.16.129`).

| Componente | Local |
|---|---|
| NRF, AMF, AUSF, UDM, UDR, SMF, UPF e demais NFs do free5GC | VM free5GC |
| **AAnF** | **VM free5GC** |
| *Patches* de AUSF, UDM e `free5gc/openapi` | VM free5GC |
| AF de teste e Ua\* | VM do AF |
| Oráculo Python | VM do AF |
| Repositório e documentação | VM do AF |
| UERANSIM (gNB e UE) | instalado hoje na VM do free5GC — ver seção 6 |

---

## 3. Motivo

1. **O AAnF é parte do que o leitor do tutorial instala no free5GC dele.** `[3GPP]` TS 33.535 §4.2.1 define o AAnF como *"the anchor function in the HPLMN"*; a NOTA de §4.1 admite implantação autônoma ou colocalizada com o AUSF ou com o NEF. Um AAnF autônomo executado na mesma máquina das demais NFs continua sendo implantação autônoma. `[MVP]` No tutorial, a instalação do AAnF é um passo do lado free5GC, ao lado dos *patches* de AUSF e UDM: mesmo *toolchain* Go, mesmo grafo de módulos (`free5gc/openapi`, `free5gc/util`), mesma convenção de configuração, mesmo NRF e mesmo certificado OAuth2. Separá-lo em outra máquina faria o tutorial ensinar uma topologia que o leitor não tem.
2. **O AF é, por definição, externo ao núcleo.** O que o tutorial ensina do lado do AF — o AF, o Ua\*, o lado UE — é o que um desenvolvedor de aplicação faria fora do 5GC. A fronteira física entre as VMs passa a coincidir com a fronteira de documentação (`docs/free5gc-side/` × `docs/af-side/`).
3. **Contribuição *upstream*.** `[FREE5GC]` As NFs do free5GC são módulos Go em `NFs/`. Um AAnF destinado a *upstream* precisa compilar e rodar no mesmo ambiente das demais NFs, que é o da VM do free5GC (Go 1.26.2, `openapi` v1.2.4, `util` v1.3.2 — ver `docs/baseline/ambiente-verificado.md`).

O que o ADR-0001 chamou de vantagem da topologia B — isolamento de dependências entre o código do projeto e o do free5GC — deixa de ser vantagem: para o AAnF, compartilhar o grafo de dependências do free5GC é exatamente o requisito.

---

## 4. O que permanece válido do ADR-0001

- **§2 — correção de premissa.** A restrição de kernel é do UPF (`gtp5g`), não do free5GC, e não alcança o AAnF. `[FREE5GC]` Verificado: kernel `5.4.0-150-generic` com `gtp5g` v0.9.14 carregado; Go `1.26.2`.
- **§4.3 — mecanismos de endereçamento.** `[3GPP]` TS 33.535 §6.7: seleção do AAnF via NRF, ou por configuração local. Ambos continuam admitidos.
- **§6.6 — persistência.** O contexto AKMA continua em memória no MVP, sem uso do MongoDB do free5GC.

---

## 5. Reavaliação dos requisitos derivados

A mudança divide os pares de comunicação em dois grupos: os que passam a ser locais à VM do free5GC e os que continuam atravessando a rede entre as VMs.

### 5.1 Requisitos de rede (revê ADR-0001 §6.1)

| Origem | Destino | Finalidade | Antes (ADR-0001) | Agora |
|---|---|---|---|---|
| AAnF | NRF | registro, *heartbeat*, token OAuth2 | entre VMs | **local** |
| AUSF | AAnF | `register-anchorkey` | entre VMs | **local** |
| AAnF | UDM | `id-translation-result`, `ee-subscriptions`, `trigger-auth` (fases posteriores) | entre VMs | **local** |
| AF | AAnF | `retrieve-applicationkey` | local | **entre VMs** |
| AAnF | AF | `POST {notifUri}` (fase posterior) | local | **entre VMs** |
| UE | AF | Ua\* | — | depende da seção 6 |

### 5.2 `registerIPv4` × `bindingIPv4` (revê ADR-0001 §6.2)

**Para o par AUSF↔AAnF, a armadilha deixa de valer.** `[FREE5GC]` As NFs do free5GC escutam em endereços de *loopback* distintos (`127.0.0.9` o AUSF, `127.0.0.10` o NRF, `127.0.0.3` o UDM). Um AAnF que seguisse a convenção e registrasse um `127.0.0.x` seria alcançável pelo AUSF, na mesma máquina.

**A armadilha passa para o par AAnF↔AF.** Se o AAnF seguir essa convenção, o registro no NRF é bem-sucedido, o AUSF o alcança, o fluxo de autenticação primária e `register-anchorkey` funcionam — e a falha só aparece quando o AF, na outra VM, tenta `retrieve-applicationkey`. É a mesma assinatura de falha descrita no ADR-0001, deslocada para o último passo do fluxo, o que a torna mais enganosa.

Decisão de configuração `[MVP]`:

- `bindingIPv4` do AAnF: `192.168.16.130` (ou `0.0.0.0`);
- `registerIPv4` do AAnF: `192.168.16.130`.

Assim, o URI publicado no NRF (usado pelo AUSF) e o URI configurado no AF são o mesmo, alcançável das duas máquinas. A porta escolhida não pode colidir com outra NF que escute em `0.0.0.0`.

`[FREE5GC]` O NRF escuta apenas em `127.0.0.10` e não é alcançável a partir da VM do AF. O AF seleciona o AAnF por **configuração local**, admitida por TS 33.535 §6.7. Expor o NRF à outra VM seria alteração de configuração do free5GC sem necessidade para o MVP.

Verificação obrigatória ao subir o AAnF (MVP-1): requisição a partir de `192.168.16.129` à porta do AAnF, e ausência de regra de *firewall* que a bloqueie.

### 5.3 Sincronização de relógio (revê ADR-0001 §6.3)

**Deixa de valer para AUSF↔AAnF**, que compartilham o relógio da mesma máquina. O AUSF não transporta instante absoluto ao AAnF no fluxo do MVP.

**Continua requisito para AAnF↔AF.** `[3GPP]` `expiry` de `AkmaAfKeyData` é `DateTime` — instante absoluto, emitido pelo AAnF e avaliado pelo AF. Um desvio de relógio entre as duas VMs faz o AF aceitar K_AF expirada ou rejeitar K_AF válida.

Estado medido em 15/10/2026: as duas VMs relatam relógio sincronizado; desvio aproximado de +0,16 s (medida por SSH, grosseira). As VMs estão em fusos horários diferentes; `[MVP]` o AAnF serializa `expiry` em UTC com sufixo `Z`, e o AF compara instantes, nunca horas locais.

O teste negativo com desvio deliberado de relógio permanece no plano, agora aplicado apenas ao AF.

### 5.4 Desempenho (revê ADR-0001 §6.4)

- A latência AUSF→AAnF passa a ser medida em condição local. O custo de transporte desse par deixa de ser variável do experimento.
- A latência AF→AAnF continua dependente da rede entre as VMs; os requisitos metodológicos do ADR-0001 §6.4 (declarar hipervisor, rede virtual, MTU e latência de base) passam a valer só para esse par. Latência de base medida: cerca de 0,5 ms (ICMP).
- **Requisito novo:** o AAnF disputa CPU e memória com as NFs do free5GC e com o UERANSIM, se este permanecer na mesma VM. A configuração da VM (vCPUs, memória) deve ser declarada junto com qualquer medição.

### 5.5 Segurança de transporte (revê ADR-0001 §6.5)

`[3GPP]` TS 33.535 §4.4.0: *"The SBA interface between the AAnF and the AUSF shall be confidentiality, integrity and replay protected"*; o mesmo para AAnF↔AF/NEF e AAnF↔UDM. TS 33.535 §4.2.2 e TS 29.535 §4.2.1.3.2: o AF *"shall be authenticated and authorized by the operator network"* antes de receber K_AF.

`[3GPP]` TS 33.501 V19.7.0 §13.1.0: todas as NFs devem suportar TLS mutuamente autenticado; *"TLS shall be used for transport protection within a PLMN unless network security is provided by other means"*. NOTA 2: *"If interfaces are trusted (e.g. physically protected), it is for the PLMN-operator to decide whether to use cryptographic protection."*

- **AUSF↔AAnF:** o tráfego deixa de atravessar a rede entre VMs e fica em *loopback*, dentro de um único hospedeiro. `[INTERPRETAÇÃO]` É o caso mais próximo de "interface confiável" da NOTA 2 (que é informativa): o uso de `http` (h2c) passa a ser decisão de operador, a registrar como tal — não deixa de ser decisão.
- **AF↔AAnF:** passa a atravessar a rede virtual entre VMs. O uso de h2c nesse par só se apoia na NOTA 2 se o projeto declarar essa rede como confiável — e, mesmo assim, deixa o AF sem autenticação de transporte, o que §4.2.2 exige. Sem ela, um AF pode pedir K_AF usando o AF_ID de outro: é o ataque de falsificação descrito por Akman et al. (Computers 12(1):2, 2023; ver `docs/literatura/`). Precisa ser decidido explicitamente antes do MVP-1 e registrado em `docs/conformidade/desvios.md` (desvio DV-1 do plano).
- O suporte a TLS continua obrigatório para a NF (§13.1.0, primeira frase), ainda que o MVP opere sem ele.

A observação do ADR-0001 sobre certificados — SAN correspondente ao endereço publicado no NRF — passa a valer para o certificado do AAnF apresentado ao AF.

### 5.6 Grafo de dependências (revê ADR-0001 §4.2)

O AAnF passa a compartilhar o grafo de módulos do free5GC. Se os modelos do Naanf_AKMA forem acrescentados ao `free5gc/openapi` por `replace`, a mesma diretiva vale para o AUSF (consumidor) e para o AAnF (produtor). `[FREE5GC]` `openapi` v1.2.4 já contém `AANF`, `naanf-akma`, `AanfInfo`, `AkmaInd` e `RoutingId`; faltam apenas os modelos e o cliente do Naanf_AKMA (`docs/baseline/ambiente-verificado.md` §3).

---

## 6. Ponto em aberto: localização do lado UE

`[FREE5GC]` O UERANSIM está instalado e funcionando **na VM do free5GC**, com o UE em *network namespace* próprio (`ueransim-ue`). A K_AUSF do UE existe apenas dentro do processo `nr-ue`. `[3GPP]` TS 33.535 §6.1 atribui ao ME a derivação de K_AKMA e do A-KID a partir de K_AUSF.

Isso conflita com a alocação "lado UE na VM do AF": o código do lado UE que deriva as chaves precisa rodar onde o `nr-ue` roda. As alternativas e a recomendação estão no plano de desenvolvimento (decisão D5). Esta revisão não decide o ponto; registra apenas que a topologia do lado UE depende dele.

---

## 7. Verificações do ADR-0001 §7 — resultado

| # | Verificação | Resultado |
|---|---|---|
| 1 | Kernel da VM do free5GC | `5.4.0-150-generic` |
| 2 | Go da VM do free5GC | `go1.26.2 linux/amd64` |
| 3 | *Commit* do free5GC e submódulos | v4.2.3, `3b34a08e`; submódulos em `docs/baseline/ambiente-verificado.md` §2.1 |
| 4 | `AANF` e `naanf-akma` no `free5gc/openapi` | presentes (v1.2.4) |
| 5 | `AkmaInd` e `RoutingId` em `AuthenticationInfoResult` | presentes (v1.2.4) |
| 6 | Alcançabilidade entre as VMs | ICMP ok. O critério muda: deixa de ser "AAnF alcança o NRF na outra VM" e passa a ser "AF alcança a porta do AAnF" — a verificar no MVP-1 |
| 7 | Sincronização de relógio | ambas sincronizadas; desvio aproximado de 0,16 s, abaixo de 1 s |

---

## 8. Consequências

- A documentação do AAnF passa para `docs/free5gc-side/`; a do AF, do Ua\* e do lado UE, para `docs/af-side/`.
- A recomendação do ADR-0001 §5 de usar contêineres (`free5gc-compose`) como alvo de reprodutibilidade fica adiada. Se retomada, o AAnF entra como mais um serviço do `free5gc-compose`, coerente com esta revisão.
- A matriz de rastreabilidade deve trocar a linha "Endereço publicado alcançável entre VMs" (origem ADR-0001 §6.2) por "Endereço do AAnF alcançável a partir da VM do AF" (origem: esta revisão, §5.2), e restringir o requisito de NTP ao par AAnF↔AF.

---

## 9. Referências

- TS 33.535 V19.0.0, cláusulas 4.1, 4.2.1, 4.2.2, 4.4.0, 6.1, 6.7 — conferidas
- TS 29.535 V19.5.0, cláusula 4.2.1.3.2 — conferida
- TS 29.522 V19.7.0, Tabela 5.14.5.3.3-1 — `AkmaAfKeyData.expiry` (`DateTime`) — conferida
- TS 29.122 V19.7.0 e TS 29.571 V19.7.0 — `DateTime` com `format: date-time` em ambas — conferido
- TS 33.501 V19.7.0, cláusula 13.1.0 e NOTA 2 — conferida
- [ADR-0001](ADR-0001-topologia-de-implantacao.md)
- `docs/baseline/ambiente-verificado.md`
