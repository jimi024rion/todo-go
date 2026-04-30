import { AppShell } from "@/components/layout/app-shell"

export default function TodosLayout({ children }: { children: React.ReactNode }) {
  return <AppShell title="todos">{children}</AppShell>
}
