---
# autoconfig.yml will use these settings. config.yml will override.
Type: Checkpoint
UID: 2df27a2d-6db3-4dd6-aa82-2e6169e5b77f
# DefaultVisibility: hidden # Uncomment this line to default Checkpoint to hidden
MaxCheckpointSubmissions: 1
# EmailOnCompletion: true #  Uncomment this line to send instructors an email once a student has completed a checkpoint
TimeLimit: 60
# Autoscore: true # Uncomment this line to finalize checkpoint scores without instructor review
---

# Checkpoint Example

Check the top of this file to see how things like the time allowed and maximum number of submissions can be configured. You do not require a maximum number of submissions or a time limit if you so choose.

Saving and Exiting will not pause the timer, but will allow you to return to the test if you wish.

The checkpoint will be graded if the time elapses without submitting, even if you close their browser.

<!--BEGIN CHALLENGE-->

### !challenge

* type: multiple-choice
* id: 06f68dd8-0819-4504-8d53-96246dc2d83f
* title: How did it go?
* points: 3
<!--Other optional fields (checkpoints only) -->
<!--`points: 1`: the number of points for scoring as a checkpoint-->
<!--`topics: python, pandas`: the topics for analyzing points-->

##### !question

1. How well do you know the material from this lesson?

##### !end-question

##### !options

* I got it!
* More practice please!
* I need some help!

##### !end-options

##### !answer

*

##### !end-answer

### !end-challenge

<!--END CHALLENGE-->

<!--BEGIN CHALLENGE-->

### !challenge

* type: checkbox
* id: b73598d5-4c11-4900-a711-cc30a9d57f21
* title: JSX components
* points: 10
<!--Other optional fields (checkpoints only) -->
<!--`points: 1`: the number of points for scoring as a checkpoint-->
<!--`topics: python, pandas`: the topics for analyzing points-->

##### !question

What are the three primary parts of a JSX component?

##### !end-question

##### !options

* Props
* Globals
* Lifecycle & Business Logic
* Rendering
* Objects
* Compilers

##### !end-options

##### !answer

* Props
* Lifecycle & Business Logic
* Rendering

##### !end-answer

### !end-challenge

<!--END CHALLENGE-->

<!--BEGIN CHALLENGE-->

### !challenge

* type: short-answer
* id: bc7be392-0529-4846-a7c7-09c356297fea
* title: Method triggers
<!--Other optional fields (checkpoints only) -->
<!--`points: 1`: the number of points for scoring as a checkpoint-->
<!--`topics: python, pandas`: the topics for analyzing points-->

##### !question

What method triggers a render of your component and updates its state?

##### !end-question

##### !placeholder
your answer
##### !end-placeholder

##### !answer

/(r|R)ender/

##### !end-answer

### !end-challenge

<!--END CHALLENGE-->

<!--BEGIN CHALLENGE-->

### !challenge

* type: paragraph
* id: b83e2e0c-02d1-4a50-b72e-247331249191
* title: Using class components
<!--Other optional fields (checkpoints only) -->
<!--`points: 1`: the number of points for scoring as a checkpoint-->
<!--`topics: python, pandas`: the topics for analyzing points-->

##### !question

When should you use a class component vs a functional component?

##### !end-question

##### !placeholder

your answer

##### !end-placeholder

### !end-challenge

<!--END CHALLENGE-->
