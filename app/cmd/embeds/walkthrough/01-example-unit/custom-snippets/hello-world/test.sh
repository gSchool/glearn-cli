#!/bin/bash
submission_file=${1:-submission.txt}

fail () {
    echo "Didn't find the exact text 'Hello world!'."
    exit 1
}

grep -F 'Hello world!' "$submission_file" >/dev/null || fail

echo "✅ Found 'Hello world!'"
exit 0
