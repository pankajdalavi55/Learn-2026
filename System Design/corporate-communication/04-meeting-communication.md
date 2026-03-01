# Meeting Communication

## The Meeting Problem

"This meeting could have been an email" - Every engineer ever.

Meetings are expensive. A 1-hour meeting with 10 people costs 10 hours of productivity.

## When to Have a Meeting

### ✅ Good Reasons for a Meeting

- **Complex decisions** requiring discussion and debate
- **Brainstorming** where collaboration sparks ideas
- **Building relationships** and team cohesion
- **Resolving conflicts** that need real-time dialogue
- **Sensitive topics** requiring nuance and empathy
- **Alignment** across multiple stakeholders
- **Urgent issues** requiring immediate coordination

### ❌ Bad Reasons for a Meeting

- Status updates (use email/doc)
- Information that could be written down
- Decisions already made (just announce them)
- Topics that only affect 2 people (chat 1-on-1)
- "We always meet on Mondays"

### The Email Test

Before scheduling a meeting, ask:
> Could this be handled with:
> - An email?
> - A shared document?
> - A Slack thread?
> - A quick 1-on-1?

If yes → Don't meet.

## Before the Meeting

### 1. Create a Clear Agenda

**Bad agenda:**
```
Weekly Sync
- Updates
- Discussion
```

**Good agenda:**
```
Weekly Planning - Jan 5, 2pm-3pm

1. Quick wins from last week (5 min)
2. Sprint planning for next 2 weeks (30 min)
   - Review backlog priorities
   - DECISION: Which features for Sprint 23?
   - Assign tasks
3. Tech debt discussion (15 min)
   - DECISION: Refactor now or later?
4. Blockers & questions (10 min)

Goal: Leave with clear sprint plan and task assignments
```

### 2. Invite the Right People

**Required:** People who must be there for decisions
**Optional:** People who should know but don't need to attend (send notes)
**FYI:** People who can read notes later

**Don't invite people "just in case"**

### 3. Send Prep Materials

Give people time to prepare:
```
Meeting: Q1 Architecture Review
Date: Friday, Jan 10, 2pm-4pm

Please review before the meeting:
- Current architecture doc: [link]
- Proposed changes: [link]
- Performance data: [link]

We'll make decisions on:
1. Database migration strategy
2. Caching layer implementation
3. API versioning approach

Come prepared with questions and concerns.
```

### 4. Time Box Appropriately

- **15 min:** Quick sync, single decision
- **30 min:** Standard meeting, couple of topics
- **60 min:** Deep discussion, multiple topics
- **2 hr:** Workshop, brainstorming, planning

**Default to 25 or 50 minutes** to give people buffer between meetings.

## During the Meeting

### 1. Start on Time

Don't punish people who arrived on time by waiting for late arrivers.

### 2. Assign Roles

**Facilitator:** Keeps discussion on track
**Note-taker:** Documents decisions and actions
**Timekeeper:** Watches the clock

### 3. Follow the Agenda

**When discussion goes off-track:**
```
"This is important, but off-topic for today. Let's schedule 
a separate discussion and get back to the agenda."
```

**Parking lot technique:**
Keep a list of "important but not now" topics to address later.

### 4. Make Decisions

Meetings should produce outcomes.

**Unclear ending:**
```
"Thanks everyone. Good discussion."
```

**Clear ending:**
```
"To summarize our decisions:
1. We're going with PostgreSQL (not MongoDB)
2. Alex will lead the migration
3. Target completion: March 1

Action items:
- Alex: Create migration plan by Friday
- Sarah: Get budget approval by Tuesday
- Team: Review and comment on plan

Next meeting: Jan 20 to review progress.
```

### 5. Encourage Participation

**For quiet people:**
- "Jordan, you have experience with this. What do you think?"
- "Let's go around the room"

**For dominant voices:**
- "Thanks Alex. Let's hear from others."
- "We've heard your perspective. Other views?"

### 6. Remote Meeting Best Practices

- **Camera on** when possible (builds connection)
- **Mute when not speaking** (reduce noise)
- **Use chat** for questions/links
- **Screen share** for visual topics
- **Record** if people can't attend (with permission)

## Meeting Types

### 1. Standup/Daily Sync

**Format:** 15 minutes, standing or quick Zoom
**Goal:** Coordination, surface blockers

**Structure:**
```
Each person (2 min max):
- Yesterday: What I did
- Today: What I'm doing
- Blockers: What's in my way
```

**Anti-patterns:**
- ❌ Turns into problem-solving session
- ❌ Detailed technical discussions
- ❌ Going over 15 minutes
- ❌ Becoming a status report to manager

**Fix:**
"Let's take that offline and discuss after standup."

### 2. Brainstorming

**Format:** 30-60 minutes, creative environment
**Goal:** Generate lots of ideas

**Rules:**
1. **No criticism** during idea generation
2. **Quantity over quality** initially
3. **Build on others' ideas**
4. **Wild ideas welcome**

