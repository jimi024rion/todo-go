# Design System — Todo App

A clean, modern task management UI with a warm, calm aesthetic. The design philosophy prioritizes legibility and focus: every visual decision removes friction from the act of managing tasks.

## 1. Visual Theme & Atmosphere

The interface uses a warm off-white background with white card surfaces and a single orange accent. The palette is warm and calm — not clinical white, not flashy. Think of a well-crafted analog notebook: structured, purposeful, with just enough warmth to feel human.

Status colors are the only additional chromatic elements: amber for pending, blue for in-progress, and emerald for completed. Cards use a left border strip (4px) to communicate status at a glance without relying solely on badges.

Typography is Inter — leaning on weight and size contrast rather than font variety. All spacing uses an 8px grid.

**Key Characteristics:**
- Single accent color: Orange (`#F97316`) used exclusively for primary actions and focus states
- Warm off-white page background (`#F7F6F2`) — not pure white
- White card surfaces with subtle shadow
- Status palette: amber / blue / emerald — used in badges and card left borders
- Inter font, system-native fallback chain
- `#1C1917` near-black (stone-950) for text — warm, not cool
- 8px spacing base unit
- 6px button radius, 8px card radius, 12px modal radius

---

## 2. Color Palette & Roles

### Primary
- **Orange 500** (`#F97316`): Primary CTA, active tab, focus ring, interactive accent. The only saturated non-status color in the UI.
- **Orange 600** (`#EA580C`): Button hover / pressed state.
- **Orange 50** (`#FFF7ED`): Subtle orange tint for selected/hover backgrounds.

### Background & Surface
- **Page Background** (`#F7F6F2`): Warm off-white. The main canvas.
- **Surface / Card** (`#FFFFFF`): White cards on the warm background.

### Text
- **Stone 950** (`#1C1917`): Headings, primary body text.
- **Stone 500** (`#78716C`): Secondary text, descriptions, metadata.
- **Stone 400** (`#A8A29E`): Placeholder text, disabled states, empty states.

### Border
- **Stone 200** (`#E7E5E0`): Standard borders — cards, inputs, dividers.
- **Stone 300** (`#D6D3D1`): Stronger borders, focus-adjacent outlines.

### Status Colors

| Status | Badge BG | Badge Text | Card Left Border |
|--------|----------|------------|-----------------|
| `pending` | `#FEF3C7` amber-100 | `#92400E` amber-800 | `#FCD34D` amber-300 |
| `in_progress` | `#DBEAFE` blue-100 | `#1D4ED8` blue-700 | `#60A5FA` blue-400 |
| `completed` | `#D1FAE5` emerald-100 | `#065F46` emerald-800 | `#34D399` emerald-400 |

### Semantic
- **Red 500** (`#EF4444`): Destructive actions (delete button hover, error messages).
- **Red 50** (`#FEF2F2`): Destructive hover background.

### Shadows
- **Card Default**: `0 1px 2px rgba(0,0,0,0.06), 0 1px 3px rgba(0,0,0,0.04)`
- **Card Hover / Lifted**: `0 4px 6px -1px rgba(0,0,0,0.08), 0 2px 4px -2px rgba(0,0,0,0.06)`
- **Modal**: `0 20px 25px -5px rgba(0,0,0,0.12), 0 8px 10px -6px rgba(0,0,0,0.08)`

---

## 3. Typography Rules

### Font Family
- **Primary**: `Inter`, with fallbacks: `ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif`

### Hierarchy

| Role | Size | Weight | Line Height | Color |
|------|------|--------|-------------|-------|
| Page Title | 24px (1.5rem) | 700 | 1.25 | `#1C1917` |
| Section Heading | 18px (1.125rem) | 600 | 1.4 | `#1C1917` |
| Todo Title | 16px (1rem) | 500 | 1.5 | `#1C1917` |
| Body / Description | 14px (0.875rem) | 400 | 1.6 | `#78716C` |
| Label / Badge | 12px (0.75rem) | 500 | 1.33 | Varies by status |
| Caption / Meta | 12px (0.75rem) | 400 | 1.33 | `#A8A29E` |
| Button | 14px (0.875rem) | 500 | 1 | `#FFFFFF` / `#1C1917` |
| Input | 14px (0.875rem) | 400 | 1.5 | `#1C1917` |

