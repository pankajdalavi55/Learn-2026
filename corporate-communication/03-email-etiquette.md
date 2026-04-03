# Email Etiquette

## The Email Essentials

### When to Use Email

✅ **Good for:**
- Formal requests or announcements
- Documentation that needs to be searchable
- External stakeholder communication
- Detailed information with attachments
- Non-urgent matters
- Communication across time zones

❌ **Not good for:**
- Urgent issues (use Slack/call)
- Complex back-and-forth discussions (use meeting)
- Sensitive/emotional topics (use video call)
- Quick questions (use chat)

## Subject Lines

The subject line determines if your email gets read.

### Best Practices

✅ **Good Subject Lines:**
- `[Action Required] API Deprecation - Update by March 1`
- `Q1 Performance Review Schedule`
- `Question: Database Migration Timeline`
- `[FYI] Deploy Window This Friday 6-8 PM`

❌ **Bad Subject Lines:**
- `Quick question`
- `Hi`
- `Follow up`
- `Important!!!`

### Prefixes to Use

- **[Action Required]** - Needs recipient to do something
- **[FYI]** - Information only, no action needed
- **[Urgent]** - Genuinely urgent matters only
- **[Decision Needed]** - Requires their input
- **Question:** - You're asking something

## Email Structure

### The 30-Second Rule
Your email should be scannable in 30 seconds or less.

### Effective Template

```
Subject: [Action Required] Review Q1 Planning Doc by Friday

Hi [Name],

[1-2 sentence context/purpose]
I've prepared the Q1 planning document for our team's initiatives.

[Key information or request]
Please review by Friday and provide feedback on:
1. Prioritization of features
2. Resource allocation
3. Timeline feasibility

[Link or attachment]
Document: [link]

[Clear next steps]
I'll finalize the plan after incorporating your feedback and share 
with leadership next Monday.

[Closing]
Thanks,
[Your Name]
```

### Key Components

1. **Clear purpose** in first sentence
2. **Specific request** or information
3. **Action items** if any
4. **Deadline** if applicable
5. **Relevant links** or attachments

## Writing Style

### Be Clear and Concise

✅ **Good:**
```
The deployment is scheduled for Friday at 6 PM. 
We expect 15 minutes of downtime.
```

❌ **Bad:**
```
So, I was thinking that maybe we should probably deploy sometime 
this week, and it might be best if we did it on Friday evening 
because there's usually less traffic then, though I'm not entirely 
sure about the exact time...
```

### Use Formatting for Clarity

**Techniques:**
- **Bullet points** for lists
- **Bold** for key information
- **Numbers** for sequential steps
- **Short paragraphs** (2-3 sentences max)
- **Headings** for long emails

**Example:**
```
**What changed:** API v1 will be deprecated

**Impact:** 
- Users must migrate to API v2
- Old endpoints stop working on March 1

**Action needed:**
1. Review migration guide: [link]
2. Update your integration
3. Test in staging
4. Confirm completion by Feb 15

**Support:** Email api-support@company.com with questions
```

## Professional Tone

### Greetings

**Formal:**
- Dear Mr./Ms. [Last Name]
- Hello [First Name]

**Semi-formal (most common):**
- Hi [First Name]
- Hello team

**Informal (close colleagues):**
- Hey [Name]
- [Name]

### Closings

**Formal:**
- Best regards
- Sincerely
- Respectfully

**Semi-formal:**
- Best
- Thanks
- Regards

**Informal:**
- Cheers
- Talk soon

### Watch Your Tone

Remember: Email lacks tone. Be extra careful.

✅ **Good:**
```
I noticed the deadline was moved up. Could we discuss the 
feasibility? I want to make sure we deliver quality work.
```

❌ **Risky (might sound harsh):**
```
This deadline is impossible. We can't do this.
```

## Common Scenarios

### 1. Requesting Something

```
Subject: Request: SSH Access to Production Server

Hi [Manager Name],

I need SSH access to the production server to debug the 
memory leak issue reported this morning.

Access needed:
- Server: prod-app-01
- Duration: Today only
- Reason: Memory profiling for incident #1234

I've completed the security training (cert #5678) and will 
follow the runbook for production access.

Could you approve this request? The issue is impacting 10% 
of users.

Thanks,
[Your Name]
```

### 2. Following Up

**Wait 2-3 business days before following up**

