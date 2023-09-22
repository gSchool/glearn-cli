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

If you just ran `learn preview .` from the root of the walkthrough directory, you'll notice new content in the side bar. This content is rendered from the walkthrough files:
```
├── 01-example-unit
│   ├── 00-hello-world.md
│   ├── 01-markdown-examples.md
│   ├── 02-challenges.md
│   ├── 03-checkpoint.md
```

Click on 'Markdown Examples' to see more of what Learn is capable of creating.

The content file 'Challenge Examples' will show you how Learn can create interactive challenges.

## Generating Content

The `learn` CLI tool can generate boilerplate markdown for challenges and other custom markdown content (like the callout above) with the command `learn md`. Try it now and see what is available.

If you use the `-o` flag the content is sent to `STDOUT`, while the `-m` flag produces a minimal version of the content.

## Checkpoints

Finally see Learn's checkpoint delivery capability with the 'Checkpoint Example'. Each unit can have only one Checkpoint, and the checkpoint content file must have challenges. They can be configured to score points automatically, grant a limited number of attempts, require a certain time limit to complete, and email the instructors upon completion.

See the header on the `01-example-unit/03-checkpoint.md` walkthrough file to configure these options.


