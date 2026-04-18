import { cookies } from "next/headers";
import { NextResponse } from "next/server";
import { createAPIKey, createUser } from "@/lib/api";

const BACKEND_URL = process.env.BACKEND_URL ?? "http://localhost:8080";

export async function POST(req: Request) {
  const { name, email } = await req.json();

  if (!name || !email) {
    return NextResponse.json({ error: "名前とメールアドレスは必須です" }, { status: 400 });
  }

  let user;
  try {
    user = await createUser(name, email, BACKEND_URL);
  } catch (e) {
    return NextResponse.json({ error: "ユーザーの作成に失敗しました" }, { status: 400 });
  }

  let apiKey;
  try {
    apiKey = await createAPIKey(user.id, "default", BACKEND_URL);
  } catch (e) {
    return NextResponse.json({ error: "APIキーの発行に失敗しました" }, { status: 500 });
  }

  const cookieStore = await cookies();
  cookieStore.set("api_key", apiKey.key, {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
  });

  return NextResponse.json({ ok: true });
}
