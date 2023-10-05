---
# This is yaml frontmatter which will set values for a generated `autoconfig.yaml` file
# It is not displayed in Learn.
Type: Lesson # Options: Lesson, Checkpoint, Survey, Instructor, Resource
# Lesson is a normal content file
# Checkpoint is a test, which has special configuration options; limit one per unit
# Surveys hold a set of ungraded challenges, and are submitted all at once
# Instructor files are only visible to instructors
# Resources can be linked within content, but will not show up in the navigation side bar
UID: unique-identifier # can be set to any string, must be unique in the repository
# DefaultVisibility: hidden # Uncomment this line to default Lesson to hidden when used
---

# Configuration

If look at this file in your preferred text editor, you'll see content at the top that is not displayed in Learn. This yaml frontmatter is used to set attributes on the content file for Learn.

You can see an example header by running
```
learn md fh -o
```
from your command line.

### !callout-info

## Setting Title

You'll notice the title is not configured in the front matter- it is always derived from the first h1 header in the content file. If none can be found, the file name defines the title.

### !end-callout

## Tracking configuration in a file

After you ran `learn preview .` you'll see a new file in the root of the walkthrough directory named `autoconfig.yaml`. The yaml frontmatter at the top of each file sets the values for the `ContentFiles` key in the `autoconfig.yaml`.

You can see the format and more information for configuration options by running
```
learn md cfy -o
```
from your command line.

## Controlling Configuration from one file

Each time the preview is created, the `autoconfig.yaml` file is overwritten.

However, running `learn preview` or `learn publish` won't generate an `autoconfig.yaml` file if there is already a file named `config.yaml` at the root of the project. Both files are used to define content; `config.yaml` takes precedence over `autoconfig.yaml`, and the latter will rebuild itself on each `preview` or `publish`.

Go ahead and change the `autoconfig.yaml` file to `config.yaml`, then comment out the `ContentFiles` entry that looks like
```
  - Type: Lesson
    Path: /01-example-unit/01-configuration.md
    UID: unique-identifier
```
Then from the root of the project run `learn preview .` again. You'll notice that this content file is now missing!

Delete the `config.yaml` file entirely, and run `learn preview .` again. It will recreate the `autoconfig.yaml` file again as no config file was discovered.

## Controlling Unit attributes

In the configuration file there is a single entry for `Standards`. When the autoconfig was generated, its attributes were read from the file `01-example-unit/description.yaml`.

You can see an example `description.yaml` file by running
```
learn md dsy -o
```
from your command line.

Each unit directory like `01-example-unit/` should have one `description.yaml` file.

## Adding more Units

Create a new Unit directory as a sibling to `01-example-unit/` and name it `02-my-unit/`. Populate it with a single markdown file called `00-playground.md`. From the root of the walkthrough

```
mkdir 02-my-unit/
touch 02-my-unit/00-playground.md
```

Write the following contents to your new file:

```
---
Type: Lesson
UID: playground
---
# Playground

Use this file to experiment!
```

Then run `learn preview .` again from the root of the walkthrough directory. Follow the link and notice that you now have _two_ units displayed in your preview.

We never wrote a `description.yaml` file for our second unit! What did it display for the unit Title and Description? Try adding your own `description.yaml` file inside your `02-my-unit/` directory and change the settings as you see fit. Re-preview to see how Title and Description are now configured.

Try renaming your new unit directory to `00-my-unit/` and preview again with `learn preview .` to see that the ordering 
