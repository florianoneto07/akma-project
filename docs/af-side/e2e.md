# Ponta a ponta e testes negativos (MVP-5 e MVP-6)

**Data:** 18 de outubro de 2026; MVP-6 em 19 de outubro de 2026 · **Resultado do MVP-6:** 83 verificações, nenhuma falha (execução `20261019-213720`, seção 5.2) · **MVP-5:** 63 verificações, nenhuma falha (execução `20261018-171510`) · **Instalação limpa:** 63 verificações, nenhuma falha, em 19 de outubro de 2026 (execução `20261019-160723`, seção 5.1)

## 1. O que o roteiro faz

`af/scripts/e2e.sh` executa o fluxo AKMA inteiro no laboratório e os testes negativos do plano, de forma reprodutível:

- autenticação primária real;
- K_AKMA e A-KID no AUSF e no UE;
- registro no AAnF;
- K_AF entregue ao AF;
- autenticação Ua\* do UE pelo plano de usuário.

Ele roda na máquina do AF e conduz a do free5GC por SSH. Lá, pede a senha de `sudo` **uma vez**: a parte remota (`af/scripts/e2e-remote.sh`) roda numa única sessão `ssh -t`, porque o free5GC (UPF), o *namespace* e o UE precisam de `root`.

As duas máquinas se sincronizam por arquivos-marcadores em três pontos (A, B, C). Em cada um, a parte remota espera enquanto a máquina do AF faz as verificações que precisam partir dela: o AAnF autoriza a K_AF por AF_ID e endereço de origem, e só a máquina do AF é o AF autorizado.

| Grupo | Onde | O que verifica |
|---|---|---|
| S | AF, antes da execução | autoteste do ME (`nr-akma-selftest`) com os vetores do MVP-0 |
| R1 | free5GC | NAT do N6, *namespace* do UE, `linkIp` do gNB, assinatura AKMA; free5GC com `-aanf`, gNB e UE no ar |
| R2 | free5GC | registro; A-KID do UE = do AUSF = do AAnF; Ua\* do UE pelo plano de usuário, 3 requisições |
| N (ponto A) | AF | testes negativos da API do AAnF, do `register-anchorkey`, do RID e do AF (seção 4); acesso anônimo direto ao AAnF (N01c, MVP-6) |
| E (ponto A) | AF | MVP-6: túnel até o NRF; três AF extras (externo na 8083, externo anônimo na 8084, interno anônimo na 8085); API do NEF (seção 4.1) |
| R2n | free5GC | MVP-6: Ua\* do UE com os três AF extras; o AAnF entrega via NEF; o NEF repassa a K_AF e não registra o SUPI |
| R3 | free5GC | nova autenticação: A-KID novo, contexto substituído, Ua\* com o novo; Ua\* com o A-KID e a K_AF anteriores, ainda aceitos pelo AF até `expiry` |
| B (ponto B) | AF | os AF extras se descadastram do NRF (E12); o AAnF recusa o A-KID anterior e aceita o novo |
| R4 | free5GC | AAnF parado durante uma autenticação (D9) |
| R5 | free5GC | assinante sem assinatura AKMA; depois restaurada |
| R6 | free5GC | `remove-context`, e de novo (404) |
| C (ponto C) | AF | `retrieve` depois do `remove-context` |
| R7, L1 | os dois | nenhuma cadeia de 64 hexadecimais nos logs do free5GC, do UE, do gNB e do AF |

## 2. Pré-requisitos

Na máquina do free5GC, as partes 3 a 7 do tutorial, isto é:

- o repositório em `~/akma-free5gc`;
- os *patches* do NRF, do UDM, do AUSF, da configuração e do `run.sh`;
- o AAnF implantado (`bin/aanf`, `config/aanfcfg.yaml`) e o `seed` em `~/akma-tools/`;
- o UERANSIM com os *patches*, o `linkIp: 10.100.0.1` do gNB e o `uaclient` em `~/akma-free5gc/af/build/bin/`.

Na máquina do AF: `af/scripts/build-af.sh` e `free5gc/scripts/build-aanf.sh`; acesso SSH por chave à máquina do free5GC; e o `af/config/afcfg.yaml` com o `aanfUri` do AAnF da máquina do free5GC.

