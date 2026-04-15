"use client";

import { useState } from "react";
import { createTodo } from "@/lib/api";

type Props = {
  onCreated: () => void;
};

export function TodoForm({ onCreated }: Props) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmedTitle = title.trim();
    if (!trimmedTitle) return;

    setLoading(true);
    setError(null);
    try {
      await createTodo({
        title: trimmedTitle,
        description: description.trim() || undefined,
      });
      setTitle("");
      setDescription("");
      onCreated();
    } catch (err) {
      setError(err instanceof Error ? err.message : "作成に失敗しました");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="rounded-lg border border-[#E5E7EB] bg-white p-5">
      <h2 className="text-sm font-semibold text-[#111827]">新しいタスクを追加</h2>
      <form onSubmit={handleSubmit} className="mt-3 space-y-3">
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="タイトル"
          maxLength={100}
          required
          className="w-full rounded-md border border-[#D1D5DB] bg-white px-3 py-2 text-sm text-[#111827] placeholder-[#9CA3AF] outline-none transition focus:border-[#4F46E5] focus:ring-2 focus:ring-[rgba(79,70,229,0.15)]"
        />
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="説明（任意）"
          maxLength={1000}
          rows={2}
          className="w-full resize-none rounded-md border border-[#D1D5DB] bg-white px-3 py-2 text-sm text-[#111827] placeholder-[#9CA3AF] outline-none transition focus:border-[#4F46E5] focus:ring-2 focus:ring-[rgba(79,70,229,0.15)]"
        />
        {error && <p className="text-xs text-[#EF4444]">{error}</p>}
        <button
          type="submit"
          disabled={!title.trim() || loading}
          className="w-full rounded-md bg-[#4F46E5] px-4 py-2 text-sm font-medium text-white transition hover:bg-[#4338CA] disabled:cursor-not-allowed disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-[#4F46E5] focus:ring-offset-2"
        >
          {loading ? "追加中…" : "追加する"}
        </button>
      </form>
    </div>
  );
}
