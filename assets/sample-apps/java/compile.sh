#!/bin/sh
if ! javac -target 8 -source 8 App.java; then
    echo "An error occured while compiling"
    exit 1
fi

if ! jar cfe App.jar App App.class; then
    echo "An error occured while compiling"
    exit 1
fi

echo "Compiled successfully!"
exit 0