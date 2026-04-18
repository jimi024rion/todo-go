import { type Todo } from "@/types/todo";

const STATUS_CONFIG: Record<
  Todo["Status"],
  { label: string; bg: string; text: string; dot: string }
> = {
  pending: {
    label: "未着手",
    bg: "bg-[#FEF3C7]",
    text: "text-[#92400E]",
    dot: "bg-[#FCD34D]",
  },
  in_progress: {
    label: "進行中",
    bg: "bg-[#DBEAFE]",
    text: "text-[#1D4ED8]",
    dot: "bg-[#60A5FA]",
  },
  completed: {
    label: "完了",
    bg: "bg-[#D1FAE5]",
    text: "text-[#065F46]",
    dot: "bg-[#34D399]",
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
