import { NextRequest, NextResponse } from "next/server"

const BACKEND = process.env.BACKEND_URL ?? "http://localhost:8080"
const API_KEY = process.env.BACKEND_API_KEY ?? ""
const headers: HeadersInit = { "Content-Type": "application/json", "X-API-Key": API_KEY }

export async function GET() {
  const res = await fetch(`${BACKEND}/v1/tags`, { headers })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}

export async function POST(req: NextRequest) {
  const body = await req.json()
  const res = await fetch(`${BACKEND}/v1/tags`, { method: "POST", headers, body: JSON.stringify(body) })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}
