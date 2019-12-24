- dry run flag
- be its own make file through the yml config file
- if user doesn't provide version we should lookup which one was the last and prompt user to enter new one
- read in github token from a file, then fall back to an env, and then fail if not found
- Create a test for pretext to make sure it can always be compiled.
- Preserve changelog in case of error. (Make changelog a predictable location and then attempt to recover if release is called again with the
  right version)
- deploy should have a skip download/upload command
