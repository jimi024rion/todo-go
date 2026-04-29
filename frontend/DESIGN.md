---
version: "1.0.0"
name: "Todo"
description: "ミニマル・クリーンなTodo管理Webアプリ。スレート/ニュートラルカラーベースのライトモード固定デザイン。"

colors:
  background: "#f8fafc"
  surface: "#ffffff"
  surface-raised: "#f1f5f9"
  border: "#e2e8f0"
  border-subtle: "#f1f5f9"
  foreground: "#0f172b"
  foreground-secondary: "#45556c"
  foreground-muted: "#90a1b9"
  primary: "#314158"
  primary-foreground: "#ffffff"
  primary-hover: "#1d293d"
  destructive: "#ef4444"
  destructive-foreground: "#ffffff"
  success: "#22c55e"
  focus-ring: "#62748e"

typography:
  fontFamily:
    primary: "Inter, -apple-system, BlinkMacSystemFont, sans-serif"
  fontSize:
    xs: "12px"
    sm: "14px"
    base: "16px"
    lg: "18px"
    xl: "20px"
    2xl: "24px"
    3xl: "30px"
  fontWeight:
    regular: 400
    medium: 500
    semibold: 600
    bold: 700
  lineHeight:
    tight: 1.25
    normal: 1.5
    relaxed: 1.75

spacing:
  "1": "4px"
  "2": "8px"
  "3": "12px"
  "4": "16px"
  "5": "20px"
  "6": "24px"
  "8": "32px"
  "10": "40px"
  "12": "48px"
  "16": "64px"
  sidebar-width: "240px"
  content-max-width: "720px"
  card-padding: "16px"
  item-gap: "12px"

rounded:
  sm: "4px"
  md: "6px"
  lg: "8px"
  full: "9999px"

elevation:
  card:
    box-shadow: "0 1px 3px rgba(0, 0, 0, 0.06), 0 1px 2px rgba(0, 0, 0, 0.04)"
  panel:
    box-shadow: "0 4px 6px rgba(0, 0, 0, 0.07), 0 2px 4px rgba(0, 0, 0, 0.06)"
  modal:
    box-shadow: "0 20px 25px rgba(0, 0, 0, 0.1), 0 10px 10px rgba(0, 0, 0, 0.04)"

components:
  button-primary:
    backgroundColor: "{colors.primary}"
    color: "{colors.primary-foreground}"
    padding: "8px 16px"
    borderRadius: "{rounded.md}"
    fontSize: "{typography.fontSize.sm}"
    fontWeight: "{typography.fontWeight.medium}"
    hover:
      backgroundColor: "{colors.primary-hover}"
  button-ghost:
    backgroundColor: "transparent"
    color: "{colors.foreground-secondary}"
    padding: "8px 12px"
    borderRadius: "{rounded.md}"
    hover:
      backgroundColor: "{colors.surface-raised}"
  button-destructive:
    backgroundColor: "{colors.destructive}"
    color: "{colors.destructive-foreground}"
    padding: "8px 16px"
    borderRadius: "{rounded.md}"
  todo-card:
    backgroundColor: "{colors.surface}"
    borderRadius: "{rounded.lg}"
    padding: "{spacing.card-padding}"
    border: "1px solid {colors.border}"
    boxShadow: "{elevation.card.box-shadow}"
    gap: "{spacing.item-gap}"
    hover:
      boxShadow: "{elevation.panel.box-shadow}"
      borderColor: "#cad5e2"
  checkbox:
    size: "16px"
    borderRadius: "{rounded.sm}"
    checked:
      backgroundColor: "{colors.success}"
      borderColor: "{colors.success}"
  sidebar:
    width: "{spacing.sidebar-width}"
    backgroundColor: "{colors.background}"
    borderRight: "1px solid {colors.border}"
    padding: "16px 12px"
  badge:
    backgroundColor: "{colors.surface-raised}"
    color: "{colors.foreground-muted}"
    borderRadius: "{rounded.full}"
    padding: "2px 8px"
    fontSize: "{typography.fontSize.xs}"
    fontWeight: "{typography.fontWeight.medium}"
  input:
    backgroundColor: "{colors.surface}"
    border: "1px solid {colors.border}"
    borderRadius: "{rounded.md}"
    padding: "8px 12px"
    fontSize: "{typography.fontSize.sm}"
    focus:
      borderColor: "{colors.primary}"
      outline: "2px solid {colors.focus-ring}"
      outlineOffset: "2px"
---

## Overview

**Todo** はミニマル・クリーンな個人用タスク管理アプリ。
余白を広く取り、装飾を最小限に抑えた設計。
Linear や Notion にインスパイアされた洗練されたインターフェース。

