# Agent Delivery Model — Concepts and Project Setup Guide

## 1. Purpose of this guide

This guide explains the concepts behind the PayFlow agent-delivery model:

```text
BA → Architecture Review (when required) → Developer → Tester
```

It explains:

- What each agent and artifact is responsible for.
- Why stages, gates, and handoffs exist.
- How pipeline status differs from an ADR decision status.
- How project conventions, agent instructions, skills, and rules work together.
- What alternatives exist and when they may be preferable.
- How to adapt the model to another project without copying PayFlow-specific
  requirements.

The goal is to help a project owner design an understandable delivery process
quickly, while preserving clear responsibility, verifiable evidence, and
explicit decisions.

---

## 2. The model in one view

A feature moves through a sequence of defined responsibilities:

```text
Requester
   │ feature request
   ▼
Business Analyst
   │ requirements.md
   ▼
Architecture Review, when required
   │ architecture-review.md
   ▼
Developer
   │ code changes + build-report.md
   ▼
Tester
   │ test-report.md
   ▼
Run closed if Gate 3 passes
```

Architecture review is conditional. If a change does not affect architecture,
the flow is:

```text
Requester → BA → Developer → Tester
```

If a change materially affects architecture, the flow is:

```text
Requester → BA → Architecture Review → Developer → Tester
```

Every stage has a defined job, a defined output, and a condition that must be
satisfied before the next stage starts.

### Why use stages?

A single agent asked to analyse requirements, make architectural decisions,
implement code, and approve the result combines incompatible responsibilities.
That makes it difficult to know:

- Which input was treated as authoritative.
- Who made a decision.
- Whether a requirement was omitted.
- Whether tests were actually run.
- Whether the final approval was independent.

Stages make responsibilities and evidence visible. They do not guarantee
correctness, but they make omissions and disagreements easier to identify.

---

## 3. Core terms

### Agent

An **agent** is a role with defined responsibilities, inputs, outputs, rules,
and success criteria. It is not simply a persona or tone of voice.

For example, the Developer Agent is responsible for implementation and build
evidence. It is not responsible for deciding that the Tester’s independent
acceptance criteria passed.

### Stage

A **stage** is a position in the delivery sequence. Each stage consumes
upstream information and produces an artifact for the next applicable stage.

A stage number describes process order, not team seniority. The current PayFlow
sequence uses:

| Stage | Agent | Main output |
|---|---|---|
| 1 | BA Agent | `requirements.md` |
| 2, when needed | Architecture Review Agent | `architecture-review.md` |
| 3 | Developer Agent | Code changes and `build-report.md` |
| 4 | Tester Agent | `test-report.md` |

Because Architecture Review is conditional, Developer is still Stage 3 in the
project’s definitions even when Architecture Review is skipped. Keep stage
identifiers consistent across agent files, audit entries, and reports.

### Handoff

A **handoff** is the transfer of work from one stage to another using defined
artifacts. A message such as “looks good, please continue” is not a complete
handoff when the process requires a requirements or build report.

A handoff tells the next stage:

- What to treat as authoritative.
- What work is complete.
- What evidence exists.
- What remains unknown.
- Whether the next stage may start.

### Gate

A **gate** is a set of checks the receiving stage performs before proceeding.
The producer prepares the artifact; the downstream stage checks it.

For example:

- The BA writes requirements.
- Architecture Review or Developer checks that the requirements satisfy Gate 1.
- The Developer produces implementation and build evidence.
- The Tester checks the Developer handoff against Gate 2.

A gate is not just a checklist for the producer. It is a receiving-stage
decision about whether the work is ready to continue.

### Artifact

An **artifact** is a durable file that records work or a decision. Examples
include requirements, an ADR, a build report, and a test report.

Artifacts are useful because they can be read later, reviewed independently,
and linked to audit entries. They also reduce reliance on conversational memory.

### Run

A **run** is one tracked attempt to deliver a ticket or feature. Its run-scoped
artifacts and audit history are stored together, for example:

```text
runs/20261025-1432-PAY-123/
```

The run ID distinguishes one delivery attempt from another. A ticket ID
identifies the business work; a run ID identifies a specific pipeline execution.

