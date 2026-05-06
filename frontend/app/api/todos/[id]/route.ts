import { NextRequest, NextResponse } from "next/server"

const BACKEND = process.env.BACKEND_URL ?? "http://localhost:8080"

type Params = { params: Promise<{ id: string }> }

export async function GET(req: NextRequest, { params }: Params) {
  const { id } = await params
  const authorization = req.headers.get("Authorization") ?? ""
  const res = await fetch(`${BACKEND}/v1/todos/${id}`, {
    headers: { "Content-Type": "application/json", Authorization: authorization },
  })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}

export async function PUT(req: NextRequest, { params }: Params) {
  const { id } = await params
  const authorization = req.headers.get("Authorization") ?? ""
  const body = await req.json()
  const res = await fetch(`${BACKEND}/v1/todos/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json", Authorization: authorization },
    body: JSON.stringify(body),
  })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}

export async function DELETE(req: NextRequest, { params }: Params) {
  const { id } = await params
  const authorization = req.headers.get("Authorization") ?? ""
  const res = await fetch(`${BACKEND}/v1/todos/${id}`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json", Authorization: authorization },
  })
  if (res.status === 204) return new NextResponse(null, { status: 204 })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}
