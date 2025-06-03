import { useState } from 'react'

// src/App.tsx
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

// 接続したいサービスをインポート
import { GreetService, TodoService } from "../gen/greet/v1/greet_pb";
import type { Todo } from "../gen/entpb/entpb_pb"


// transportではどのタイプのエンドポイントを使うか定義します
// 今回はConnect endpointを使います。
// エンドポイントが`g-RPC-web`しか対応していない場合は`createGrpcWebTransport`を使ってください
const transport = createConnectTransport({
    baseUrl: "http://localhost:8080"
})

// サービス定義とtransportを組み合わせてクライアントを作ります
const client = createClient(GreetService, transport)
const todoClient = createClient(TodoService, transport)

function App() {
    const[requestState, setRequest] = useState("")
    const [responseState, setResponse] = useState<string[]>([]);
    const [todos, setTodos] = useState<Todo[]>([]);
    return <>
        <form onSubmit={async (e) => {
            e.preventDefault(); // ページリロードを避ける
            const response = await client.greet({
                name: requestState
            })
            setResponse([...responseState, response.greeting])
        }}>
            <input value={requestState} onChange={e =>setRequest(e.target.value)}/>
            <button type="submit">Send</button>
        </form>

        <>
            {
                responseState.map( item => (
                    <li>{item}</li>
                ))
            }
        </>

        <button onClick={async () => {
            const response = await todoClient.listTodos({})
            setTodos(response.todos)
        }}>
            get todos
        </button>

        <div>
            {
                todos.map( todo => (
                    <li key={todo.id}>
                        ID: {todo.id}, タイトル: {todo.text}, 状態: {todo.status}, 優先度: {todo.priority}
                    </li>
                ))
            }
        </div>
    </>
}
export default App
