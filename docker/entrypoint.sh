#!/usr/bin/env sh
set -e

echo "starting SMPP simulator"
java -Djava.net.preferIPv4Stack=true -Djava.util.logging.config.file=/app/smpp/smppsim_logging.properties -jar /app/smpp/smppsim.jar /app/smpp/smppsim.1.props &
java -Djava.net.preferIPv4Stack=true -Djava.util.logging.config.file=/app/smpp/smppsim_logging.properties -jar /app/smpp/smppsim.jar /app/smpp/smppsim.2.props &
java -Djava.net.preferIPv4Stack=true -Djava.util.logging.config.file=/app/smpp/smppsim_logging.properties -jar /app/smpp/smppsim.jar /app/smpp/smppsim.3.props
