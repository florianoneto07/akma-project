# ADR-0001 — Topologia de implantação do AAnF em relação ao free5GC

**Estado:** Revisto pelo [ADR-0001-R2](ADR-0001-R2-aanf-na-vm-do-free5gc.md) em 16/10/2026 — seções 3, 5, 6 e 8 substituídas; seções 2 e 4.3 mantidas
**Data:** 02 de setembro de 2026
**Contexto:** decisão anterior ao início da fase F3 (AAnF autônomo)

*Nota (16/10/2026):* as fases citadas neste ADR correspondem ao plano de desenvolvimento aprovado (`docs/plano/plano-de-desenvolvimento.md`): F2 = MVP-0, F3 = MVP-1, F4 = MVP-2, F5 = MVP-3, F6 = MVP-4, F7 = MVP-5. F11 (validação de segurança, TLS), F12 (avaliação de desempenho) e F15 (contribuição *upstream*) não constam do README; estão no MVP-6+.

---

## 1. Questão

O AAnF, o AF de teste e os componentes auxiliares deste projeto devem ser executados na mesma máquina em que roda o free5GC, ou podem ser executados em máquina separada, com o endereço apontado na configuração do free5GC?

A motivação declarada para a pergunta é a preocupação de que a exigência de versão de kernel do free5GC restrinja as bibliotecas e ferramentas utilizáveis na implementação do AAnF.

---

## 2. Correção de premissa

A premissa de que "o free5GC usa um kernel muito antigo" precisa de qualificação. Duas correções, ambas verificáveis:

### 2.1 A restrição de kernel é do UPF, não do free5GC

O guia oficial de instalação do free5GC condiciona a exigência de versão de kernel ao UPF: a versão 5.0.0-23-generic ou 5.4.x é requerida *"para usar o elemento UPF"*, porque o free5GC emprega o módulo de kernel `gtp5g`, testado e compilado contra essas versões. O README do repositório `free5gc/gtp5g` reitera que o módulo não funciona com toda versão de kernel e recomenda 5.0.0-23-generic, acima de 5.4 (Ubuntu 20.04), ou RHEL8.

O `gtp5g` implementa o plano de usuário GTP-U, tratando pacotes segundo PDR e FAR de PFCP (TS 29.281 e TS 29.244). **Nenhuma outra função de rede do free5GC o utiliza.** AMF, AUSF, UDM, UDR, NRF, PCF, NSSF e SMF são processos Go em espaço de usuário que se comunicam por HTTP/2 sobre TCP. Não possuem dependência de kernel além da pilha de rede padrão.

A restrição, além disso, vem sendo relaxada: a versão v0.8.5 do `gtp5g` registra alteração para suportar kernel 6.2 ou superior. Relatos de implantação em Kubernetes indicam operação bem-sucedida em 5.15, com ressalvas para a família 6.x.

### 2.2 O toolchain Go não é antigo

O guia oficial de instalação corrente instrui a instalação do **Go 1.26.2**. Referências a Go 1.12 ou 1.14 provêm de páginas de wiki arquivadas dos estágios 2 e 3 do projeto e não descrevem o estado atual.

### 2.3 Consequência

A preocupação de incompatibilidade de bibliotecas por causa do kernel **não se aplica ao AAnF**. Ela se aplicaria a um componente que dependesse de módulo de kernel, de `netlink` específico ou de instruções de plano de dados — o que não é o caso de uma função de rede SBI.

Existe, contudo, uma restrição real, de outra natureza, tratada na seção 4.

---

## 3. Decisão

**O AAnF, o AF de teste e o oráculo de derivação serão executados em máquina separada da que executa o free5GC.**

**Os patches do AUSF e do UDM serão obrigatoriamente construídos e executados na máquina do free5GC.**

---

## 4. Fundamentação

### 4.1 O AAnF pode ficar em qualquer lugar

O AAnF é uma NF de plano de controle. Sua superfície de integração é inteiramente SBI:

| Interface | Papel | Transporte |
|---|---|---|
| N61 | AUSF → AAnF (`register-anchorkey`) | HTTP/2 |
| N62 | AF interno → AAnF (`retrieve-applicationkey`) | HTTP/2 |
| N63 | NEF → AAnF | HTTP/2 |
| — | AAnF ↔ NRF (registro, descoberta, OAuth2) | HTTP/2 |
| — | AAnF → UDM (`id-translation-result`, `ee-subscriptions`, `trigger-auth`) | HTTP/2 |

Nenhuma delas envolve plano de usuário, módulo de kernel ou interface de túnel. O único requisito de infraestrutura é conectividade IP bidirecional.

### 4.2 Os patches do AUSF e do UDM não podem

Esta é a restrição real, e ela **não é de kernel — é de grafo de dependências**.