O roteiro confere esses itens no início (*preflight*) e mostra o estado dos *patches*. O NAT do N6 e o *namespace* do UE ele mesmo aplica.

## 3. Como rodar

```bash
af/scripts/build-af.sh
```

```bash
af/scripts/e2e.sh woca@192.168.16.130
```

O argumento é a máquina do free5GC, `usuario@endereço`. Sem ele, vale a do laboratório, `woca@192.168.16.130`.

Os endereços vêm de duas configurações:

- o do AAnF, do `aanfUri` de `af/config/afcfg.yaml`. No *preflight*, o roteiro confere que ele é o `sbi.bindingIPv4` do `~/free5gc/config/aanfcfg.yaml` instalado e para se não for. A parte remota toma o endereço desse mesmo `aanfcfg.yaml`. O `seed` recebe o endereço em `--aanf` nas duas máquinas;
- o dos AF extras, de relógio deslocado, do `listen` de `af/config/afcfg.yaml`, nas portas 8081 e 8082.

Com `--keep`, o free5GC, o gNB e o UE continuam no ar no fim. Sem ele, o roteiro os encerra. Os arquivos ficam em `af/build/e2e/<data-hora>/`:

| Arquivo | Conteúdo |
|---|---|
| `report.md` | resumo, A-KID de cada autenticação, tabela de verificações e as linhas AKMA do free5GC, do UE e do AF |
| `local-results.txt`, `remote/remote-results.txt` | uma linha `PASS` ou `FAIL` por verificação |
| `remote/` | `free5gc.log`, `ue.log`, `gnb.log`, saídas do cliente Ua\* e do `seed`, `state` (os A-KID) |
| `af.log`, `af-*.log` | o AF, e os AF de relógio deslocado |

O código de saída é 0 só se todas as verificações passarem. Leva cerca de dois minutos.

## 4. Testes negativos

A lista do plano (MVP-5), com o resultado esperado de cada um:

| Caso do plano | Verificações | Esperado |
|---|---|---|
| A-KID desconhecido | N06; R5.ua | `403 K_AKMA_NOT_PRESENT` no AAnF; `403 AKMA_KEY_NOT_AVAILABLE` no AF |
| rótulo `akid` aceito | N02 | 200 (ADR-0002) |
| RID `"012"` × `"12"` | N19, N20, N21 | o mesmo A-TID com `rid12` é outro A-KID: 403 |
| AF não autorizado | N03 | 403 sem causa (P3) |
| `afId` malformado (9 dígitos, sem ponto) | N04 | `400 MANDATORY_IE_INCORRECT` |
| `AkmaKeyInfo` com `supi` e `gpsi` | N17 | `400 OPTIONAL_IE_INCORRECT` |
| `kAkma` fora do padrão | N18 | `400 MANDATORY_IE_INCORRECT` |
| `Content-Type` incorreto | N09 | 415 |
| corpo acima do limite | N10 | `413 MAX_JSON_SIZE_EXCEEDED` |
| demais códigos do MVP-1 | N05, N07, N08, N11 a N16 | `400 MANDATORY_IE_INCORRECT`, `MANDATORY_IE_MISSING`, `INVALID_MSG_FORMAT`; `411 INCORRECT_LENGTH`; 501; `404 RESOURCE_URI_STRUCTURE_NOT_FOUND`; `400 INVALID_API`; 401 sem *token*; nunca 204 (DV-7) |
| relógio do AF deslocado | N22, N23 | 2 h adiantado: K_AF válida tratada como vencida, `503 AANF_UNAVAILABLE`; 2 h atrasado: K_AF aceita, que ficaria em uso 2 h além do `expiry` |
| reautenticação | R3, B1, B2 | A-KID antigo: 403 no AAnF; a K_AF antiga continua aceita no AF até `expiry` (TS 33.535 §6.4.2, §5.2) |
| assinante sem AKMA | R5 | sem indicação do UDM, sem contexto no AAnF, `403 AKMA_KEY_NOT_AVAILABLE` no AF |
| AAnF indisponível na autenticação | R4 | registro aceito, erro só no log do AUSF (D9), `503 AANF_UNAVAILABLE` no AF |
| `remove-context` seguido de `retrieve` | R6, C1 | 204, depois 404 `AKMA_CONTEXT_NOT_FOUND`; `retrieve`: 403 |
| — | N24, N25 | Ua\* com K_AF errada: 401; `Host` de outro AF: 421 |

