import { NextRequest, NextResponse } from "next/server"

const BACKEND = process.env.BACKEND_URL ?? "http://localhost:8080"
const API_KEY = process.env.BACKEND_API_KEY ?? ""

const backendHeaders: HeadersInit = {
  "Content-Type": "application/json",
  "X-API-Key": API_KEY,
}

export async function GET() {
  const res = await fetch(`${BACKEND}/v1/todos`, { headers: backendHeaders })
  if (!res.ok) {
    const data = await res.json().catch(() => ({}))
    return NextResponse.json({ error: data.message ?? res.statusText }, { status: res.status })
  }
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}

export async function POST(req: NextRequest) {
  const body = await req.json()
  const res = await fetch(`${BACKEND}/v1/todos`, {
    method: "POST",
    headers: backendHeaders,
    body: JSON.stringify(body),
  })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}
