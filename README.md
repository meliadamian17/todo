# Cheap Simple Terminal TODO App (original ik)

I wanted to have a reminder of what I want to do everytime i launch my terminal, so i made this.

## Installation
Simply just run `./install.sh`
This will place the todos.json into `$HOME/.config/todo` and move the built binary to `/bin` so you can run it.

## My Setup
I want this to prompt me on launch of my terminal, so, in my `.zshrc` I added the following line:
```bash
todo l
```