**設計原則:**
- 余白は多めに、テキストを主役にする
- 色に頼らず、形と余白で階層を表現する
- アニメーションは控えめ（150ms, ease-out）
- アクセシビリティ: WCAG AA 準拠を目指す

## Colors

スレート/ニュートラルを基調としたモノクロベースのカラーパレット。
プライマリアクションも濃いスレートを使用し、鮮やかな色は使わない。
唯一の例外は完了状態の緑（`success`）と削除の赤（`destructive`）。

**使用ガイド:**
- `background`: ページ全体の背景（薄いスレート `#f8fafc`）
- `surface`: カード・モーダル・パネルの背景（白 `#ffffff`）
- `primary`: ボタン・アクション要素（濃いスレート `#314158`）
- カラフルなアクセントカラーは使用しない

## Typography

**Inter** を全体で使用。システムフォントへのフォールバックあり。
フォントサイズは実質4段階（`xs`, `sm`, `base`, `xl`）に絞る。

**使用ガイド:**
- ページタイトル: `xl` / `semibold`
- セクション見出し: `base` / `semibold`
- 本文・Todoタイトル: `sm` / `regular`
- サブテキスト・メタ情報: `xs` / `regular` + `foreground-muted`

## Layout

**8px グリッドシステム**。全スペーシング値は 4px または 8px の倍数。

**デスクトップ（768px以上）:** 左サイドバー 240px 固定 + メインコンテンツ
**モバイル（768px未満）:** ハンバーガーアイコン（☰）+ `Sheet` ドロワー

コンテンツ最大幅: 720px（中央寄せ）

## Elevation & Depth

影は3レベルのみ:
1. `card`: 通常カード（非常に薄い）
2. `panel`: ホバー時・編集パネル
3. `modal`: モーダル・ドロワー

ボーダーと影を組み合わせて奥行きを表現する。

## Shapes

`border-radius: 6px`（Tailwind `rounded-md`）をデフォルトとして全コンポーネントに統一。
カードのみ `8px`（`rounded-lg`）。バッジ・アバターは `rounded-full`。

## Components

### Todo カード
- ホワイト背景 + `border` 色のボーダー + `card` 影
- 左: チェックボックス（16px, `rounded-sm`）
- 中央: タイトル（`sm`/`regular`）+ 説明文先頭2行（`xs`, `foreground-muted`, `line-clamp-2`）
- 完了時: タイトルに打ち消し線 + 全テキストが `foreground-muted` に
- ホバー時: `panel` 影 + ボーダーが `slate-300` に
- クリックで編集パネル（Sheet）が開く

### チェックボックス
- 16×16px, `rounded-sm`
- 完了時: `success` 色で塗りつぶし + 白チェックマーク
- Tailwind の `peer` パターンで実装

### サイドバー（デスクトップ）
- 幅: 240px 固定, `background` 色, 右ボーダーあり
- 上部: アプリ名（`xl`/`semibold`）
- ナビ項目: アイコン + ラベル + 未完了数バッジ
  - アクティブ: `surface-raised` 背景 + `foreground` テキスト
  - 非アクティブ: 透明背景 + `foreground-secondary` テキスト
- 下部固定: ユーザーアバター + 名前

### インライン Todo 作成フォーム
- リストの最上部に常時表示
- 入力欄: `border` 色のボーダー, `rounded-md`
- プレースホルダー: 「タスクを追加...」
- Enter で追加 / Escape でクリア

### 編集パネル（Sheet）
- 右からスライドイン（幅 400px）
- タイトル・説明文: クリックで即インライン編集
- `blur` イベントで自動保存
- フッター: 「完了にする」ボタン / 「削除」ボタン（destructive）

### 空の状態
- 中央配置の SVG イラスト
- 見出し: 「タスクがありません」（`base`/`semibold`）
- サブテキスト: 「上の入力欄からタスクを追加してください」（`sm`, `foreground-muted`）

### 設定画面
- ユーザーカード: アバター + 名前 + メール（仮データ）
- 「ログアウト」ボタン（Firebase Auth 追加後に機能実装）

## Do's and Don'ts

**Do:**
- 余白を惜しまず使う（`padding: 16px`, `gap: 12px` を基本）
- テキストの階層は `fontWeight` と `color` のみで表現
- フォーカス状態を必ず視覚化（`outline: 2px solid {colors.focus-ring}`）
- トランジション: `transition-all duration-150 ease-out`

**Don't:**
- 鮮やかなアクセントカラー（ブルー・オレンジ・パープル等）を使わない
- `150ms` を超えるアニメーション
- フォントサイズを5種類以上使う
- ボックスシャドウを重ねがけする
