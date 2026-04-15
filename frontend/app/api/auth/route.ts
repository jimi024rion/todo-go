import { NextResponse } from "next/server";

// POST /api/auth
// APIキーを受け取り httpOnly Cookie にセットする
// httpOnly = true にすることで JavaScript からアクセス不可能になる（XSS対策）
export async function POST(req: Request) {
  const body = await req.json().catch(() => ({}));
  const key = typeof body.key === "string" ? body.key.trim() : "";

  if (!key) {
    return NextResponse.json({ error: "key is required" }, { status: 400 });
  }

  const res = NextResponse.json({ ok: true });
  res.cookies.set("api_key", key, {
    httpOnly: true,   // JavaScriptの document.cookie からは見えない
    sameSite: "lax",  // CSRF対策（同一サイトのリクエストのみCookieを送信）
    path: "/",
    // secure: true は本番環境（HTTPS）では必須。開発環境（HTTP）では外す
    secure: process.env.NODE_ENV === "production",
  });
  return res;
}

// DELETE /api/auth
// Cookie を削除する（APIキーの変更・リセット時に使用）
export async function DELETE() {
  const res = NextResponse.json({ ok: true });
  res.cookies.delete("api_key");
  return res;
}
