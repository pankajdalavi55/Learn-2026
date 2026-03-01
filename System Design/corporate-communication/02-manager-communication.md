# Manager Communication

## Understanding the Relationship

Your manager is not your adversary - they're your advocate, coach, and support system. Their success depends on your success.

### What Managers Care About
1. **Delivery:** Are you completing your commitments?
2. **Proactivity:** Do you identify and solve problems independently?
3. **Growth:** Are you developing your skills?
4. **Team health:** Are you a positive team member?
5. **Communication:** Do they have visibility into your work?

## Regular 1-on-1s

### Preparation is Key

**Before Each 1-on-1:**
- Review your previous notes
- List accomplishments since last meeting
- Identify blockers or concerns
- Prepare questions
- Think about growth goals

### Effective Agenda Structure

```
1. Wins & Accomplishments (5 min)
2. Current Work & Blockers (10 min)
3. Feedback & Development (10 min)
4. Strategic Topics / Career (5 min)
```

### Sample Topics

**Regular Topics:**
- Progress on current projects
- Blockers you need help with
- Feedback on your work
- Team dynamics or concerns
- Process improvements

**Periodic Topics:**
- Career goals and progression
- Learning opportunities
- Project interests
- Work-life balance
- Compensation

### Real 1-on-1 Conversation Examples

**Example 1 - Discussing a Blocker:**
```
You: "I'm stuck on the payment integration. The vendor's API 
documentation is incomplete, and their support is slow to respond."

Manager: "How long have you been blocked?"

You: "Three days. I've tried their sample code, read forums, and 
opened two support tickets. Still can't get the webhook verification 
to work."

Manager: "That's frustrating. Have you considered reaching out to 
their engineering team directly? I know someone there."

You: "I didn't know we had that connection. That would be really helpful."

Manager: "I'll make an intro today. In the meantime, can you work on 
another part of the feature?"

You: "Yes, I can build the UI and mock the API responses."
```

**Example 2 - Seeking Feedback:**
```
You: "I'd like feedback on my technical leadership. I've been trying 
to step up in design reviews and mentoring, but I'm not sure if I'm 
striking the right balance."

Manager: "What specifically are you concerned about?"

You: "I don't want to come across as know-it-all in reviews, but I 
also don't want to hold back valuable input."

Manager: "I've seen your reviews - they're actually really good. You 
ask questions instead of dictating, which is great. One thing: sometimes 
you could be more decisive. Like in last week's architecture review, 
we needed your opinion and you were too diplomatic."

You: "That's helpful. So be more opinionated when my expertise is needed?"

Manager: "Exactly. You're the expert on our frontend - own that."
```

**Example 3 - Career Discussion:**
```
You: "I've been thinking about my career path. I'm not sure if I want 
to go into management or stay on the technical track."

Manager: "What appeals to you about each path?"

You: "I love coding and architecture, but I also enjoy helping teammates 
grow. I just led the intern's project and found that really rewarding."

Manager: "Good news - you don't have to decide right now. At your level, 
you can develop both skill sets. Why don't you try tech leading the next 
major project? It's a hybrid role that involves both technical work and 
people coordination."

You: "That sounds perfect. Which project were you thinking?"

Manager: "The mobile app rewrite. It's high visibility, cross-functional, 
and you'd work with 3 other engineers. Interested?"
```

### What NOT to Do

❌ Cancel frequently
❌ Come unprepared
❌ Only talk about tactical work
❌ Complain without solutions
❌ Hide problems until they're critical

## Status Updates

### Weekly Update Email Template

```
Subject: Weekly Update - [Your Name] - [Date]

## Completed This Week
- [Specific achievement with impact]
- [Feature shipped / Bug fixed with ticket #]
- [Helped team member with X]

## In Progress
- [Task 1] - 60% complete, on track for Friday
- [Task 2] - Blocked on design review (see blockers)

## Coming Up Next Week
- [Planned task 1]
- [Planned task 2]

## Blockers / Needs
- Need design approval for feature X by Wednesday
- Unclear on priority between Task A and Task B

## Metrics (if applicable)
- API latency reduced from 500ms to 200ms
- Test coverage increased to 85%

## Other Notes
- Attending React conference on Thursday
```

### Key Principles

1. **Be specific:** "Reduced load time by 40%" not "Made it faster"
2. **Show impact:** Connect work to business goals
3. **Flag risks early:** Don't surprise your manager
4. **Suggest solutions:** Not just problems

## Asking for Help

### Frame Requests Effectively

✅ **Good:**
```
I'm working on the checkout optimization. I've identified three approaches:

1. Server-side caching - fastest to implement, moderate impact
2. Database query optimization - medium effort, high impact
3. Architecture refactor - high effort, highest impact

Given our Q1 goals, I recommend #2. Do you agree, or should we 
discuss tradeoffs?
```

❌ **Bad:**
```
The checkout is slow. What should I do?
```

### When You're Overwhelmed

**Be honest and specific:**
```
I'm currently assigned to:
- Feature A (8 hours/day) - due Friday
- Bug triage rotation (2 hours/day)
- Feature B (requested yesterday) - unclear deadline

This is 10 hours of work daily. Can we:
1. Delay Feature B until next sprint, or
2. Find someone else for bug rotation this week?

What's your preference?
```

