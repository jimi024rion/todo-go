# Design System — Todo App

A clean, focused task management UI. The design philosophy prioritizes legibility and calm focus: every visual decision removes friction from the act of managing tasks, not adding to it.

## 1. Visual Theme & Atmosphere

The interface is built on pure white surfaces with a single indigo accent — there are no competing colors, no gradients, no decorative noise. The aesthetic is closer to a well-designed notebook than a feature-rich dashboard. Surfaces feel flat and airy; interactions are signaled through subtle shadows and understated color transitions rather than bold animations.

Status colors are the only chromatic elements allowed: a muted slate for pending tasks, a calm blue for in-progress, and a soft green for completed. These colors are desaturated enough to read as data, not decoration.

Typography is system-native (`Inter`), leaning on weight and size contrast rather than font variety. Body text is dark but not pure black — `#111827` — giving a slightly warmer, less clinical feel on white. All spacing uses an 8px grid, creating consistent visual rhythm without rigidity.

**Key Characteristics:**
- Single accent color: Indigo (`#4F46E5`) used exclusively for primary actions and focus states
- Flat surfaces: no decorative shadows, depth only on interactive overlays and modals
- Status palette: slate / blue / green — muted, not saturated
- Inter font, system-native fallback chain
- `#111827` near-black text instead of pure black
- 8px spacing base unit
- 6px button radius for a slightly rounded but functional feel
- Borders: `1px solid #E5E7EB` — light gray, not whisper-weight

---

## 2. Color Palette & Roles

### Primary
- **Indigo 600** (`#4F46E5`): Primary CTA, active state, focus ring, interactive accent. The only saturated color in the UI.
- **Indigo 700** (`#4338CA`): Button hover / pressed state.
- **Indigo 50** (`#EEF2FF`): Subtle indigo tint for selected backgrounds.
- **White** (`#FFFFFF`): Page background, card surfaces, button text on indigo.

### Text
- **Gray 900** (`#111827`): Headings, primary body text.
- **Gray 600** (`#4B5563`): Secondary text, descriptions, metadata.
- **Gray 400** (`#9CA3AF`): Placeholder text, disabled states, empty states.

### Surface & Border
- **Gray 50** (`#F9FAFB`): Page background tint, input backgrounds.
- **Gray 100** (`#F3F4F6`): Hover backgrounds on list items, subtle section fills.
- **Gray 200** (`#E5E7EB`): Standard borders — cards, inputs, dividers.
- **Gray 300** (`#D1D5DB`): Stronger borders, focus-adjacent outlines.

### Status Colors
- **Pending** — Background: `#F3F4F6`, Text: `#4B5563`, Dot: `#9CA3AF`
- **In Progress** — Background: `#DBEAFE`, Text: `#1D4ED8`, Dot: `#3B82F6`
- **Completed** — Background: `#DCFCE7`, Text: `#15803D`, Dot: `#22C55E`

### Semantic
- **Red 500** (`#EF4444`): Destructive actions (delete button hover, error messages).
- **Red 50** (`#FEF2F2`): Destructive hover background.

### Shadows
- **Overlay Shadow** (`0 4px 6px -1px rgba(0,0,0,0.1), 0 2px 4px -2px rgba(0,0,0,0.1)`): Modals, dropdowns.
- **Card Lift** (`0 1px 3px rgba(0,0,0,0.07), 0 1px 2px rgba(0,0,0,0.06)`): Hover state on cards.

---

## 3. Typography Rules

### Font Family
- **Primary**: `Inter`, with fallbacks: `ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif`

### Hierarchy

| Role | Size | Weight | Line Height | Letter Spacing | Color |
|------|------|--------|-------------|----------------|-------|
| Page Title | 24px (1.5rem) | 700 | 1.25 | -0.025em | `#111827` |
| Section Heading | 18px (1.125rem) | 600 | 1.4 | -0.015em | `#111827` |
| Todo Title | 16px (1rem) | 500 | 1.5 | normal | `#111827` |
| Body / Description | 14px (0.875rem) | 400 | 1.6 | normal | `#4B5563` |
| Label / Badge | 12px (0.75rem) | 500 | 1.33 | 0.025em | Varies by status |
| Caption / Meta | 12px (0.75rem) | 400 | 1.33 | normal | `#9CA3AF` |
| Button | 14px (0.875rem) | 500 | 1 | normal | `#FFFFFF` / `#111827` |
| Input | 14px (0.875rem) | 400 | 1.5 | normal | `#111827` |

### Principles
- **Weight over variety**: Use weight (400/500/600/700) to create hierarchy. Never use more than two font sizes in a single component.
- **Negative tracking on headings only**: Page title at -0.025em, section heading at -0.015em. Body text always at normal tracking.
- **Completed todo treatment**: strike-through (`line-through`) + color shift to `#9CA3AF`. Weight stays the same.

---

## 4. Component Stylings

### Buttons

