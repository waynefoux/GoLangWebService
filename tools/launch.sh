#!/bin/bash

echo "Starting launch script."

umask 000

# Set applications home
if [[ ! -n "$APP_HOME" ]]; then
  export APP_HOME=".."
else
  export APP_HOME="/opt/app/locker-hw-consumer"
fi

echo "$appconf" > ${APP_HOME}/conf/application.properties

LOGGING_HOME="${APP_HOME}/logs"

STDOUT_FILE="${LOGGING_HOME}/${HOSTNAME}_csg_std.out"

EXEC_COMMAND="exec /opt/app/locker-hw-consumer/bin/locker-hw-consumer -properties=${APP_HOME}/conf/application.properties"

echo "Starting service."
echo ${EXEC_COMMAND}

if [[ ! -n ${LOG_TO_DISK} ]]
then
	$EXEC_COMMAND
else
	$EXEC_COMMAND > ${STDOUT_FILE} 2>&1
fi
