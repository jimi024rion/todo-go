"use client";

import { useEffect, useRef, useState } from "react";
import { updateTodo } from "@/lib/api";
import { type Todo } from "@/types/todo";

type Props = {
  todo: Todo;
  onClose: () => void;
  onSaved: (updated: Todo) => void;
};

export function EditTodoModal({ todo, onClose, onSaved }: Props) {
  const [title, setTitle] = useState(todo.Title);
  const [description, setDescription] = useState(todo.Description);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const titleRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    titleRef.current?.focus();
  }, []);

  function handleOverlayClick(e: React.MouseEvent) {
    if (e.target === e.currentTarget) onClose();
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (!title.trim()) {
      setError("タイトルを入力してください");
      return;
    }

    setLoading(true);
    try {
      const updated = await updateTodo(todo.ID, {
        title: title.trim(),
        description: description.trim(),
      });
      onSaved(updated);
      onClose();
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存に失敗しました");
    } finally {
      setLoading(false);
    }
  }

  const inputClass =
    "w-full px-3 py-2 text-sm text-[#1C1917] bg-white border border-[#D6D3D1] rounded-[6px] placeholder-[#A8A29E] " +
    "focus:outline-none focus:border-[#F97316] focus:ring-2 focus:ring-[rgba(249,115,22,0.15)] transition-colors";

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm px-4"
      onClick={handleOverlayClick}
    >
      <div className="w-full max-w-[480px] bg-white rounded-xl shadow-2xl p-6">
        <h2 className="text-lg font-semibold text-[#1C1917] mb-5">タスクを編集</h2>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-[#1C1917]">タイトル</label>
            <input
              ref={titleRef}
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              maxLength={100}
              disabled={loading}
              className={inputClass}
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-[#1C1917]">説明</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={4}
              maxLength={1000}
              disabled={loading}
              className={inputClass + " resize-none"}
            />
          </div>

          {error && <p className="text-sm text-[#EF4444]">{error}</p>}

          <div className="flex justify-end gap-2 mt-1">
            <button
              type="button"
              onClick={onClose}
              disabled={loading}
              className="px-4 py-2 text-sm font-medium text-[#1C1917] bg-white border border-[#E7E5E0] rounded-[6px] hover:bg-[#F7F6F2] disabled:opacity-60 transition-colors min-h-[44px]"
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 text-sm font-medium text-white bg-[#F97316] rounded-[6px] hover:bg-[#EA580C] disabled:opacity-60 transition-colors min-h-[44px]"
            >
              {loading ? "保存中..." : "保存"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
