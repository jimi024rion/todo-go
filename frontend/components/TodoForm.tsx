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

  const inputClass =
    "w-full rounded-[6px] border border-[#D6D3D1] bg-white px-3 py-2 text-sm text-[#1C1917] placeholder-[#A8A29E] outline-none transition focus:border-[#F97316] focus:ring-2 focus:ring-[rgba(249,115,22,0.15)]";

  return (
    <div className="rounded-lg border border-[#E7E5E0] bg-white p-5 shadow-sm">
      <h2 className="text-sm font-semibold text-[#1C1917]">新しいタスクを追加</h2>
      <form onSubmit={handleSubmit} className="mt-3 space-y-3">
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="タイトル"
          maxLength={100}
          required
          className={inputClass}
        />
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="説明（任意）"
          maxLength={1000}
          rows={2}
          className={inputClass + " resize-none"}
        />
        {error && <p className="text-xs text-[#EF4444]">{error}</p>}
        <button
          type="submit"
          disabled={!title.trim() || loading}
          className="w-full rounded-[6px] bg-[#F97316] px-4 py-2.5 text-sm font-medium text-white transition hover:bg-[#EA580C] disabled:cursor-not-allowed disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-[#F97316] focus:ring-offset-2 min-h-[44px]"
        >
          {loading ? "追加中…" : "追加する"}
        </button>
      </form>
    </div>
  );
}
