## Gokbilgin Wiki Generator

### Motivation

M. Tayyib Gökbilgin is one of the most important Ottoman historians of the 20th century; he is known for his researches on the Ottoman Empire from its foundation to its collapse, especially on urbanism, administrative organization, and religion.

https://wiki.gokbilgin.com is dedicated to provides resources about Tayyib. One page in particular is the correspondence has had with other professors and historians at the time. 

Tayyib received letters from over 600 authors, each of which sent multiple letters. The manually create pages for each of these letters would be tedious and incredibly time consuming. 

Instead, this script will iterate through the directories and generate the pages, saving significant time. 

### Usage

Various scripts under the `./scripts/` directory can be used to ensure the letter folder structure is formatted correctly. 

Process: 
- `./scripts/rename_letter_dirs.sh` will rename the letter directories to the format `YYYY.MM.dd`
- `./scripts/check.sh` will detail which folders require attention
    - Manually review the directory and ensure the following convention:
        - a pdf ending with *original.pdf
        - an option pdf ending with *tercume.pdf
    - `runbook.md`provides insights on how combine images into pdfs
- `./scripts/extract_translation.sh` will extract translations from the pdf-image into a text file; this will be used for the content of the page
- `./scripts/generate_letters_yaml.sh` will generate a `letter.yaml` file
- Manually create an `author.yaml` file for the author in the author directory