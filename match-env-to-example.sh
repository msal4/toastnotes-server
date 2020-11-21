#!/bin/sh
# Use this script to match the .env file with .env.example 

grep -o "\(^#.*\|.*=\|^\$\)" .env > .env.example
