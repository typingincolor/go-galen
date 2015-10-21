#!/usr/bin/env bash

gotags -tag-relative=true -R=true -sort=true -f="tags" -fields=+l .