**Structure:**
```
1. State the problem clearly (5 min)
2. Individual ideation (5 min - silent, write ideas)
3. Share all ideas (15 min - no debate)
4. Group and discuss (15 min)
5. Vote/prioritize (10 min)
6. Next steps (5 min)
```

### 3. Retrospective

**Format:** 60 minutes, end of sprint/project
**Goal:** Learn and improve

**Structure:**
```
1. What went well? (15 min)
2. What didn't go well? (15 min)
3. What should we change? (20 min)
4. Action items (10 min)
```

**Ground rules:**
- Blame the process, not people
- Focus on actionable improvements
- Everyone participates

### 4. 1-on-1 with Manager

**Format:** 30-60 minutes, weekly or bi-weekly
**Goal:** Feedback, growth, alignment

See [Manager Communication](02-manager-communication.md) for detailed guidance.

### 5. All-Hands

**Format:** 60 minutes, monthly/quarterly
**Goal:** Company updates, transparency

**Good all-hands:**
- Clear structure and agenda
- Mix of updates and Q&A
- Engaging presentation
- Leave time for questions
- Record for those who can't attend

### 6. Architecture Review

**Format:** 60-90 minutes
**Goal:** Validate technical decisions

**Structure:**
```
1. Context & problem (10 min)
2. Proposed solution (20 min)
3. Alternative approaches considered (10 min)
4. Tradeoffs analysis (15 min)
5. Discussion & questions (20 min)
6. Decision or next steps (5 min)
```

**Pre-work:** Share architecture doc 2 days before

### 7. Incident Postmortem

**Format:** 60 minutes
**Goal:** Learn from outages, prevent recurrence

**Structure:**
```
1. Timeline of events (10 min)
2. Root cause analysis (15 min)
3. What worked well in response (10 min)
4. What we should improve (15 min)
5. Action items (10 min)
```

**Rules:**
- No blame - blameless postmortems
- Focus on systems and processes
- Everyone can speak freely

**Example Postmortem Conversation:**
```
Facilitator: "Let's walk through the timeline. Alex, you detected 
the issue first?"

Alex: "Yes, at 2:15 PM I got alerts about API errors spiking to 40%."

Facilitator: "What happened next?"

Alex: "I checked the logs and saw database connection timeouts. 
I pinged the on-call DBA in Slack at 2:17 PM."

Sarah (DBA): "I saw the message at 2:25 PM - I was in another meeting. 
I immediately checked the database and saw we'd hit max connections."

Facilitator: "What was the root cause?"

Sarah: "A new service deployed that morning didn't properly close 
connections. Each request leaked a connection until we hit the limit."

Facilitator: "How did you fix it?"

Sarah: "I restarted the service at 2:30 PM and increased connection 
limits as a temporary fix. Jordan deployed a code fix at 3:15 PM."

Facilitator: "What worked well?"

Alex: "Alerts caught it quickly. Communication was clear once we 
were all online."

Facilitator: "What should we improve?"

Sarah: "I should have seen the Slack ping sooner. Maybe we need 
PagerDuty for critical alerts."

Jordan: "We should have caught the connection leak in testing."

Facilitator: "Action items:
1. Jordan: Add connection leak testing to CI/CD
2. Sarah: Set up PagerDuty for database alerts  
3. Alex: Document runbook for database connection issues

Anything else?"
```

### 8. Design Review Meeting

**Example Design Review Flow:**
```
Presenter: "Today I'm proposing we move from REST to GraphQL for 
our mobile API."

[Shares architecture diagram]

Presenter: "The main benefits are:
- Reduces API calls from 5 to 1 for the home screen
- Eliminates over-fetching - mobile bandwidth savings
- Better tooling and type safety"

Reviewer 1: "What's the learning curve for the mobile team?"

Presenter: "Good question. I've created a sample app and training 
plan. Estimate 2 weeks to get up to speed."

Reviewer 2: "What about caching? We heavily cache REST responses."

Presenter: "GraphQL has Apollo Client cache. I've tested it with 
our use case - it's actually better than what we have now."

Reviewer 3: "Did you consider just optimizing the REST endpoints?"

Presenter: "Yes, slide 5 shows that comparison. We'd need to create 
12 custom endpoints vs 1 GraphQL schema. Maintenance cost is lower 
with GraphQL."

Reviewer 1: "What's the migration strategy?"

Presenter: "Run both in parallel for 6 months. Migrate features 
gradually. No big-bang rewrite."

Facilitator: "Any other concerns? [pause] Okay, let's vote. 
Thumbsup for approve, thumbs down for more discussion needed."

[Team gives thumbs up]

Facilitator: "Approved. Alex, next step is RFC document by Friday?"
```

## Running an Effective Meeting

### As the Organizer

**Before:**
- [ ] Clear agenda sent 24 hours ahead
- [ ] Right people invited
- [ ] Prep materials shared
- [ ] Goals/decisions clearly stated

