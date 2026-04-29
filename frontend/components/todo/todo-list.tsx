"use client"

import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { ArrowUpDown } from "lucide-react"
import { todoApi } from "@/lib/api"
import { TodoCard } from "./todo-card"
import { TodoCreate } from "./todo-create"
import { TodoPanel } from "./todo-panel"
import { EmptyState } from "./empty-state"
import type { Todo } from "@/types/todo"

const QUERY_KEY = ["todos"] as const

type SortKey = "created_desc" | "created_asc" | "status" | "title"

const SORT_OPTIONS: { value: SortKey; label: string }[] = [
  { value: "created_desc", label: "新しい順" },
  { value: "created_asc",  label: "古い順" },
  { value: "status",       label: "未完了優先" },
  { value: "title",        label: "タイトル順" },
]

function sortTodos(todos: Todo[], key: SortKey): Todo[] {
  return [...todos].sort((a, b) => {
    switch (key) {
      case "created_desc":
        return b.created_at.localeCompare(a.created_at)
      case "created_asc":
        return a.created_at.localeCompare(b.created_at)
      case "status": {
        // pending → in_progress → completed の順
        const order = { pending: 0, in_progress: 1, completed: 2 }
        return (order[a.status] ?? 0) - (order[b.status] ?? 0)
      }
      case "title":
        return a.title.localeCompare(b.title, "ja")
    }
  })
}

export function TodoList() {
  const queryClient = useQueryClient()
  const [selectedTodo, setSelectedTodo] = useState<Todo | null>(null)
  const [sortKey, setSortKey] = useState<SortKey>("created_desc")

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
      const optimistic: Todo = {
        id: `__optimistic__${Date.now()}`,
        title: input.title,
        description: input.description ?? "",
        status: "pending",
        tags: [],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) => [optimistic, ...old])
      return { previous }
    },
    onError: (_err, _input, ctx) => {
      queryClient.setQueryData(QUERY_KEY, ctx?.previous)
    },
    onSettled: () => {
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
      setSelectedTodo((prev) =>
        prev?.id === id
          ? { ...prev, title: input.title ?? prev.title, description: input.description ?? prev.description, status: input.status ?? prev.status }
          : prev
      )
      return { previous }
    },
    onError: (_err, _input, ctx) => {
      queryClient.setQueryData(QUERY_KEY, ctx?.previous)
    },
    onSuccess: (serverTodo) => {
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
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) =>
        old.filter((t) => t.id !== id)
      )
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

  const sorted = sortTodos(todos, sortKey)

  return (
    <>
      <div className="space-y-4">
        {/* 作成フォーム */}
        <TodoCreate onSubmit={handleCreate} />

        {/* ソートセレクター */}
        {todos.length > 1 && (
          <div className="flex items-center justify-end gap-1.5">
            <ArrowUpDown className="h-3.5 w-3.5 text-muted-foreground" />
            <div className="flex gap-1">
              {SORT_OPTIONS.map((opt) => (
                <button
                  key={opt.value}
                  onClick={() => setSortKey(opt.value)}
                  className={
                    sortKey === opt.value
                      ? "rounded px-2 py-1 text-xs font-medium bg-secondary text-foreground"
                      : "rounded px-2 py-1 text-xs text-muted-foreground hover:bg-secondary/60 transition-colors"
                  }
                >
                  {opt.label}
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Todo リスト */}
        {todos.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="space-y-3">
            {sorted.map((todo) => (
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
