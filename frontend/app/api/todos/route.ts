import { NextRequest, NextResponse } from "next/server"

const BACKEND = process.env.BACKEND_URL ?? "http://localhost:8080"

export async function GET(req: NextRequest) {
  const authorization = req.headers.get("Authorization") ?? ""
  const res = await fetch(`${BACKEND}/v1/todos`, {
    headers: { "Content-Type": "application/json", Authorization: authorization },
  })
  if (!res.ok) {
    const data = await res.json().catch(() => ({}))
    return NextResponse.json({ error: data.message ?? res.statusText }, { status: res.status })
  }
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}

export async function POST(req: NextRequest) {
  const authorization = req.headers.get("Authorization") ?? ""
  const body = await req.json()
  const res = await fetch(`${BACKEND}/v1/todos`, {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: authorization },
    body: JSON.stringify(body),
  })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}
