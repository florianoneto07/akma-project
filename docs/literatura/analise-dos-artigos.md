# Análise da literatura sobre AKMA

**Data:** 16 de outubro de 2026
**Acervo:** `papers/` (fora do controle de versão) — 15 PDF, dos quais dois são idênticos (mesmo SHA-256): 14 trabalhos distintos. Os SHA-256 conferem com `INDICE.txt`.
**Papel no projeto:** contexto, trabalhos relacionados, propriedades para verificação formal futura e antecipação de limitações. **Nenhuma modificação de protocolo proposta nesses trabalhos entra no baseline.**

Os artigos não são redistribuíveis: vários são cópias licenciadas de editoras. Ficam fora do controle de versão (`.gitignore`, `/papers/`).

---

## 1. Inventário

| # | Trabalho | Publicação | Tipo | Versão de TS 33.535 analisada |
|---|---|---|---|---|
| 1 | Khan, Ginzboorg, Niemi — *AKMA: Delegated Authentication System of 5G* | IEEE Communications Standards Magazine, set. 2021 | visão geral; questões de privacidade | Rel-17 em elaboração |
| 2 | Khan, Ginzboorg, Niemi — *Privacy Preserving AKMA in 5G* | ACM SSR '19 | proposta de "modo privacidade" | 16.1.0 |
| 3 | Akman, Ginzboorg, Niemi — *AKMA for Secure Multi-access Edge Computing Mobility in 5G* | ICCSA 2022 Workshops, LNCS 13380 | protocolo para MEC sobre AKMA | 17.1.0–17.4.0 |
| 4 | Akman, Ginzboorg, Damir, Niemi — *Privacy-Enhanced AKMA for Multi-Access Edge Computing Mobility* | Computers 12(1):2, 2023 (dois arquivos idênticos no acervo) | análise formal (ProVerif) + protocolo modificado | 17.1.0–17.6.0 |
| 5 | Yang et al. — *Formal Analysis of 5G AKMA* (preprint) | submetido ao J. Systems Architecture, 2022 | análise formal (Tamarin) | 17.4.0 |
| 6 | Yang et al. — *Formal Analysis of 5G Authentication and Key Management for Applications (AKMA)* | J. Systems Architecture 126 (2022) 102478 | versão publicada do item 5 | 17.4.0 |
| 7 | Boureanu, Wesemeyer, Rajaona, Schneider, Treharne — *A Systematic Study of Practical & Formal Privacy in the 5G AKMA Procedure* | IEEE EuroS&P 2025 (DOI 10.1109/EuroSP63326.2025.00010) | análise formal de privacidade (Tamarin) + patch "AKMA_p" | não indicada |
| 8 | Yang et al. — *AKMA+: Security and Privacy-Enhanced and Standard-Compatible AKMA for 5G Communication* | não identificada no texto; artefato Zenodo DOI 10.5281/zenodo.14649847 | protocolo modificado + Tamarin + implementação em C (OpenSSL) | 18.2.0–19.0.0 |
| 9 | Boureanu et al. — *Post-Compromise Security with Application-Level Key-Controls – with a comprehensive study of the 5G AKMA protocol* | ACM ASIA CCS '25 | arcabouço de PCS + análise (ProVerif, Tamarin) + implementação no Open5GCore | 16.0.0 (citada) |
| 10 | Li, Chen, Gao, Zhang — *E2E-AKMA* | não identificada no texto | protocolo modificado contra comprometimento do AAnF; implementação em Go | não indicada |
| 11 | *Efficient Privacy-Preserving 5G Authentication and Key Agreement for Applications (5G-AKMA) in MEC* | IEEE Trans. Network and Service Management 23 (2026) | protocolo modificado; Scyther, ROR; custos com MIRACL e NS-3 | 17.3.0 |
| 12 | *An Enhanced Application Authentication and Key Management in 5G* | IOP Publishing (EEICE, 2023) | propostas de melhoria (renovação de chave, domínio de confiança) | 16.1.0–17.7.0 |
| 13 | *Toward an Era of Secure 5G Convergence Applications: Formal Security Verification of 3GPP AKMA with TLS 1.3 PSK Option* | Applied Sciences 14 (2024) 11152 | análise formal (ProVerif) e emulação do Ua\* TLS 1.3 PSK com OpenSSL 3.2.1 | 18.3.0–18.5.0 |
| 14 | Mohan, I V — *Selection of 5G AKMA Authentication Policy and Application Function Based on a Tenant Identifier Obtained from a UE eSIM* | Technical Disclosure Commons, Defensive Publication 7161, 2024 | publicação defensiva (multi-inquilino) | — |

`[INTERPRETAÇÃO]` Todas as análises formais tratam versões de TS 33.535 anteriores ao baseline (17.x e 18.x, com exceção do item 8). TS 33.535 V19.0.0 é a 18.8.0 republicada sem mudança técnica; resultados obtidos sobre 17.x precisam ser reconferidos antes de citados como válidos para V19.0.0.

---

## 2. Implementações existentes

