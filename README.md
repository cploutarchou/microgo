## MicroGO the Go version of Laravel Framework

In MicroGO, I take some of the most valuable features in Laravel and implement similar functionality in Go.

Since Go is compiled and type-safe, web applications written in this language are typically much faster and far less
error-prone than an equivalent application, Laravel, written in PHP.

________________________________
# Work-in-Progress ....................................
________________________________

### MicroGO Terminal Commands:

* **help**                           - Show the help commands
* **version**                        - Print application version
* **make auth**                      - Create and runs migrations for auth tables, create models and middleware.
* **migrate**                        - Runs all up migrations that have not been run previously
* **migrate down**                   - Reverses the most recent migration
* **migrate reset**                  - Runs all down migrations in reverse order, and then all up migrations
* **make migration migration_name**  - Create two new up and down migrations in the migrations folder
* **make handler handler_name**      - Create a stub handler on handlers directory
* **make model  model_name**         - Create a new mode in the models directory
* **make key**                       - Create a random key of 32 characters.
* **make mail**                      - Create two starter mail templates in the mail directory.

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate?hosted_button_id=EH6BNRFVPZ63N)