### Principles
- **Weight over variety**: Use weight (400/500/600/700) to create hierarchy.
- **Completed todo treatment**: strike-through (`line-through`) + color shift to `#A8A29E`. Weight stays the same.

---

## 4. Component Stylings

### Buttons

**Primary (Orange)**
- Background: `#F97316`
- Hover: `#EA580C`
- Active: `#C2410C`
- Text: `#FFFFFF`, 14px weight 500
- Padding: `8px 16px`
- Radius: `6px`
- Focus: `2px solid #F97316` outline + `2px` offset

**Secondary (Outline)**
- Background: `#FFFFFF`
- Border: `1px solid #E7E5E0`
- Hover bg: `#F7F6F2`
- Text: `#1C1917`, 14px weight 500
- Padding: `8px 16px`
- Radius: `6px`

**Ghost (Icon Button)**
- Background: transparent
- Hover bg: `#F5F5F4` (stone-100)
- Icon color: `#A8A29E` → hover `#78716C`
- Padding: `6px`
- Radius: `6px`
- Use: Edit icon button, settings

**Destructive Ghost**
- Same as Ghost, but hover: bg `#FEF2F2`, icon `#EF4444`
- Use: Delete action

### Cards & Containers

**Todo Card**
- Background: `#FFFFFF`
- Border: `1px solid #E7E5E0`
- Border left: `4px solid <status-color>` (amber-300 / blue-400 / emerald-400)
- Radius: `8px`
- Padding: `16px`
- Shadow: `0 1px 2px rgba(0,0,0,0.06), 0 1px 3px rgba(0,0,0,0.04)`
- Hover: shadow lifted + border `#D6D3D1`
- Completed state: opacity-75, title with line-through

**Form Container**
- Background: `#FFFFFF`
- Border: `1px solid #E7E5E0`
- Shadow: same as Todo Card
- Radius: `8px`
- Padding: `20px`

**Register / Setup Card**
- Background: `#FFFFFF`
- Shadow: `0 10px 15px -3px rgba(0,0,0,0.08), 0 4px 6px -4px rgba(0,0,0,0.06)`
- Radius: `16px`
- Padding: `32px`
- Max-width: `400px`

### Inputs & Forms

**Text Input / Textarea**
- Background: `#FFFFFF`
- Border: `1px solid #D6D3D1`
- Radius: `6px`
- Padding: `8px 12px`
- Text: 14px `#1C1917`
- Placeholder: `#A8A29E`
- Focus: border `#F97316` + `0 0 0 2px rgba(249,115,22,0.15)` ring
- Disabled: bg `#F5F5F4`, text `#A8A29E`

### Status Badge

**Pill shape** (radius: `9999px`), padding: `2px 10px`

| Status | Background | Text | Dot |
|--------|-----------|------|-----|
| `pending` | `#FEF3C7` | `#92400E` | `#FCD34D` |
| `in_progress` | `#DBEAFE` | `#1D4ED8` | `#60A5FA` |
| `completed` | `#D1FAE5` | `#065F46` | `#34D399` |

Structure: `● label` — 6px dot (●) + space + text at 12px weight 500

### Filter Tabs

Horizontal tab bar below the page header. Tabs: **すべて** / **未着手** / **進行中**

- Container: `border-b border-stone-200 bg-[#F7F6F2]`
- Tab item padding: `8px 16px`
- **Active**: `border-b-2 border-orange-500 text-orange-600 font-semibold` — no bg change
- **Inactive**: `text-stone-500 hover:text-stone-700`
- Tab font: 14px weight 400 (500 when active)