**Primary (Indigo)**
- Background: `#4F46E5`
- Hover: `#4338CA`
- Active: `#3730A3`
- Text: `#FFFFFF`, 14px weight 500
- Padding: `8px 16px`
- Radius: `6px`
- Focus: `2px solid #4F46E5` outline + `2px` offset

**Secondary (Outline)**
- Background: `#FFFFFF`
- Border: `1px solid #E5E7EB`
- Hover bg: `#F9FAFB`
- Text: `#374151`, 14px weight 500
- Padding: `8px 16px`
- Radius: `6px`

**Ghost (Icon Button)**
- Background: transparent
- Hover bg: `#F3F4F6`
- Icon color: `#9CA3AF` → hover `#4B5563`
- Padding: `6px`
- Radius: `6px`
- Use: Delete icon button, settings

**Destructive Ghost**
- Same as Ghost, but hover: bg `#FEF2F2`, icon `#EF4444`
- Use: Delete action in confirmed context

### Cards & Containers

**Todo Card**
- Background: `#FFFFFF`
- Border: `1px solid #E5E7EB`
- Radius: `8px`
- Padding: `16px`
- Hover: border `#D1D5DB` + shadow `0 1px 3px rgba(0,0,0,0.07), 0 1px 2px rgba(0,0,0,0.06)`
- Completed state: bg `#F9FAFB`

**Form Container**
- Background: `#FFFFFF`
- Border: `1px solid #E5E7EB`
- Radius: `8px`
- Padding: `20px`

### Inputs & Forms

**Text Input / Textarea**
- Background: `#FFFFFF`
- Border: `1px solid #D1D5DB`
- Radius: `6px`
- Padding: `8px 12px`
- Text: 14px `#111827`
- Placeholder: `#9CA3AF`
- Focus: border `#4F46E5` + `0 0 0 2px rgba(79,70,229,0.15)` ring
- Disabled: bg `#F9FAFB`, text `#9CA3AF`

### Status Badge

**Pill shape** (radius: `9999px`), padding: `2px 10px`

| Status | Background | Text | Dot |
|--------|-----------|------|-----|
| `pending` | `#F3F4F6` | `#4B5563` | `#9CA3AF` |
| `in_progress` | `#DBEAFE` | `#1D4ED8` | `#3B82F6` |
| `completed` | `#DCFCE7` | `#15803D` | `#22C55E` |

Structure: `● label` — 6px dot (●) + space + text at 12px weight 500

### Modal (API Key Setup)

- Overlay: `rgba(0,0,0,0.4)` backdrop, `backdrop-filter: blur(2px)`
- Panel: `#FFFFFF`, radius `12px`, padding `32px`, max-width `400px`
- Shadow: `0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1)`
- Title: 20px weight 700 `#111827`
- Description: 14px weight 400 `#4B5563`

### Divider
- `1px solid #F3F4F6`
- Use to separate todo list from form, or group sections

---

## 5. Layout Principles

### Spacing System (8px base)
- `4px` — tight spacing, icon gaps, badge padding vertical
- `8px` — component-level internal spacing (input padding)
- `12px` — card padding tight, gap between form elements
- `16px` — card padding standard, section gap
- `20px` — form container padding
- `24px` — between cards in list
- `32px` — section vertical margin, modal padding
- `48px` — page top padding

### Page Layout
- Max content width: `640px` (single-column, task-focused)
- Centered with horizontal auto margins
- Horizontal padding: `16px` on mobile, `24px` on tablet+
- No sidebars, no multi-column — task management is a focused activity

### Component Arrangement
```
[Page Header: "Todo" title + API key reset button]
─────────────────────────────────────────
[TodoForm: input fields + submit button]
─────────────────────────────────────────
[TodoList]
  [TodoCard]
  [TodoCard]
  ...
```

### Whitespace Philosophy
- List items have `8px` gap — dense enough to see many tasks, loose enough to distinguish them.
- The form and list are separated by a `32px` margin.
- Empty state uses center-aligned text with `48px` top margin.

### Border Radius Scale
- `4px` — small chips, tight elements
- `6px` — buttons, inputs
- `8px` — cards, containers
- `12px` — modal panel
- `9999px` — status badges (full pill)

---

## 6. Depth & Elevation

| Level | Treatment | Use |
|-------|-----------|-----|
| Flat (0) | No shadow, `1px solid #E5E7EB` border | Page bg, standard cards, inputs |
| Lifted (1) | `0 1px 3px rgba(0,0,0,0.07), 0 1px 2px rgba(0,0,0,0.06)` | Card hover state |
| Overlay (2) | `0 4px 6px -1px rgba(0,0,0,0.1), 0 2px 4px -2px rgba(0,0,0,0.1)` | Dropdowns |
| Modal (3) | `0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1)` | API key modal |
| Focus | `0 0 0 2px rgba(79,70,229,0.15)` + border `#4F46E5` | Keyboard focus on inputs |

