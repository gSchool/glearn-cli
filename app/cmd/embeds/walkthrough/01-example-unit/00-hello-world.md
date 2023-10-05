---
Type: Lesson
UID: hello-world
---

# Hello World

### !callout-info
If you just ran `learn preview 00-hello-world.md`, you’re now looking at a temporary preview of this file.

Go back to the command line and run `learn preview .` to preview all example materials.

Visit the new link and click on the 'Example Unit' to return to this page.
### !end-callout

This file is written in markdown. You can change any part of it and run
```
learn preview .
```
to see your update. You may preview individual files by specifying them for faster previewing. This is the part of the recommended workflow when developing new content and modifying existing content.

## Your Preview Cohort

If you look at the URL Learn created where your content can be found, you'll notice a cohort id which belongs to you! Learn cohorts are where instructors and students come together to learn subjects and montior progress, but this cohort is just for your curriculum development. Content you preview will always attach itself to your preview cohort and remain there until your next preview.

The urls are built using the file names in the project, so each time you preview the same project you can simply refresh the browser to see changes you have made.

If you just ran `learn preview .` from the root of the walkthrough directory, you might notice new content in the side bar. This content is produced from the walkthrough files:
```
├── 01-example-unit
│   ├── 00-hello-world.md
│   ├── 01-configuration.md
│   ├── 02-publishing.md
│   ├── 03-markdown-examples.md
│   ├── 04-challenges.md
│   ├── 05-checkpoint.md
```

* The 'Configuration' file explains how a repository can be organized into units of content.
* 'Publishing' shows you how to make your materials available for use in a cohort.
* Explore rendering options in 'Markdown Examples'.
* See how Learn enables inline checks for understanding with 'Challenge Examples'.
* Each unit can assess a student's understanding with 'Checkpoint Example'.

## Generating Content

The `learn` CLI tool can generate boilerplate markdown for challenges and other custom markdown content (like the callout above) with the command 
```
learn md
```

so for example, you can see how a callout markdown can be rendered in a content file with

```
learn md co -o
```

Try it now and compare it to the content of the `00-hello-world.md` file in your editor. The callout markdown you generated appears much like the callout at the top of the file.

If you use the `-o` flag the content is sent to `STDOUT`, while the `-m` flag produces a minimal version of the content.

Each lesson will explore different options with this command.

