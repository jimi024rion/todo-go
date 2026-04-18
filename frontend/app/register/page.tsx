import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import RegisterForm from "@/components/RegisterForm";

export default async function RegisterPage() {
  const cookieStore = await cookies();
  if (cookieStore.get("api_key")) {
    redirect("/");
  }

  return (
    <main className="min-h-screen bg-[#F7F6F2] flex items-center justify-center px-4">
      <div className="w-full max-w-[400px] bg-white rounded-2xl shadow-lg p-8">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold text-[#1C1917] tracking-tight">Todo</h1>
          <p className="mt-2 text-sm text-[#78716C]">
            アカウントを作成してはじめましょう
          </p>
        </div>
        <RegisterForm />
      </div>
    </main>
  );
}
