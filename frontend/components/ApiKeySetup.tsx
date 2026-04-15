"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

// onSave はキーを受け取らない（キーはサーバー側のCookieに保存されるため）
type Props = {
  onSave?: () => void;
};

export function ApiKeySetup({ onSave }: Props) {
  const [value, setValue] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = value.trim();
    if (!trimmed) return;

    setLoading(true);
    setError(null);
    try {
      // POST /api/auth → サーバー側で httpOnly Cookie をセット
      const res = await fetch("/api/auth", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ key: trimmed }),
      });
      if (!res.ok) throw new Error("保存に失敗しました");

      onSave?.();
      // router.refresh() でServer Componentを再実行させる
      // → page.tsxがCookieを読んでTodoListを表示する
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm px-4">
      <div className="w-full max-w-sm rounded-xl bg-white p-8 shadow-[0_20px_25px_-5px_rgba(0,0,0,0.1),0_8px_10px_-6px_rgba(0,0,0,0.1)]">
        <h1 className="text-xl font-bold tracking-tight text-[#111827]">
          Todo アプリへようこそ
        </h1>
        <p className="mt-2 text-sm text-[#4B5563]">
          このアプリを使用するには API キーが必要です。
          バックエンドで発行したキーを入力してください。
        </p>
        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <input
            type="text"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder="todo_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
            className="w-full rounded-md border border-[#D1D5DB] bg-white px-3 py-2 text-sm text-[#111827] placeholder-[#9CA3AF] outline-none transition focus:border-[#4F46E5] focus:ring-2 focus:ring-[rgba(79,70,229,0.15)]"
          />
          {error && <p className="text-xs text-[#EF4444]">{error}</p>}
          <button
            type="submit"
            disabled={!value.trim() || loading}
            className="w-full rounded-md bg-[#4F46E5] px-4 py-2 text-sm font-medium text-white transition hover:bg-[#4338CA] disabled:cursor-not-allowed disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-[#4F46E5] focus:ring-offset-2"
          >
            {loading ? "保存中…" : "保存して開始する"}
          </button>
        </form>
      </div>
    </div>
  );
}