---

## 4. Project instructions, agent definitions, skills, and rules

These files solve different problems. Avoid combining them into one large
instruction file.

### `AGENTS.md` — project conventions

**Purpose:** Explain the project’s architecture and non-negotiable conventions.

It should answer questions such as:

- What systems and services exist?
- Which service owns which data?
- What security and compliance constraints apply?
- How is money represented?
- Where are canonical ADRs stored?
- What build and test commands should be used?
- What is the project’s ADR approval policy?

In PayFlow, `AGENTS.md` says that the ledger is authoritative, money uses
integer minor units, services do not access one another’s databases, and raw
PAN/CVV must not enter PayFlow systems.

**Why keep this separate from agent definitions?** The same project facts apply
to multiple agents. Duplicating them across each role creates conflicting
versions.

**Alternative:** Put all instructions in each agent file. This may work for a
tiny project, but changes must then be repeated consistently in several places.

### Agent definition files — role-specific instructions

Each agent definition describes one role’s work:

- Inputs and whether they are required.
- Outputs and allowed statuses.
- Start conditions.
- Procedure.
- Handoff requirements.
- Evaluation criteria.

**Why?** Each agent should know what it owns and what it must not do.

### `pipeline-rule.md` — workflow and gates

The pipeline rule defines:

- Stage order.
- Conditional routing.
- Artifact contract.
- Gate requirements.
- Rework routing.
- Retry limits and escalation.
- Audit-entry format.
- Run closure.

**Why separate it from agent files?** Orchestration rules must be consistent
across every stage. If each agent independently describes stage order or retry
limits, they can disagree.

### `frontend-skill.md` — reusable implementation expertise

A **skill** is reusable technical guidance for a specific kind of work. The
frontend skill applies to dashboard changes but not to backend-only changes.

In PayFlow it covers:

- Zod validation and inferred types.
- TanStack Query.
- UI loading, empty, error, and success states.
- Accessible forms and dialogs.
- Integer minor-unit handling.
- PSP tokenization and avoidance of raw card data.
- Idempotent payment submissions.
- Frontend testing.

**Why not put all frontend details in the pipeline rule?** The pipeline rule
should explain when frontend guidance applies. The skill should explain how to
build the UI.

---

## 5. BA Agent and requirements

### What the BA does

The Business Analyst translates a request into business requirements. The BA
defines **what** users need and **why**, not the technical solution.

The BA output should make it possible for a developer to implement the work and
a tester to verify it without guessing business behavior.

### User story

A user story expresses a capability from the perspective of a user or role:

> As a merchant admin, I want to issue a partial refund, so that I can resolve
> a partial dispute without refunding the entire payment.

The common form is:

> As a **[user or role]**, I want **[capability]**, so that **[benefit]**.

**Purpose:** Keep the requirement connected to a user or business outcome.

**Alternative:** A list of system functions may be adequate for technical
maintenance work, but user stories help expose who benefits and why.

### Acceptance criteria

Acceptance criteria define observable outcomes. Given/When/Then makes the
starting condition, action, and expected result explicit:

- **Given:** relevant starting state.
- **When:** action or event.
- **Then:** outcome that can be checked.

Example:

```text
Given a captured payment with 100.00 USD refundable,
When a merchant submits a refund for 30.00 USD,
Then the refund is created for 3000 minor units and 70.00 USD remains refundable.
```

**Purpose:** Give developers an implementation target and testers a verification
target.

A criterion such as “the refund flow works correctly” is not sufficiently
testable. It does not say which inputs, states, or outcomes matter.

### Scope

Scope specifies which behavior belongs in the change.

- **In scope:** Work the team agrees to deliver.
- **Out of scope:** Related work that is explicitly excluded.

**Purpose:** Avoid accidental expansion and disagreements about whether an
unrequested feature was required.

### Non-functional requirements

A non-functional requirement describes a system quality or constraint rather
than a user-visible capability. Examples include latency, auditability,
availability, security, and retention.

For payment-related work, PayFlow asks the BA to address:

- PCI-DSS scope and card-data handling.
- Idempotency.
- Auditability.
- Money and currency precision.
- Latency.
- Failure and timeout behavior.

