"use client"

import { useEffect, useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { ArrowUpDown } from "lucide-react"
import {
  DndContext,
  closestCenter,
  PointerSensor,
  TouchSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core"
import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
  arrayMove,
} from "@dnd-kit/sortable"
import { CSS } from "@dnd-kit/utilities"
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
      case "created_desc": return b.created_at.localeCompare(a.created_at)
      case "created_asc":  return a.created_at.localeCompare(b.created_at)
      case "status": {
        const order = { pending: 0, in_progress: 1, completed: 2 }
        return (order[a.status] ?? 0) - (order[b.status] ?? 0)
      }
      case "title": return a.title.localeCompare(b.title, "ja")
    }
  })
}

// ── DnD ドラッグ可能なカードラッパー ──────────────────────────────────
// カード全体を長押し（モバイル 300ms / デスクトップ 8px 移動）でドラッグ開始

function SortableTodoCard(props: React.ComponentProps<typeof TodoCard>) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({ id: props.todo.id })

  return (
    <div
      ref={setNodeRef}
      style={{ transform: CSS.Transform.toString(transform), transition }}
      className={isDragging ? "opacity-50 z-50 relative cursor-grabbing" : "cursor-grab"}
      {...attributes}
      {...listeners}
    >
      <TodoCard {...props} />
    </div>
  )
}

// ── メインコンポーネント ──────────────────────────────────────────────

