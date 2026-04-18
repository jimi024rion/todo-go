"use client";

import { useState } from "react";
import { deleteTodo, updateTodo } from "@/lib/api";
import { type Todo } from "@/types/todo";
import { EditTodoModal } from "./EditTodoModal";
import { StatusBadge } from "./StatusBadge";

const NEXT_STATUS: Record<Todo["Status"], Todo["Status"] | null> = {
  pending: "in_progress",
  in_progress: "completed",
  completed: null,
};

const NEXT_STATUS_LABEL: Record<Todo["Status"], string | null> = {
  pending: "進行中にする",
  in_progress: "完了にする",
  completed: null,
};

const STATUS_BORDER: Record<Todo["Status"], string> = {
  pending: "border-l-[#FCD34D]",
  in_progress: "border-l-[#60A5FA]",
  completed: "border-l-[#34D399]",
};

type Props = {
  todo: Todo;
  onUpdated: (updated: Todo) => void;
  onDeleted: (id: string) => void;
};

export function TodoCard({ todo, onUpdated, onDeleted }: Props) {
  const [loading, setLoading] = useState(false);
  const [editing, setEditing] = useState(false);

  const nextStatus = NEXT_STATUS[todo.Status];
  const isCompleted = todo.Status === "completed";

  const handleStatusChange = async () => {
    if (!nextStatus || loading) return;
    setLoading(true);
    try {
      const updated = await updateTodo(todo.ID, { status: nextStatus });
      onUpdated(updated);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (loading) return;
    setLoading(true);
    try {
      await deleteTodo(todo.ID);
      onDeleted(todo.ID);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <div
        className={`group rounded-lg border border-[#E7E5E0] border-l-4 ${STATUS_BORDER[todo.Status]} bg-white p-4 shadow-sm transition-shadow hover:shadow-md ${isCompleted ? "opacity-75" : ""}`}
      >
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2">
              <StatusBadge status={todo.Status} />
            </div>
            <p
              className={`mt-2 text-base font-medium leading-snug ${isCompleted ? "text-[#A8A29E] line-through" : "text-[#1C1917]"}`}
            >
              {todo.Title}
            </p>
            {todo.Description && (
              <p className="mt-1 text-sm leading-relaxed text-[#78716C]">
                {todo.Description}
              </p>
            )}
          </div>

          <div className="flex shrink-0 gap-1">
            <button
              onClick={() => setEditing(true)}
              disabled={loading}
              aria-label="編集"
              className="rounded-md p-2 text-[#A8A29E] transition hover:bg-[#F5F5F4] hover:text-[#78716C] disabled:opacity-50 min-h-[40px] min-w-[40px] flex items-center justify-center"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
              </svg>
            </button>
            <button
              onClick={handleDelete}
              disabled={loading}
              aria-label="削除"
              className="rounded-md p-2 text-[#A8A29E] transition hover:bg-[#FEF2F2] hover:text-[#EF4444] disabled:opacity-50 min-h-[40px] min-w-[40px] flex items-center justify-center"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M3 6h18" />
                <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
                <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
              </svg>
            </button>
          </div>
        </div>

        {nextStatus && (
          <div className="mt-3 flex justify-end">
            <button
              onClick={handleStatusChange}
              disabled={loading}
              className="rounded-md border border-[#E7E5E0] bg-white px-3 py-1.5 text-xs font-medium text-[#1C1917] transition hover:bg-[#F7F6F2] disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-[#F97316] focus:ring-offset-1 min-h-[36px]"
            >
              {loading ? "更新中…" : NEXT_STATUS_LABEL[todo.Status]}
            </button>
          </div>
        )}
      </div>

      {editing && (
        <EditTodoModal
          todo={todo}
          onClose={() => setEditing(false)}
          onSaved={(updated) => {
            onUpdated(updated);
            setEditing(false);
          }}
        />
      )}
    </>
  );
}
