
# ReadSON

Goal of this project is to develop a fast CLI utility tool that quickly allows to convert a large amount of complex *JSON* files into *markdown*, according to a 
template. This is particularly useful for use in knowledge base systems where users have a lot of structred data they want to turn into *readable* text.

Commands:
- Turn a single JSON into a Markdown file with a certain template 
`readson file.json -t template.md`
- Turn several files - defined by either a folder or an expression - into Markdown files, following a single template, applied to each file in the directory or each file matching the expression
`readson [dir/ | file_*.json] -t template.md`
- More complex structure, through a config file.
`readson dir/ -t template.md -c config.yml`

> By default, these files are stored with the same name but with a different extension - `.md`. Customize with `-o` flag. In the case of multiple files, a sequential number is added to each, in no particular order.

Config file is useful for different situations and works in blocks, the goal is to allow further custimzation over what is applied to what. You may apply expressions to the files as well - the return of that expression must either be an object or an array.

```yaml
blocks:
    - files: people*
      template: person.md
    - files: people*.affiliates
      template: afilliate.md
    - files: books*
      template: book.md 
```


