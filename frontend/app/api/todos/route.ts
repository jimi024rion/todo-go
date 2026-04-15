import { type NextRequest, NextResponse } from "next/server";

// バックエンドのURLは環境変数から読む（ブラウザには公開されない）
const BACKEND = process.env.BACKEND_URL ?? "http://localhost:8080";

// リクエストの Cookie から APIキーを取り出してバックエンド向けのヘッダーを作る
// Goのミドルウェアで言う「ヘッダーの変換処理」に相当する
function backendHeaders(req: NextRequest): HeadersInit {
  const apiKey = req.cookies.get("api_key")?.value ?? "";
  return {
    "Content-Type": "application/json",
    "X-API-Key": apiKey,
  };
}

// GET /api/todos → GET /v1/todos（バックエンドに転送）
export async function GET(req: NextRequest) {
  const res = await fetch(`${BACKEND}/v1/todos`, {
    headers: backendHeaders(req),
  });
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}

// POST /api/todos → POST /v1/todos（バックエンドに転送）
export async function POST(req: NextRequest) {
  const body = await req.json();
  const res = await fetch(`${BACKEND}/v1/todos`, {
    method: "POST",
    headers: backendHeaders(req),
    body: JSON.stringify(body),
  });
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
