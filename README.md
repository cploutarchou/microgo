# MicroGO

### MicroGO Terminal Commands:

* **help**                  - show the help commands
* **version**               - print application version
* **make auth**             - Create and runs migrations for auth tables, create models and middleware.
* **migrate**               - runs all up migrations that have not been run previously
* **migrate down**          - reverses the most recent migration
* **migrate reset**         - runs all down migrations in reverse order, and then all up migrations
* **make migration** <name> - creates two new up and down migrations in the migrations folder