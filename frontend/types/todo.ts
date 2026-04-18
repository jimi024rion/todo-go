// バックエンドのTodoResponseはPascalCase（jsonタグに合わせる）
export type Todo = {
  ID: string;
  Title: string;
  Description: string;
  Status: "pending" | "in_progress" | "completed";
  CreatedAt: string;
  UpdatedAt: string;
};

export type CreateTodoRequest = {
  title: string;
  description?: string;
};

export type UpdateTodoRequest = {
  title?: string;
  description?: string;
  status?: Todo["Status"];
};

export type User = {
  id: string;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
};

export type APIKey = {
  id: string;
  key: string;
  user_id: string;
  name: string;
  created_at: string;
};
