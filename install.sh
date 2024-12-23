#!/bin/bash

TODO_DIR="$HOME/.config/todo"
TODO_FILE="$TODO_DIR/todos.json"
BINARY_NAME="todo"
INSTALL_PATH="/bin/$BINARY_NAME"

if [ ! -d "$TODO_DIR" ]; then
  mkdir -p "$TODO_DIR"
  echo "Created directory: $TODO_DIR"
fi

if [ ! -f "$TODO_FILE" ]; then
  echo "[]" > "$TODO_FILE"
  echo "Created file: $TODO_FILE"
fi

echo "Building the binary..."
go build -o $BINARY_NAME

chmod +x $BINARY_NAME
sudo mv $BINARY_NAME $INSTALL_PATH

echo "Installed $BINARY_NAME to $INSTALL_PATH"
echo "You can now run the app using '$BINARY_NAME'"

