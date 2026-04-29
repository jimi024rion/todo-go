import type { CreateTodoInput, Todo, UpdateTodoInput } from "@/types/todo"

// Next.js の API Routes を経由してバックエンドを呼ぶ（クライアント側から呼ぶ）
// BACKEND_API_KEY は API Routes のサーバー側のみで使用される

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: { "Content-Type": "application/json", ...init?.headers },
  })
  if (res.status === 204) return undefined as T
  const data = await res.json()
  if (!res.ok) throw new Error(data.error ?? res.statusText)
  return data as T
}

export const todoApi = {
  list: () => request<Todo[]>("/api/todos"),
  create: (input: CreateTodoInput) =>
    request<{ id: string }>("/api/todos", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  update: (id: string, input: UpdateTodoInput) =>
    request<Todo>(`/api/todos/${id}`, {
      method: "PUT",
      body: JSON.stringify(input),
    }),
  delete: (id: string) =>
    request<void>(`/api/todos/${id}`, { method: "DELETE" }),
}
