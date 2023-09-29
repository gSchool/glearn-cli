---
# This is yaml frontmatter which will set values for a generated `autoconfig.yaml` file
# It is not displayed in Learn
Type: Lesson # Options: Lesson, Checkpoint, Survey, Instructor, Resource
UID: unique-in-the-repository # can be set to any string, must be unique
# DefaultVisibility: hidden # Uncomment this line to default Lesson to hidden
---

# Configuration

If look at this file in your preferred text editor, you'll notice content at the top that is not displayed in Learn. This yaml frontmatter is used to set attributes on the content file.

You can see an example header by running
```
learn md fh -o
```
from your command line.

### !callout-info

## Setting Title

You'll notice the title is not configured in the front matter- it is always derived from the first h1 header in the content file.

### !end-callout

## Tracking configuration in a file

After you ran `learn preview .` you may notice a new file in the root of the walkthrough directory named `autoconfig.yaml`. The yaml frontmatter on each file sets the values within this file.

You can see the format and more information on configuration options by running
```
learn md cfy -o
```
from your command line.

## Controlling Configuration from one file

Each time the preview is created, the `autoconfig.yaml` file is overwritten.

However, running `learn preview` or `learn publish` won't generate an `autoconfig.yaml` file if there is already a file named `config.yaml` at the root of the project. This file is used to define content, and takes precedence over `autoconfig.yaml`.

Go ahead and change the `autoconfig.yaml` file to `config.yaml`. Then comment out the `ContentFiles` entry that looks like
```
  - Type: Lesson
    Path: /01-example-unit/01-configuration.md
    UID: unique-in-the-repository
```
Then from the root of the project run `learn preview .` again. You'll notice that this content file is now missing!

## Controlling Unit attributes

You'll notice in the configuration file there is a single standard. When the autoconfig was generated, its attributes were read from the file `01-example-unit/description.yaml`. 

You can see an example `description.yaml` file by running
```
learn md dsy -o
```
from your command line.
