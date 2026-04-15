"use client";

import { useCallback, useEffect, useState } from "react";
import { listTodos } from "@/lib/api";
import { type Todo } from "@/types/todo";
import { TodoCard } from "./TodoCard";
import { TodoForm } from "./TodoForm";

type Props = {
  // サーバーサイドで取得済みの初期データ（SSRで渡される）
  // 指定された場合、初期ローディングスピナーを表示しない
  initialTodos?: Todo[];
};

export function TodoList({ initialTodos }: Props) {
  const [todos, setTodos] = useState<Todo[]>(initialTodos ?? []);
  const [loading, setLoading] = useState(initialTodos === undefined);
  const [error, setError] = useState<string | null>(null);

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
    // initialTodosが渡された場合はSSRで取得済みなので再取得しない
    if (initialTodos !== undefined) return;
    fetchTodos();
  }, [fetchTodos, initialTodos]);

  const handleUpdated = (updated: Todo) => {
    setTodos((prev) => prev.map((t) => (t.ID === updated.ID ? updated : t)));
  };

  const handleDeleted = (id: string) => {
    setTodos((prev) => prev.filter((t) => t.ID !== id));
  };

  return (
    <div className="space-y-6">
      <TodoForm onCreated={fetchTodos} />

      <div className="h-px bg-[#F3F4F6]" />

      {loading && (
        <div className="flex justify-center py-12">
          <svg
            className="h-6 w-6 animate-spin text-[#4F46E5]"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
            />
          </svg>
        </div>
      )}

      {!loading && error && (
        <p className="text-center text-sm text-[#EF4444]">{error}</p>
      )}

      {!loading && !error && todos.length === 0 && (
        <p className="mt-12 text-center text-sm text-[#9CA3AF]">
          タスクはまだありません
        </p>
      )}

      {!loading && !error && todos.length > 0 && (
        <div className="space-y-2">
          {todos.map((todo) => (
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
  );
}
