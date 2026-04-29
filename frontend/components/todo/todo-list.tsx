"use client"

import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { todoApi } from "@/lib/api"
import { TodoCard } from "./todo-card"
import { TodoCreate } from "./todo-create"
import { TodoPanel } from "./todo-panel"
import { EmptyState } from "./empty-state"
import type { Todo } from "@/types/todo"

export function TodoList() {
  const queryClient = useQueryClient()
  const [selectedTodo, setSelectedTodo] = useState<Todo | null>(null)

  const { data: todos = [], isLoading, error } = useQuery({
    queryKey: ["todos"],
    queryFn: todoApi.list,
  })

  const createMutation = useMutation({
    mutationFn: todoApi.create,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["todos"] }),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, ...input }: { id: string } & Parameters<typeof todoApi.update>[1]) =>
      todoApi.update(id, input),
    onSuccess: (updatedTodo) => {
      queryClient.setQueryData<Todo[]>(["todos"], (old = []) =>
        old.map((t) => (t.id === updatedTodo.id ? updatedTodo : t))
      )
      // パネルの Todo も最新に
      if (selectedTodo?.id === updatedTodo.id) setSelectedTodo(updatedTodo)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: todoApi.delete,
    onSuccess: (_, id) => {
      queryClient.setQueryData<Todo[]>(["todos"], (old = []) =>
        old.filter((t) => t.id !== id)
      )
    },
  })

  async function handleCreate(title: string) {
    await createMutation.mutateAsync({ title })
  }

  async function handleToggle(id: string, done: boolean) {
    await updateMutation.mutateAsync({ id, status: done ? "done" : "pending" })
  }

  async function handleUpdate(id: string, input: Parameters<typeof todoApi.update>[1]) {
    await updateMutation.mutateAsync({ id, ...input })
  }

  async function handleDelete(id: string) {
    await deleteMutation.mutateAsync(id)
  }

  if (isLoading) {
    return (
      <div className="space-y-3">
        <div className="h-10 animate-pulse rounded-md bg-secondary" />
        {[1, 2, 3].map((i) => (
          <div key={i} className="h-16 animate-pulse rounded-lg bg-secondary" />
        ))}
      </div>
    )
  }

  if (error) {
    return (
      <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4 text-sm text-destructive">
        データの取得に失敗しました。BACKEND_API_KEY が設定されているか確認してください。
      </div>
    )
  }

  return (
    <>
      <div className="space-y-4">
        {/* インライン作成フォーム */}
        <TodoCreate onSubmit={handleCreate} />

        {/* Todo リスト */}
        {todos.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="space-y-3">
            {todos.map((todo) => (
              <TodoCard
                key={todo.id}
                todo={todo}
                onToggle={handleToggle}
                onClick={setSelectedTodo}
              />
            ))}
          </div>
        )}
      </div>

      {/* 編集パネル */}
      <TodoPanel
        todo={selectedTodo}
        onClose={() => setSelectedTodo(null)}
        onUpdate={handleUpdate}
        onDelete={handleDelete}
      />
    </>
  )
}
