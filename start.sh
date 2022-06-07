#!/bin/sh

set -e
echo "exec database migration"
/workspace/garden migrate -s 2
exec "$@"

