import { TodoList } from "@/components/todo/todo-list"

export default function TodosPage() {
  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <h1 className="mb-6 text-xl font-semibold text-foreground">Todos</h1>
      <TodoList />
    </div>
  )
}
