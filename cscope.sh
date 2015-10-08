#!/bin/sh

find . -name "*.go">cscope.files

cscope -bkq -i cscope.files

ctags -R
