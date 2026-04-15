import { type NextRequest, NextResponse } from "next/server";

const BACKEND = process.env.BACKEND_URL ?? "http://localhost:8080";

function backendHeaders(req: NextRequest): HeadersInit {
  const apiKey = req.cookies.get("api_key")?.value ?? "";
  return {
    "Content-Type": "application/json",
    "X-API-Key": apiKey,
  };
}

type Params = { params: Promise<{ id: string }> };

// GET /api/todos/:id → GET /v1/todos/:id
export async function GET(req: NextRequest, { params }: Params) {
  const { id } = await params;
  const res = await fetch(`${BACKEND}/v1/todos/${id}`, {
    headers: backendHeaders(req),
  });
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}

// PUT /api/todos/:id → PUT /v1/todos/:id
export async function PUT(req: NextRequest, { params }: Params) {
  const { id } = await params;
  const body = await req.json();
  const res = await fetch(`${BACKEND}/v1/todos/${id}`, {
    method: "PUT",
    headers: backendHeaders(req),
    body: JSON.stringify(body),
  });
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}

// DELETE /api/todos/:id → DELETE /v1/todos/:id
export async function DELETE(req: NextRequest, { params }: Params) {
  const { id } = await params;
  const res = await fetch(`${BACKEND}/v1/todos/${id}`, {
    method: "DELETE",
    headers: backendHeaders(req),
  });
  if (res.status === 204) {
    return new NextResponse(null, { status: 204 });
  }
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