O AUSF e o UDM modificados continuam sendo NFs do free5GC. Compilam contra `github.com/free5gc/openapi`, `github.com/free5gc/util` e demais módulos fixados no `go.mod` de cada NF. As alterações previstas — leitura de `akmaInd` e `routingId`, derivação de K_AKMA e A-TID, consumidor `Naanf_AKMA` — vivem dentro desses módulos.

Consequências práticas:

- O desenvolvimento dos patches pode ocorrer em qualquer máquina; a **construção e a execução** ficam atadas ao ambiente do free5GC.
- Alterações em `free5gc/openapi` — necessárias para `AkmaInd`, `RoutingId`, `NFType` `AANF` e `ServiceName` `naanf-akma` — devem ser feitas por `replace` no `go.mod` ou por *fork* do módulo, e valem para todas as NFs que o importam.
- O AAnF, sendo binário separado, **não precisa** compartilhar esse grafo. Pode usar versão de Go distinta e dependências próprias. Recomenda-se, ainda assim, alinhá-lo ao free5GC para facilitar a eventual contribuição *upstream* prevista na fase F15.

### 4.3 O padrão admite ambos os mecanismos de endereçamento

`[3GPP]` TS 33.535 V19.0.0, cláusula 6.7: o AUSF, o AF interno e o NEF selecionam a instância de AAnF utilizando o NRF, *"salvo se a informação do AAnF estiver disponível por outros meios, por exemplo configuração local"*.

Portanto, apontar um endereço IP remoto não é desvio da norma sob nenhum dos dois mecanismos:

**Mecanismo 1 — descoberta via NRF (preferido).** O AAnF registra-se no NRF publicando seu próprio endereço alcançável. O AUSF e o AF descobrem-no por consulta a `target-nf-type=AANF`, com o parâmetro `routing-indicator` (ver Adendo B §B.1.2). É o mecanismo padrão do free5GC e o que será exercitado na fase F3.

**Mecanismo 2 — configuração local.** URI do AAnF fixada no arquivo de configuração do AUSF e do AF. Admitido explicitamente pela cláusula 6.7. Útil como caminho de contorno durante a depuração inicial.

---

## 5. Topologias consideradas

| | **A — Máquina única** | **B — Duas VMs** (adotada) | **C — Containers** |
|---|---|---|---|
| free5GC | *bare metal* ou VM | VM 1 | `free5gc-compose` |
| AAnF, AF, oráculo | mesma máquina | VM 2 | container adicional na mesma rede |
| Isolamento de dependências | nenhum | completo | completo |
| Complexidade de rede | mínima | moderada | baixa |
| Latência AUSF→AAnF | *loopback* | rede entre VMs | rede de containers |
| Reprodutibilidade | baixa | média | alta |
| Efeito da restrição de kernel | UPF apenas | UPF apenas, VM 1 | UPF apenas, kernel do hospedeiro |

**Topologia B é adotada** para as fases F3 a F7, pelas razões: isolamento de dependências entre o código do projeto e o do free5GC; possibilidade de reinstalar ou reverter qualquer dos dois lados sem afetar o outro; e correspondência com a arquitetura de implantação real, em que o AAnF é uma NF autônoma (TS 33.535 §4.1 admite implantação autônoma ou colocada).

A topologia C é recomendada como alvo da fase F7, para reprodutibilidade dos experimentos e para eventual publicação do ambiente junto com o trabalho.

---

## 6. Consequências e requisitos derivados

### 6.1 Requisitos de rede

| Origem | Destino | Finalidade |
|---|---|---|
| AAnF | NRF | Registro do perfil, *heartbeat*, descoberta, token OAuth2 |
| NRF | AAnF | Notificações de gerenciamento de NF, se subscritas |
| AUSF | AAnF | `POST /naanf-akma/v1/register-anchorkey` |
| AF | AAnF | `POST /naanf-akma/v1/retrieve-applicationkey` |
| AAnF | UDM | `id-translation-result`, `ee-subscriptions`, `trigger-auth` (fases posteriores) |
| AAnF | AF | `POST {notifUri}` — notificação de desabilitação (fase posterior) |

### 6.2 Armadilha principal: `registerIPv4` versus `bindingIPv4`

As NFs do free5GC distinguem o endereço em que o servidor escuta do endereço publicado no NRF. Registrar `127.0.0.1` — o valor padrão em várias configurações de exemplo — faz o registro no NRF ser bem-sucedido e a descoberta retornar um URI inalcançável a partir da outra máquina. A falha manifesta-se apenas no momento da primeira requisição do AUSF, com erro de conexão, e é frequentemente diagnosticada como problema do AAnF.

Verificação obrigatória ao subir a topologia B: o valor publicado no perfil do AAnF deve ser o endereço da VM 2 alcançável a partir da VM 1, e vice-versa para o NRF.

### 6.3 Sincronização de relógio — requisito, não recomendação

O atributo `expiry` de `AkmaAfKeyData` é do tipo `DateTime`, isto é, um instante absoluto (Adendos B §B.5.3 e C §C.2). Com AAnF e AF em máquinas distintas, um desvio de relógio produz aceitação de K_AF expirada ou rejeição de K_AF válida.