| Trabalho | Plataforma | Abrangência | Código |
|---|---|---|---|
| 9 (Boureanu et al.) | Open5GCore (Fraunhofer, em C), Ubuntu 18.04 | AKMA acrescentado ao NEF do Open5GCore, seguindo o YAML `Nnef_AKMA`; derivações de K_AKMA no núcleo | Open5GCore não é aberto; os autores publicaram partes do código. Afirmam ser *"the first implementation of AKMA which is publicly-available"* |
| 8 (AKMA+) | C + OpenSSL; HPLMN em notebook, UE em Raspberry Pi 5 | emulação das entidades (UDM, AUSF, AAnF, AF, UE), sem núcleo 5G real | Zenodo 10.5281/zenodo.14649847 (cerca de 11 290 linhas) |
| 10 (E2E-AKMA) | Go (biblioteca ligada a Swift no iOS); servidores em nuvem | emulação da rede doméstica e do AF; UE em telefone real | não localizado no texto |
| 13 (TLS 1.3 PSK) | OpenSSL 3.2.1 em VM Ubuntu 20.04 | só o Ua\* (modos PSK-only, PSK-(EC)DHE, 0-RTT) | não localizado no texto |
| 11 | MIRACL (custos criptográficos) e NS-3 | simulação | — |

`[INTERPRETAÇÃO]` Nenhum trabalho do acervo implementa AKMA sobre o free5GC, nem com o AAnF como NF própria registrada no NRF, nem com o Ua\* sobre um núcleo real. A afirmação "primeira implementação aberta sobre o free5GC" fica restrita ao acervo examinado; não deve ser feita sem busca mais ampla.

---

## 3. O que os trabalhos dizem sobre o protocolo padrão

| Constatação | Trabalhos | Efeito neste projeto |
|---|---|---|
| **Falsificação de AF:** a norma não manda verificar explicitamente que o AF tem direito ao AF_ID que apresenta; um AF que capture o A-KID pode pedir a K_AF de outro AF e se passar por ele | 4 | O AAnF deve associar cada AF_ID autorizado à identidade autenticada do AF. É verificação de implementação dentro da "política local" de TS 33.535 §6.2.1 passo 2, não modificação de protocolo. Sem autenticação de transporte (DV-1) a associação fica fraca — risco residual a declarar |
| **Sem confirmação de chave, a autenticação falha:** os objetivos de autenticação UE↔AF não se cumprem se o Ua\* não confirmar a posse de K_AF dos dois lados | 5, 6 | O Ua\* simplificado tem de confirmar a chave nos dois sentidos (no Digest, `rspauth`) — critério de aceitação do MVP-3 |
| **Rastreabilidade:** o mesmo A-KID é apresentado a vários AF; a rede doméstica sabe todos os AF que o assinante usa; o SUPI vai a AF internos | 1, 2, 4, 7, 8, 10, 11 | Limitações do padrão, a declarar no tutorial. Resolver exige modificar o protocolo — fora do escopo |
| **Comprometimento do AAnF expõe tudo** (todas as K_AF deriváveis) | 8, 10 | Limitação do padrão; reforça TLS e isolamento do AAnF fora do MVP |
| **Uso prolongado da mesma K_AF; segurança pós-comprometimento fraca** — K_AF só muda com nova autenticação primária | 9, 12 | Tempo de vida de K_AF finito e configurável (TS 33.535 §4.4.0, §5.2); medir em MVP-6+ |
| **TLS 1.3 PSK:** 0-RTT sem sigilo futuro e sujeito a repetição; PSK-(EC)DHE troca eficiência por segurança | 13 | Para o perfil K.3.3 (MVP-6+): `psk_dhe_ke`, sem 0-RTT |

---

## 4. Coisas a não confundir

- **Dois "AKMA+" diferentes.** O AKMA+ de Yang et al. (item 8, Zenodo 10.5281/zenodo.14649847) é o que o README e o baseline citam como referência de *testbed*. O "AKMA+" de Boureanu et al. (item 9) é outro: um patch de segurança pós-comprometimento que força nova autenticação primária a cada K_AF. A documentação deve dizer qual dos dois cita.
- Os itens 5 e 6 são o mesmo trabalho (preprint e versão publicada). Citar a versão publicada.
- O item 14 é publicação defensiva, não revisada por pares.

---

## 5. Propriedades para verificação formal futura

Derivadas dos trabalhos acima, para uma fase posterior; nenhuma é requisito do MVP.

1. Sigilo de K_AKMA e de K_AF diante de atacante Dolev-Yao no Ua\* e nas SBI sem TLS.
2. Acordo injetivo UE↔AF sobre (A-KID, K_AF), com e sem confirmação de chave no Ua\*.
3. Associação AF↔AF_ID: o AAnF só entrega K_AF ao AF dono do AF_ID (ataque do item 4).
4. Não-rastreabilidade do UE entre AF e entre sessões (itens 7 e 8), registrada como propriedade que o padrão **não** satisfaz.
5. Segurança pós-comprometimento sob o arcabouço do item 9.

Modelos de partida: Tamarin (itens 5–8) e ProVerif (itens 4, 9, 13). Precisam ser atualizados para TS 33.535 V19.0.0 e incluir os parâmetros concretos deste projeto (identificador Ua\*, formato do A-KID).

---

## 6. Codificação do A-TID

Nenhum trabalho fixa a codificação textual do A-TID na NAI. O item 8 descreve as derivações *"following 3GPP TS 33.535"* sem tratar da serialização. A lacuna C continua aberta; a decisão do projeto (base64, ADR-0002) não é afetada.
