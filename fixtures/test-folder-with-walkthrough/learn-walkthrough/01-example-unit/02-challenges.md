# Challenge Examples

Interactive questions called "challenges" can be added to any markdown lesson to check for understanding. These same challenges can also be used to construct practice assignments, quizzes, and assessments.

## Examples of Each Type

* multiple choice (select one)
* checkbox (select multiple)
* short answer
* number
* paragraph
* code snippet--js, py, java, sql (code directly in Learn)
* project, testable project (code locally, submit a repo)



## Multiple Choice

Multiple Choice challenges allow a student to submit a single answer to a multiple-choice question.

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

<!--'type' is required-->
<!--'id' is required, string, must be unique within a branch-->
<!--'title' is required, string, used when displaying results-->

* type: multiple-choice
* id: 4fb73ef6-1492-45c9-b937-fe718cb80fea
* title: Example Multiple Choice Challenge

<!--'question' is required, markdown, the question to be answered-->

### !question

Which of these cities is home to a Galvanize campus?

### !end-question

<!--'options' is required, a bulleted markdown list, the options the student selects the correct answer from-->

### !options

* Chicago
* Denver
* Fort Collins
* Miami
* $$x = {-b \pm \sqrt{b^2-4ac} \over 2a}.$$
* Inline LaTeX \\( z^2 = x^2 + y^2 \\)
* \\( z^2 = x^2 + y^2 \\)
* `public void moveTo(int x, int y, int z)`

### !end-options

<!--'answer' is required, the correct answer, must exist in the options-->

### !answer

Denver

### !end-answer

<!--'explanation' is optional. Shown after the student correctly answers the question.-->

### !end-challenge

<!----------------------END CHALLENGE----------------------------->




## Checkbox

Checkbox challenges allow a student to submit multiple answers to a multiple-choice question.

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

* type: checkbox
* id: a2b0dd37-4d72-490f-ba26-c6e7b08d21a0
* title: Example Checkbox Challenge

##### !question

Mark all of the cities that are home to a Galvanize campus.

##### !end-question

##### !options

* Portland
* Denver
* Los Angeles
* Miami

##### !end-options

##### !answer

* Denver
* Los Angeles

##### !end-answer

### !end-challenge

<!----------------------END CHALLENGE----------------------------->



## Short Answer

Short-answer challenges allow a student to submit a short answer, usually a single word, as the answer to a question. By default, the answer is evaluated as a case-insensitive exact match, but can also be evaluated as a regex.

### !challenge

<!--'type' is required-->
<!--'id' is required, string, must be unique within a single markdown file-->
<!--'title' is required, string, used when displaying results-->

* type: short-answer
* id: 1a3cc34f-ea21-4533-bed4-6da3c8c95a4e
* title: Example Short Answer Challenge

<!--'question' is required, markdown, the question to be answered-->

### !question

What will the following code produce?

```javascript
var myArray = ["Elie", "Janey", "Matt", "Parker", "Tim"];
myArray[3]
```

### !end-question

<!--'placeholder' is optional, the placeholder text in the input box. Removed when the user starts typing-->

#### !placeholder

What does myArray[3] equal?

#### !end-placeholder