**Depth Philosophy**: Depth only appears when the user is interacting (hover) or when an element demands attention (modal). Static cards are flat — the border alone defines their boundary. This keeps the list scannable without visual noise.

---

## 7. Design Guidelines

### Do
- Use Indigo (`#4F46E5`) exclusively for primary actions. Never apply it to borders, text, or decorative elements.
- Apply `line-through` + `#9CA3AF` text color to completed todo titles for immediate status recognition.
- Show loading state inline (spinner icon replacing the button label) rather than disabling the whole card.
- Use the status badge in every todo card — it's the primary information hierarchy element.
- Animate only transitions: 150ms ease for color/shadow/border changes. No transform animations on list items.

### Don't
- Don't use more than one primary button per visible screen section.
- Don't add icons to buttons unless they are icon-only ghost buttons. Text buttons stay text-only.
- Don't use red for anything other than destructive confirmation (delete). Errors use red text, not red backgrounds.
- Don't show skeleton loaders for the todo list — use a single centered spinner instead. The list is short.
- Don't truncate todo titles with ellipsis — allow wrapping. Tasks need to be fully readable at a glance.
- Don't change card width or layout on status change. The list should feel stable.

---

## 8. Responsive Behavior

### Breakpoints
| Name | Width | Key Changes |
|------|-------|-------------|
| Mobile | < 640px | Full-width cards, 16px horizontal padding, stacked form fields |
| Tablet+ | ≥ 640px | Centered 640px column, 24px horizontal padding |

### Collapsing Strategy
- The layout is already single-column — no column collapsing needed.
- Form: on mobile, title and description fields stack vertically with the submit button full-width below.
- Todo card action buttons (status change, delete): always visible, never hidden behind a "more" menu — the list is compact enough.
- API key modal: full-width with 16px horizontal margin on mobile; centered 400px panel on tablet+.

### Touch Targets
- All buttons minimum `36px` height
- Delete icon button: `32px × 32px` touch area
- Status change button: `32px` height minimum

---

## 9. Agent Prompt Guide

### Quick Color Reference
```
Primary CTA:        #4F46E5  (Indigo 600)
CTA Hover:          #4338CA  (Indigo 700)
Page background:    #FFFFFF
Alt background:     #F9FAFB  (Gray 50)
Primary text:       #111827  (Gray 900)
Secondary text:     #4B5563  (Gray 600)
Muted text:         #9CA3AF  (Gray 400)
Border:             #E5E7EB  (Gray 200)
Strong border:      #D1D5DB  (Gray 300)
Focus ring:         rgba(79,70,229,0.15)

Status — Pending:     bg #F3F4F6  text #4B5563
Status — In Progress: bg #DBEAFE  text #1D4ED8
Status — Completed:   bg #DCFCE7  text #15803D
Destructive:          #EF4444  hover-bg #FEF2F2
```

### Example Component Prompts

- **Todo Card**: `White bg, 1px solid #E5E7EB border, 8px radius, 16px padding. Title: 16px Inter weight 500, color #111827. Description: 14px weight 400, color #4B5563. Status badge: pill 9999px, 12px weight 500, colors per status table. Delete: ghost icon button, hover bg #FEF2F2, icon #EF4444.`
- **Primary Button**: `#4F46E5 bg, white text, 14px weight 500, 6px radius, 8px 16px padding. Hover: #4338CA. Focus: 2px solid #4F46E5 outline + rgba(79,70,229,0.15) ring.`
- **Text Input**: `White bg, 1px solid #D1D5DB border, 6px radius, 8px 12px padding, 14px #111827 text, placeholder #9CA3AF. Focus: border #4F46E5 + 0 0 0 2px rgba(79,70,229,0.15).`
- **Status Badge**: `Pill (9999px), 2px 10px padding, 12px Inter weight 500. Pending: bg #F3F4F6 text #4B5563. In Progress: bg #DBEAFE text #1D4ED8. Completed: bg #DCFCE7 text #15803D. Prefix with colored dot ●.`
- **API Key Modal**: `White panel, 12px radius, 32px padding, max-width 400px. Backdrop: rgba(0,0,0,0.4) + blur(2px). Shadow: 0 20px 25px -5px rgba(0,0,0,0.1). Title: 20px weight 700 #111827.`

### Iteration Guide
1. Single-column always — this is a focused task manager, not a dashboard. Resist the urge to add sidebars.
2. Indigo is the only saturated color. If you're reaching for another bright color, use a status color or reconsider.
3. Status badges are pills (`9999px`). Buttons and inputs are `6px`. Cards are `8px`. Modal is `12px`. Don't mix these.
4. Completed todos get `line-through` + `#9CA3AF`. That's the only style change for completion — no fade, no collapse.
5. Hover states on cards add a subtle shadow. The border doesn't change color — the shadow signals interactivity.
6. Keep transitions at 150ms. This app is used repeatedly; snappy feedback matters more than smooth animation.
7. Empty state: centered `#9CA3AF` text at 14px, `48px` top margin. No illustrations, no calls-to-action.