```
Subject: Re: Review Q1 Planning Doc by Friday

Hi [Name],

Following up on my email from Tuesday about the Q1 planning doc.

The deadline is tomorrow (Friday). Could you let me know if 
you need more time or have any questions about the document?

Happy to hop on a call if that's easier.

Thanks,
[Your Name]
```

### 3. Delivering Bad News

```
Subject: Delay: User Dashboard Feature

Hi [Product Manager],

I need to inform you that the user dashboard feature will be 
delayed by one week.

**Reason:** 
We discovered critical performance issues during load testing. 
The page takes 8 seconds to load with 1000+ users, which fails 
our <2s standard.

**Solution in progress:**
- Implementing database query optimization (2 days)
- Adding caching layer (2 days)
- Re-testing (1 day)

**New delivery date:** March 15 (was March 8)

I should have caught this earlier in development. I've added 
performance testing to our standard checklist.

Let me know if you need to discuss the impact with stakeholders.

Best,
[Your Name]
```

### 4. Saying No Professionally

```
Subject: Re: Additional Feature Request for Sprint 23

Hi [Name],

Thanks for the feature request. I understand why this would be 
valuable for users.

However, I don't think we can include it in Sprint 23 because:
1. We're at capacity with committed items
2. This feature needs design input (2-3 days)
3. It risks our deadline for the higher-priority checkout fix

**Alternative options:**
- Add it to Sprint 24 backlog (my recommendation)
- Reduce scope of current feature X to make room
- Assign to another engineer if available

What works best from your perspective?

Best,
[Your Name]
```

### 5. Asking for Clarification

```
Subject: Question: Requirements for Mobile App Feature

Hi [Product Manager],

I'm starting work on the mobile app notification feature and 
need clarification on a few points:

1. Should notifications work when the app is closed, or only 
   when it's in the background?
2. Are we supporting both iOS and Android in this phase?
3. What's the priority: delivery speed or rich content 
   (images, actions)?

The answers will significantly impact the technical approach 
and timeline.

Could we briefly discuss this, or would you prefer to respond 
via email?

Thanks,
[Your Name]
```

### 6. Requesting Meeting

```
Subject: Request: 30-Minute Sync on API Redesign

Hi Sarah,

I'd like to schedule 30 minutes to discuss the API redesign 
approach before I start implementation.

Specific topics:
- Authentication strategy (OAuth vs JWT)
- Versioning approach (URL vs header)
- Backward compatibility requirements

I've drafted a proposal here: [link]

Are you available this week? I'm free:
- Tuesday 2-4 PM
- Wednesday 10 AM - 12 PM  
- Thursday after 2 PM

Let me know what works!

Thanks,
Alex
```

### 7. Thanking Someone

```
Subject: Thank You - Code Review Help

Hi Jordan,

Quick note to say thanks for the thorough code review on my 
PR yesterday. Your suggestions about the caching strategy were 
spot-on - implementing them improved performance by 40%.

I really appreciate you taking the time to explain the reasoning 
behind each suggestion. Learned a lot!

Thanks again,
Alex
```

### 8. Announcing Completion

```
Subject: ✅ User Dashboard Feature - Shipped to Production

Hi team,

Good news - the user dashboard feature is now live in production!

**What shipped:**
- Real-time activity feed
- Customizable widgets
- Export to CSV functionality
- Mobile-responsive design

**Metrics:**
- Page load time: 1.2s (under our 2s target)
- Test coverage: 87%
- Zero critical bugs in QA

**Next steps:**
- Monitoring performance for 48 hours
- Gathering user feedback
- Planning phase 2 enhancements

Thanks to Sarah (design), Jordan (backend), and Morgan (QA) 
for their collaboration!

Dashboard: [link]
Documentation: [link]

Let me know if you see any issues.

Cheers,
Alex
```

### 9. Escalating an Issue

```
Subject: [Urgent] Production Database at 95% Capacity

Hi [Manager],

We have a critical issue that needs immediate attention.

**Problem:** Production database is at 95% capacity and growing 
5% per day.

**Impact:** 
- Will hit 100% in ~24 hours
- Will cause service outage when full
- Affects all customers

**Immediate actions I've taken:**
- Archived old logs (freed 3%)
- Identified largest tables
- Contacted AWS support (ticket #12345)

**What I need:**
- Approval to scale database instance ($500/month increase)
- OR decision to delete data older than 2 years
- Decision needed within 4 hours

**Recommendation:** Scale the instance now, plan data archival 
strategy for next week.

I'm available to discuss immediately.

Alex
[Phone number]
```

