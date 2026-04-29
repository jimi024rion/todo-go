export function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center py-24 text-center">
      {/* SVG イラスト: チェックリストのシルエット */}
      <svg
        className="mb-6 h-24 w-24 text-muted-foreground/30"
        fill="none"
        viewBox="0 0 96 96"
        xmlns="http://www.w3.org/2000/svg"
        aria-hidden="true"
      >
        <rect x="16" y="12" width="64" height="72" rx="6" stroke="currentColor" strokeWidth="3" />
        <line x1="28" y1="34" x2="68" y2="34" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
        <line x1="28" y1="48" x2="68" y2="48" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
        <line x1="28" y1="62" x2="52" y2="62" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
        <circle cx="22" cy="34" r="4" fill="currentColor" />
        <circle cx="22" cy="48" r="4" fill="currentColor" />
        <circle cx="22" cy="62" r="4" fill="currentColor" />
      </svg>

      <p className="text-base font-semibold text-foreground">タスクがありません</p>
      <p className="mt-1 text-sm text-muted-foreground">
        上の入力欄からタスクを追加してください
      </p>
    </div>
  )
}