**Why record unknowns?** An unknown requirement can affect design and testing.
Writing “unknown” makes the gap visible. Guessing a value creates a false
requirement.

### Hard constraint

A **hard constraint** is a limit or obligation an acceptable design must meet,
such as:

- Monthly infrastructure cost at or below a specified ceiling.
- p95 latency below a specified value at a stated workload.
- A specific compliance or data-residency obligation.
- A team’s available skills or operational capacity.

A constraint should include its source and meaning. “Fast” is not a measurable
latency constraint. “p95 time to first delivery attempt below two seconds at
5,000 jobs per second” is measurable.

**Why require sources?** The source tells reviewers whether a value is a
business requirement, a project convention, an estimate, or an assumption.

### Blocker

A **BLOCKER** is missing or unresolved information that prevents safe,
testable, or defensible downstream work.

Not every unknown is a blocker. For example, the exact dashboard locale may be
non-blocking if the feature does not involve formatted amounts. But an unknown
latency budget may block architecture review if the decision is specifically
about meeting latency.

**Rule:** Mark an unknown as a blocker when it prevents the next required
decision. Do not mark it as a blocker merely because it is unknown.

---

## 6. Architecture Review and ADRs

### When architecture review is required

Architecture review is required when a change may materially affect:

- Service boundaries or ownership.
- Data storage, ownership, or movement.
- Event or API contracts.
- Security or compliance controls.
- Infrastructure.
- Availability or scalability.
- Another material architecture decision.

For a local change within established boundaries, the BA may record
`architecture_review: not_required` with a reason.

If applicability is unclear, PayFlow treats review as required until the
uncertainty is resolved.

### What an ADR is

An **Architecture Decision Record (ADR)** is a durable record of an architecture
decision and its reasoning.

A useful ADR normally states:

1. **Context:** Why a decision is needed.
2. **Decision drivers:** Requirements and constraints that matter.
3. **Options:** The proposed design and credible alternatives.
4. **Decision:** What is recommended or accepted.
5. **Consequences:** What follows from the decision, including trade-offs.

An ADR is not just a design description. Its main value is preserving why the
team chose one option instead of another.

### Proposed design

The **proposed design** is the candidate design to assess. It can come from the
requester or an authorized design owner. The BA records its reference but does
not invent the technical design.

The Architecture Review Agent needs the proposal and constraints to make a
meaningful comparison. If either is missing in a way that prevents comparison,
the BA blocks the handoff.

### Alternatives

An ADR must compare credible alternatives. A rejected alternative should be
plausible enough that a reviewer can understand why it was considered.

Weak rejection:

> PostgreSQL polling was rejected because it has downsides.

Constraint-based rejection:

> PostgreSQL polling was rejected because the supplied p95 result is 4.8
> seconds against a 2-second latency budget, and its estimated monthly cost is
> $14,100 against a $12,000 ceiling.

**Purpose:** Show that the decision was selected by evaluating the scenario,
not by presenting a preferred solution and dismissing other options generically.

### Trade-offs

A **trade-off** is something the chosen option gives up or makes harder in
exchange for its benefits.

Example:

> Kafka workers isolate slow merchant endpoints, but delivery is at least once.
> Merchants must deduplicate notifications using the stable delivery ID.

**Why require trade-offs?** Every real design has costs. A decision with only
benefits is usually missing consequences or has not been examined critically.

### Concrete consequences

A consequence is a specific downstream effect that follows from the decision.

Weak:

> This design may create operational risk.

Concrete:

> The webhook service must alert on Kafka consumer lag, retry depth, exhausted
> deliveries, and delivery-record database growth. Its on-call runbook must
> describe safe replay using the original delivery ID.

Concrete consequences help teams estimate, plan ownership, implement, and test
the decision.

### Evidence, estimates, assumptions, and unknowns

These terms distinguish the confidence of claims:

- **Evidence:** A supplied measurement, executed test, or established project
  fact.
- **Estimate:** A forecast that has not been validated.
- **Assumption:** A premise used to compare options.
- **Unknown:** Information not currently available.

Example:

