import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { ResetKeyButton } from "@/components/ResetKeyButton";
import { TodoList } from "@/components/TodoList";
import { type Todo } from "@/types/todo";

async function fetchInitialTodos(apiKey: string): Promise<Todo[]> {
  try {
    const res = await fetch(
      `${process.env.BACKEND_URL ?? "http://localhost:8080"}/v1/todos`,
      {
        headers: {
          "Content-Type": "application/json",
          "X-API-Key": apiKey,
        },
        cache: "no-store",
      },
    );
    if (!res.ok) return [];
    const data = await res.json();
    return data.resultBody ?? [];
  } catch {
    return [];
  }
}

export default async function Home() {
  const apiKey = (await cookies()).get("api_key")?.value;

  if (!apiKey) {
    redirect("/register");
  }

  const initialTodos = await fetchInitialTodos(apiKey);

  return (
    <div className="min-h-screen bg-[#F7F6F2] px-4 py-10 sm:px-6">
      <div className="mx-auto max-w-xl">
        <header className="mb-6 flex items-center justify-between">
          <h1 className="text-2xl font-bold tracking-tight text-[#1C1917]">Todo</h1>
          <ResetKeyButton />
        </header>

        <TodoList initialTodos={initialTodos} />
      </div>
    </div>
  );
}
