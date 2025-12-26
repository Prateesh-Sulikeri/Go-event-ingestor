#!/bin/bash
echo "Resetting event_ingestor.events table..."
sudo -iu postgres psql event_ingestor -c "TRUNCATE TABLE events;"
echo "Done."