### Edit Modal

Triggered by the edit icon button on a Todo Card.

- Overlay: `rgba(0,0,0,0.4)` backdrop, `backdrop-filter: blur(2px)`
- Panel: `#FFFFFF`, radius `12px`, padding `24px`, max-width `480px`
- Shadow: `0 20px 25px -5px rgba(0,0,0,0.12), 0 8px 10px -6px rgba(0,0,0,0.08)`
- Title: 18px weight 600 `#1C1917`
- Fields: title input + description textarea
- Buttons: [キャンセル] secondary + [保存] primary, right-aligned

### Completed Section

Appears below the active todo list when there are completed tasks.

- Header: `border-t border-stone-200 pt-4 mt-4`
- Toggle text: `text-stone-500 text-sm font-medium` + chevron icon
- Collapsed by default (show count in header)
- Inner cards: same card style but with `opacity-75`

### Divider
- `1px solid #E7E5E0`
- Use to separate sections

---

## 5. Layout Principles

### Spacing System (8px base)
- `4px` — tight spacing, badge padding vertical
- `8px` — internal component spacing
- `12px` — gap between form elements
- `16px` — card padding, section gap, tab padding
- `20px` — form container padding
- `24px` — modal padding, between cards in list
- `32px` — section vertical margin, register card padding
- `48px` — page top padding

### Page Layout
- Max content width: `640px` (single-column)
- Centered with horizontal auto margins
- Horizontal padding: `16px` on mobile, `24px` on tablet+
- Page background: `#F7F6F2`

### Component Arrangement
```
[Page Header: "Todo" title + ResetKey button]
─────────────────────────────────────────────
[Filter Tabs: すべて | 未着手 | 進行中]
─────────────────────────────────────────────
[TodoForm: input fields + submit button]
─────────────────────────────────────────────
[Active TodoList]
  [TodoCard] ← left border: status color
  [TodoCard]
  ...
─────────────────────────────────────────────
[完了済み (N) ▶]  ← collapsible section
  [TodoCard (completed)]
  ...
```

### Register Page Layout
```
[Full screen: bg #F7F6F2]
  [Centered card: white, 400px, rounded-2xl, shadow-lg]
    [App title]
    [RegisterForm: name + email + submit]
```

### Whitespace Philosophy
- List items have `8px` gap.
- Form and list separated by `24px` margin.
- Empty state: centered `#A8A29E` text, `48px` top margin.

### Border Radius Scale
- `4px` — small chips
- `6px` — buttons, inputs
- `8px` — cards
- `12px` — modal panel
- `16px` — register card
- `9999px` — status badges (full pill)

---

## 6. Depth & Elevation

| Level | Treatment | Use |
|-------|-----------|-----|
| Flat (0) | No shadow, border `#E7E5E0` | Inputs standard |
| Card (1) | shadow-sm | Todo cards, form containers |
| Lifted (2) | shadow-md | Card hover state |
| Modal (3) | shadow-xl | Edit modal, register card |
| Focus | `0 0 0 2px rgba(249,115,22,0.15)` + border `#F97316` | Keyboard focus on inputs |

**Depth Philosophy**: Depth appears on interactive hover and attention-demanding elements (modals, registration). Static cards have a subtle shadow on the warm background — the contrast itself creates separation without heavy borders.

---

## 7. Design Guidelines

### Do
- Use Orange (`#F97316`) exclusively for primary actions, active tabs, and focus rings.
- Use the 4px left border on cards to signal status — it's faster to scan than badges alone.
- Apply `line-through` + `#A8A29E` text color to completed todo titles.
- Show loading state inline (spinner replacing button label).
- Collapse the completed section by default — completed tasks are history, not current work.
- Minimum touch target height `44px` for all interactive elements on mobile.