<!--'answer' is required, the correct answer to the question.
By default, answers evaluated as case-insensitive exact match. If answer wrapped in '/', evaluated as regex. The answer to this question could have been defined as /\A['"]Parker['"]$/ to allow students to enter single or double quotes. Also supports regex ending in '/i' for case-insensitive regex. To test your regex, you can use http://rubular.com/-->

### !answer
"Parker"
### !end-answer

### !end-challenge

<!----------------------END CHALLENGE----------------------------->



## Number

Number challenges allow a student to submit a number as the answer to a question. The answer is evaluated numerically, and the student can answer with a decimal or a fraction--so that things like 3/10, 30/100, .3, 0.3, 0.300 etc. are all equivalent. There also also an option when creating the challenge to define the precision used when the answer is scored -- see the code below for details.

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

<!--'type' is required-->
<!--'id' is required, string, must be unique within a branch-->
<!--'title' is required, string, used when displaying results-->
<!--'decimal' is optional, the number of decimal places that the answer will be rounded to before evaluating for correctness. If not specified, answer must be exact.-->

* type: number
* id: e88f0c35-7f2d-4990-81c1-a529d5e81dbb
* title: Example Number Challenge
* decimal: 5

<!--'question' is required, markdown, the question to be answered-->

### !question
Suppose a card is drawn from a standard 52 card deck. What's the probability that the card is a queen?
### !end-question

<!--'placeholder' is optional, the placeholder text in the input box. Removed when the user starts typing-->

#### !placeholder
Write your answer as a decimal to 5 places
#### !end-placeholder

<!--'answer' is required, the correct answer. Can be specified as a decimal like 0.07692 or a fraction like 1/13.-->

### !answer
1/13
### !end-answer

### !end-challenge

<!----------------------END CHALLENGE----------------------------->



## Paragraph

Paragraph Challenges allow a student to submit a long free-form text answer to a question, such as a definition of explanation. The answers to these Challenges are not evaluated by Learn, but are available for the instructor to view.

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

<!--'type' is required-->
<!--'id' is required, string, must be unique within a branch-->
<!--'title' is required, string, used when displaying results-->

* type: paragraph
* id: f397a35a-2d2a-42e2-a8aa-9bc4be353e59
* title: Example Paragraph Challenge

<!--'question' is required, markdown, the question to be answered-->

### !question
Explain at least 2 benefits of writing semantic HTML.
### !end-question

<!--'placeholder' is optional, the placeholder text in the input box. Removed when the user starts typing-->

#### !placeholder
Write your answer here
#### !end-placeholder

<!--'explanation' is optional. Shown after the student correctly answers the question.-->

### !explanation
Your answer may have covered accessibility, SEO, and human readability.
### !end-explanation

### !end-challenge

<!----------------------END CHALLENGE----------------------------->



## Javascript Code Snippet

Code Snippet Challenges allow a student to write code directly in Learn. The submission is evaluated against unit tests that are setup as part of the Challenge. The student then sees the standard output from the test runner in Learn.

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

<!--'type' is required-->
<!--'id' is required, string, must be unique within a branch-->
<!--'language' is required for type: code-snippet. For javascript, use 'javascript' -->
<!--'title' is required, string, used when displaying results-->

* type: code-snippet
* id: dd9c31af-0fe8-440d-b4ec-bab3e8dc8a1d
* language: javascript
* title: Javascript `Repeats` Function

<!--'question' is required, markdown, the question to be answered-->

### !question

## Repeats

Write a function named `repeats`
* `repeats` should take one argument, `str`, the string to test.
* For this exercise, you can assume that `str` is a string.
* If the first half of `str` equals the second half, return true.
* If `str` is an empty string, return true.
* Otherwise, return false.
* Do not use the `.repeat` method.

### !end-question

<!--'placeholder' is optional, the starting code in the editor. Not removed when the user starts typing-->

#### !placeholder

```js
function repeats(str) {
   // return str.substring(0, str.length/2) === str.substring(str.length/2, str.length);
}
```

#### !end-placeholder

<!--'tests' is required for code-snippets and contains the unit tests that will be run to evaluate the student submission. Tests should be broken down into as many 'it(should...)' blocks as possible, because this is how students will get results. -->

### !tests

```js
describe('repeats', function() {

    it("should return true when given an empty string (which seems strange, but go with it :) )", function() {
      expect(repeats(""), "Default value is incorrect").to.deep.eq(true)
    })

    it("should return true when the second half of the string equals the first", function() {
      expect(repeats("bahbah")).to.deep.eq(true)
      expect(repeats("nananananananana")).to.deep.eq(true)
    })

    it("should return false when the second half of the string does not equal the first", function() {
      expect(repeats("bahba")).to.deep.eq(false)
      expect(repeats("nananananann")).to.deep.eq(false)
    })

    it("should not use .repeat", function() {
      expect(repeats.toString()).to.not.match(/\.repeat/)
    })

})
```
### !end-tests

### !end-challenge

<!----------------------END CHALLENGE----------------------------->



## Python Code Snippet

Code Snippet Challenges allow a student to write Python directly in Learn. The submission is evaluated against unit tests that are setup as part of the Challenge. The student then sees the standard output from the test runner in Learn.

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

<!-- type is required, type of challenge -->
<!--'language' is required for type: code-snippet. For python, options are 'python2.7' and 'python3.6' -->
<!-- id is required -->
<!-- title is required -->

* type: code-snippet
* language: python3.6
* id: 6f1a61d2-a3b4-402d-83be-38d443f03e72
* title: filter by class

<!-- !question is required -->

### !question

Implement the function `filter_by_class`: It takes a feature matrix, `X`, an array of classes, `y`, and a class label, `label`. It should return all of the rows from X whose label is the given label.

```python
>>> X = np.array([[1, 2, 3], [4, 5, 6], [7, 8, 9], [10, 11, 12]])
>>> y = np.array(["a", "c", "a", "b"])
>>> filter_by_class(X, y, "a")
array([[1, 2, 3],
       [7, 8, 9]])
>>> filter_by_class(X, y, "b")
array([[10, 11, 12]])
```
### !end-question

<!-- !placeholder is optional, the starter code that will be added to the editor -->

### !placeholder

```python
def filter_by_class(X, y, label):
    '''
    INPUT: 2 dimensional numpy array, numpy array, object
    OUTPUT: 2 dimensional numpy array

    Return the rows from X whose corresponding label from y is the given label.
    '''
    # return X[y == label]
```
### !end-placeholder

 <!-- !tests are required -->

### !tests
```python
import unittest
import main as p		
import numpy as np

class TestChallenge(unittest.TestCase):
  def test_filter_by_class1(self):
      X = np.array([[1, 2, 3], [4, 5, 6], [7, 8, 9], [10, 11, 12]])
      y = np.array(["a", "c", "a", "b"])
      result = p.filter_by_class(X, y, "a")
      answer = np.array([[1, 2, 3], [7, 8, 9]])
      assert np.array_equal(result, answer)

  def test_filter_by_class2(self):
      X = np.array([[1, 2, 3], [4, 5, 6], [7, 8, 9], [10, 11, 12]])
      y = np.array(["a", "c", "a", "b"])
      result = p.filter_by_class(X, y, "b")
      answer = np.array([[10, 11, 12]])
      assert np.array_equal(result, answer)
```
### !end-tests

### !end-challenge

<!----------------------END CHALLENGE----------------------------->



## Java Code Snippet

Code Snippet Challenges allow a student to write Java directly in Learn. The submission is evaluated against unit tests that are setup as part of the Challenge. The student then sees the standard output from the test runner in Learn.

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

* type: code-snippet
* language: java
* id: 48900333-21b1-47e8-bb04-9666cf3b4a03
* title: Single comparison

### !question

In the space given below, define and implement a method called `isActive`. It takes as input a `String` and returns `true` if the passed in string is "active", `false` if it is any other string.

### !end-question

### !setup

// [to allow student to submit simple statements, wrap the submission
//  using the !setup and !tests sections; example below]
class ChallengeSolution {

### !end-setup

### !placeholder
boolean isActive(String status) {
    // Implement your solution
}
### !end-placeholder

### !tests

}

public class SnippetTest {

    ChallengeSolution solution = new ChallengeSolution();

    @Test
    public void emptyStringReturnsFalse() {
        String input = new String("");
        assertEquals(false, solution.isActive(input), "For empty string");
    }

    @Test
    public void singleLetterReturnsFalse() {
        String input = new String("a");
        assertEquals(false, solution.isActive(input), "For single letter string");
    }

    @Test
    public void activeReturnsTrue() {
        String input = new String("active");
        assertEquals(true, solution.isActive(input), "For new \"active\" string");
    }

    @Test
    public void substringActiveReturnsFalse() {
        String input = new String("ctive");
        assertEquals(false, solution.isActive(input), "For substrings of \"active\"");
    }

    @Test
    public void activeSubstringReturnsFalse() {
        String input = new String("superactive");
        assertEquals(false, solution.isActive(input), "For input where a substring is \"active\"");
    }
}
### !end-tests

### !end-challenge

<!----------------------END CHALLENGE----------------------------->




## SQL Code Snippet

SQL Snippet Challenges allow a student to write SQL `SELECT` queries against a database defined by a curriculum developer. The submission is evaluated against a supplied test query. The student sees a truncated set of rows as if their query was run from the `psql` interpreter, and the challenge is graded against row and column matches, along with data matches. `ORDER BY` clauses will be respected if they are supplied in the test query.

The database used is:

`PostgreSQL 10.6 on x86_64-pc-linux-gnu, compiled by gcc (GCC) 4.8.3 20140911 (Red Hat 4.8.3-9), 64-bit`

### !callout-info
## Known issue with SQL challenge preview
SQL challenges will render in preview, but will not function if you submit an answer. To test them, publish the repo and try them out attached to an example cohort.
### !end-callout

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

<!--'type' is required-->
<!--'id' is required, string, must be unique; use a uuid generator-->
<!--'language' is required for type: code-snippet. For sql, use 'sql' -->
<!--'title' is required, string, used when displaying results-->
<!--'data_path' is required, string, must specify a sql file from the root of the repo. It contains everything necessary to create the database and populate the data-->

* type: code-snippet
* language: sql
* id: 7264b4e8-1a0d-44c5-b019-601305617ff8
* title: Absence of a join with Left Outer Join
* data_path: /01-example-unit/sql-files/foodtruck.sql

<!--'question' is required, markdown, the question to be answered-->

##### !question

Given a `users` table with a primary key of `id` and a `trucks` table with a foreign key that references the `users.id` column with `trucks.owner_id`, select all user information from users who do not own trucks, ordered by their last name descending, `users.last`.

##### !end-question

<!--'placeholder' is optional, the starting query in the editor. Not removed when the user starts typing-->

##### !placeholder

```sql
-- Write your query to select users who do not own trucks
select users.* from users
left outer join trucks on users.id = trucks.owner_id
where trucks.owner_id IS NULL
order by users.last desc
```

##### !end-placeholder

<!--'tests' is required for sql code-snippets and contains the query that should answer the question. Student exact-matches are correct, along with data matches. -->

##### !tests

SELECT users.*
FROM   users
       LEFT OUTER JOIN trucks
                    ON trucks.owner_id = users.id
WHERE  owner_id IS NULL
ORDER BY users.last DESC

##### !end-tests

### !end-challenge

<!----------------------END CHALLENGE----------------------------->



## Project Challenge

Project Challenges allow a student to do work outside of Learn. After completing the work, the student submits a link (typically from Github) to the work that they did for tracking and review within Learn.

<!-- Let Learn run tests automatically on repo submissions! See [Testable Project Challenges](Testable-Project-Challenge.md).-->

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

<!--'type' is required-->
<!--'id' is required, string, must be unique within a branch-->
<!--'title' is required, string, used when displaying results-->

* type: project
* id: 372d5f1f-e432-47f1-8de6-d0643f3773dc
* title: JS Native Array Methods

<!--'question' is required, markdown, the question to be answered-->

##### !question

### JS Native Array Methods

Submit the github repo containing your work on the js-native-array-methods below.

##### !end-question

<!--'placeholder' is optional, the placeholder text in the input field. -->

##### !placeholder

https://github.com/<username>/js-native-array-methods

##### !end-placeholder

### !end-challenge

<!----------------------END CHALLENGE----------------------------->



## Testable Project

Just like Project Challenges, but Learn will automatically run tests included in the repo and make the results available to the instructor.

To do this, you need two things --

1. An exercise repo that is setup to run in Docker (covered [here](https://learn-2.galvanize.com/cohorts/667/blocks/13/content_files/Testing-Project-Challenges.md)).
2. A project challenge where the students can submit their fork of the exercise repo (this doc).

<!----------------------BEGIN CHALLENGE----------------------------->

### !challenge

<!--'type' is required-->
<!--'id' is required, string, must be unique within a branch-->
<!--'title' is required, string, used when displaying results-->
<!--'upstream' is required, the upstream repo that will be forked. Set to the branch with the correct Dockerfile and test.sh for running tests-->
<!--'validate_fork' is options, set to true to require that the student submission is a fork of the upstream. If not defined, default is false.-->

* type: testable-project
* id: 358b7964-f097-4c11-a29c-7a3d72ea1098
* title: JS Native Array Methods
* upstream: https://github.com/gSchool/js-native-array-methods/
* validate_fork: false

<!--'question' is required, markdown, the question to be answered-->

##### !question

### JS Native Array Methods

Submit the github repo containing your work on the js-native-array-methods below.

incorrect answer: https://github.com/sperella/js-native-array-methods

correct answer: https://github.com/sperella/js-native-array-methods/tree/correct-answer

**Note: Fix correct example repo**

##### !end-question

<!--'placeholder' is optional, the placeholder text in the input field. -->

##### !placeholder

https://github.com/<username>/js-native-array-methods

##### !end-placeholder


### !end-challenge

<!----------------------END CHALLENGE----------------------------->
