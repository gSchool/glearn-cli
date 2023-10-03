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

Each unit can have only one Checkpoint, and the checkpoint must have challenges. Checkpoints can be configured to score points automatically, grant a limited number of attempts, require a certain time limit to complete, and email the instructors upon completion.

See the header on this file `01-example-unit/05-checkpoint.md` to configure these options.

Saving and Exiting will not pause the timer, but will allow you to return to the test if you wish.

If a time limit is set, the checkpoint will be graded if the time elapses without submitting, even if you close the browser.

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

1. How well do you know the material from this lesson? Check the answer in the markdown; this is how you make all options valid in a multiple choice challenge.

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

Which yaml files can be used when configuring your repository contents for a block?

##### !end-question

##### !options

* `description.yaml`
* `course.yaml`
* `manifest.yaml`
* `config.yaml`
* `yak.yaml`
* `autoconfig.yaml`

##### !end-options

##### !answer

* `description.yaml`
* `config.yaml`
* `autoconfig.yaml`

##### !end-answer

### !end-challenge

<!--END CHALLENGE-->

<!--BEGIN CHALLENGE-->

### !challenge

* type: short-answer
* id: bc7be392-0529-4846-a7c7-09c356297fea
* title: Iterating
<!--Other optional fields (checkpoints only) -->
<!--`points: 1`: the number of points for scoring as a checkpoint-->
<!--`topics: python, pandas`: the topics for analyzing points-->

##### !question

Which `learn` CLI command lets you view content in Learn to iterate quickly before publishing?

Check this short answer to see how regex is used to match answers. You can use [Rubular](https://rubular.com/) to quickly test your own regular expressions.

##### !end-question

##### !placeholder
your answer
##### !end-placeholder

##### !answer

/(preview|learn preview)/

##### !end-answer

### !end-challenge

<!--END CHALLENGE-->

### !challenge

* type: checkbox
* id: fb4e6a97-ee62-4ffa-80a6-860f1654353c
* title: Explanations when incorrect
* points: 3
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

Which of the following are Learn callout colors?  Try and get the answer wrong, see how you are guided to the correct answer, then check the markdown to see how explanations are used to create this interactive challenge.

The `!explanation` can be configured with different variants to supply custom responses to checkbox and multiple-choice challenges. The options are:

* `!explanation-correct` -> The response when the answer is correct.
* `!explanation-incorrect` -> The default incorrect response if no other matches are found.
* `!explanation: <OPTION>` -> The `<OPTION>` here must match one of the challenge options. Use with incorrect selected options.
* `!explanation-not: <OPTION>` -> The `<OPTION>` here must match one of the challenge options. Use with correct selected options.

##### !end-question

##### !options

* info
* false
* warning
* danger
* moon
* star

##### !end-options

##### !answer

* info
* warning
* danger
* star

##### !end-answer

#### !explanation-correct:
That's right! The full list is info, success, warning, danger, secondary, star
#### !end-explanation


#### !explanation: false
False is a boolean value but not a callout.
#### !end-explanation

#### !explanation: moon
While 'star' may be an option, moon is not.
#### !end-explanation

#### !explanation-not: info
Sometimes you want an extra tidbit of slightly tangential details.
#### !end-explanation

#### !explanation-not: warning
What about the color of trees and grass?
#### !end-explanation

#### !explanation-incorrect:
Try the hint.
#### !end-explanation

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
#### !hint
You can see the list when you run `learn md co -o`.
#### !end-hint
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->

### !end-challenge