### Don't
- Don't use Indigo or any blue except for `in_progress` status.
- Don't use pure white (`#FFFFFF`) as the page background — use `#F7F6F2`.
- Don't use more than one primary (orange) button per visible screen section.
- Don't add icons to text buttons — icon-only ghost buttons only.
- Don't use red for anything other than destructive actions.
- Don't show skeleton loaders — use a single centered spinner.
- Don't truncate todo titles with ellipsis — allow wrapping.
- Don't change card width or layout on status change.

---

## 8. Responsive Behavior

### Breakpoints
| Name | Width | Key Changes |
|------|-------|-------------|
| Mobile | < 640px | Full-width cards, 16px horizontal padding, stacked form fields |
| Tablet+ | ≥ 640px | Centered 640px column, 24px horizontal padding |

### Collapsing Strategy
- Single-column always — no column collapsing needed.
- Form: title and description fields stack vertically; submit button full-width on mobile.
- Todo card action buttons: always visible.
- Edit modal: full-width with 16px margin on mobile; centered 480px panel on tablet+.

### Touch Targets
- All buttons minimum `44px` height on mobile (use `min-h-[44px]`)
- Edit and delete icon buttons: `40px × 40px` touch area
- Filter tabs: `44px` height on mobile

---

## 9. Agent Prompt Guide

### Quick Color Reference
```
Primary CTA:        #F97316  (Orange 500)
CTA Hover:          #EA580C  (Orange 600)
Page background:    #F7F6F2  (Warm off-white)
Card surface:       #FFFFFF
Primary text:       #1C1917  (Stone 950)
Secondary text:     #78716C  (Stone 500)
Muted text:         #A8A29E  (Stone 400)
Border:             #E7E5E0  (Stone 200)
Strong border:      #D6D3D1  (Stone 300)
Focus ring:         rgba(249,115,22,0.15)

Status — Pending:     badge bg #FEF3C7  text #92400E  border #FCD34D
Status — In Progress: badge bg #DBEAFE  text #1D4ED8  border #60A5FA
Status — Completed:   badge bg #D1FAE5  text #065F46  border #34D399
Destructive:          #EF4444  hover-bg #FEF2F2
```

### Example Component Prompts

- **Todo Card**: `White bg, border-l-4 (status color), 1px solid #E7E5E0 border, 8px radius, 16px padding, shadow-sm. Title: 16px weight 500 #1C1917. Description: 14px #78716C. Edit: ghost icon. Delete: destructive ghost icon. Status badge: pill.`
- **Primary Button**: `#F97316 bg, white text, 14px weight 500, 6px radius, 8px 16px padding. Hover: #EA580C. Focus: 2px orange outline + rgba(249,115,22,0.15) ring.`
- **Text Input**: `White bg, 1px solid #D6D3D1, 6px radius, 8px 12px padding, 14px #1C1917, placeholder #A8A29E. Focus: border #F97316 + orange ring.`
- **Filter Tab (active)**: `border-b-2 border-orange-500 text-orange-600 font-semibold, 8px 16px padding.`
- **Edit Modal**: `White panel, 12px radius, 24px padding, max-width 480px. Backdrop: rgba(0,0,0,0.4) + blur(2px). Shadow: shadow-xl.`
- **Register Card**: `White, 16px radius, 32px padding, max-width 400px, shadow-lg on #F7F6F2 background.`

### Iteration Guide
1. Single-column always. Resist sidebars.
2. Orange is the only saturated non-status color. If you reach for another, reconsider.
3. Status badges are pills (`9999px`). Buttons/inputs are `6px`. Cards are `8px`. Modal is `12px`. Register card is `16px`. Don't mix these.
4. Completed todos get `line-through` + `#A8A29E`. That's the only style change for completion.
5. Hover states on cards add a lifted shadow. The border subtly darkens.
6. Keep transitions at 150ms. Snappy feedback over smooth animation.
7. Empty state: centered `#A8A29E` text at 14px, `48px` top margin. No illustrations.
8. The left `4px` border on cards is always the status color — it's the fastest visual signal.
