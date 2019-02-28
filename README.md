A simple Job executor.
---------------------


This project gives a server that will allow you to: 

1. Create a task either time based or event based.
2. Execute a event based task by posting an event.
3. Stop a scehduled task.

Each task conatins a shell script file, which needs to be executed when executing the task.

Work pending:
1. Check the task's shell script before executing it, for any dangerous commands.
2. Endpoint to get all jobs.

