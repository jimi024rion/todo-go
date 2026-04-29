"use client"

import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { todoApi } from "@/lib/api"
import { TodoCard } from "./todo-card"
import { TodoCreate } from "./todo-create"
import { TodoPanel } from "./todo-panel"
import { EmptyState } from "./empty-state"
import type { Todo } from "@/types/todo"

const QUERY_KEY = ["todos"] as const

export function TodoList() {
  const queryClient = useQueryClient()
  const [selectedTodo, setSelectedTodo] = useState<Todo | null>(null)

  const { data: todos = [], isLoading, error } = useQuery({
    queryKey: QUERY_KEY,
    queryFn: todoApi.list,
  })

  // ── 作成 ──────────────────────────────────────────────────────────
  const createMutation = useMutation({
    mutationFn: todoApi.create,
    onMutate: async (input) => {
      await queryClient.cancelQueries({ queryKey: QUERY_KEY })
      const previous = queryClient.getQueryData<Todo[]>(QUERY_KEY)

      // 仮IDで楽観的に追加（リスト先頭に挿入）
      const optimistic: Todo = {
        id: `__optimistic__${Date.now()}`,
        title: input.title,
        description: input.description ?? "",
        status: "pending",
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) => [optimistic, ...old])
      return { previous }
    },
    onError: (_err, _input, ctx) => {
      // ロールバック
      queryClient.setQueryData(QUERY_KEY, ctx?.previous)
    },
    onSettled: () => {
      // 成功・失敗どちらでもサーバーと同期
      queryClient.invalidateQueries({ queryKey: QUERY_KEY })
    },
  })

  // ── 更新 ──────────────────────────────────────────────────────────
  const updateMutation = useMutation({
    mutationFn: ({ id, ...input }: { id: string } & Parameters<typeof todoApi.update>[1]) =>
      todoApi.update(id, input),
    onMutate: async ({ id, ...input }) => {
      await queryClient.cancelQueries({ queryKey: QUERY_KEY })
      const previous = queryClient.getQueryData<Todo[]>(QUERY_KEY)

      // キャッシュを楽観的に更新
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) =>
        old.map((t) =>
          t.id === id
            ? {
                ...t,
                title: input.title ?? t.title,
                description: input.description ?? t.description,
                status: input.status ?? t.status,
                updated_at: new Date().toISOString(),
              }
            : t
        )
      )

      // 開いているパネルも即時反映
      setSelectedTodo((prev) =>
        prev?.id === id
          ? {
              ...prev,
              title: input.title ?? prev.title,
              description: input.description ?? prev.description,
              status: input.status ?? prev.status,
            }
          : prev
      )

      return { previous }
    },
    onError: (_err, _input, ctx) => {
      queryClient.setQueryData(QUERY_KEY, ctx?.previous)
    },
    onSuccess: (serverTodo) => {
      // サーバーの正確なデータで上書き
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) =>
        old.map((t) => (t.id === serverTodo.id ? serverTodo : t))
      )
      setSelectedTodo((prev) => (prev?.id === serverTodo.id ? serverTodo : prev))
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: QUERY_KEY })
    },
  })

  // ── 削除 ──────────────────────────────────────────────────────────
  const deleteMutation = useMutation({
    mutationFn: todoApi.delete,
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: QUERY_KEY })
      const previous = queryClient.getQueryData<Todo[]>(QUERY_KEY)

      // 即時削除
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) =>
        old.filter((t) => t.id !== id)
      )
      // 削除対象がパネルで開いていたら閉じる
      setSelectedTodo((prev) => (prev?.id === id ? null : prev))

      return { previous }
    },
    onError: (_err, _id, ctx) => {
      queryClient.setQueryData(QUERY_KEY, ctx?.previous)
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: QUERY_KEY })
    },
  })

  // ── ハンドラー ────────────────────────────────────────────────────

  async function handleCreate(input: Parameters<typeof todoApi.create>[0]) {
    await createMutation.mutateAsync(input)
  }

  async function handleToggle(todo: Todo, done: boolean) {
    await updateMutation.mutateAsync({
      id: todo.id,
      title: todo.title,
      description: todo.description,
      status: done ? "completed" : "pending",
    })
  }

  async function handleUpdate(id: string, input: Parameters<typeof todoApi.update>[1]) {
    await updateMutation.mutateAsync({ id, ...input })
  }

  async function handleDelete(id: string) {
    await deleteMutation.mutateAsync(id)
  }

  // ── レンダー ──────────────────────────────────────────────────────

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
        <TodoCreate onSubmit={handleCreate} />

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

      <TodoPanel
        todo={selectedTodo}
        onClose={() => setSelectedTodo(null)}
        onUpdate={handleUpdate}
        onDelete={handleDelete}
      />
    </>
  )
}
