package mysqldiff

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost     string `split_words:"true"`
	DBPort     string `split_words:"true"`
	DBUser     string `split_words:"true"`
	DBPassword string `split_words:"true"`
	DBName     string `split_words:"true"`
	FilePath   string `split_words:"true"`
}

func mergeConfig(config *Config, env Config) {
	if env.DBHost != "" {
		config.DBHost = env.DBHost
	}
	if env.DBPort != "" {
		config.DBPort = env.DBPort
	}
	if env.DBUser != "" {
		config.DBUser = env.DBUser
	}
	if env.DBPassword != "" {
		config.DBPassword = env.DBPassword
	}
}

func InitializeConfig() (Config, Config, Config) {
	var envLocal, envSrc, envDst Config
	if err := envconfig.Process("GO_MYSQL_DIFF_LOCAL", &envLocal); err != nil {
		panic(err)
	}
	if err := envconfig.Process("GO_MYSQL_DIFF_SRC", &envSrc); err != nil {
		panic(err)
	}
	if err := envconfig.Process("GO_MYSQL_DIFF_DST", &envDst); err != nil {
		panic(err)
	}
	var (
		argLocalDBHost     = flag.String("local-db-host", "127.0.0.1", "your local db host.")
		argLocalDBPort     = flag.String("local-db-port", "3306", "your local db port.")
		argLocalDBUser     = flag.String("local-db-user", "root", "your local db user.")
		argLocalDBPassword = flag.String("local-db-password", "", "your local db password.")

		argSrcDBHost     = flag.String("src-db-host", "", "src db host.")
		argSrcDBPort     = flag.String("src-db-port", "", "src db port.")
		argSrcDBUser     = flag.String("src-db-user", "", "src db user.")
		argSrcDBPassword = flag.String("src-db-password", "", "src db password.")
		argSrcDBName     = flag.String("src-db-name", "", "src db name.")
		argSrcFilePath   = flag.String("src-file-path", "", "src sql file.")

		argDstDBHost     = flag.String("dst-db-host", "", "dst db host.")
		argDstDBPort     = flag.String("dst-db-port", "", "dst db port.")
		argDstDBUser     = flag.String("dst-db-user", "", "dst db user.")
		argDstDBPassword = flag.String("dst-db-password", "", "dst db password.")
		argDstDBName     = flag.String("dst-db-name", "", "dst db name.")
		argDstFilePath   = flag.String("dst-file-path", "", "dst sql file.")
	)
	flag.Parse()

	// Set default Value
	localConfig := Config{
		DBHost:     *argLocalDBHost,
		DBPort:     *argLocalDBPort,
		DBUser:     *argLocalDBUser,
		DBPassword: *argLocalDBPassword,
	}
	srcConfig := Config{
		DBHost:     *argSrcDBHost,
		DBPort:     *argSrcDBPort,
		DBUser:     *argSrcDBUser,
		DBPassword: *argSrcDBPassword,
		DBName:     *argSrcDBName,
		FilePath:   *argSrcFilePath,
	}
	dstConfig := Config{
		DBHost:     *argDstDBHost,
		DBPort:     *argDstDBPort,
		DBUser:     *argDstDBUser,
		DBPassword: *argDstDBPassword,
		DBName:     *argDstDBName,
		FilePath:   *argDstFilePath,
	}
	// Set env value
	mergeConfig(&localConfig, envLocal)
	mergeConfig(&srcConfig, envSrc)
	mergeConfig(&dstConfig, envDst)

	// Set args
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "local-db-host":
			localConfig.DBHost = f.Value.String()
		case "local-db-port":
			localConfig.DBPort = f.Value.String()
		case "local-db-user":
			localConfig.DBUser = f.Value.String()
		case "local-db-password":
			localConfig.DBPassword = f.Value.String()

		case "src-db-host":
			srcConfig.DBHost = f.Value.String()
		case "src-db-port":
			srcConfig.DBPort = f.Value.String()
		case "src-db-user":
			srcConfig.DBUser = f.Value.String()
		case "src-db-password":
			srcConfig.DBPassword = f.Value.String()
		case "src-db-name":
			srcConfig.DBName = f.Value.String()
		case "src-file-path":
			srcConfig.FilePath = f.Value.String()

		case "dst-db-host":
			dstConfig.DBHost = f.Value.String()
		case "dst-db-port":
			dstConfig.DBPort = f.Value.String()
		case "dst-db-user":
			dstConfig.DBUser = f.Value.String()
		case "dst-db-password":
			dstConfig.DBPassword = f.Value.String()
		case "dst-db-name":
			dstConfig.DBName = f.Value.String()
		case "dst-file-path":
			dstConfig.FilePath = f.Value.String()
		}
	})

	return localConfig, srcConfig, dstConfig
}
