# WIBX LABS — OPERATING CHARTER
**Version 1.1 · Governing ruleset · Machine-readable**

> **What this file is.** The governing ruleset for WibX Labs. It states rules; it does not argue for them. Every rule here is binding on Labs, on Labs' agents, automations, and projects.
>
> **What this file is not.** It is not the reasoning behind the rules, and it contains no assumptions, warnings, or open questions. Those live in internal context documentation. **A rule believed to be wrong is amended under §11, never hedged, annotated, or worked around.**
>
> **Section numbers are stable.** Cite them. They do not renumber between versions; superseded sections are marked void.

---

## §1 — Objective

**1.1** WibX Labs exists to remove the structural limits that prevent the company from operating at scale.

**1.2** It does this by embedding inside departments, transferring capability to the people closest to the problem, and leaving behind systems those departments own and never require Labs to return to.

**1.3** Labs output is measured in capability transferred, not work delivered.

**1.4** Labs is not a service desk, not an internal consultancy, and not a permanent owner of what it builds.

---

## §2 — North Star and Primitives

**2.1** Company North Star: **scale and profitability.** Prepare the organization to serve millions, then billions, of users. Scale is the forcing function.

**2.2** An initiative that cannot state its path to §2.1 does not enter the portfolio.

**2.3** There are exactly two strategic outcome primitives:

| Primitive | Definition |
|---|---|
| **Autonomy** | Reduce dependency on a specific individual or gatekeeper. |
| **Scalability** | Survive 10x–1000x volume without structural failure. |

**2.4** Efficiency is not a primitive and is not scoreable. It is a lever, not a destination. *Rationale retained because its absence causes misapplication: "efficiency" previously covered six unrelated department pathologies, producing false equivalence across departments and destroying the ranking axis.*

**2.5** Communication discipline — departments surfacing dependencies and work-in-progress to each other — is a company operating norm taught during embeds using existing tools. **It is not a Labs build.**

---

## §3 — *Void (v1.1). Removed.*

---

## §4 — Unit of Work

**4.1** Two classes of obstacle exist: problems inside a department (**boxes**) and problems in the dependency between two departments (**edges**).

**4.2** The default unit is the department. Diagnosis and embedding happen inside the box, every time.

**4.3 Discriminator.** Does the problem survive an internal fix?
- **No** → it was a box problem. Closed.
- **Yes** → it is an edge.

**4.4** The §4.3 test is retrospective and evidence-based. It is never applied at intake. *Rationale retained: applying it at intake permits the framing to be chosen to suit a desired conclusion.*

**4.5** An edge claim must record the evidence that the problem regenerated after the box was fixed. An edge asserted without that record is not an edge.

**4.6** The box is fixed first, to **the department's own standard** — not to a speculative future edge. A stable box is the precondition for closing an edge.

**4.7** While fixing a box, Labs watches for the edge but builds nothing for it. Preparation begins only once an edge is observed under §4.3. Expecting an edge is attention, not construction.

**4.8** A confirmed edge is worked from both ends concurrently — one Labs member embedded in each of the two departments. Edges do not require sequential passes.

**4.9** Edges currently confirmed: **Product ↔ IT** · **Design ↔ CS / Marketing / Commercial** · **Marketing ↔ Product**.

---

## §5 — Prioritization

**5.1** A request is scored, not judged.

**5.2 Modules.** There are exactly two scoring modules, one per primitive: **Autonomy** and **Scalability**.

**5.3** A department **loads** the module matching its assigned primitive. The module is never selected by discretion, and no weight is chosen during an evaluation.

**5.4** Criteria — three, and only three:

| Criterion | Definition |
|---|---|
| **Autonomy served** | How much the request reduces dependency on an individual or gatekeeper. |
| **North-Star service** | How much the request serves the company and department North Star. |
| **Readiness** | Existing stack and knowledge versus learning on the job. |

**5.5** **Effort is not a criterion.** The three criteria score value only. Effort is a **sequencing input**: among requests already judged worth doing, effort orders the queue. Effort never raises or lowers a score. *Rationale retained: value and cost in one sum lets a cheap trivial request outrank an expensive critical one.*

**5.6 Scoring standard.** A score is a **measured delta**, never a rating. The score is how far the request is projected to move that department's unit against its captured baseline.

**5.7** **No baseline, no score.** This is the same gate as §6.1.

**5.8** **Weight values are not yet set.** Until they are set by §11 amendment, no numeric score may be produced or reported. Requests are ordered by explicit comparison against §5.4, and any ordering so produced is marked **provisional**.

**5.8.1** While §5.8 holds, §5.2–5.3, §5.6 and §5.9–5.11 are dormant: they bind from the amendment that sets the weights.

**5.9** Weights, once set, are fixed and versioned. They are unchangeable inside an evaluation.

**5.10** **Labs is strictly bound by its own scores.** There is no discretionary override, for Labs or for anyone.

**5.11** A refusal never means *this scored low*. It means *these scored higher*. A refusal requires the whole current candidate field scored and live; an unscored or stale field cannot produce one.

**5.12 Grandfather clause.** Bifrost, Investigator, Nexus PR, and the CRM notification dashboard predate this standard, lack captured baselines, and are **grandfathered as pre-charter.** They are not scored retroactively for admission. The rubric applies to engagements opened after v1.0.

---

## §6 — Engagement Lifecycle

**6.1** Stages:

