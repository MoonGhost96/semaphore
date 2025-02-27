package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ansible-semaphore/semaphore/api"
	"github.com/ansible-semaphore/semaphore/api/sockets"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/db/factory"
	"github.com/ansible-semaphore/semaphore/services/schedules"
	"github.com/ansible-semaphore/semaphore/services/tasks"
	"github.com/ansible-semaphore/semaphore/util"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "semaphore",
	Short: "Ansible Semaphore is a beautiful web UI for Ansible",
	Long: `Ansible Semaphore is a beautiful web UI for Ansible.
Source code is available at https://github.com/ansible-semaphore/semaphore.
Complete documentation is available at https://ansible-semaphore.com.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Configuration file path")
	if err := rootCmd.Execute(); err != nil {
		log.Error(os.Stderr, err)
		os.Exit(1)
	}
}

func runService() {
	logFormatter := &log.TextFormatter{
		ForceColors: true,
	}
	log.SetFormatter(logFormatter)
	log.Println("服务启动...")
	store := createStore("root")
	taskPool := tasks.CreateTaskPool(store)
	schedulePool := schedules.CreateSchedulePool(store, &taskPool)
	tasks.ClearPlaybookTmpFile()

	defer schedulePool.Destroy()

	util.Config.PrintDbInfo()

	log.Printf("Tmp Path (projects home) %v", util.Config.TmpPath)
	log.Printf("Semaphore %v", util.Version)
	log.Printf("Interface %v", util.Config.Interface)
	log.Printf("Port %v", util.Config.Port)

	go sockets.StartWS()
	go schedulePool.Run()
	go taskPool.Run()

	route := api.Route()

	route.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.Set(r, "store", store)
			context.Set(r, "schedule_pool", schedulePool)
			context.Set(r, "task_pool", &taskPool)
			next.ServeHTTP(w, r)
		})
	})

	var router http.Handler = route

	router = handlers.ProxyHeaders(router)
	http.Handle("/", router)

	log.Println("Server is running")

	if store.PermanentConnection() {
		defer store.Close("root")
	} else {
		store.Close("root")
	}

	err := http.ListenAndServe(util.Config.Interface+util.Config.Port, cropTrailingSlashMiddleware(router))

	if err != nil {
		log.Panic(err)
	}
}

func createStore(token string) db.Store {
	util.ConfigInit(configPath)

	store := factory.CreateStore()

	store.Connect(token)

	//if err := store.Connect(token); err != nil {
	//	switch err {
	//	case bbolt.ErrTimeout:
	//		fmt.Println("\n BoltDB supports only one connection at a time. You should stop Semaphore to use CLI.")
	//	default:
	//		fmt.Println("\n Have you run `semaphore setup`?")
	//	}
	//	os.Exit(1)
	//}

	err := db.Migrate(store)

	if err != nil {
		panic(err)
	}

	return store
}
