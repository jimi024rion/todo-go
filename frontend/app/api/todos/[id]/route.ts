import { NextRequest, NextResponse } from "next/server"

const BACKEND = process.env.BACKEND_URL ?? "http://localhost:8080"
const API_KEY = process.env.BACKEND_API_KEY ?? ""

const backendHeaders: HeadersInit = {
  "Content-Type": "application/json",
  "X-API-Key": API_KEY,
}

type Params = { params: Promise<{ id: string }> }

export async function GET(_req: NextRequest, { params }: Params) {
  const { id } = await params
  const res = await fetch(`${BACKEND}/v1/todos/${id}`, { headers: backendHeaders })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}

export async function PUT(req: NextRequest, { params }: Params) {
  const { id } = await params
  const body = await req.json()
  const res = await fetch(`${BACKEND}/v1/todos/${id}`, {
    method: "PUT",
    headers: backendHeaders,
    body: JSON.stringify(body),
  })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}

export async function DELETE(_req: NextRequest, { params }: Params) {
  const { id } = await params
  const res = await fetch(`${BACKEND}/v1/todos/${id}`, {
    method: "DELETE",
    headers: backendHeaders,
  })
  if (res.status === 204) return new NextResponse(null, { status: 204 })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}