- Evidence: A load test measured 1.1 seconds p95 at 5,000 jobs per second.
- Estimate: Monthly infrastructure cost is projected to be $10,800.
- Assumption: The existing Kafka cluster has capacity for this service.
- Unknown: The target region’s retention controls have not been verified.

Do not present an estimate or assumption as a verified fact.

### ADR approval versus pipeline handoff

An ADR has a **decision status**, for example:

- `proposed`: recommended but not accepted by the authorized owner.
- `accepted`: accepted, with approval evidence.
- `rejected`: explicitly rejected.
- `superseded`: replaced by a later decision.

The pipeline artifact separately has a **pipeline status**, such as `pass` or
`blocked`.

For example:

```yaml
status: pass
decision_status: proposed
```

This means the Architecture Review artifact is complete enough for Developer
validation, but the decision has not been accepted.

Under the supplied PayFlow policy:

- Implementation may begin after the architecture handoff passes validation,
  even if the ADR remains `proposed`.
- Production launch requires acceptance by the architecture owner.
- The owner’s acceptance must be recorded.

These policies are project-specific. Another project may require ADR
acceptance before implementation. If so, record that policy in its project
conventions and block implementation until approval exists.

### Canonical ADR versus run-scoped review

The run-scoped `architecture-review.md` records the current delivery’s review
and handoff. It is not automatically a canonical project decision.

The canonical, accepted ADR belongs under:

```text
docs/architecture/adr/
```

Publish it there according to the project’s approval process. This separates
working review evidence from the lasting record of accepted architecture.

---

## 7. Developer Agent and build report

### What the Developer does

The Developer implements the passed requirements and follows the applicable
architecture decision. It adds or updates tests and records actual build and
unit-test evidence.

The Developer does not change passed requirements or ADRs silently. If
implementation reveals that an upstream decision must change, it requests
rework through the pipeline.

### Traceability table

A traceability table maps:

```text
Story → Acceptance criterion → Implementation → Test
```

Example:

| Story | Criterion | Implementation | Test |
|---|---|---|---|
| S1 | S1-AC1 | `services/refunds/create.ts` | `create.test.ts::partial_refund` |

**Purpose:** Make omissions visible. A missing row shows that a criterion has
not been connected to implementation and verification.

If an item is not implemented, state the reason. An explicit
“not implemented” is transparency; it is not evidence that a required
criterion is satisfied.

### Build report

The build report records:

- What changed.
- How each criterion maps to code and tests.
- Which local build and test commands actually ran.
- Their actual output.
- Applicable security checks.
- Frontend skill evidence, or “No UI impact.”

**Why require actual output?** “Tests passed” is a claim. A command and output
provide evidence that someone can reproduce or investigate.

### Frontend skill

For dashboard changes, the Developer applies the frontend skill. That includes
validation, asynchronous UI states, accessibility, safe money handling,
tokenized card handling, idempotency, and relevant tests.

For backend-only work, the Developer records “No UI impact” rather than
pretending that frontend requirements apply.

---

## 8. Tester Agent and test report

### What the Tester does

The Tester independently checks the implementation against the acceptance
criteria in `requirements.md`.

The Developer’s tests are useful evidence, but they do not replace the
Tester’s verification. The Tester must not mark a criterion as passing merely
because the Developer says it passed.

### Criterion-level verdict

Every acceptance criterion gets:

- A stable identifier.
- A test executed.
- Expected result.
- Actual result.
- Evidence reference.
- `PASS` or `FAIL`.

A criterion that was not executed or lacks evidence cannot be marked `PASS`.

### Negative and boundary testing

- **Positive test:** Valid case produces the required result.
- **Negative test:** Invalid or unauthorized action is rejected correctly.
- **Boundary test:** Minimum, maximum, and just-beyond-limit values behave
  correctly.
- **Regression test:** Existing behavior affected by the change still works.
- **Payment-specific test:** Idempotency, financial effect, currency precision,
  provider failure, or other relevant payment behavior is verified.

Not every category applies to every story. Record why a category is not
applicable rather than adding meaningless tests.

### Gate 3 authority

The Tester is the final Gate 3 decision-maker. This is not self-approval of its
own evidence: the Tester independently tests the implementation and records
the outcome. A downstream agent is not required to approve the test report.

