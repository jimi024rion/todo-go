"use client";

import { useState } from "react";
import { deleteTodo, updateTodo } from "@/lib/api";
import { type Todo } from "@/types/todo";
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

type Props = {
  todo: Todo;
  onUpdated: (updated: Todo) => void;
  onDeleted: (id: string) => void;
};

export function TodoCard({ todo, onUpdated, onDeleted }: Props) {
  const [loading, setLoading] = useState(false);

  const nextStatus = NEXT_STATUS[todo.Status];

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

  const isCompleted = todo.Status === "completed";

  return (
    <div
      className={`group rounded-lg border border-[#E5E7EB] bg-white p-4 transition hover:border-[#D1D5DB] hover:shadow-[0_1px_3px_rgba(0,0,0,0.07),0_1px_2px_rgba(0,0,0,0.06)] ${isCompleted ? "bg-[#F9FAFB]" : ""}`}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-center gap-2">
            <StatusBadge status={todo.Status} />
          </div>
          <p
            className={`mt-2 text-base font-medium leading-snug ${isCompleted ? "text-[#9CA3AF] line-through" : "text-[#111827]"}`}
          >
            {todo.Title}
          </p>
          {todo.Description && (
            <p className="mt-1 text-sm leading-relaxed text-[#4B5563]">
              {todo.Description}
            </p>
          )}
        </div>
        <button
          onClick={handleDelete}
          disabled={loading}
          aria-label="削除"
          className="shrink-0 rounded-md p-1.5 text-[#9CA3AF] transition hover:bg-[#FEF2F2] hover:text-[#EF4444] disabled:opacity-50"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M3 6h18" />
            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
          </svg>
        </button>
      </div>

      {nextStatus && (
        <div className="mt-3 flex justify-end">
          <button
            onClick={handleStatusChange}
            disabled={loading}
            className="rounded-md border border-[#E5E7EB] bg-white px-3 py-1.5 text-xs font-medium text-[#374151] transition hover:bg-[#F9FAFB] disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-[#4F46E5] focus:ring-offset-1"
          >
            {loading ? "更新中…" : NEXT_STATUS_LABEL[todo.Status]}
          </button>
        </div>
      )}
    </div>
  );
}
