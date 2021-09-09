#!/bin/bash

git config core.hooksPath scripts/git_hooks && echo ">> local githook has been set" && git config --get core.hooksPath
