#!/bin/sh

set -e
echo "exec database migration"
/workspace/garden migrate
exec "$@"
