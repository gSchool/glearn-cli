# Walkthrough Materials

Use this directory when learning how to publish curriculum to Galvanize Learn.

Start by previewing an individual file with `learn preview 00-hello-world.md`. Modify the file and preview again to see your changes. Follow the directions in the lesson to learn more.

Running `learn preview .` from the root of this project to preview the full directory contents.

## Access Requirements

In order to preview curriculum, a Learn user is required with the ability to [generate an api token](https://learn-2.galvanize.com/api_token). Only admins and instructors can generate an API token. To publish blocks, your user will also need the role `forge.blocks_manager`; this role is not necessary to preview curriculum.

Once an API token is obtained, it can be stored locally with `learn set --api_token=<YOUR_TOKEN>` which writes to a dotfile in your home directory called `.glearn-config.yaml`. The token is then used for authorizing CLI access to Learn.

## More details

Visit the [Learn Documentation](https://learn-2.galvanize.com/cohorts/667) cohort for a more comprehensive explanation of Learn and its features, including richer [challenge examples](https://learn-2.galvanize.com/cohorts/667/blocks/13/content_files/Testable-Project-Challenge.md), explanations for [content file types](https://learn-2.galvanize.com/cohorts/667/blocks/13/content_files/content-file-types/10-Lesson-Content-File.md), [attaching curriculum to cohorts](https://learn-2.galvanize.com/cohorts/667/blocks/13/content_files/create-cohort-course.md), viewing [additional markdown examples](https://learn-2.galvanize.com/cohorts/667/blocks/13/content_files/walkthrough/markdown-examples.md) and even using [MathJax to render LaTex](https://learn-2.galvanize.com/cohorts/667/blocks/13/content_files/walkthrough/math-jax-examples.md).

