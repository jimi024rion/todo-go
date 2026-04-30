"use client"

import { useRef, useState } from "react"
import { Plus, X } from "lucide-react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { cn } from "@/lib/utils"
import { tagApi } from "@/lib/api"
import type { CreateTodoInput } from "@/types/todo"

interface TodoCreateProps {
  onSubmit: (input: CreateTodoInput & { tagIds: string[] }) => Promise<void>
}

export function TodoCreate({ onSubmit }: TodoCreateProps) {
  const [title, setTitle] = useState("")
  const [description, setDescription] = useState("")
  const [selectedTagIds, setSelectedTagIds] = useState<string[]>([])
  const [dueDate, setDueDate] = useState("")
  const [expanded, setExpanded] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [showNewTagInput, setShowNewTagInput] = useState(false)
  const [newTagName, setNewTagName] = useState("")
  const queryClient = useQueryClient()
  const titleRef = useRef<HTMLInputElement>(null)
  const descriptionRef = useRef<HTMLTextAreaElement>(null)
  const formRef = useRef<HTMLDivElement>(null)

  const { data: allTags = [] } = useQuery({
    queryKey: ["tags"],
    queryFn: tagApi.list,
    enabled: expanded,
  })

  const createTagMutation = useMutation({
    mutationFn: tagApi.create,
    onSuccess: (newTag) => {
      queryClient.invalidateQueries({ queryKey: ["tags"] })
      setSelectedTagIds((prev) => [...prev, newTag.id])
      setNewTagName("")
      setShowNewTagInput(false)
    },
  })

  function handleFocus() {
    setExpanded(true)
  }

  function handleBlur(e: React.FocusEvent) {
    if (formRef.current?.contains(e.relatedTarget as Node)) return
    if (!title.trim() && !description.trim() && selectedTagIds.length === 0) {
      setExpanded(false)
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const trimmedTitle = title.trim()
    if (!trimmedTitle || isSubmitting) return

    // フォームを即時クリア（楽観的更新と同期）
    const submittedTitle = trimmedTitle
    const submittedDescription = description.trim() || undefined
    const submittedTagIds = [...selectedTagIds]
    const submittedDueDate = dueDate ? new Date(dueDate).toISOString() : undefined
    setTitle("")
    setDescription("")
    setSelectedTagIds([])
    setDueDate("")
    setExpanded(false)

    setIsSubmitting(true)
    try {
      await onSubmit({ title: submittedTitle, description: submittedDescription, due_date: submittedDueDate, tagIds: submittedTagIds })
    } catch {
      setTitle(submittedTitle)
      setDescription(submittedDescription ?? "")
      setSelectedTagIds(submittedTagIds)
      setDueDate(dueDate)
      setExpanded(true)
    } finally {
      setIsSubmitting(false)
    }
  }

  function handleCancel() {
    setTitle("")
    setDescription("")
    setSelectedTagIds([])
    setDueDate("")
    setShowNewTagInput(false)
    setNewTagName("")
    setExpanded(false)
    titleRef.current?.blur()
  }

  function toggleTag(id: string) {
    setSelectedTagIds((prev) =>
      prev.includes(id) ? prev.filter((t) => t !== id) : [...prev, id]
    )
  }

  function handleTitleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === "Escape") { handleCancel(); return }
    if (e.key === "Enter") { e.preventDefault(); descriptionRef.current?.focus(); return }
    if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) { e.preventDefault(); handleSubmit(e as unknown as React.FormEvent) }
  }

  function handleDescriptionKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === "Escape") { handleCancel(); return }
    if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) { e.preventDefault(); handleSubmit(e as unknown as React.FormEvent) }
  }

  return (
    <div ref={formRef} onBlur={handleBlur}>
      <form onSubmit={handleSubmit}>
        <div className={cn(
          "rounded-lg border bg-card shadow-sm transition-all duration-150",
          expanded ? "border-primary/50 shadow-md" : "border-border"
        )}>
          {/* タイトル入力 */}
          <div className="flex items-center gap-2 px-3 py-2.5">
            <Plus className={cn("h-4 w-4 shrink-0 transition-colors duration-150", expanded ? "text-primary" : "text-muted-foreground")} />
            <input
              ref={titleRef}
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              onFocus={handleFocus}
              onKeyDown={handleTitleKeyDown}
              placeholder="タスクを追加..."
              disabled={isSubmitting}
              className={cn("flex-1 bg-transparent text-base text-foreground", "placeholder:text-muted-foreground focus:outline-none disabled:opacity-50")}
            />
          </div>

          {expanded && (
            <>
              {/* 説明欄 */}
              <div className="border-t border-border/60 px-3 py-2">
                <textarea
                  ref={descriptionRef}
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  onKeyDown={handleDescriptionKeyDown}
                  placeholder="説明を追加（任意）"
                  rows={2}
                  disabled={isSubmitting}
                  className={cn("w-full resize-none bg-transparent text-base text-foreground", "placeholder:text-muted-foreground focus:outline-none disabled:opacity-50")}
                />
              </div>

              {/* 期限設定 */}
              <div className="border-t border-border/60 px-3 py-2">
                <p className="mb-1.5 text-xs text-muted-foreground">期限（任意）</p>
                <input
                  type="date"
                  value={dueDate}
                  onChange={(e) => setDueDate(e.target.value)}
                  min={new Date().toISOString().split("T")[0]}
                  disabled={isSubmitting}
                  className={cn(
                    "rounded-md border border-border bg-background px-2.5 py-1 text-sm text-foreground",
                    "focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50"
                  )}
                />
              </div>

              {/* タグ選択 + 新規作成 */}
              <div className="border-t border-border/60 px-3 py-2">
                <p className="mb-1.5 text-xs text-muted-foreground">タグ</p>
                <div className="flex flex-wrap gap-1.5">
                  {allTags.map((tag) => {
                    const selected = selectedTagIds.includes(tag.id)
                    return (
                      <button
                        key={tag.id}
                        type="button"
                        onClick={() => toggleTag(tag.id)}
                        className={cn(
                          "inline-flex items-center rounded-full px-2.5 py-1 text-xs transition-colors duration-150",
                          selected
                            ? "bg-foreground/10 text-foreground font-medium"
                            : "bg-secondary text-muted-foreground hover:bg-secondary/80"
                        )}
                      >
                        {selected && <span className="mr-1">✓</span>}
                        {tag.name}
                      </button>
                    )
                  })}
                  {/* 新規タグ作成ボタン */}
                  {!showNewTagInput && (
                    <button
                      type="button"
                      onClick={() => setShowNewTagInput(true)}
                      className="inline-flex items-center gap-1 rounded-full border border-dashed border-border px-2.5 py-1 text-xs text-muted-foreground hover:border-primary hover:text-foreground transition-colors"
                    >
                      <Plus className="h-3 w-3" />
                      新規
                    </button>
                  )}
                </div>
                {showNewTagInput && (
                  <div className="mt-2 flex items-center gap-2">
                    <input
                      autoFocus
                      type="text"
                      value={newTagName}
                      onChange={(e) => setNewTagName(e.target.value)}
                      onKeyDown={(e) => {
                        if (e.key === "Enter") { e.preventDefault(); if (newTagName.trim()) createTagMutation.mutate(newTagName.trim()) }
                        if (e.key === "Escape") { setShowNewTagInput(false); setNewTagName("") }
                      }}
                      placeholder="タグ名..."
                      className={cn(
                        "flex-1 rounded-md border border-border bg-background px-2.5 py-1 text-base",
                        "focus:outline-none focus:ring-2 focus:ring-ring"
                      )}
                    />
                    <button
                      type="button"
                      onClick={() => { if (newTagName.trim()) createTagMutation.mutate(newTagName.trim()) }}
                      disabled={!newTagName.trim()}
                      className="rounded-md bg-primary px-2.5 py-1 text-xs text-primary-foreground disabled:opacity-40"
                    >
                      追加
                    </button>
                    <button
                      type="button"
                      onClick={() => { setShowNewTagInput(false); setNewTagName("") }}
                      className="text-muted-foreground hover:text-foreground"
                    >
                      <X className="h-4 w-4" />
                    </button>
                  </div>
                )}
              </div>

              {/* フッター */}
              <div className="flex items-center justify-between border-t border-border/60 px-3 py-2">
                <p className="hidden md:block text-xs text-muted-foreground">Ctrl+Enter で追加</p>
                <p className="md:hidden text-xs invisible" aria-hidden="true" />
                <div className="flex items-center gap-2">
                  <button type="button" onClick={handleCancel} disabled={isSubmitting}
                    className={cn("rounded-md px-3 py-1.5 text-sm text-muted-foreground", "hover:bg-secondary hover:text-foreground transition-colors duration-150 disabled:opacity-50")}>
                    キャンセル
                  </button>
                  <button type="submit" disabled={!title.trim() || isSubmitting}
                    className={cn("rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground", "hover:bg-primary/90 transition-colors duration-150 disabled:opacity-40 disabled:cursor-not-allowed")}>
                    {isSubmitting ? "追加中..." : "追加"}
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </form>
    </div>
  )
}
