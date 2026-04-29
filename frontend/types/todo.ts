export type TodoStatus = "pending" | "done"

export interface Todo {
  id: string
  title: string
  description: string
  status: TodoStatus
  created_at: string
  updated_at: string
}

export interface CreateTodoInput {
  title: string
  description?: string
}

export interface UpdateTodoInput {
  title?: string
  description?: string
  status?: TodoStatus
}
