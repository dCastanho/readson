
# ReadSON


ReadSON is a [command line utility](#command) with the purpose of converting structred JSON data into a more readable format, such as markdown or plain text. The main use case is to convert large databases, exported as JSON, into personal knowledge bases, such as [Obsidian](https://obsidian.md/) or [Foam](https://github.com/foambubble/foam).

This conversion is done through templates, the detailed syntax of which is described in [Templates](#templates). These templates are textual files, in which blocks of text are mixed with blocks of logic and variable accesses.  This conversion can be performed at large scale, allowing you to generate files from several files or array elements of a file in a single command - so long as they all use the same template. For example, take the following JSON.

```json
{
    "name": "John",
    "age" : 28,
    "positions" : [
        { 
            "position" : "Junior Dev",
            "start" : "10-12-2018",
            "end" : "05-08-2021"
        },
        { 
            "position" : "Mid-Level Dev",
            "start" : "05-08-2021"
        }
    ]
}
```

If you wanted to convert this into a more readable format, you could apply the following template, which would convert this data. 

```markdown
# $name$
Age: $age$

Positions:
$ for i, pos = range positions $
- **$ pos->position $**. From: $pos->start$$ if exists pos->end $, To: $pos->end$ $end$ $ end $
```

Which would result in the following file. 

```markdown
# John
Age: 28

Positions:
- **Junior Dev**. From: 10-12-2018, To: 05-08-2021 
- **Mid-Level Dev**. From: 05-08-2021
```


## Templates

Templates are simple text files composed of *text blocks* and *logic blocks*. Text blocks are kept as is, whereas logic blocks are either flow control constructs (such as if and for) or data variable accesses, which change depending on the data. Logic blocks are defined by opening and closing `$`.

### Variable Accesses and JSON paths

The simplest logic block is a variable access. This block - the opening and closing of `$` and everything in it - is replaced by the variable they are accessing. So the block `$ age $` is replaced by the age value, for example `28`.

The templates support indexes and JSON paths too. If you have nested objects, you can traverse down with the `->` symbol. So if `name` is an object with the properties `first` and `last`, you could access the first name by doing `name->first`.

Indexes are accessed by using `[i]`. So accessing the element `0` of an array is as simple as `array[0]`. These components can be stacked together.

> There is no type checking in templates, so if you try to access a variable that is not there or a non-existent index, it will give an error on evaluation.

### If-then-else

The if construct is also supported by templates with the following syntax:

```
$ if <condition> $
 true clause
$ else $
 false clause
$ end $
```

The `$ else $` and false clause are optional, and in the case of several sequential "if else"s, the alternate syntax is also available (where once again the final else is optional):

```
$ if <condition one> $
 first clase
$ else if <condition> $
 second clause
$ else $ 
 third clause
$ end $
```

In an if, when a condition is true, the whole construct is replaced by the clause matching it - so if `condition one` was true, the whole thing would be replaced by "first clause". In each clause, **every** symbol between the two `$` is kept, including spaces, tabs, and **line breaks**. Clauses are evaluated as normal, so they may contain anything a template would, including other logic blocks.

#### Conditions

##### Boolean operators

Conditions are the gates of if, which determine the clause to evaluate and use in the final result. The boolean operators of `and` and `or` are supported (using those keywords), and with *ors having priority over ands*. Negation `!` is also possible as well as changing priority/grouping conditions with parenthesis `()`. 

Short-circuit evaluation is implemented, which means that in `and` the condition stops evaluating on the first (left-most) false, and in `or` it stops on the first true. 

##### Comparison

Base conditions can be comparisons between variables and constants (i.e. `age > 18`) and the following operators are supported:
- `=` equals
- `!=` difference
- `>` greater and `>=` greater or equal
- `<` lesser and `<=` lesser or equal

Boolean variables can also be used as conditions (i.e `married and age > 18`). 

> Variables/constants of different types cannot be compared (a string with a number). Strings and booleans can be compared to each other only with the `=` operator. 


##### Exists and Isa

There are two built-in constructs that can be used in conditions:
- `exists` checks whether a given variable exists in the data, allowing for safe access. For example, `$ if exists spouse $ $ spouse->name $ $end$` gets the name of the spouse, if there is one. Notice that accessing the name of the spouse is done as a normal access to a variable.

- `isa` checks whether a variable is of a given type. For example, `$ if spouse isa object $ $ spouse->name$ $ else if spouse isa string $ $ spouse $ $ end $` gets the name of the spouse if it is stored as an object, or just prints it if it is a string. This construct gives an error if the variable does not exist.


### For loops

For loops are a way to iterate over arrays or object properties. They are essentially *for each* loops, and not associated with a condition. There are two types of for loops in templates.

- `range` loops iterate over an array's items and indexes. 
```
$ for i, job = range positions$ 
 - $ job->name $ $end$ 
```
- `props` loops iterate over objects and their property/value pairs.

```
$ for name, grade = props grades $
- **$name$:** $grade$ $end$
```

In both loops it is possible to access general fields of the data. It is also necessary to assign both names, even if only one is used. ***Indexes start at 1***.

As with if clauses, every character between the `$ for ... $` and the `$ end $` are kept, including spaces and line breaks.

### Pre-processing

ReadSON supports basic pre-processing of templates, in this case, the only function that is executed is a defines-like replacement. Every line at the start of the template that begins with `$$$ <name> text` is a defines clause. Every block `$<name>$` further in the template is thus replaced by text. This allows for some simple refactorings - **linebreaks in `text` are not yet supported**.

Pre-processing generates a new template file called `processed_<template name>`, which is deleted after execution. For debugging purposes, the file may be kept by passing the `-k` option.


### Functions

Templates support functions in the form of function calls done as `function(args,)`. Functions can be used in conditions or normal text blocks - for example the function  `callHello` returns "Hello world!", and thus the block `$callHello()$` will be replaced by "Hello World!". 

There are no default functions, but an external Javascript file with function definitions may be defined, which are inserted into the templates. Thanks to the [otto package](https://pkg.go.dev/github.com/robertkrimen/otto#section-readme). The function file is provided to the command with the `-f <filename>` option.

## Command

The readson command looks like this:

```
readson [OPTIONS] data-path
```

### Data path

The datapath may be a glob (linux path expressions) and the last element may have data accesses. Each resulting "object" from the data path is passed by the given template independently. For example, if have the following file tree, and pass `people/person_*` or `people/`. Then each file will be passed by the template and generate a resulting file.

```
.
└── people/
    ├── person_a.json
    ├── person_b.json
    └── person_c.json
```

Now imagine that each file has a field `children : [ ... ]`. If we wanted to pass each child through a template and have each one generate a text file for themselves, we could instead pass the following datapath `people/person_*->children` or `people/->children`. **By default the resulting files have the name of the glob used to generate them followed by a sequential number.** Resulting files have the same extension as the template used to generate them.

### Options

`-h` 

Shows the help message

`-t TEMPLATE`, `--templ TEMPLATE`

This option is **mandatory** and is used to tell readson what template to use for this execution.

`-name PATTERN`, `-p pattern` 

This option is used to define the name of each resulting file according to a variable within the data. For example, if each block of data has a `name` field, then `-p name` would make each file have the value in `name` as their name.

`-o OUTPUT`, `--output OUTPUT`

Will give all generated files the same *constant* name that is in OUTPUT. Useful when there is only one file being generated.

`-f FILE`, `--functions FILE`

Define the path to the javascript functions file

`-k`, `--keep`

Tells ReadSON to keep the post-processed template

`-v`, `--verbose` 

Prints logs