export function TodoList() {
  const queryClient = useQueryClient()
  const [selectedTodo, setSelectedTodo] = useState<Todo | null>(null)
  const [sortKey, setSortKey] = useState<SortKey>("created_desc")
  // 表示順序をローカルで管理（DnD でドラッグすると更新）
  const [orderedIds, setOrderedIds] = useState<string[]>([])

  const { data: todos = [], isLoading, error } = useQuery({
    queryKey: QUERY_KEY,
    queryFn: todoApi.list,
  })

  // サーバーからデータが届いたとき、orderedIds を初期化
  useEffect(() => {
    if (todos.length > 0 && orderedIds.length === 0) {
      setOrderedIds(sortTodos(todos, sortKey).map((t) => t.id))
    }
  }, [todos]) // eslint-disable-line react-hooks/exhaustive-deps

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
    // モバイル: 300ms 長押しでドラッグ開始。横スワイプ（削除）との競合を避けるため tolerance を広めに
    useSensor(TouchSensor, { activationConstraint: { delay: 300, tolerance: 10 } })
  )

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event
    if (!over || active.id === over.id) return
    setOrderedIds((ids) => {
      const oldIndex = ids.indexOf(active.id as string)
      const newIndex = ids.indexOf(over.id as string)
      return arrayMove(ids, oldIndex, newIndex)
    })
  }

  // ソートボタンを押したとき → 指定した順序を orderedIds に適用
  function applySort(key: SortKey) {
    setSortKey(key)
    setOrderedIds(sortTodos(todos, key).map((t) => t.id))
  }

  // orderedIds に従って表示リストを構築
  function getDisplayList(): Todo[] {
    if (orderedIds.length === 0) return sortTodos(todos, sortKey)
    const todoMap = new Map(todos.map((t) => [t.id, t]))
    const ordered = orderedIds.flatMap((id) => todoMap.has(id) ? [todoMap.get(id)!] : [])
    // orderedIds にない新しい todo（楽観的追加など）を末尾に追加
    const inOrder = new Set(orderedIds)
    const rest = todos.filter((t) => !inOrder.has(t.id))
    return [...ordered, ...rest]
  }

  // ── Mutations ─────────────────────────────────────────────────────

  const createMutation = useMutation({
    mutationFn: todoApi.create,
    onMutate: async (input) => {
      await queryClient.cancelQueries({ queryKey: QUERY_KEY })
      const previous = queryClient.getQueryData<Todo[]>(QUERY_KEY)
      const optimisticId = `__optimistic__${Date.now()}`
      const optimistic: Todo = {
        id: optimisticId,
        title: input.title,
        description: input.description ?? "",
        status: "pending",
        tags: [],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) => [optimistic, ...old])
      setOrderedIds((ids) => [optimisticId, ...ids])
      return { previous }
    },
    onError: (_e, _i, ctx) => {
      queryClient.setQueryData(QUERY_KEY, ctx?.previous)
      setOrderedIds((ids) => ids.filter((id) => !id.startsWith("__optimistic__")))
    },
    onSettled: () => queryClient.invalidateQueries({ queryKey: QUERY_KEY }),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, ...input }: { id: string } & Parameters<typeof todoApi.update>[1]) =>
      todoApi.update(id, input),
    onMutate: async ({ id, ...input }) => {
      await queryClient.cancelQueries({ queryKey: QUERY_KEY })
      const previous = queryClient.getQueryData<Todo[]>(QUERY_KEY)
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) =>
        old.map((t) => t.id === id ? {
          ...t,
          title: input.title ?? t.title,
          description: input.description ?? t.description,
          status: input.status ?? t.status,
          updated_at: new Date().toISOString(),
        } : t)
      )
      setSelectedTodo((prev) => prev?.id === id ? {
        ...prev,
        title: input.title ?? prev.title,
        description: input.description ?? prev.description,
        status: input.status ?? prev.status,
      } : prev)
      return { previous }
    },
    onError: (_e, _i, ctx) => queryClient.setQueryData(QUERY_KEY, ctx?.previous),
    onSuccess: (serverTodo) => {
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) =>
        old.map((t) => t.id === serverTodo.id ? serverTodo : t)
      )
      setSelectedTodo((prev) => prev?.id === serverTodo.id ? serverTodo : prev)
    },
    onSettled: () => queryClient.invalidateQueries({ queryKey: QUERY_KEY }),
  })

  const deleteMutation = useMutation({
    mutationFn: todoApi.delete,
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: QUERY_KEY })
      const previous = queryClient.getQueryData<Todo[]>(QUERY_KEY)
      queryClient.setQueryData<Todo[]>(QUERY_KEY, (old = []) => old.filter((t) => t.id !== id))
      setOrderedIds((ids) => ids.filter((i) => i !== id))
      setSelectedTodo((prev) => prev?.id === id ? null : prev)
      return { previous }
    },
    onError: (_e, _i, ctx) => queryClient.setQueryData(QUERY_KEY, ctx?.previous),
    onSettled: () => queryClient.invalidateQueries({ queryKey: QUERY_KEY }),
  })

  // ── ハンドラー ────────────────────────────────────────────────────

  async function handleCreate(input: Parameters<typeof todoApi.create>[0]) {
    await createMutation.mutateAsync(input)
  }

  async function handleToggle(todo: Todo, done: boolean) {
    await updateMutation.mutateAsync({ id: todo.id, title: todo.title, description: todo.description, status: done ? "completed" : "pending" })
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

  const displayList = getDisplayList()

  return (
    <>
      <div className="space-y-4">
        <TodoCreate onSubmit={handleCreate} />

        {/* ソートボタン（初期順序を適用する） */}
        {todos.length > 1 && (
          <div className="flex items-center justify-end gap-1.5">
            <ArrowUpDown className="h-3.5 w-3.5 text-muted-foreground" />
            <div className="flex gap-1">
              {SORT_OPTIONS.map((opt) => (
                <button
                  key={opt.value}
                  onClick={() => applySort(opt.value)}
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

        {/* Todo リスト（常に DnD 有効） */}
        {todos.length === 0 ? (
          <EmptyState />
        ) : (
          <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
            <SortableContext items={displayList.map((t) => t.id)} strategy={verticalListSortingStrategy}>
              <div className="space-y-3">
                {displayList.map((todo) => (
                  <SortableTodoCard
                    key={todo.id}
                    todo={todo}
                    onToggle={handleToggle}
                    onClick={setSelectedTodo}
                    onDelete={handleDelete}
                  />
                ))}
              </div>
            </SortableContext>
          </DndContext>
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