If every required criterion passes with evidence, the Tester closes the run.
If a criterion fails or required verification is missing, the Tester returns
the work to the Developer.

---

## 9. Statuses: do not confuse the meanings

### Pipeline artifact status

The artifact status describes the stage handoff:

- `pass`: The producing stage considers the artifact ready for downstream gate
  validation. It does not mean the downstream stage approved it.
- `fail`: The producing stage could not produce a valid artifact in this
  attempt.
- `rework`: The artifact or implementation needs correction.
- `blocked`: A decision-critical issue or retry limit prevents progress.

The downstream stage still validates the handoff.

### ADR decision status

The ADR decision status describes the state of an architecture decision:

- `proposed`
- `accepted`
- `rejected`
- `superseded`

Example:

```yaml
status: pass
decision_status: proposed
```

This combination is valid if the review artifact is complete but project
approval is not yet recorded. Whether implementation may proceed depends on the
project’s approval policy.

### Test criterion verdict

A test criterion uses:

- `PASS`: The tested result matches the acceptance criterion and evidence is
  recorded.
- `FAIL`: The result is wrong, untested, or unsupported by evidence.

### Why separate these?

A single “status” field cannot accurately answer all of these questions:

- Is the artifact complete?
- Was the ADR accepted?
- Did a specific test pass?
- Is the run closed?

Separate fields make each decision explicit.

---

## 10. Rework, attempts, and audit

### Rework

**Rework** means a stage must correct an artifact or implementation before the
pipeline continues. It goes to the stage responsible for the failed work.

Examples:

| Problem | Responsible stage |
|---|---|
| Requirements omit a decision-critical constraint | BA |
| ADR lacks credible alternatives | Architecture Review |
| Build report lacks actual test output | Developer |
| Acceptance criterion fails in testing | Developer |

A stage must not silently modify a passed artifact produced by another stage.

### Retry limit

The PayFlow rule allows three attempts total per responsible stage:

```text
Initial attempt + at most two retries = three attempts
```

After the third failed attempt, the run stops, the artifact is marked
`blocked`, and the issue is escalated to a human owner.

**Purpose:** Avoid unbounded agent loops and ensure repeated failures receive
human attention.

### Artifact revisions

Passed artifacts are immutable. A rework attempt must produce a new artifact
revision rather than overwriting the prior one. The rule defines attempt
directories such as:

```text
runs/<run-id>/developer/attempt-01/build-report.md
runs/<run-id>/developer/attempt-02/build-report.md
```

The audit log identifies which artifact version is current. If an upstream
artifact changes, downstream artifacts based on its old version must be
revalidated.

### Audit log

The `pipeline-log.md` records stage attempts, gate decisions, rework, retries,
escalations, and closure:

```text
<timestamp> | <agent> | <stage> | <status> | <run/ticket> | <summary> | <next action>
```

The log is append-only: add new facts; do not erase or rewrite prior history.

**Purpose:** Provide an audit trail of what happened and why the run moved,
stopped, or returned for rework.

---

## 11. Why this order?

### BA before Developer

**Reason:** Implementation requires business behavior and testable acceptance
criteria. Starting without them invites assumptions and repeated clarification.

**Alternative:** Let the Developer clarify requirements during coding. This
can be efficient for a small internal change, but it mixes requirements
decisions with implementation and makes independent testing harder.

### Architecture Review before Developer when needed

**Reason:** A design decision can affect multiple services, data ownership,
event contracts, compliance, cost, or operational responsibility. Reviewing
those decisions before implementation can avoid building against an unsuitable
architecture.

**Alternative:** Let developers make the architecture decision while
implementing. This can be faster for small, low-impact changes. It is less
suitable when a decision has cross-service or compliance consequences.

### Developer before Tester

**Reason:** The Developer builds the feature and provides implementation and
build evidence. The Tester independently verifies that implementation against
the accepted requirements.

**Alternative:** Ask the Developer to self-test and close the work. This is
faster but reduces independence and can miss differences between intended
behavior and actual behavior.

### Tester as final gate

