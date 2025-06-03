package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"example/ent"
	"entgo.io/ent/dialect"
    _ "github.com/mattn/go-sqlite3"

	"connectrpc.com/connect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	greetv1 "example/gen/greet/v1"
	"example/gen/greet/v1/greetv1connect"

	"example/gen/entpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/* Greet */
type GreetServer struct{}

func (s *GreetServer) SaveRequest(
	ctx context.Context,
	name string,
) (string, error) {
	// ファイルベースのSQLiteデータベースを持つent.Clientを作成します。
	client, err := ent.Open(dialect.SQLite, "file:ent.db?_fk=1")

	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// 自動マイグレーションツールを実行して、すべてのスキーマリソースを作成します。
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// レコード追加
	_, err = client.Todo.Create().SetText(fmt.Sprintf("Hello from %s", name)).Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}

	// 全てのTODOを取得する
	items, err := client.Todo.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed querying todos: %v", err)
	}

	// 全てのTODOを接続して返す
	var combinedText string
	for _, t := range items {
		combinedText += fmt.Sprintf("%d: %q\n", t.ID, t.Text)
	}
	return combinedText, nil
}

func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())

	// リクエストを保存して結合した結果を返す
	all_request, err := s.SaveRequest(ctx, req.Msg.Name)
	if err != nil {
		log.Fatalf("failed SavedRequest: %v", err)
	}

	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: all_request,
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

/* Todo */
type TodoServer struct{}

func (s *TodoServer) ListTodos(
	ctx context.Context,
	req *connect.Request[greetv1.ListTodosRequest],
) (*connect.Response[greetv1.ListTodosResponse], error) {
	// ファイルベースのSQLiteデータベースを持つent.Clientを作成します。
	client, err := ent.Open(dialect.SQLite, "file:ent.db?_fk=1")

	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// 自動マイグレーションツールを実行して、すべてのスキーマリソースを作成します。
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	todos, err := client.Todo.Query().All(ctx)
	if err != nil {
	return nil, err
	}

	var pbTodos []*entpb.Todo
	for _, todo := range todos {
		pbTodo := &entpb.Todo{
			Id: int64(todo.ID),
			Text: todo.Text,
			CreatedAt: timestamppb.New(todo.CreatedAt),
			Status: entpb.Todo_Status(entpb.Todo_Status_value[string(todo.Status)]),
			Priority: int64(todo.Priority),
		}
		if err != nil {
			return nil, err
		}
		pbTodos = append(pbTodos, pbTodo)
	}

	res := connect.NewResponse(&greetv1.ListTodosResponse{
		Todos: pbTodos,
	})

	return res, nil
}

func withCORS(h http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-User-Agent", "Connect-Protocol-Version"},
	}).Handler(h)
}

func main() {
	// マルチプレクサ(ルータ)を生成
	mux := http.NewServeMux()

	// マルチプレクサ(ルータ)にパスとハンドラを追加
	mux.Handle(greetv1connect.NewGreetServiceHandler(&GreetServer{}))
	mux.Handle(greetv1connect.NewTodoServiceHandler(&TodoServer{}))

	// httpサーバーを起動
	http.ListenAndServe(
		"0.0.0.0:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		withCORS(h2c.NewHandler(mux, &http2.Server{})),
	)
}