## Delivering Bad News

### The "No Surprises" Rule
Never let your manager be blindsided. Alert them to problems early.

### Structure for Bad News

1. **State the problem clearly**
2. **Explain what happened** (not who to blame)
3. **Present your proposed solution**
4. **Ask for guidance** if needed

**Example:**
```
Heads up - the payment integration is delayed.

What happened: The third-party API had undocumented rate limits that 
we hit during testing. It's been down for 6 hours.

My plan: 
1. Implement exponential backoff (today)
2. Add circuit breaker pattern (tomorrow)
3. Coordinate with vendor on rate limit increase

This pushes our launch from Monday to Wednesday.

Do you want to align with Product, or should I?
```

## Receiving Feedback

### How to Respond

✅ **Good Response:**
```
Manager: "Your code reviews have been too brief lately."

You: "Thanks for flagging that. You're right - I've been rushing them 
due to my feature deadline. I'll block off dedicated review time daily 
going forward. Can you share an example of the depth you'd like to see?"
```

❌ **Bad Response:**
```
"I've been really busy though. Others don't review thoroughly either."
```

### Key Principles
- **Don't get defensive:** Listen fully first
- **Ask for specifics:** "Can you give me an example?"
- **Show commitment:** Explain how you'll improve
- **Follow up:** Report on progress later

## Career Conversations

### Performance Review Prep

**Document throughout the year:**
- Major projects and their impact
- Skills you've developed
- Ways you've helped teammates
- Positive feedback received
- Goals achieved

**Example Self-Assessment:**
```
## Key Accomplishments
1. Led migration to microservices architecture
   - Reduced deployment time by 60%
   - Mentored 3 junior engineers through the process
   
2. Improved code quality metrics
   - Increased test coverage from 60% to 85%
   - Reduced production bugs by 40%

## Growth Areas
- Want to improve system design skills
- Need more experience with ML pipelines
- Working on presentation skills

## Goals for Next Period
- Lead a cross-team initiative
- Mentor a new hire
- Get AWS Solutions Architect certification
```

### Asking for Promotion

**Don't say:**
"I've been here for 2 years, so I should be promoted."

**Do say:**
```
I'd like to discuss progression to Senior Engineer. 

I believe I'm demonstrating the next level by:
1. Leading the payment system redesign autonomously
2. Mentoring 2 junior engineers - both shipped their first features
3. Improved team velocity by introducing better testing practices

What gaps do you see? What can I focus on to get there?
```

## Managing Up

### Understand Their Communication Style

- **Detail-oriented:** Give them data and thorough updates
- **Big-picture:** Lead with impact and high-level summary
- **Hands-on:** Invite them to code reviews or technical discussions
- **Hands-off:** Send regular updates so they don't have to ask

### Make Their Job Easier

1. **Bring solutions, not just problems**
2. **Give them visibility** into your work
3. **Meet commitments** or communicate early if you can't
4. **Support team goals** beyond just your tasks
5. **Share wins** they can report upward

### Asking Good Questions

✅ **Strategic questions:**
- "What's the most important thing I can work on right now?"
- "How does my work connect to the team's quarterly goals?"
- "What would a successful outcome look like for this project?"

❌ **Questions that could be self-solved:**
- "What should I do about this small bug?" (Use your judgment)
- "Is the code review approved?" (Check the PR yourself)

## Red Flags in Manager Relationships

🚩 You're afraid to share bad news
🚩 You feel micromanaged but haven't discussed it
🚩 Your 1-on-1s keep getting cancelled
🚩 You're unclear on expectations
🚩 You haven't discussed growth in 6+ months

**Action:** Have an honest conversation about improving the relationship

## Sample Scenarios

### Scenario 1: Requesting Time Off

❌ **Bad:**
"I'm taking next week off. Out of office is set."

✅ **Good:**
```
I'd like to take next week off (Aug 7-11) for a family trip.

Current status:
- Feature X will be complete by Aug 4
- Sarah can cover the on-call shift
- No critical meetings scheduled

Does this work with the team's calendar?
```

### Scenario 2: Disagreeing with Direction

❌ **Bad:**
"This won't work. We shouldn't do it this way."

✅ **Good:**
```
I have concerns about the proposed approach. 

My worry is [specific technical/business risk].

Could we consider [alternative] which addresses this while still 
meeting the goal of [objective]?

Happy to prototype both approaches if that helps the decision.
```

### Scenario 3: Asking for Development Opportunity

✅ **Good:**
```
I'm interested in developing my backend skills. I noticed the 
API redesign project starting next quarter.

Would it make sense for me to join that project? I could:
- Contribute to the design phase given my product knowledge
- Learn from the backend team
- Help with integration since I know the frontend

What do you think?
```

## Quick Tips

1. **Respect their calendar** - be punctual and prepared
2. **Use their preferred communication channel**
3. **Don't save up problems** - surface them as they arise
4. **Keep them informed** on visible projects
5. **Assume good intent** when receiving feedback
6. **Share credit** - highlight team contributions
7. **Be solution-oriented** - show initiative
8. **Follow through** on commitments

---

**Remember:** Your manager's job is to support your success. Build a relationship based on trust, transparency, and mutual respect.
