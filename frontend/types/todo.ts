export type TodoStatus = "pending" | "in_progress" | "completed"

export interface Tag {
  id: string
  name: string
}

export interface Todo {
  id: string
  title: string
  description: string
  status: TodoStatus
  due_date?: string | null
  tags: Tag[]
  created_at: string
  updated_at: string
}

export interface CreateTodoInput {
  title: string
  description?: string
  due_date?: string | null
}

export interface UpdateTodoInput {
  title?: string
  description?: string
  status?: TodoStatus
  due_date?: string | null
  clear_due_date?: boolean
}
