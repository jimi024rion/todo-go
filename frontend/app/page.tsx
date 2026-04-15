// "use client" なし → Server Component（サーバー側で実行される）
// Goのハンドラと同じ感覚：リクエストのたびにサーバーで実行され、HTMLを返す

import { cookies } from "next/headers";
import { ApiKeySetup } from "@/components/ApiKeySetup";
import { ResetKeyButton } from "@/components/ResetKeyButton";
import { TodoList } from "@/components/TodoList";
import { type Todo } from "@/types/todo";

// サーバー側でTodo一覧を取得する関数
// ブラウザを介さないのでCORSが発生しない（サーバー同士の通信）
async function fetchInitialTodos(apiKey: string): Promise<Todo[]> {
  try {
    const res = await fetch(
      `${process.env.BACKEND_URL ?? "http://localhost:8080"}/v1/todos`,
      {
        headers: {
          "Content-Type": "application/json",
          "X-API-Key": apiKey,
        },
        // SSRなので毎回最新データを取得する（Next.jsのキャッシュを無効化）
        cache: "no-store",
      },
    );
    if (!res.ok) return [];
    return res.json();
  } catch {
    return [];
  }
}

// async function → サーバー側で await が使える（useEffect不要）
export default async function Home() {
  // cookies() はサーバー側でのみ使える関数
  // httpOnly CookieはJavaScriptからは読めないが、サーバーからは読める
  const apiKey = (await cookies()).get("api_key")?.value;

  // サーバー側でTodoを取得（Cookieが設定されている場合）
  const initialTodos = apiKey ? await fetchInitialTodos(apiKey) : undefined;

  return (
    <div className="min-h-screen bg-white px-4 py-12 sm:px-6">
      {/* APIキー未設定時はセットアップモーダルを表示 */}
      {!apiKey && <ApiKeySetup />}

      <div className="mx-auto max-w-xl">
        <header className="mb-8 flex items-center justify-between">
          <h1 className="text-2xl font-bold tracking-tight text-[#111827]">Todo</h1>
          {/* ResetKeyButtonはクリックイベントが必要なのでClient Component */}
          {apiKey && <ResetKeyButton />}
        </header>

        {/* initialTodosをpropsで渡す → 画面表示と同時にデータが表示される */}
        {apiKey && <TodoList initialTodos={initialTodos} />}
      </div>
    </div>
  );
}
