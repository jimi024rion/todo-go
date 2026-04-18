"use client";

import { useCallback, useEffect, useState } from "react";
import { listTodos } from "@/lib/api";
import { type Todo } from "@/types/todo";
import { TodoCard } from "./TodoCard";
import { TodoForm } from "./TodoForm";

type Filter = "all" | "pending" | "in_progress";

const FILTER_LABELS: Record<Filter, string> = {
  all: "すべて",
  pending: "未着手",
  in_progress: "進行中",
};

type Props = {
  initialTodos?: Todo[];
};

export function TodoList({ initialTodos }: Props) {
  const [todos, setTodos] = useState<Todo[]>(initialTodos ?? []);
  const [loading, setLoading] = useState(initialTodos === undefined);
  const [error, setError] = useState<string | null>(null);
  const [activeFilter, setActiveFilter] = useState<Filter>("all");
  const [completedOpen, setCompletedOpen] = useState(false);

  const fetchTodos = useCallback(async () => {
    setError(null);
    try {
      const data = await listTodos();
      setTodos(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "取得に失敗しました");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (initialTodos !== undefined) return;
    fetchTodos();
  }, [fetchTodos, initialTodos]);

  const handleUpdated = (updated: Todo) => {
    setTodos((prev) => prev.map((t) => (t.ID === updated.ID ? updated : t)));
  };

  const handleDeleted = (id: string) => {
    setTodos((prev) => prev.filter((t) => t.ID !== id));
  };

  const activeTodos = todos.filter((t) => t.Status !== "completed");
  const completedTodos = todos.filter((t) => t.Status === "completed");

  const filteredTodos =
    activeFilter === "all"
      ? activeTodos
      : activeTodos.filter((t) => t.Status === activeFilter);

  return (
    <div className="space-y-5">
      <TodoForm onCreated={fetchTodos} />

      {/* フィルタータブ */}
      <div className="border-b border-[#E7E5E0] flex gap-1">
        {(Object.keys(FILTER_LABELS) as Filter[]).map((f) => (
          <button
            key={f}
            onClick={() => setActiveFilter(f)}
            className={`px-4 py-2.5 text-sm transition-colors border-b-2 -mb-px ${
              activeFilter === f
                ? "border-[#F97316] text-[#F97316] font-semibold"
                : "border-transparent text-[#78716C] hover:text-[#1C1917]"
            }`}
          >
            {FILTER_LABELS[f]}
          </button>
        ))}
      </div>

      {loading && (
        <div className="flex justify-center py-12">
          <svg
            className="h-6 w-6 animate-spin text-[#F97316]"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
        </div>
      )}

      {!loading && error && (
        <p className="text-center text-sm text-[#EF4444]">{error}</p>
      )}

      {!loading && !error && (
        <>
          {filteredTodos.length === 0 ? (
            <p className="mt-12 text-center text-sm text-[#A8A29E]">
              {activeFilter === "all" ? "タスクはまだありません" : "該当するタスクはありません"}
            </p>
          ) : (
            <div className="space-y-2">
              {filteredTodos.map((todo) => (
                <TodoCard
                  key={todo.ID}
                  todo={todo}
                  onUpdated={handleUpdated}
                  onDeleted={handleDeleted}
                />
              ))}
            </div>
          )}

          {/* 完了済みセクション */}
          {completedTodos.length > 0 && (
            <div className="border-t border-[#E7E5E0] pt-4 mt-4">
              <button
                onClick={() => setCompletedOpen((v) => !v)}
                className="flex w-full items-center gap-2 text-sm font-medium text-[#78716C] hover:text-[#1C1917] transition-colors py-1"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className={`transition-transform ${completedOpen ? "rotate-90" : ""}`}
                >
                  <polyline points="9 18 15 12 9 6" />
                </svg>
                完了済み ({completedTodos.length})
              </button>

              {completedOpen && (
                <div className="mt-2 space-y-2">
                  {completedTodos.map((todo) => (
                    <TodoCard
                      key={todo.ID}
                      todo={todo}
                      onUpdated={handleUpdated}
                      onDeleted={handleDeleted}
                    />
                  ))}
                </div>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
}
