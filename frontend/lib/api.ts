import { type APIKey, type CreateTodoRequest, type Todo, type UpdateTodoRequest, type User } from "@/types/todo";

// クライアント側からは /api/* (Next.js API Route) を呼ぶ
// API RouteがCookieからAPIキーを取り出してバックエンドに転送する
// バックエンドのレスポンス形式: { resultHeader: { resultCode }, resultBody: <payload> }

async function handleBackendResponse<T>(res: Response): Promise<T> {
  if (res.status === 204) return undefined as T;
  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.error ?? res.statusText);
  }
  return data.resultBody as T;
}

export async function listTodos(): Promise<Todo[]> {
  const res = await fetch("/api/todos");
  const body = await handleBackendResponse<Todo[]>(res);
  return body ?? [];
}

export async function createTodo(body: CreateTodoRequest): Promise<{ id: string }> {
  const res = await fetch("/api/todos", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return handleBackendResponse<{ id: string }>(res);
}

export async function updateTodo(id: string, body: UpdateTodoRequest): Promise<Todo> {
  const res = await fetch(`/api/todos/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return handleBackendResponse<Todo>(res);
}

export async function deleteTodo(id: string): Promise<void> {
  const res = await fetch(`/api/todos/${id}`, { method: "DELETE" });
  return handleBackendResponse<void>(res);
}

// サーバーサイド専用（Next.js API Route から呼ぶ用）
export async function createUser(
  name: string,
  email: string,
  backendUrl: string
): Promise<User> {
  const res = await fetch(`${backendUrl}/v1/users`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, email }),
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? res.statusText);
  }
  const data = await res.json();
  return data.resultBody as User;
}

export async function createAPIKey(
  userId: string,
  name: string,
  backendUrl: string
): Promise<APIKey> {
  const res = await fetch(`${backendUrl}/v1/api-keys`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ user_id: userId, name }),
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? res.statusText);
  }
  const data = await res.json();
  return data.resultBody as APIKey;
}
