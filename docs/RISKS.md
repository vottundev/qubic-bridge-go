# Crypto Token Bridge Project: Risks and Dependencies Document

## Abstract
The purpose of this document is to outline potential risks and dependencies for the development of the Vottun Qubic Bridge, including mitigation strategies and responsible parties.

**The document is 'live' and will be updated in the coming days as new dependencies or risks are discovered, so check back regularly for updates.**

---

## Risks

| Risk ID | Risk Description | Impact Level | Probability | Mitigation Strategy | Responsible |
|---------|-------------------|--------------|-------------|----------------------|-------------|
| R1      | EVM Smart contract vulnerabilities leading to token loss or theft | High          | Medium      | Conduct comprehensive audits with third-party security experts before launch; perform extensive testing. | ceseshi |
| R2      | Qubic Smart contract vulnerabilities leading to token loss or theft | High          | Medium      | Conduct comprehensive audits with third-party security experts before launch; perform extensive testing. | anarojoagusti |
| R3      | Cross-chain bridging communication can fail in the middleware, resulting in delayed or lost token transactions. | High          | Medium      | Implement a retry mechanism for failed transactions as well as a refund strategy | alexlopezt |
| R4      | Rising gas prices make bridges economically unviable to use | Medium        | High        | Implementation of an adaptive fee structure; consideration of batching of transactions in order to reduce costs. | alexlopezt, rasito99 |
| R5      | Network congestion affecting transaction confirmation times | Medium        | Medium      | Notify users of potential delays and provide flexible transaction timeout settings. | alexlopezt |
| R6      | Poor user experience due to complex transaction processes | Medium        | Medium      | Design a user-friendly interface; conduct usability testing and gather feedback. | bmora-vottun, vicargo |
| R7      | Inadequate liquidity for token swaps across the bridge | Medium        | Low         | Establish partnerships with liquidity providers and incentivize early liquidity provision. | qubic team |
| R8     | Loss of trust as a result of security incidents or miscommunication with users | High          | Medium      | Develop a robust incident response plan and maintain clear, transparent communication with users. | TBD |

---

## Dependencies

| Dependency ID | Dependency Description | Impact Level | Mitigation Strategy | Responsible |
|---------------|------------------------|--------------|----------------------|-------------|
| D1            | Dependency on Redis pubsub strategy to be implemented and released to get smart contract events | High          | Release the strategy at the development stage as soon as possible, and later on the main network before releasing the bridge. | qubic team |
| D2            | Integration with third-party wallet providers | Medium        | Ensure compatibility with leading wallet providers; test integrations for stability on a regular basis. | vicargo - alexlopezt |
| D3            | We will need to work with liquidity providers and exchanges for token swaps, assuming liquidity pools are used. | High          | Maintain active communication with partners early in the development process. | qubic Team |
| D4            | Reliance on the stability of blockchain networks and APIs (e.g. Ethereum, Binance Smart Chain) | High          | Maintain communication with the blockchain providers; establish a fallback mechanism for API failures. | vottun team, qubic team |
| D5            | Security audits and review processes by 3rd party | High          | Schedule audits well in advance. | rasito99 |
| D6            | Dependency on reliable server infrastructure for bridge operations | High          | Use scalable, decentralized infrastructure options where possible; have disaster recovery plans. | vottun Team |
| D7           | Secure and compliant storage for user transaction data | Medium        | Implement data encryption and access control policies; work with compliance to align on data retention. | alexlopezt |

---