Os testes da API do AAnF (N01 a N16) estão em `af/scripts/negative/aanf_api.py`. São requisições HTTP/1.1 cruas, sem o código Go do projeto. Os do `register-anchorkey` precisam de *token* e usam o `seed` na máquina do free5GC, com as opções de teste `--gpsi` e `--k-akma`.

`[MVP]` Os casos de RID, relógio e Ua\* usam contextos semeados com os vetores `rid-012` e `mnc-3-digits` do MVP-0, de outros SUPI, para não tocar o contexto do UE real; o roteiro os remove depois.

`[MVP]` O relógio deslocado é simulado pela opção de teste `clockOffset` do AF (`af.md`). O AF não muda o relógio do sistema.

### 4.1 API do NEF e acesso anônimo (MVP-6)

`af/scripts/negative/nef_api.py` faz requisições HTTP/1.1 cruas, sem o código Go do projeto. Ele registra no NRF, pelo túnel, um AF temporário para obter os *tokens* e o descadastra no fim (`docs/free5gc-side/nef.md` §2.1).

| Verificação | Caso | Esperado |
|---|---|---|
| E00 | os dois AF externos se registram no NRF pelo túnel | 2 registros |
| E01 | K_AF via NEF | 200 com `kaf` e `expiry`, sem `supi` e sem `gpsi` (TS 33.535 §6.3 passo 5) |
| E02 | a mesma K_AF via NEF e direto do AAnF, para o mesmo A-KID e AF_ID | iguais; nenhuma das duas é impressa |
| E03 | acesso anônimo via NEF | 200 sem `supi` e sem `gpsi` (§6.2.2; TS 29.522 tab. 5.14.5.3.3-1 NOTA 3) |
| E04 | sem *token* | 401 (*middleware* do free5GC, sem `ProblemDetails`: F12) |
| E05 | *token* de outro serviço (`naanf-akma`) | 401 |
| E06 | AF_ID fora da política do NEF | 403 sem causa (§6.3 passo 2) |
| E07 | AF_ID aceito pelo NEF e não pelo AAnF (`af3`, entrada de teste do `nefcfg.yaml`) | 403 sem causa, repassado (P11) |
| E08 | A-KID desconhecido | 403 `K_AKMA_NOT_PRESENT`, repassado (TS 29.522 §4.4.23.2); nunca 204 |
| E09, E10 | `afId` malformado; `aKId` ausente | 400 com `invalidParams` |
| E11 | `Content-Type` `text/plain` | 415 (TS 29.122 tab. 5.2.6-1) |
| E12 | os AF externos se descadastram do NRF ao encerrar | 2 descadastros |
| N01c | acesso anônimo direto ao AAnF (`anonInd`) | 200 sem `supi` |

## 5. Resultado em 18/10/2026

`63 passed, 0 failed`: 30 verificações na máquina do free5GC e 33 na do AF. Os cinco A-KID da execução, um por autenticação primária:

| Autenticação | A-KID | Situação |
|---|---|---|
| 1 | `rid0000.atidmq7LS8IJTQYlRh7lrP2XD1SpvrLD1uIU5AOoM7RJzYo=@…` | registrado; Ua\* aceito |
| 2 (nova autenticação) | `rid0000.atidAdq+4U+8HPV+Bw1n5F+EfWvjyoEqAE9rE6nMU8nP3Ic=@…` | substituiu o 1; o 1 passou a 403 no AAnF |
| 3 (AAnF parado) | `rid0000.atidqU2g296EbEun+yEPMArMbIpkKogSw5ZRmhaz8eTLNOI=@…` | derivado pelo AUSF; registro falhou (D9) |
| 4 (sem assinatura) | `rid0000.atidATIugkCuUJTyczuW/cDgcdldY0jOravQOMEMcmEUShk=@…` | derivado só pelo UE; desconhecido do AAnF |
| 5 (restaurada) | `rid0000.atidy6lXvgC9uQe8p/NNUb0B+rQtNOJTI2VLGBtqfV/bDeQ=@…` | registrado; removido no R6 |