| Stage | Rule |
|---|---|
| **0 · Baseline** | Pre-engagement measurement captured and verified. **No baseline, no engagement.** |
| **1 · Proof** | Pre-handoff validation. Nothing is handed off unproven. |
| **2 · Handoff** | Deliverable transfers. **Labs extracts here.** Observation clock starts. |
| **3 · Observation** | **One month**, unsupervised. Telemetry collected remotely. |
| **4 · Certification** | Review data, revise if required, confirm maintenance requirements and owner. Engagement closed. |

**6.2** Baseline collection is an operational task owned by **Caio**. The charter states the standard; it does not perform the collection.

**6.3 Proof method by artifact type.**
- **6.3.1 Technical artifacts → stress test.** Simulate 10x–100x load before handoff.
- **6.3.2 Process and organizational interventions → scoped live pilot.** One squad, or one client cohort. Load cannot be simulated against an organization; shrink the blast radius instead.

**6.4** Extraction occurs at **Handoff**, not Certification. *Rationale retained: an embed that remains makes the period supervised and the telemetry invalid.*

**6.5** Observation is **one month**, not one week. *Rationale retained: process change reverts in weeks three and four; a week measures compliance only.*

**6.6** At Stage 0, the engagement's **principal maintenance actions are enumerated by name.** An engagement whose maintenance actions cannot be listed in advance is not sufficiently understood to open.

**6.7 Observation test.** The named owner either performed the §6.6 actions unassisted during the observation month, or did not. Failure is tested, not judged.

**6.8** **Any Labs assistance during the observation month fails that observation.** There is no trivial threshold. Labs is unavailable for the duration by design; that unavailability is the measurement.

**6.9** All assistance is logged, including informal assistance. **Unlogged assistance discovered after Certification voids that Certification retroactively.**

---

## §7 — Failed Observation

**7.1** A failed observation permits **one retry**: a single re-embed and one fresh observation month.

**7.2** Re-entry at Stage 3 is not a violation of §10.2. "Never return" governs closed engagements; an uncertified engagement has not closed.

**7.3** **A second failed observation terminates the engagement.** It is declared failed. There is no third attempt.

**7.4** A second failure is diagnostic information: the pathology was misidentified, or the department was never the correct owner. It is not a request for further help.

**7.5** On termination, the artifact routes to IT as production code or is decommissioned.

**7.6** The original pathology returns as unresolved and **re-enters through §5 as a new candidate**, scored against the current field, and may lose.

---

## §8 — Execution Model

**8.1** The execution model is concentration of force at a single point rather than dispersal across a front, speed of penetration, and refusal to settle into attrition.

**8.2** Applied: few departments in flight at once, decisive intervention rather than sustained presence, extraction at Handoff, never entrenched. Labs' failure mode is not losing — it is getting stuck.

**8.3** *Void (v1.1). The execution model carries no name.*

**8.4 Embedded Liaison.** Labs does not operate externally. Members embed inside the target department, observe workflow directly, and translate findings into systemic fixes. Communication between embeds is itself an observation target.

**8.5** Part of every embed is teaching the department to surface its own dependencies and progress to other departments (§2.5).

**8.6 Fusion Team.** A member of the target department is temporarily designated part of Labs for one operation, reporting to Labs while embedded in their home department.

**8.7** The Fusion Team's primary value is **ownership transfer**, not capacity. The participant returns to their department as the system's owner, requiring no handoff ritual.

**8.8** A Fusion Team is **mandatory for any engagement producing a maintainable artifact.**

---

## §9 — *Void (v1.1). Removed.*

---

## §10 — Hard Limits

**10.1 Diplomatic constraint.** Every engagement is framed as *capability unlocked for the department*, never as *deficiency exposed*. A technically successful engagement that damages a department head's standing has failed.

**10.2 Completion bar — completion over coverage.** An engagement succeeds only if the job is done fast, with no fallbacks and no need to return. A partial fix requiring Labs to re-enter later is recorded as a failure, not a phase.

**10.3** No baseline, no engagement. Stage 0 is a gate, not a formality.

**10.4** Labs does not replicate existing infrastructure. Where a tool exists and is underused, the intervention is discipline, not a rebuild.

**10.5** *Void (v1.1).*

**10.6** *Void (v1.1).*

**10.7** The charter is headcount-agnostic. It defines principles, not resource math, and holds at two people and at ten. Current capacity is an operating reality that constrains scheduling and is **never encoded as a decision rule.**

**10.8** *Void (v1.1).*

**10.9** The charter governs the organization, not a list.

---

## §11 — Amendment

**11.1** A rule believed to be wrong is **amended**, never overridden, hedged, or worked around. This binds Labs and every agent operating under this charter.

**11.2** An amendment acts on the **instrument** — a rule, criterion, scale, module, or weight — never on a live verdict or a live engagement.

**11.3** The trigger is **evidence of instrument fault**, not disagreement with an outcome. The misfiring component must be named and the fault stated. *"This ranking is wrong"* is not a trigger. *"The Readiness criterion does not account for a dependency on an unstaffed role"* is.

**11.4** Amendment requests are logged, including those declined. Repeated requests against the same component are themselves evidence that the component is faulty. The log is this repository's pull requests labelled `amendment`: merged is adopted, closed unmerged is declined. Each names the component and the fault (§11.3). v1.1 is the baseline and predates the log.

**11.5** **The live verdict or engagement executes as decided while an amendment is argued.** There is no pause. An amendment applies to the next verdict, never the one in flight.

**11.6** Superseded sections are marked void. Section numbers never renumber.

---

## §12 — *Void (v1.1). Removed.*

