package main

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	// run the migration command
	switch arg2 {
	case "up":
		err := micro.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":
		if arg3 == "all" {
			err := micro.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := micro.Steps(-1, dsn)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := micro.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = micro.MigrateUp(dsn)
		if err != nil {
			return err
		}
	default:
		help()
	}

	return nil
}