No console do `run.sh` aparece uma cadeia de 64 hexadecimais: o nome do arquivo de CDR que o CHF do free5GC apaga ao encerrar. Não é material AKMA, e a verificação R7 não a alcança, porque copia o `free5gc.log` antes do encerramento.

### 5.1 Instalação limpa, em 19/10/2026 (critério 4)

`63 passed, 0 failed` (30 e 33), na execução `20261019-160723`. Ambiente:

- a VM 192.168.16.131, com o free5GC v4.2.3 e o UERANSIM instalados pelo guia oficial e sem nenhum *patch* do projeto;
- o tutorial seguido do início, partes 1 a 8, a partir de um clone limpo do repositório na máquina do AF;
- `woca@192.168.16.131` como argumento do roteiro.

A primeira passagem achou três problemas que impediam o roteiro fora da VM 192.168.16.130:

- o `e2e.sh` tinha o endereço do AAnF fixo;
- o `e2e-remote.sh` esperava o AAnF em `192.168.16.130:8000`;
- os dois chamavam o `seed` sem `--aanf`, e o padrão dele é o AAnF do laboratório.

Os endereços agora vêm das configurações (seção 3). As divergências do tutorial e a correção de cada uma estão no plano, MVP-5, critério 4.

Os cinco A-KID da execução:

| Autenticação | A-KID |
|---|---|
| 1 | `rid0000.atidvfVWKPf69ZWBMVQGrvOkLaDQFdGba3KyQ5U0gCgNJFg=@…` |
| 2 (nova autenticação) | `rid0000.atidEnHVYP6TrlaX6wXXUX9Y7QNGnhfI3zcLH/mb/s6INFA=@…` |
| 3 (AAnF parado) | `rid0000.atidSvpTFXFecorSiZI+PGTCEXoMHM9XMgnQAoi3LAJiebw=@…` |
| 4 (sem assinatura) | `rid0000.atidHGpmuka9FkKJw7ogzBXHnbejSfqB1XfXUHVpT4iHchc=@…` |
| 5 (restaurada) | `rid0000.atidSaWLZ7zwKidV8wciRplt7HcrhqFnt2YAV+Mu5eBlhs4=@…` |

### 5.2 MVP-6, em 19/10/2026

`83 passed, 0 failed` na execução `20261019-213720`: 36 verificações na máquina do free5GC e 47 na do AF. A configuração é a do MVP-6:

- AAnF em 192.168.16.130:8010;
- NEF corrigido em 192.168.16.130:8000;
- N2 e N3 em 192.168.16.130;
- gNB local com `ngapIp` 192.168.16.130 e `gtpIp` 10.100.0.1.

As 20 verificações novas são E00 a E12, N01c e R2n.ua-nef, R2n.ua-nef-anon, R2n.ua-anon, R2n.aanf-via-nef, R2n.nef-relayed e R2n.nef-no-supi. As 63 do MVP-5 passaram com a configuração nova.

Antes dela, duas execuções no mesmo dia:

- **`20261019-203849` (63/0):** a regressão do MVP-5 com a configuração nova.
- **`20261019-213431` (80/3):** as três falhas eram do roteiro, não do sistema. A função que gera as configurações dos AF extras falhava com `set -u` quando chamada com três argumentos, e o AF externo da porta 8083 não subiu. Corrigido.

O teste com o UERANSIM em outra máquina (critério 6) está em `docs/af-side/teste-nef.md`. O roteiro supõe o gNB e o UE na máquina do free5GC.

## 6. Limitações

- **Senha de `sudo`.** É pedida uma vez por execução; o roteiro não altera o `sudoers`.
- **Endereços.** O roteiro não lê a máquina do free5GC das configurações. Ela é o argumento, e o padrão é a do laboratório.
- **Ordem fixa.** Uma falha no meio não interrompe o roteiro: as verificações seguintes rodam no estado que houver e podem falhar em cascata. A primeira `FAIL` do relatório é a que importa.
