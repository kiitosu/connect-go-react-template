package example // go.modで定義しているpackageにあわせる

import (
    "context"
    "log"
    "example/ent" // package名を修正
    "example/ent/todo"
	"fmt"

    "entgo.io/ent/dialect"
    _ "github.com/mattn/go-sqlite3"
)

func ExampleTodo() {
    // インメモリーのSQLiteデータベースを持つent.Clientを作成します。
    client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
    if err != nil {
        log.Fatalf("failed opening connection to sqlite: %v", err)
    }
    defer client.Close()
    ctx := context.Background()
    // 自動マイグレーションツールを実行して、すべてのスキーマリソースを作成します。
    if err := client.Schema.Create(ctx); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }
	// レコード追加
	task1, err := client.Todo.Create().SetText("Add GraphQL Example").Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	fmt.Printf("%d: %q\n", task1.ID, task1.Text)

	task2, err := client.Todo.Create().SetText("Add Tracing Example").Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	fmt.Printf("%d: %q\n", task2.ID, task2.Text)

	// リレーションシップ追加
	if err := task2.Update().SetParent(task1).Exec(ctx); err != nil {
        log.Fatalf("failed connecting todo2 to its parent: %v", err)
    }

	// 全てのTODOを取得する
	items, err := client.Todo.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed querying todos: %v", err)
	}
	for _, t := range items {
		fmt.Printf("%d: %q\n", t.ID, t.Text)
	}

	// 他のTODOを親に持つTODOを取得する
    items, err = client.Todo.Query().Where(todo.HasParent()).All(ctx)
    if err != nil {
        log.Fatalf("failed querying todos: %v", err)
    }
    for _, t := range items {
        fmt.Printf("%d: %q\n", t.ID, t.Text)
	}

	// 他のTODOを親に持たず、かつ、他のTODOを子に持つTODOを取得する
	items, err = client.Todo.Query().
        Where(
            todo.Not(
                todo.HasParent(),
            ),
            todo.HasChildren(),
        ).
        All(ctx)
    if err != nil {
        log.Fatalf("failed querying todos: %v", err)
    }
    for _, t := range items {
        fmt.Printf("%d: %q\n", t.ID, t.Text)
    }

	// 子TODOを通じて親TODOを取得し、
    // クエリが正確に1つのTODOを返すことを期待します。
    parent, err_parent := client.Todo.Query(). // すべてのtodoアイテムを取得する
        Where(todo.HasParent()).        // 親todoアイテムを持つtodoアイテムのみにフィルタリング
        QueryParent().                  // 親todoアイテムについて走査を続ける
        Only(ctx)                       // 1つのtodoアイテムのみ取得する
    if err != nil {
        log.Fatalf("failed querying todos: %v", err_parent)
    }
    fmt.Printf("%d: %q\n", parent.ID, parent.Text)


	// Output:
	// 1: "Add GraphQL Example"
	// 2: "Add Tracing Example"
	// 1: "Add GraphQL Example"
	// 2: "Add Tracing Example"
	// 2: "Add Tracing Example"
	// 1: "Add GraphQL Example"
	// 1: "Add GraphQL Example"
}


