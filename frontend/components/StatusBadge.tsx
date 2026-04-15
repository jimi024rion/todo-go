import { type Todo } from "@/types/todo";

const STATUS_CONFIG: Record<
  Todo["Status"],
  { label: string; bg: string; text: string; dot: string }
> = {
  pending: {
    label: "未着手",
    bg: "bg-[#F3F4F6]",
    text: "text-[#4B5563]",
    dot: "bg-[#9CA3AF]",
  },
  in_progress: {
    label: "進行中",
    bg: "bg-[#DBEAFE]",
    text: "text-[#1D4ED8]",
    dot: "bg-[#3B82F6]",
  },
  completed: {
    label: "完了",
    bg: "bg-[#DCFCE7]",
    text: "text-[#15803D]",
    dot: "bg-[#22C55E]",
  },
};

type Props = {
  status: Todo["Status"];
};

export function StatusBadge({ status }: Props) {
  const config = STATUS_CONFIG[status];
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium ${config.bg} ${config.text}`}
    >
      <span className={`h-1.5 w-1.5 rounded-full ${config.dot}`} />
      {config.label}
    </span>
  );
}
