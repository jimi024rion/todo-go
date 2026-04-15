import { type CreateTodoRequest, type Todo, type UpdateTodoRequest } from "@/types/todo";

// クライアント側からは /api/* (Next.js API Route) を呼ぶ
// API RouteがCookieからAPIキーを取り出してバックエンドに転送する
// → ブラウザのJavaScriptはAPIキーを一切扱わない

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? res.statusText);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export async function listTodos(): Promise<Todo[]> {
  const res = await fetch("/api/todos");
  return handleResponse<Todo[]>(res);
}

export async function createTodo(body: CreateTodoRequest): Promise<{ id: string }> {
  const res = await fetch("/api/todos", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return handleResponse<{ id: string }>(res);
}

export async function updateTodo(id: string, body: UpdateTodoRequest): Promise<Todo> {
  const res = await fetch(`/api/todos/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return handleResponse<Todo>(res);
}

export async function deleteTodo(id: string): Promise<void> {
  const res = await fetch(`/api/todos/${id}`, { method: "DELETE" });
  return handleResponse<void>(res);
}