**NTP é requisito de configuração da topologia B.** A fase F7 deve incluir um teste negativo com desvio deliberado de relógio, conforme já registrado no Adendo B.

### 6.4 Efeito sobre a fase F12 (avaliação de desempenho)

A separação em VMs introduz latência de rede nas medições de AUSF→AAnF e AF→AAnF. Isto **não invalida** as medições, mas as torna dependentes da topologia.

Requisitos metodológicos derivados:

1. A topologia utilizada deve ser declarada com todos os parâmetros relevantes: hipervisor, tipo de rede virtual, MTU, e latência de base medida entre as VMs.
2. Recomenda-se medir a derivação da KDF isoladamente — que é insensível à topologia — em separado da latência ponta a ponta.
3. Recomenda-se repetir as medições de latência em topologia A, para separar o custo de processamento do custo de transporte. A diferença entre as duas topologias é, em si, um resultado publicável.

### 6.5 Segurança de transporte

`[3GPP]` TS 33.535 §4.4.0 exige proteção de confidencialidade, integridade e repetição nas interfaces SBI entre AAnF e AUSF, AAnF e AF/NEF, e AAnF e UDM.

Em topologia A, o tráfego permanece em *loopback*. Em topologia B, atravessa a rede entre VMs. O MVP pode operar em `http` (h2c), conforme o padrão de laboratório do free5GC, mas **isto passa a ser uma simplificação de ambiente com efeito sobre um requisito normativo**, e não mais uma escolha indiferente. Deve ser registrada como tal e revertida na fase F11 (validação de segurança), quando TLS será habilitado.

Ao habilitar TLS, os certificados devem conter SAN correspondente ao endereço ou FQDN efetivamente publicado no NRF — armadilha correlata à da seção 6.2.

### 6.6 Persistência

O contexto AKMA do AAnF será mantido em memória durante o MVP. Não há dependência do MongoDB do free5GC, o que remove uma razão adicional para colocalização. Caso a persistência venha a ser necessária em fase posterior, o AAnF terá instância própria na VM 2.

---

## 7. Verificações a executar antes de fixar a decisão

| # | Verificação | Comando ou ação | Critério |
|---|---|---|---|
| 1 | Versão de kernel da máquina do free5GC | `uname -r` | Registrar; relevante apenas para o UPF |
| 2 | Versão de Go da máquina do free5GC | `go version` | Registrar no baseline |
| 3 | *Commit hash* do free5GC e submódulos | `git rev-parse HEAD` e `git submodule status` | Registrar no baseline |
| 4 | Presença de `AANF` e `naanf-akma` em `free5gc/openapi` | Inspeção do repositório | Determina se há *patch* no módulo (INCERTEZA 5) |
| 5 | Presença de `AkmaInd` e `RoutingId` em `AuthenticationInfoResult` | Inspeção do repositório | Idem |
| 6 | Alcançabilidade bidirecional VM1 ↔ VM2 | `ping`, `curl` ao NRF | Pré-requisito da F3 |
| 7 | Sincronização de relógio entre VMs | `chronyc tracking` ou `timedatectl` | Desvio abaixo de 1 s |

---

## 8. Decisão resumida

| Componente | Local | Motivo |
|---|---|---|
| AAnF | **VM separada** | Sem dependência de kernel; isolamento de dependências; corresponde à implantação autônoma de TS 33.535 §4.1 |
| AF de teste | **VM separada** | Idem |
| Oráculo de derivação (Python) | **VM separada** ou estação de trabalho | Sem dependência de rede |
| Patch do AUSF | **Máquina do free5GC** | Compila no grafo de módulos do free5GC |
| Patch do UDM | **Máquina do free5GC** | Idem |
| Patch de `free5gc/openapi` | **Máquina do free5GC** | Afeta todas as NFs que o importam |
| UPF | **Máquina do free5GC** | Única com dependência de kernel (`gtp5g`) |

**A implementação em máquina separada é viável e é a opção recomendada.** A restrição de kernel que motivou a dúvida não alcança o AAnF. A restrição que efetivamente existe é de grafo de dependências Go e atinge apenas os patches do AUSF e do UDM, que por natureza pertencem ao free5GC.

---

## 9. Referências

- TS 33.535 V19.0.0, cláusulas 4.1, 4.4.0, 6.7
- TS 29.510 V19.7.0 — `NFType` `AANF`, `AanfInfo`, descoberta por `routing-indicator`
- TS 29.522 V19.7.0 — `AkmaAfKeyData.expiry` (`DateTime`)
- Guia oficial de instalação do free5GC — requisito de kernel do UPF e versão de Go
- `free5gc/gtp5g` — README e histórico de versões
- Documentos internos: `00-baseline-normativo.md` §6, `02-adendo-B-sba-nrf.md` §§B.1, B.5.3, `03-adendo-C-ts29122.md` §C.2
