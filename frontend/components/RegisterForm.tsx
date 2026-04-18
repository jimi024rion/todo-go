"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

export default function RegisterForm() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (!name.trim()) {
      setError("名前を入力してください");
      return;
    }
    if (!email.trim() || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setError("有効なメールアドレスを入力してください");
      return;
    }

    setLoading(true);
    try {
      const res = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: name.trim(), email: email.trim() }),
      });
      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.error ?? "登録に失敗しました");
      }
      router.push("/");
      router.refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "登録に失敗しました");
    } finally {
      setLoading(false);
    }
  }

  const inputClass =
    "w-full px-3 py-2 text-sm text-[#1C1917] bg-white border border-[#D6D3D1] rounded-[6px] placeholder-[#A8A29E] " +
    "focus:outline-none focus:border-[#F97316] focus:ring-2 focus:ring-[rgba(249,115,22,0.15)] transition-colors";

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <div className="flex flex-col gap-1.5">
        <label className="text-sm font-medium text-[#1C1917]">名前</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="山田 太郎"
          maxLength={100}
          disabled={loading}
          className={inputClass}
        />
      </div>

      <div className="flex flex-col gap-1.5">
        <label className="text-sm font-medium text-[#1C1917]">メールアドレス</label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="taro@example.com"
          disabled={loading}
          className={inputClass}
        />
      </div>

      {error && (
        <p className="text-sm text-[#EF4444]">{error}</p>
      )}

      <button
        type="submit"
        disabled={loading}
        className="mt-2 w-full py-2.5 px-4 text-sm font-medium text-white bg-[#F97316] rounded-[6px] hover:bg-[#EA580C] disabled:opacity-60 disabled:cursor-not-allowed transition-colors min-h-[44px]"
      >
        {loading ? "登録中..." : "はじめる"}
      </button>
    </form>
  );
}
