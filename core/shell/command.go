package shell

import "github.com/abiosoft/ishell"

func RunShell() {
	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "cleanup-emails",
		Help: "Find email duplicates & allow to clean them up.",
		Func: CleanupDuplicatedEmails,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "test-events",
		Help: "Test events abstraction.",
		Func: TestEventHandler,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "migrate-comments",
		Help: "Migrate legacy comments (before anzu).",
		Func: MigrateComments,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "gc",
		Help: "Run garbage collector and timed events.",
		Func: RunAnzuGarbageCollector,
	})

	// start shell
	shell.Start()

	select {}
}


