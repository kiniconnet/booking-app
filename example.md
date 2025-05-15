dbPass := flag.String("dbpass", "npg_9d4fAcePNJZx", "Database password")
dbName := flag.String("dbname", "bookings", "Database name")
	dbUser := flag.String("dbuser", "bookings_owner", "Database user")


    DATABASE_URL='postgresql://bookings_owner:npg_9d4fAcePNJZx@ep-wild-fog-a4zo4n4z-pooler.us-east-1.aws.neon.tech/bookings?sslmode=require'