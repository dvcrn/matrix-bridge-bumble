package ai

import "fmt"

const welcomePrompt = `
You will be given a Bumble profile to analyze. Your task is to extract key points about the person and generate potential conversation starters based on their profile.

Here is the Bumble profile:
<profile>
%s
</profile>

For reference, here is MY bumble profile. DO NOT ANALYZE THIS
<profile>
%s
</profile>

Carefully read through the profile and follow these steps:

1. Write out their bio and prompts

2. Identify key points about the person and provide a data sheet. These should be short, concise bullet points that highlight anything interesting or unique about them. This could include their interests, hobbies, occupation, values, or any STANDOUT information. Standout meaning, it's something unique or interesting about hat person. 
Things like favorite food or random replies to prompts are not standout. Keep it on interesting points and not things like "star sign" or their education, but keep what they're looking for, drinking/smoking habits, age, drugs, height and plans for children.

Point out interesting things, standout points or unique points about this profile, cut out fluff that is clearly not relevant to be repeated and just lighthearted. 

Point out any warning signs or red flags if there are any (and explanation why it's a warning/red flag), and a subjective take on other things that I should look out for with this person.

Provide a 1-2 sentence short but straight to the point without BS subjective take on this person. 

Provide a non-BS take on whether I would get along with this person.

Provide a quick statement how to talk to this person, for example, whether directly asking them out is a good idea, or if they wouldn't like that, or how to get a reply from them. No bullshit or fluff, keep it real and concise. 


3. Generate 5 SHORT conversation starters or questions that could get the conversation with this person going
- The first message that is sent is the most important, it decides whether the other person will answer or not, so ask questions
- Use information in the profile as context but DO NOT FOCUS OVERLY ON IT. If coffee is mentioned, don't make everything about coffee. Focus on their BIO/About me mostly.
- max 1-2 SHORT sentences
- if the prompts have an interesting question, or if the profile asks a question, try to reply to it
- if there are matches between my profile and theirs, try to reference that
- If the profile barely has anything in it, don't reference much from it
- Don't ask questions about things like their favorite food, favorite spots or where they're from, noone cares about that. 
- If they want to skip small talk and just meet, directly suggest to meet with a few ideas
- don't overuse emojis, add them when it makes sense
- KEEP IT SHORT!! FOLLOW THE EXAMPLES AS CLOSE AS YOU CAN IN WRITING STYLE!!
- use casual puncutation and capitalization 

try some of the following if it applies: 
- a dumb wordplay joke that plays on their profile
- a quick snappy message "aggressive roasting? I'm so down"
- a cheesy pickup line
- a quick invitation to a drink/coffee/date 
- a open question, like 'would you rather', but more strange questions
- ask a weird unique open ended question


No go's (NEVER DO THIS):
- Talking about food
- Asking about favorite spots
- Asking about craziest things they have done
- long messages, keep it as short as possible
- use lingo like "challenge accepted"
- use language like "showdown this week" or "who would win in a battle of wits?" 


Good examples (FOLLOW THIS STYLE, INCLUDING WRITING AND NUANCE): 
- man I'd be so down for a hot chocolate date
- you're into roasting? you go first, then my turn 👀
- ok tell me your first impression of my profile and I'll tell you mine.
- you're my personal cheerleader? I could really use one right now lol how do I sign up? 
- you like analyzing people? try me lol what stood out to you? 
- so, would you rather always be a little thirsty, or always a little hungry? 
- hi, I want to apply for the position as [xx] you mentioned, where do i send my resume in?? 
- can I hire you as my [xxx]? I think I could need one 😆
- soo how about checking out xxx this week sometime if you're free? 
- sooo you wanna go running together sometime then? 😂
- I pick the place, you pick what we'll do? Go to a park. 
- What's the worst Tinder opener you've gotten so far?
- Are you one of the pioneers of flight? You seem just Wright for me
- Your eyes are like Ikea. I'm totally lost in them
- I'm researching important dates in history. Do you want to be mine

Format your response as follows:

**Bio**:
Their bio

Prompts: 

- [Their 1st prompt]
- [Their 2nd prompt]
- [Their 3rd prompt]

**Key Points**:
- [First key point]
- [Second key point]
- [Third key point]
... (continue for all relevant points)

Unique Points & Warning Signs
... (continue with standout + unique points, if any)
... (continue with red flags & warning signs, if any)

Subjective take: 

Fit: 

Tips for dealing with this person:
- How to get a reply from them? 
- Ask out directly? Yes/no, why: 


**Conversation Starters**:

1. [First potential question or conversation starter]
2. [Second potential question or conversation starter]
3. [Third potential question or conversation starter]

Ensure that your key points are concise and relevant, and that your conversation starters are engaging and tailored to the specific details in the profile. Avoid generic statements or questions that could apply to anyone.

No preamble, just print out the answer. 

Output everything in valid markdown format
`

func createWelcomePrompt(ownUserProfile, userProfile string) string {
	return fmt.Sprintf(welcomePrompt, userProfile, ownUserProfile)
}
