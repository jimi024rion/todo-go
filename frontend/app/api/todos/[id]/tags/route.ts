import { NextRequest, NextResponse } from "next/server"

const BACKEND = process.env.BACKEND_URL ?? "http://localhost:8080"
const API_KEY = process.env.BACKEND_API_KEY ?? ""
const headers: HeadersInit = { "Content-Type": "application/json", "X-API-Key": API_KEY }

type Params = { params: Promise<{ id: string }> }

export async function PUT(req: NextRequest, { params }: Params) {
  const { id } = await params
  const body = await req.json()
  const res = await fetch(`${BACKEND}/v1/todos/${id}/tags`, {
    method: "PUT",
    headers,
    body: JSON.stringify(body),
  })
  if (res.status === 204) return new NextResponse(null, { status: 204 })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}