**Reason:** The Tester has direct evidence from independent test execution and
can compare it against the requirements. The pipeline rule makes the Tester
the final Gate 3 decision-maker.

**Alternative:** Add a separate human release approver after testing. This may
be appropriate in regulated release processes, but it adds another stage and
must be documented as part of the pipeline rather than implied.

---

## 12. Common design choices and alternatives

### Conditional versus mandatory architecture review

**PayFlow choice:** Conditional review, based on BA-recorded impact.

**Why:** A UI-only change within existing service boundaries may not need a full
ADR. A change to event contracts, data ownership, or service boundaries likely
does.

**Mandatory review alternative:** Every change receives architecture review.
This maximizes review consistency but costs time and can turn the review into a
routine formality.

**No formal review alternative:** Developers make all design decisions.
This minimizes process overhead but leaves cross-service decisions less
consistent and less documented.

### Run-scoped ADR versus direct canonical ADR

**PayFlow choice:** Review in the run folder; publish accepted decisions to
`docs/architecture/adr/`.

**Why:** A proposed decision needs review evidence, but not every proposal
should become a canonical project decision.

**Direct canonical ADR alternative:** Write all decisions directly into the
canonical ADR directory. This is simpler, but the directory can fill with
unaccepted proposals unless its status and approval process are carefully
managed.

### Conditional constraints versus numeric scoring

**PayFlow choice:** Treat hard constraints as feasibility gates, then explain
the decision among feasible choices.

**Why:** A weighted score must not make an option that violates a mandatory
compliance or budget limit look acceptable.

**Scoring alternative:** Score options across dimensions and select the highest
score. This can help compare preferences once mandatory constraints are met,
but it should not replace explicit constraint checks.

### One centralized pipeline rule versus duplicated rules

**PayFlow choice:** Store orchestration rules in one `pipeline-rule.md`.

**Why:** One source of truth reduces drift in stage order, status values, gates,
and retry behavior.

**Duplicated-in-agent-files alternative:** Each agent describes the entire
pipeline. This is easier to read in isolation but requires synchronized edits
whenever the process changes.

---

## 13. Adapting the model to a new project

Use this sequence to create a project-specific version.

### Step 1: Write project conventions

Create a project guidance file such as `AGENTS.md`. Record:

- Project purpose and system boundaries.
- Data ownership.
- Security and compliance rules.
- Coding conventions.
- Build and test commands.
- Canonical documentation locations.
- ADR approval policy.

Do not copy PayFlow-specific rules such as the ledger model or PCI-DSS scope
unless they apply to the new project.

### Step 2: Decide whether architecture review is needed

Choose whether review is:

- Required for every change.
- Conditional based on architectural impact.
- Optional and separate from the delivery pipeline.

Write the trigger conditions and who records applicability.

### Step 3: Define the stages

For every stage, define:

- Role and responsibility.
- Required inputs.
- Output artifacts.
- Start conditions.
- Completion criteria.
- Downstream receiver.

Keep responsibilities distinct. For example, avoid asking the same agent to
write a requirement and independently approve its own requirement handoff.

### Step 4: Define the gates

Make each gate testable. A gate should say what evidence the receiving stage
must find, not just “review for quality.”

Include:

- Required contents.
- Accepted statuses.
- Who validates the gate.
- Rework target when it fails.

### Step 5: Define statuses

Keep artifact workflow status distinct from domain-specific decisions and test
results. For example:

```yaml
status: pass
decision_status: proposed
```

Do not reuse one field to mean “complete,” “approved,” and “test passed.”

### Step 6: Define handoff storage and revisions

Choose:

- Run-folder naming.
- Artifact names.
- Attempt/revision locations.
- How the current artifact is identified.
- How the audit log is preserved.

Ensure every agent definition uses the same paths as the pipeline rule.

### Step 7: Set approval and escalation policy

Record:

- Who may accept an ADR or release.
- Whether approval is needed before implementation or production.
- Retry count.
- What happens after retries are exhausted.
- Who receives blocked work.

Do not leave an agent to infer these policies.

### Step 8: Add examples that demonstrate reasoning

An example should contain real scenario values and show:

