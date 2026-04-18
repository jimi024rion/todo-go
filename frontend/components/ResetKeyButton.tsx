"use client";

import { useRouter } from "next/navigation";

export function ResetKeyButton() {
  const router = useRouter();

  const handleReset = async () => {
    await fetch("/api/auth", { method: "DELETE" });
    router.push("/register");
    router.refresh();
  };

  return (
    <button
      onClick={handleReset}
      className="rounded-[6px] border border-[#E7E5E0] bg-white px-3 py-1.5 text-xs font-medium text-[#78716C] transition hover:bg-[#F7F6F2] focus:outline-none focus:ring-2 focus:ring-[#F97316] focus:ring-offset-1"
    >
      ログアウト
    </button>
  );
}