## Email Etiquette Rules

### DO:

✅ Reply within 24 hours (even if just to say "I'll get back to you by X")
✅ Use CC and BCC appropriately
✅ Proofread before sending
✅ Use descriptive attachment names
✅ Include context when forwarding
✅ Use "Reply All" when everyone needs to know
✅ Double-check recipients before sending
✅ Put the ask/action at the top for busy people

### DON'T:

❌ Reply all to company-wide emails unnecessarily
❌ Use all caps (LOOKS LIKE SHOUTING)
❌ Use excessive exclamation marks!!!
❌ Send huge attachments (use cloud links)
❌ Write novels (keep it under 200 words when possible)
❌ Use emojis in formal communication
❌ Forward without permission if it contains sensitive info
❌ Mark everything as urgent/high priority

## CC vs BCC vs TO

### TO:
People who need to take action or are primary recipients

### CC (Carbon Copy):
People who should be informed but don't need to act
- Your manager (for visibility)
- Team members (for awareness)

### BCC (Blind Carbon Copy):
- Large mailing lists (to hide email addresses)
- When you don't want recipients to see each other
- **Don't use to secretly include people** - it's unethical

## Mobile Email

### Keep It Even Shorter

People read email on phones. Make it skimmable.

✅ **Mobile-friendly:**
```
Subject: [Action Needed] Approve Deploy by 2 PM

Hi Sarah,

Need approval to deploy hotfix for payment bug.

Changes: [link to PR #345]

Risk: Low (config change only)

Approve? Reply with 👍

Thanks,
Alex
```

## Email Mistakes to Avoid

### 1. The "Reply All" Disaster
Responding to 500 people when you meant to reply to one.

**Prevention:** Always double-check recipients.

### 2. The Emotional Email
Sending an angry email in the heat of the moment.

**Prevention:** Save as draft. Review in 1 hour. Usually delete it.

### 3. The Missing Attachment
"Please see attached" but you forgot to attach.

**Prevention:** Attach files FIRST, then write email.

### 4. The Wrong Recipient
Sending to the wrong John/Sarah/Alex.

**Prevention:** Verify email addresses, especially for sensitive content.

### 5. The Vague Ask
Unclear what you want the recipient to do.

**Prevention:** Always include clear action items or questions.

## Advanced Tips

### 1. The TL;DR Technique

For long emails, add a summary at the top:

```
**TL;DR:** Server migration scheduled for Saturday 2-6 AM. 
No action needed from you, just expect 4 hours downtime.

[Detailed information follows...]
```

### 2. Boomerang/Schedule Send

Use scheduled sending for:
- Emails written at night (send at 8 AM instead)
- Following up automatically if no response
- Respecting recipient's timezone

### 3. Email Templates

Create templates for common scenarios:
- Status updates
- Meeting follow-ups
- Pull request notifications
- Deployment announcements

### 4. Signature Best Practices

```
[Your Name]
[Your Title]
[Company]
[Email] | [Work Phone]
[LinkedIn or GitHub - optional]

[Company tagline or legal notice if required]
```

Keep it simple. Don't add inspirational quotes.

## Cross-Cultural Considerations

### Global Teams

- **Be explicit:** Some cultures are more direct than others
- **Avoid idioms:** "Touch base" might confuse non-native speakers
- **Mind time zones:** Don't expect immediate responses
- **Be patient:** Writing in a second language takes time
- **Use simple language:** Clear beats clever

### Formality Varies

- Some cultures prefer formal titles (Dr., Mr., Ms.)
- Others prefer first names immediately
- When in doubt, match their style

## Quick Reference

| Situation | Response Time | Formality | Length |
|-----------|---------------|-----------|--------|
| Manager request | Same day | Semi-formal | Brief |
| Peer question | 1 business day | Informal | Brief |
| External stakeholder | Same day | Formal | Medium |
| Bug report | 2-4 hours | Semi-formal | Detailed |
| Team update | As scheduled | Semi-formal | Medium |

---

**Remember:** Email is permanent. Write every email as if it might be forwarded to the entire company (or shown in court). Be professional, clear, and kind.