- A specific decision.
- At least two credible rejected alternatives.
- Rejection reasons tied to stated constraints.
- The selected option’s costs.
- Concrete downstream consequences.
- Which claims are evidence, estimates, assumptions, or unknowns.

An example is instructional data. Label hypothetical figures clearly so they
cannot be mistaken for project facts.

### Step 9: Cross-check the files

Before using the pipeline, verify that:

- Agent names and filenames match.
- Stage numbers match across files.
- Every required artifact has the same name and path everywhere.
- Required inputs can be produced by the preceding stage.
- Status values are consistent.
- Rework goes to the responsible stage.
- The last stage’s authority is explicit.
- Approval policy is stated once in a canonical location and referenced
  consistently.

---

## 14. PayFlow-specific concepts to generalize carefully

The process structure is reusable; many PayFlow rules are not.

| PayFlow concept | Generalize as |
|---|---|
| Ledger is the only writer of money | Define the authoritative owner for critical data or state |
| Integer minor units | Define the project’s exact representation for critical numeric data |
| Idempotency keys | Define retry safety for externally visible or irreversible operations |
| No raw PAN/CVV | Define prohibited sensitive data and permitted tokenized references |
| Kafka events | Define the project’s event or messaging contracts |
| PostgreSQL per service | Define data ownership and database access boundaries |
| PCI-DSS SAQ-A | Replace with the new project’s actual compliance scope |
| Merchant webhooks | Use a real subsystem and workload from the new project |

Do not copy a constraint just because it appears in a sample. A project-specific
constraint must come from an authorized requirement, project convention,
verified regulation, or documented decision.

---

## 15. Consistency checks for the current PayFlow files

Before running the pipeline, verify these points in the repository:

1. `AGENTS.md` says the agent and rule files are in `artifacts/`; links should
   point to paths that actually exist.
2. Every agent’s artifact output path follows the attempt-directory structure
   in `pipeline-rule.md`.
3. The BA, Architecture Review, Developer, and Tester use the same stage IDs
   as the pipeline audit examples.
4. The BA routes to Architecture Review only when
   `architecture_review: required`; otherwise it routes to Developer.
5. The Architecture Review Agent treats missing decision-critical constraints
   as blockers rather than inventing values.
6. The Developer follows the ADR approval policy in `AGENTS.md`.
7. The Tester is explicitly the final Gate 3 decision-maker, as defined by the
   pipeline rule.
8. Example data is identified as hypothetical and is not presented as actual
   test evidence.

These checks matter because correct individual files can still produce an
inconsistent process when they disagree about artifact paths, stage numbers,
statuses, or approval authority.

---

## 16. Final mental model

Think of the model as four connected questions:

1. **What does the business need?**  
   BA Agent → `requirements.md`

2. **Which design should satisfy the need and constraints?**  
   Architecture Review Agent when required → `architecture-review.md`

3. **Was the selected scope implemented, and what build evidence exists?**  
   Developer Agent → code and `build-report.md`

4. **Does the implementation meet every requirement?**  
   Tester Agent → `test-report.md`

The rules connect the answers, the artifacts preserve them, the gates prevent
unsupported handoffs, and the audit log records how the run progressed.

The process is useful when it makes decisions and evidence clearer. If a stage
or artifact adds no distinct responsibility or decision, simplify the process
rather than keeping it merely for the appearance of rigor.
```


# MCP

```text
MCP (Model Context Protocol) — A standard protocol for exposing tools, resources, and prompts to an agent.
Tool — An action, usually with side effects, that an agent may invoke. Defined by a name, a description, an input schema, and an output schema.
Resource — Readable data an MCP server exposes by URI. No side effects.
Prompt (MCP primitive) — A reusable prompt template offered by a server to a client.
Tool contract — The full agreement: input schema, output schema, error cases, and the description text that governs when it's called.
Idempotency key — A caller-supplied identifier that makes a retry safe by returning the original result rather than repeating the effect.
Security boundary — The set of actions a tool must never perform, regardless of the request.
isError vs protocol error — MCP distinguishes a tool that failed (a result flagged as an error, which the agent sees and can reason about) from a protocol failure (a JSON-RPC error, which is a transport-level problem). Recoverable tool failures belong in the first category.

```