**During:**
- [ ] Start on time
- [ ] Follow agenda
- [ ] Keep discussion on track
- [ ] Document decisions
- [ ] Capture action items

**After:**
- [ ] Send notes within 24 hours
- [ ] Follow up on action items
- [ ] Schedule follow-ups if needed

### As a Participant

**Prepare:**
- Read materials beforehand
- Prepare questions
- Be on time

**Participate:**
- Contribute meaningfully
- Listen actively
- Stay focused (no laptop multitasking)
- Respect others' time

**Follow through:**
- Complete your action items
- Communicate if you can't
- Read the notes

## Meeting Notes Template

```
# [Meeting Title]
**Date:** Jan 5, 2026
**Attendees:** Alex, Sarah, Jordan, Morgan
**Note-taker:** Alex

## Agenda
1. Sprint planning
2. Tech debt discussion
3. Blockers

## Discussion Summary

### Sprint Planning
We reviewed the backlog and prioritized features for Sprint 23.

Key points:
- User dashboard is highest priority
- Payment retry logic is second
- Mobile app features delayed to Sprint 24

### Tech Debt
Discussed the database migration. Team agreed it's urgent.

## Decisions Made
1. ✅ User dashboard is top priority for Sprint 23
2. ✅ Database migration happens in Sprint 24
3. ✅ Morgan will lead the migration effort

## Action Items
- [ ] Alex: Create user dashboard spec by Friday
- [ ] Morgan: Draft migration plan by next Wednesday
- [ ] Sarah: Review payment retry PR by tomorrow

## Parking Lot (for future discussion)
- Mobile app architecture refactor
- Monitoring improvements

## Next Meeting
Sprint 24 Planning - Jan 20, 2pm
```

## Meeting Etiquette

### DO:

✅ Arrive on time (or 2 minutes early)
✅ Come prepared
✅ Mute when not speaking (remote)
✅ Listen actively
✅ Take notes
✅ Ask questions
✅ Contribute ideas
✅ Stay focused
✅ End on time

### DON'T:

❌ Show up late without notice
❌ Multitask on laptop
❌ Interrupt others
❌ Dominate the conversation
❌ Check your phone
❌ Go on tangents
❌ Rehash old debates
❌ Let meetings run over

## Handling Difficult Situations

### Someone is Dominating

**Politely redirect:**
```
"Thanks for sharing your thoughts, Alex. Let's hear from 
others who haven't spoken yet."
```

### Meeting Going Off-Track

**Refocus:**
```
"This is a great discussion, but we're getting away from 
today's agenda. Can we table this and get back to [topic]?"
```

### No One is Talking

**Prompt participation:**
```
"I'd like to hear everyone's perspective. Jordan, let's 
start with you."
```

### Can't Reach a Decision

**Name it and move forward:**
```
"We're not going to align in the time we have. Let's 
identify what info we need and reconvene on Friday."
```

### Running Over Time

**Options:**
1. **Extend:** "We need 10 more minutes. Anyone have a hard stop?"
2. **Pause:** "Let's stop here and continue tomorrow."
3. **Decide:** "We're out of time. Based on the discussion, I'm deciding we'll go with option A. Objections?"

## Declining Meetings

It's okay to decline if you're not needed.

**Template:**
```
Thanks for including me. Based on the agenda, I don't think 
I can add much value to this discussion. 

Would it work to send me the notes afterward? Happy to provide 
async feedback if needed.
```

**When to decline:**
- You're optional and have more urgent work
- The meeting doesn't align with your role
- A teammate can represent your perspective

## Alternative to Meetings

### 1. Async Standups

```
Daily standup - Post by 10 AM

Yesterday:
- Completed OAuth integration
- Fixed 3 QA bugs

Today:
- Starting session management refactor
- Code review for Sarah

Blockers:
- None
```

### 2. Collaborative Documents

Use Google Docs for:
- RFCs and proposals
- Decision-making (comments)
- Gathering feedback
- Status updates

### 3. Recorded Videos

Record a Loom/video for:
- Demos
- Explanations
- Tutorials
- Updates

People can watch on their schedule.

### 4. Email Threads

Good for:
- Announcements
- Status updates
- Simple questions
- Collecting feedback

## Measuring Meeting Effectiveness

**Ask yourself after each meeting:**

1. ✅ Did we achieve the stated goal?
2. ✅ Were the right people there?
3. ✅ Did we make decisions?
4. ✅ Are action items clear?
5. ✅ Could this have been an email?

**If too many ❌:** Rethink your meeting culture.

## Company-Wide Meeting Practices

### Good Practices to Adopt

- **No-meeting Fridays** (focus time)
- **25/50 minute default** (not 30/60)
- **Meeting-free hours** (e.g., mornings)
- **Default to decline** (optional means optional)
- **Record by default** (for async folks)

---

**Remember:** The best meeting is sometimes no meeting. Respect people's time, have a clear purpose, and make meetings productive. Your teammates will thank you.
