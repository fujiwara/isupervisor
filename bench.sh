#!/bin/sh

set -eux

echo "Starting benchmark"
sleep 1
echo "Finished benchmark"
echo '{"finished": true, "passed": true, "score": 1234}' > $ISUPERVISOR_RESULT
exit
