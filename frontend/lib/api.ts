import type { CreateTodoInput, Tag, Todo, UpdateTodoInput } from "@/types/todo"

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
  setTags: (todoId: string, tagIds: string[]) =>
    request<void>(`/api/todos/${todoId}/tags`, {
      method: "PUT",
      body: JSON.stringify({ tag_ids: tagIds }),
    }),
}

export const tagApi = {
  list: () => request<Tag[]>("/api/tags"),
  create: (name: string) =>
    request<Tag>("/api/tags", {
      method: "POST",
      body: JSON.stringify({ name }),
    }),
  delete: (id: string) =>
    request<void>(`/api/tags/${id}`, { method: "DELETE" }),
}
