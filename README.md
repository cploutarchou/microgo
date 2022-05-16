## MicroGO the Go version of Laravel Framework

In MicroGO, I take some of the most valuable features in Laravel and implement similar functionality in Go.

Since Go is compiled and type-safe, web applications written in this language are typically much faster and far less
error-prone than an equivalent application, Laravel, written in PHP.

________________________________
# Work-in-Progress ....................................
________________________________

### MicroGO Terminal Commands:

* **help**                           - show the help commands
* **version**                        - print application version
* **make auth**                      - Create and runs migrations for auth tables, create models and middleware.
* **migrate**                        - runs all up migrations that have not been run previously
* **migrate down**                   - reverses the most recent migration
* **migrate reset**                  - runs all down migrations in reverse order, and then all up migrations
* **make migration migration_name**  - creates two new up and down migrations in the migrations folder
* **make handler handler_name**      - create a stub handler on handlers directory
* **make model  model_name**          - create a new mode in the models directory

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate?hosted_button_id=EH6BNRFVPZ63N)
