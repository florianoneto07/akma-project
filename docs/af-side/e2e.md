# Ponta a ponta e testes negativos (MVP-5)

**Data:** 18 de outubro de 2026 · **Resultado:** 63 verificações, nenhuma falha (execução `20261018-171510`)

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
| N (ponto A) | AF | testes negativos da API do AAnF, do `register-anchorkey`, do RID e do AF (seção 4) |
| R3 | free5GC | nova autenticação: A-KID novo, contexto substituído, Ua\* com o novo; Ua\* com o A-KID e a K_AF anteriores, ainda aceitos pelo AF até `expiry` |
| B (ponto B) | AF | o AAnF recusa o A-KID anterior e aceita o novo |
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

Na máquina do AF: `af/scripts/build-af.sh` e `free5gc/scripts/build-aanf.sh`, e acesso SSH por chave à máquina do free5GC.

O roteiro confere esses itens no início (*preflight*) e mostra o estado dos *patches*. O NAT do N6 e o *namespace* do UE ele mesmo aplica.

## 3. Como rodar

```bash
af/scripts/build-af.sh
```

```bash
af/scripts/e2e.sh
```

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

## 6. Limitações

- **Senha de `sudo`.** É pedida uma vez por execução; o roteiro não altera o `sudoers`.
- **Reprodução.** Uma reprodução totalmente limpa, a partir do free5GC v4.2.3 recém-instalado e só com o tutorial, exige uma VM nova ou um *snapshot* (critério 4 do MVP-5). Ver `docs/plano/plano-de-desenvolvimento.md`.
- **Ordem fixa.** Uma falha no meio não interrompe o roteiro: as verificações seguintes rodam no estado que houver e podem falhar em cascata. A primeira `FAIL` do relatório é a que importa.
