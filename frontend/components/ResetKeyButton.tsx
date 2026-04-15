"use client";

import { useRouter } from "next/navigation";

// APIキーをリセット（Cookie削除）するボタン
// Client Componentにする理由：クリックイベントとrouter.refresh()が必要なため
export function ResetKeyButton() {
  const router = useRouter();

  const handleReset = async () => {
    await fetch("/api/auth", { method: "DELETE" });
    // Server Componentを再実行 → Cookie未設定 → ApiKeySetupが表示される
    router.refresh();
  };

  return (
    <button
      onClick={handleReset}
      className="rounded-md border border-[#E5E7EB] bg-white px-3 py-1.5 text-xs font-medium text-[#4B5563] transition hover:bg-[#F9FAFB] focus:outline-none focus:ring-2 focus:ring-[#4F46E5] focus:ring-offset-1"
    >
      API Key 変更
    </button>
  );
}
