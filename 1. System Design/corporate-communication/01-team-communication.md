# Team Communication

## Core Principles

### 1. **Transparency**
Share information openly with your team. Hiding problems only makes them worse.

**Good:** "I'm stuck on the authentication module. I've tried X and Y approaches. Can someone pair with me?"
**Bad:** *Silently struggling for days before missing the deadline*

### 2. **Respect Everyone's Time**
- Keep messages concise and purposeful
- Use threads to keep conversations organized
- Don't @channel unless truly urgent
- Record important decisions for async team members

### 3. **Assume Positive Intent**
When reading messages, assume your teammate means well. Text lacks tone.

## Daily Standups / Check-ins

### Effective Standup Updates

✅ **Good Format:**
```
Yesterday: Completed user authentication API, merged PR #234
Today: Working on password reset flow
Blockers: Need design specs for the email template
```

❌ **Poor Format:**
```
Working on stuff. No blockers.
```

### Key Elements
- **Be specific:** "Finished X" not "Made progress"
- **Be honest about blockers:** Don't hide issues
- **Keep it brief:** 1-2 minutes max

## Code Reviews

### Giving Feedback

**Principles:**
- Focus on the code, not the person
- Explain the "why" behind suggestions
- Use questions instead of commands
- Praise good solutions

**Examples:**

✅ **Good:**
```
Could we use a Map here instead of an array? It would give us O(1) lookup 
instead of O(n), which matters since this runs in a hot path.

Nice use of the factory pattern here! Makes this much more testable.
```

❌ **Bad:**
```
This is wrong.
Why would you do it this way?
```

### Receiving Feedback

- Don't take it personally
- Ask clarifying questions
- Thank reviewers for their time
- Push back respectfully if you disagree with reasoning

**Example Response:**
```
Good catch on the memory leak! Fixed in the latest commit.

Regarding the interface suggestion - I kept it as a class because we need 
the constructor logic. Open to alternatives if you have ideas.
```

**More Review Examples:**

✅ **Constructive feedback:**
```
I noticed this function is 150 lines. Could we break it into smaller 
functions? Something like:
- validateInput()
- processData()
- formatResponse()

This would make it easier to test and understand. Happy to pair on 
the refactor if helpful!
```

✅ **Suggesting alternatives:**
```
This works, but have you considered using the repository pattern here? 
It would make mocking easier in tests. Check out how we did it in 
UserService.ts for reference.
```

✅ **Asking questions:**
```
I'm trying to understand the logic here - what happens if the API 
returns null? Should we add a null check?
```

✅ **Praising good work:**
```
Really clean implementation! I especially like how you handled the 
edge cases in lines 45-52. This is exactly the pattern we should use.
```

## Asking for Help

### The XY Problem
Don't just ask about your attempted solution (Y), explain your actual goal (X).

❌ **Bad:** "How do I convert a string to uppercase in Python?"
✅ **Good:** "I need to normalize usernames for case-insensitive comparison. Should I uppercase, lowercase, or something else?"

### Effective Help Requests

**Template:**
```
**Goal:** [What you're trying to achieve]
**What I've tried:**
1. [Attempt 1] - [Result/Error]
2. [Attempt 2] - [Result/Error]

**Relevant code/logs:** [Link or snippet]
**Question:** [Specific question]
```

**Example:**
```
**Goal:** Get the user service to authenticate via OAuth

**What I've tried:**
1. Using the `oauth2` library - getting "invalid_grant" error
2. Checked token expiry - it's valid for another hour
3. Verified redirect URI matches exactly

**Error logs:** https://gist.github.com/...

**Question:** Has anyone seen this error with our OAuth provider? 
Am I missing a required scope?
```

**More Help Request Examples:**

**Example 1 - Database Issue:**
```
**Goal:** Query users by multiple filter criteria efficiently

**What I've tried:**
1. Using WHERE clauses with AND - works but slow (3 seconds)
2. Creating a compound index - helped but still 1.5 seconds
3. Looked at query execution plan - seeing full table scan

**Current query:**
SELECT * FROM users 
WHERE status = 'active' AND created_at > '2025-01-01' AND region = 'US'

**Question:** Is there a better indexing strategy? Should I denormalize 
the data? We have 2M users and this query runs 1000x/day.

**Urgency:** Medium - not blocking but impacting user experience
```

**Example 2 - Architecture Decision:**
```
**Context:** Building a notification service that needs to send emails, 
push notifications, and SMS.

**Options I'm considering:**
1. Pub/Sub with separate workers for each channel
2. Queue-based with a single worker that routes
3. Direct API calls (simplest but not scalable)

**Constraints:**
- Need to handle 10k notifications/hour
- Must track delivery status
- Should be resilient to service outages

**Question:** Which pattern have you used for similar use cases? Any 
pitfalls I should watch out for?
```

**Example 3 - Tool/Library Help:**
```
**Goal:** Add authentication to our GraphQL API

**Research done:**
- Looked at Apollo Server docs
- Found 3 different approaches (context, directives, middleware)
- Not sure which is best practice for our use case

**Question:** @backend-team - how did you implement auth in the REST API? 
Should I follow the same pattern or is GraphQL different?

Willing to do a spike and present findings if no one has experience here.
```

## Sharing Knowledge

### When You Learn Something
- Document it in team wiki/confluence
- Share in team chat with context
- Offer to do a quick demo/session

**Example Message:**
```
TIL: You can use `git commit --fixup <hash>` to mark commits for auto-squashing.
Wrote a quick guide: [link]

This will save us time during PR cleanup!
```

## Handling Disagreements

### Disagree and Commit
It's okay to disagree, but once a decision is made, commit to it.

**Steps:**
1. **State your perspective clearly** with reasoning
2. **Listen** to others' views genuinely
3. **Propose alternatives** or ask questions
4. **Accept the decision** once made
5. **Support the direction** publicly

**Example:**
```
I still think approach A is more maintainable, but I understand the 
performance concerns. Let's go with approach B and document the 
tradeoffs. I'll implement it properly.
```

## Remote/Async Communication

### Best Practices
- **Overcommunicate:** Can't read body language
- **Use video** when discussing complex topics
- **Document decisions** in writing
- **Respect timezones:** Use async-friendly processes
- **Use emoji/GIFs** appropriately to add warmth

### Async-First Mindset
Structure communication so people can catch up later:

✅ **Good:**
```
Thread: Q2 Planning Discussion
├─ Here's the proposal: [link to doc]
├─ Key decision points: [list]
├─ Please review by Friday
└─ We'll finalize in Monday's meeting
```

## Celebrating Wins

Don't forget to celebrate! It builds team morale.

**Examples:**
- "Great job shipping the payment feature @alex!"
- "Team crushed this sprint - we cleared 38 story points!"
- "Shoutout to @sara for the thorough incident postmortem"

## Red Flags to Avoid

🚩 **Passive-aggressive comments**
🚩 **Public criticism of teammates**
🚩 **Taking credit for others' work**
🚩 **Ignoring team norms (like PR review times)**
🚩 **"I told you so" when something fails**
🚩 **Ghosting conversations you started**

## Quick Tips

1. **Reply with ETA if you can't answer immediately**
2. **Use threads** to keep channels organized
3. **Search before asking** - might be documented
4. **Tag relevant people** but don't over-tag
5. **Use code blocks** for code/logs
6. **Summarize long discussions** for clarity
7. **Follow up** on conversations you started
8. **Say thanks** - appreciation matters

---

**Remember:** Your team is your first support system. Invest in these relationships through clear, kind, and timely communication.
